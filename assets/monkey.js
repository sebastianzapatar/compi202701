/**
 * Monkey en el navegador — réplica didáctica en JavaScript del intérprete
 * de "Writing An Interpreter In Go" (Thorsten Ball). Misma estructura que el
 * código Go del libro: Lexer → Parser (Pratt) → Evaluador (Tree-Walking).
 *
 * Uso: Monkey.lex(src) · Monkey.parse(src, {trace:true}) · Monkey.run(src, env)
 * Lenguajes y Compiladores 2027-1
 */
(function (global) {
    'use strict';

    /* ============================ TOKENS ============================ */
    const KEYWORDS = { fn: 'FUNCTION', let: 'LET', true: 'TRUE', false: 'FALSE', if: 'IF', else: 'ELSE', return: 'RETURN' };
    const SINGLE = {
        '+': 'PLUS', '-': 'MINUS', '*': 'ASTERISK', '/': 'SLASH', '<': 'LT', '>': 'GT',
        ',': 'COMMA', ';': 'SEMICOLON', ':': 'COLON', '(': 'LPAREN', ')': 'RPAREN',
        '{': 'LBRACE', '}': 'RBRACE', '[': 'LBRACKET', ']': 'RBRACKET'
    };

    /* Valor real de cada token.TokenType en Go (token.ASSIGN = "=", etc.).
       Los mensajes de error del libro imprimen estos valores. */
    const GO_TYPE = {
        ASSIGN: '=', PLUS: '+', MINUS: '-', BANG: '!', ASTERISK: '*', SLASH: '/', LT: '<', GT: '>',
        EQ: '==', NOT_EQ: '!=', COMMA: ',', SEMICOLON: ';', COLON: ':', LPAREN: '(', RPAREN: ')',
        LBRACE: '{', RBRACE: '}', LBRACKET: '[', RBRACKET: ']'
    };
    const goType = (t) => GO_TYPE[t] || t;

    /* ============================ LEXER ============================= */
    function Lexer(input) {
        this.input = input;
        this.position = 0;
        this.readPosition = 0;
        this.ch = '';
        this.readChar();
    }
    Lexer.prototype.readChar = function () {
        this.ch = this.readPosition >= this.input.length ? '\0' : this.input[this.readPosition];
        this.position = this.readPosition;
        this.readPosition += 1;
    };
    Lexer.prototype.peekChar = function () {
        return this.readPosition >= this.input.length ? '\0' : this.input[this.readPosition];
    };
    Lexer.prototype.skipWhitespace = function () {
        while (this.ch === ' ' || this.ch === '\t' || this.ch === '\n' || this.ch === '\r') this.readChar();
    };
    const isLetter = (c) => (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c === '_';
    const isDigit = (c) => c >= '0' && c <= '9';
    Lexer.prototype.readWhile = function (pred) {
        const start = this.position;
        while (pred(this.ch)) this.readChar();
        return this.input.slice(start, this.position);
    };
    Lexer.prototype.readString = function () {
        const start = this.position + 1;
        do { this.readChar(); } while (this.ch !== '"' && this.ch !== '\0');
        return this.input.slice(start, this.position);
    };
    Lexer.prototype.nextToken = function () {
        this.skipWhitespace();
        const ch = this.ch;
        let tok;
        if (ch === '=') {
            if (this.peekChar() === '=') { this.readChar(); tok = { type: 'EQ', literal: '==' }; }
            else tok = { type: 'ASSIGN', literal: '=' };
        } else if (ch === '!') {
            if (this.peekChar() === '=') { this.readChar(); tok = { type: 'NOT_EQ', literal: '!=' }; }
            else tok = { type: 'BANG', literal: '!' };
        } else if (ch === '"') {
            tok = { type: 'STRING', literal: this.readString() };
        } else if (SINGLE[ch]) {
            tok = { type: SINGLE[ch], literal: ch };
        } else if (ch === '\0') {
            return { type: 'EOF', literal: '' };
        } else if (isLetter(ch)) {
            const lit = this.readWhile(isLetter);
            return { type: KEYWORDS[lit] || 'IDENT', literal: lit };
        } else if (isDigit(ch)) {
            return { type: 'INT', literal: this.readWhile(isDigit) };
        } else {
            tok = { type: 'ILLEGAL', literal: ch };
        }
        this.readChar();
        return tok;
    };

    function lex(src) {
        const l = new Lexer(src);
        const out = [];
        for (;;) {
            const t = l.nextToken();
            out.push(t);
            if (t.type === 'EOF') return out;
        }
    }

    /* ========================= PRECEDENCIAS ========================= */
    const P = { LOWEST: 1, EQUALS: 2, LESSGREATER: 3, SUM: 4, PRODUCT: 5, PREFIX: 6, CALL: 7, INDEX: 8 };
    const PNAME = Object.fromEntries(Object.entries(P).map(([k, v]) => [v, k]));
    const precedences = {
        EQ: P.EQUALS, NOT_EQ: P.EQUALS, LT: P.LESSGREATER, GT: P.LESSGREATER,
        PLUS: P.SUM, MINUS: P.SUM, SLASH: P.PRODUCT, ASTERISK: P.PRODUCT,
        LPAREN: P.CALL, LBRACKET: P.INDEX
    };

    /* ======================== AST · String() ======================== */
    function str(n) {
        if (!n) return '';
        switch (n.type) {
            case 'Program': return n.statements.map(str).join('');
            case 'LetStatement': return `let ${str(n.name)} = ${str(n.value)};`;
            case 'ReturnStatement': return `return ${str(n.value)};`;
            case 'ExpressionStatement': return str(n.expression);
            case 'BlockStatement': return n.statements.map(str).join('');
            case 'Identifier': return n.value;
            case 'IntegerLiteral': return String(n.value);
            case 'StringLiteral': return n.value;
            case 'Boolean': return String(n.value);
            case 'PrefixExpression': return `(${n.operator}${str(n.right)})`;
            case 'InfixExpression': return `(${str(n.left)} ${n.operator} ${str(n.right)})`;
            case 'IfExpression': return `if${str(n.condition)} ${str(n.consequence)}` + (n.alternative ? `else ${str(n.alternative)}` : '');
            case 'FunctionLiteral': return `fn(${n.parameters.map(str).join(', ')}) ${str(n.body)}`;
            case 'CallExpression': return `${str(n.fn)}(${n.args.map(str).join(', ')})`;
            case 'ArrayLiteral': return `[${n.elements.map(str).join(', ')}]`;
            case 'IndexExpression': return `(${str(n.left)}[${str(n.index)}])`;
            case 'HashLiteral': return `{${n.pairs.map(([k, v]) => `${str(k)}:${str(v)}`).join(', ')}}`;
            default: return '?';
        }
    }

    /* ======================== PARSER (Pratt) ======================== */
    function Parser(lexer, opts) {
        this.l = lexer;
        this.errors = [];
        this.trace = opts && opts.trace ? [] : null;
        this.depth = 0;
        this.nextToken();
        this.nextToken();

        this.prefixParseFns = {
            IDENT: this.parseIdentifier, INT: this.parseIntegerLiteral, STRING: this.parseStringLiteral,
            BANG: this.parsePrefixExpression, MINUS: this.parsePrefixExpression,
            TRUE: this.parseBoolean, FALSE: this.parseBoolean, LPAREN: this.parseGroupedExpression,
            IF: this.parseIfExpression, FUNCTION: this.parseFunctionLiteral,
            LBRACKET: this.parseArrayLiteral, LBRACE: this.parseHashLiteral
        };
        this.infixParseFns = {
            PLUS: this.parseInfixExpression, MINUS: this.parseInfixExpression,
            SLASH: this.parseInfixExpression, ASTERISK: this.parseInfixExpression,
            EQ: this.parseInfixExpression, NOT_EQ: this.parseInfixExpression,
            LT: this.parseInfixExpression, GT: this.parseInfixExpression,
            LPAREN: this.parseCallExpression, LBRACKET: this.parseIndexExpression
        };
    }
    const tokStr = (t) => (t.type === 'EOF' ? 'EOF' : `${t.type} "${t.literal}"`);
    Parser.prototype.log = function (kind, text) {
        if (this.trace) this.trace.push({ depth: this.depth, kind, text, cur: tokStr(this.curToken), peek: tokStr(this.peekToken) });
    };
    Parser.prototype.nextToken = function () {
        this.curToken = this.peekToken || { type: 'EOF', literal: '' };
        this.peekToken = this.l.nextToken();
    };
    Parser.prototype.curTokenIs = function (t) { return this.curToken.type === t; };
    Parser.prototype.peekTokenIs = function (t) { return this.peekToken.type === t; };
    Parser.prototype.expectPeek = function (t) {
        if (this.peekTokenIs(t)) { this.nextToken(); return true; }
        this.peekError(t);
        return false;
    };
    Parser.prototype.peekError = function (t) {
        this.errors.push(`expected next token to be ${goType(t)}, got ${goType(this.peekToken.type)} instead`);
    };
    Parser.prototype.peekPrecedence = function () { return precedences[this.peekToken.type] || P.LOWEST; };
    Parser.prototype.curPrecedence = function () { return precedences[this.curToken.type] || P.LOWEST; };

    Parser.prototype.parseProgram = function () {
        const program = { type: 'Program', statements: [] };
        while (!this.curTokenIs('EOF')) {
            const stmt = this.parseStatement();
            if (stmt) program.statements.push(stmt);
            this.nextToken();
        }
        return program;
    };
    Parser.prototype.parseStatement = function () {
        switch (this.curToken.type) {
            case 'LET': return this.parseLetStatement();
            case 'RETURN': return this.parseReturnStatement();
            default: return this.parseExpressionStatement();
        }
    };
    Parser.prototype.parseLetStatement = function () {
        this.log('stmt', 'parseLetStatement()');
        if (!this.expectPeek('IDENT')) return null;
        const name = { type: 'Identifier', value: this.curToken.literal };
        if (!this.expectPeek('ASSIGN')) return null;
        this.nextToken();
        const value = this.parseExpression(P.LOWEST);
        if (this.peekTokenIs('SEMICOLON')) this.nextToken();
        return { type: 'LetStatement', name, value };
    };
    Parser.prototype.parseReturnStatement = function () {
        this.log('stmt', 'parseReturnStatement()');
        this.nextToken();
        const value = this.parseExpression(P.LOWEST);
        if (this.peekTokenIs('SEMICOLON')) this.nextToken();
        return { type: 'ReturnStatement', value };
    };
    Parser.prototype.parseExpressionStatement = function () {
        this.log('stmt', 'parseExpressionStatement()');
        const expression = this.parseExpression(P.LOWEST);
        if (this.peekTokenIs('SEMICOLON')) this.nextToken();
        return { type: 'ExpressionStatement', expression };
    };

    /* El corazón del Pratt parser (idéntico a parser.go del libro) */
    Parser.prototype.parseExpression = function (precedence) {
        this.log('enter', `parseExpression(${PNAME[precedence]})`);
        this.depth++;
        const prefix = this.prefixParseFns[this.curToken.type];
        if (!prefix) {
            this.errors.push(`no prefix parse function for ${goType(this.curToken.type)} found`);
            this.depth--;
            return null;
        }
        this.log('prefix', `prefijo → ${prefix.label || 'parse'}(${tokStr(this.curToken)})`);
        let leftExp = prefix.call(this);

        while (!this.peekTokenIs('SEMICOLON') && precedence < this.peekPrecedence()) {
            this.log('loop', `peek ${this.peekToken.type}: ${PNAME[this.peekPrecedence()]} > ${PNAME[precedence]} → continúa con infijo`);
            const infix = this.infixParseFns[this.peekToken.type];
            if (!infix) { this.depth--; return leftExp; }
            this.nextToken();
            leftExp = infix.call(this, leftExp);
        }
        if (this.trace) {
            const why = this.peekTokenIs('SEMICOLON') ? 'peek es SEMICOLON'
                : `peek ${this.peekToken.type}: ${PNAME[this.peekPrecedence()]} ≤ ${PNAME[precedence]}`;
            this.log('stop', `${why} → retorna ${str(leftExp)}`);
        }
        this.depth--;
        return leftExp;
    };

    Parser.prototype.parseIdentifier = function () { return { type: 'Identifier', value: this.curToken.literal }; };
    Parser.prototype.parseIntegerLiteral = function () {
        const v = Number(this.curToken.literal);
        if (!Number.isSafeInteger(v)) {
            this.errors.push(`could not parse "${this.curToken.literal}" as integer`);
            return null;
        }
        return { type: 'IntegerLiteral', value: v };
    };
    Parser.prototype.parseStringLiteral = function () { return { type: 'StringLiteral', value: this.curToken.literal }; };
    Parser.prototype.parseBoolean = function () { return { type: 'Boolean', value: this.curTokenIs('TRUE') }; };
    Parser.prototype.parsePrefixExpression = function () {
        const operator = this.curToken.literal;
        this.nextToken();
        const right = this.parseExpression(P.PREFIX);
        return { type: 'PrefixExpression', operator, right };
    };
    Parser.prototype.parseInfixExpression = function (left) {
        const operator = this.curToken.literal;
        const precedence = this.curPrecedence();
        this.log('infix', `infijo "${operator}" con left = ${str(left)}`);
        this.nextToken();
        const right = this.parseExpression(precedence);
        return { type: 'InfixExpression', left, operator, right };
    };
    Parser.prototype.parseGroupedExpression = function () {
        this.nextToken();
        const exp = this.parseExpression(P.LOWEST);
        if (!this.expectPeek('RPAREN')) return null;
        return exp;
    };
    Parser.prototype.parseBlockStatement = function () {
        const block = { type: 'BlockStatement', statements: [] };
        this.nextToken();
        while (!this.curTokenIs('RBRACE') && !this.curTokenIs('EOF')) {
            const stmt = this.parseStatement();
            if (stmt) block.statements.push(stmt);
            this.nextToken();
        }
        return block;
    };
    Parser.prototype.parseIfExpression = function () {
        if (!this.expectPeek('LPAREN')) return null;
        this.nextToken();
        const condition = this.parseExpression(P.LOWEST);
        if (!this.expectPeek('RPAREN')) return null;
        if (!this.expectPeek('LBRACE')) return null;
        const consequence = this.parseBlockStatement();
        let alternative = null;
        if (this.peekTokenIs('ELSE')) {
            this.nextToken();
            if (!this.expectPeek('LBRACE')) return null;
            alternative = this.parseBlockStatement();
        }
        return { type: 'IfExpression', condition, consequence, alternative };
    };
    Parser.prototype.parseFunctionParameters = function () {
        const ids = [];
        if (this.peekTokenIs('RPAREN')) { this.nextToken(); return ids; }
        this.nextToken();
        ids.push({ type: 'Identifier', value: this.curToken.literal });
        while (this.peekTokenIs('COMMA')) {
            this.nextToken();
            this.nextToken();
            ids.push({ type: 'Identifier', value: this.curToken.literal });
        }
        if (!this.expectPeek('RPAREN')) return null;
        return ids;
    };
    Parser.prototype.parseFunctionLiteral = function () {
        if (!this.expectPeek('LPAREN')) return null;
        const parameters = this.parseFunctionParameters() || []; // en Go queda nil y se sigue
        if (!this.expectPeek('LBRACE')) return null;
        const body = this.parseBlockStatement();
        return { type: 'FunctionLiteral', parameters, body };
    };
    Parser.prototype.parseExpressionList = function (end) {
        const list = [];
        if (this.peekTokenIs(end)) { this.nextToken(); return list; }
        this.nextToken();
        list.push(this.parseExpression(P.LOWEST));
        while (this.peekTokenIs('COMMA')) {
            this.nextToken();
            this.nextToken();
            list.push(this.parseExpression(P.LOWEST));
        }
        if (!this.expectPeek(end)) return null;
        return list;
    };
    Parser.prototype.parseCallExpression = function (fn) {
        this.log('infix', `llamada: ${str(fn)}( … )`);
        const args = this.parseExpressionList('RPAREN');
        return { type: 'CallExpression', fn, args: args || [] };
    };
    Parser.prototype.parseArrayLiteral = function () {
        const elements = this.parseExpressionList('RBRACKET');
        return { type: 'ArrayLiteral', elements: elements || [] };
    };
    Parser.prototype.parseIndexExpression = function (left) {
        this.log('infix', `índice: ${str(left)}[ … ]`);
        this.nextToken();
        const index = this.parseExpression(P.LOWEST);
        if (!this.expectPeek('RBRACKET')) return null;
        return { type: 'IndexExpression', left, index };
    };
    Parser.prototype.parseHashLiteral = function () {
        const pairs = [];
        while (!this.peekTokenIs('RBRACE')) {
            this.nextToken();
            const key = this.parseExpression(P.LOWEST);
            if (!this.expectPeek('COLON')) return null;
            this.nextToken();
            const value = this.parseExpression(P.LOWEST);
            pairs.push([key, value]);
            if (!this.peekTokenIs('RBRACE') && !this.expectPeek('COMMA')) return null;
        }
        if (!this.expectPeek('RBRACE')) return null;
        return { type: 'HashLiteral', pairs };
    };
    // Etiquetas legibles para la traza
    [['parseIdentifier'], ['parseIntegerLiteral'], ['parseStringLiteral'], ['parseBoolean'], ['parsePrefixExpression'],
        ['parseGroupedExpression'], ['parseIfExpression'], ['parseFunctionLiteral'], ['parseArrayLiteral'], ['parseHashLiteral']]
        .forEach(([k]) => { Parser.prototype[k].label = k; });

    function parse(src, opts) {
        const p = new Parser(new Lexer(src), opts);
        const program = p.parseProgram();
        return { program, errors: p.errors, trace: p.trace };
    }

    /* ========================= OBJETOS ========================= */
    const TRUE = { type: 'BOOLEAN', value: true };
    const FALSE = { type: 'BOOLEAN', value: false };
    const NULL = { type: 'NULL' };
    const nativeBool = (b) => (b ? TRUE : FALSE);
    const newError = (msg) => ({ type: 'ERROR', message: msg });
    const isError = (o) => o && o.type === 'ERROR';

    function inspect(o) {
        if (!o) return '';
        switch (o.type) {
            case 'INTEGER': return String(o.value);
            case 'BOOLEAN': return String(o.value);
            case 'NULL': return 'null';
            case 'STRING': return o.value;
            case 'RETURN_VALUE': return inspect(o.value);
            case 'ERROR': return 'ERROR: ' + o.message;
            case 'FUNCTION': return `fn(${o.parameters.map(str).join(', ')}) {\n${str(o.body)}\n}`;
            case 'BUILTIN': return 'builtin function';
            case 'ARRAY': return `[${o.elements.map(inspect).join(', ')}]`;
            case 'HASH': return `{${[...o.pairs.values()].map((p) => `${inspect(p.key)}: ${inspect(p.value)}`).join(', ')}}`;
            default: return '?';
        }
    }

    function Environment(outer) { this.store = new Map(); this.outer = outer || null; }
    Environment.prototype.get = function (name) {
        if (this.store.has(name)) return this.store.get(name);
        return this.outer ? this.outer.get(name) : undefined;
    };
    Environment.prototype.set = function (name, val) { this.store.set(name, val); return val; };
    const newEnclosedEnvironment = (outer) => new Environment(outer);

    function hashKey(o) {
        if (o.type === 'INTEGER' || o.type === 'STRING' || o.type === 'BOOLEAN') return `${o.type}:${o.value}`;
        return null;
    }

    /* ========================= BUILT-INS ========================= */
    function makeBuiltins(out) {
        const argErr = (got, want) => newError(`wrong number of arguments. got=${got}, want=${want}`);
        return {
            len: { type: 'BUILTIN', fn: (args) => {
                if (args.length !== 1) return argErr(args.length, 1);
                const a = args[0];
                if (a.type === 'STRING') return { type: 'INTEGER', value: a.value.length };
                if (a.type === 'ARRAY') return { type: 'INTEGER', value: a.elements.length };
                return newError(`argument to \`len\` not supported, got ${a.type}`);
            } },
            first: { type: 'BUILTIN', fn: (args) => {
                if (args.length !== 1) return argErr(args.length, 1);
                if (args[0].type !== 'ARRAY') return newError(`argument to \`first\` must be ARRAY, got ${args[0].type}`);
                return args[0].elements.length ? args[0].elements[0] : NULL;
            } },
            last: { type: 'BUILTIN', fn: (args) => {
                if (args.length !== 1) return argErr(args.length, 1);
                if (args[0].type !== 'ARRAY') return newError(`argument to \`last\` must be ARRAY, got ${args[0].type}`);
                const e = args[0].elements;
                return e.length ? e[e.length - 1] : NULL;
            } },
            rest: { type: 'BUILTIN', fn: (args) => {
                if (args.length !== 1) return argErr(args.length, 1);
                if (args[0].type !== 'ARRAY') return newError(`argument to \`rest\` must be ARRAY, got ${args[0].type}`);
                const e = args[0].elements;
                return e.length ? { type: 'ARRAY', elements: e.slice(1) } : NULL;
            } },
            push: { type: 'BUILTIN', fn: (args) => {
                if (args.length !== 2) return argErr(args.length, 2);
                if (args[0].type !== 'ARRAY') return newError(`argument to \`push\` must be ARRAY, got ${args[0].type}`);
                return { type: 'ARRAY', elements: [...args[0].elements, args[1]] };
            } },
            puts: { type: 'BUILTIN', fn: (args) => { args.forEach((a) => out.push(inspect(a))); return NULL; } }
        };
    }

    /* ======================= EVALUADOR ======================= */
    function Evaluator(out) { this.builtins = makeBuiltins(out); this.steps = 0; }

    Evaluator.prototype.eval = function (node, env) {
        if (++this.steps > 200000) throw new Error('Límite de pasos alcanzado (¿recursión infinita?)');
        switch (node.type) {
            case 'Program': return this.evalProgram(node, env);
            case 'ExpressionStatement': return this.eval(node.expression, env);
            case 'BlockStatement': return this.evalBlockStatement(node, env);
            case 'ReturnStatement': {
                const val = this.eval(node.value, env);
                if (isError(val)) return val;
                return { type: 'RETURN_VALUE', value: val };
            }
            case 'LetStatement': {
                const val = this.eval(node.value, env);
                if (isError(val)) return val;
                env.set(node.name.value, val);
                return undefined;
            }
            case 'IntegerLiteral': return { type: 'INTEGER', value: node.value };
            case 'StringLiteral': return { type: 'STRING', value: node.value };
            case 'Boolean': return nativeBool(node.value);
            case 'PrefixExpression': {
                const right = this.eval(node.right, env);
                if (isError(right)) return right;
                return evalPrefixExpression(node.operator, right);
            }
            case 'InfixExpression': {
                const left = this.eval(node.left, env);
                if (isError(left)) return left;
                const right = this.eval(node.right, env);
                if (isError(right)) return right;
                return evalInfixExpression(node.operator, left, right);
            }
            case 'IfExpression': {
                const cond = this.eval(node.condition, env);
                if (isError(cond)) return cond;
                if (isTruthy(cond)) return this.eval(node.consequence, env);
                if (node.alternative) return this.eval(node.alternative, env);
                return NULL;
            }
            case 'Identifier': {
                const v = env.get(node.value);
                if (v !== undefined) return v;
                if (this.builtins[node.value]) return this.builtins[node.value];
                return newError('identifier not found: ' + node.value);
            }
            case 'FunctionLiteral': return { type: 'FUNCTION', parameters: node.parameters, body: node.body, env };
            case 'CallExpression': {
                const fn = this.eval(node.fn, env);
                if (isError(fn)) return fn;
                const args = this.evalExpressions(node.args, env);
                if (args.length === 1 && isError(args[0])) return args[0];
                return this.applyFunction(fn, args);
            }
            case 'ArrayLiteral': {
                const elements = this.evalExpressions(node.elements, env);
                if (elements.length === 1 && isError(elements[0])) return elements[0];
                return { type: 'ARRAY', elements };
            }
            case 'IndexExpression': {
                const left = this.eval(node.left, env);
                if (isError(left)) return left;
                const index = this.eval(node.index, env);
                if (isError(index)) return index;
                return evalIndexExpression(left, index);
            }
            case 'HashLiteral': return this.evalHashLiteral(node, env);
            default: return NULL;
        }
    };
    Evaluator.prototype.evalProgram = function (program, env) {
        let result;
        for (const s of program.statements) {
            result = this.eval(s, env);
            if (result && result.type === 'RETURN_VALUE') return result.value;
            if (isError(result)) return result;
        }
        return result;
    };
    Evaluator.prototype.evalBlockStatement = function (block, env) {
        let result;
        for (const s of block.statements) {
            result = this.eval(s, env);
            if (result && (result.type === 'RETURN_VALUE' || result.type === 'ERROR')) return result;
        }
        return result;
    };
    Evaluator.prototype.evalExpressions = function (exps, env) {
        const result = [];
        for (const e of exps) {
            const ev = this.eval(e, env);
            if (isError(ev)) return [ev];
            result.push(ev);
        }
        return result;
    };
    Evaluator.prototype.applyFunction = function (fn, args) {
        if (fn.type === 'FUNCTION') {
            const extendedEnv = newEnclosedEnvironment(fn.env);
            fn.parameters.forEach((p, i) => extendedEnv.set(p.value, args[i]));
            const evaluated = this.eval(fn.body, extendedEnv);
            return evaluated && evaluated.type === 'RETURN_VALUE' ? evaluated.value : evaluated;
        }
        if (fn.type === 'BUILTIN') return fn.fn(args);
        return newError('not a function: ' + fn.type);
    };
    Evaluator.prototype.evalHashLiteral = function (node, env) {
        const pairs = new Map();
        for (const [kn, vn] of node.pairs) {
            const key = this.eval(kn, env);
            if (isError(key)) return key;
            const hk = hashKey(key);
            if (!hk) return newError('unusable as hash key: ' + key.type);
            const value = this.eval(vn, env);
            if (isError(value)) return value;
            pairs.set(hk, { key, value });
        }
        return { type: 'HASH', pairs };
    };

    function evalPrefixExpression(op, right) {
        if (op === '!') {
            if (right === TRUE) return FALSE;
            if (right === FALSE) return TRUE;
            if (right === NULL) return TRUE;
            return FALSE;
        }
        if (op === '-') {
            if (right.type !== 'INTEGER') return newError('unknown operator: -' + right.type);
            return { type: 'INTEGER', value: -right.value };
        }
        return newError(`unknown operator: ${op}${right.type}`);
    }
    function evalInfixExpression(op, l, r) {
        if (l.type === 'INTEGER' && r.type === 'INTEGER') {
            const a = l.value, b = r.value;
            switch (op) {
                case '+': return { type: 'INTEGER', value: a + b };
                case '-': return { type: 'INTEGER', value: a - b };
                case '*': return { type: 'INTEGER', value: a * b };
                case '/':
                    if (b === 0) return newError('division by zero (en Go esto sería un panic)');
                    return { type: 'INTEGER', value: Math.trunc(a / b) };
                case '<': return nativeBool(a < b);
                case '>': return nativeBool(a > b);
                case '==': return nativeBool(a === b);
                case '!=': return nativeBool(a !== b);
                default: return newError(`unknown operator: ${l.type} ${op} ${r.type}`);
            }
        }
        if (l.type === 'STRING' && r.type === 'STRING') {
            if (op !== '+') return newError(`unknown operator: ${l.type} ${op} ${r.type}`);
            return { type: 'STRING', value: l.value + r.value };
        }
        if (op === '==') return nativeBool(l === r);
        if (op === '!=') return nativeBool(l !== r);
        if (l.type !== r.type) return newError(`type mismatch: ${l.type} ${op} ${r.type}`);
        return newError(`unknown operator: ${l.type} ${op} ${r.type}`);
    }
    function isTruthy(o) {
        if (o === NULL) return false;
        if (o === TRUE) return true;
        if (o === FALSE) return false;
        return true;
    }
    function evalIndexExpression(left, index) {
        if (left.type === 'ARRAY' && index.type === 'INTEGER') {
            const i = index.value, max = left.elements.length - 1;
            if (i < 0 || i > max) return NULL;
            return left.elements[i];
        }
        if (left.type === 'HASH') {
            const hk = hashKey(index);
            if (!hk) return newError('unusable as hash key: ' + index.type);
            const pair = left.pairs.get(hk);
            return pair ? pair.value : NULL;
        }
        return newError('index operator not supported: ' + left.type);
    }

    /* Ejecuta un programa completo: devuelve resultado, salida de puts y errores de parser */
    function run(src, env) {
        const { program, errors } = parse(src);
        if (errors.length) return { parseErrors: errors };
        const out = [];
        const ev = new Evaluator(out);
        env = env || new Environment();
        try {
            const result = ev.eval(program, env);
            return { result, output: out, env, program };
        } catch (e) {
            return { result: newError(e.message), output: out, env, program };
        }
    }

    /* ===================== AST → HTML (árbol) ===================== */
    function esc(s) { return String(s).replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c])); }
    function astChildren(n) {
        switch (n.type) {
            case 'Program': case 'BlockStatement': return n.statements;
            case 'LetStatement': return [n.name, n.value];
            case 'ReturnStatement': return [n.value];
            case 'ExpressionStatement': return [n.expression];
            case 'PrefixExpression': return [n.right];
            case 'InfixExpression': return [n.left, n.right];
            case 'IfExpression': return [n.condition, n.consequence, n.alternative].filter(Boolean);
            case 'FunctionLiteral': return [...n.parameters, n.body];
            case 'CallExpression': return [n.fn, ...n.args];
            case 'ArrayLiteral': return n.elements;
            case 'IndexExpression': return [n.left, n.index];
            case 'HashLiteral': return n.pairs.flat();
            default: return [];
        }
    }
    function astLabel(n) {
        switch (n.type) {
            case 'Identifier': return ['lit', 'Identifier', n.value];
            case 'IntegerLiteral': return ['lit', 'Integer', n.value];
            case 'StringLiteral': return ['lit', 'String', `"${n.value}"`];
            case 'Boolean': return ['lit', 'Boolean', n.value];
            case 'PrefixExpression': return ['op', 'Prefix', n.operator];
            case 'InfixExpression': return ['op', 'Infix', n.operator];
            case 'LetStatement': return ['stmt', 'Let', ''];
            case 'ReturnStatement': return ['stmt', 'Return', ''];
            case 'ExpressionStatement': return ['stmt', 'ExprStmt', ''];
            case 'CallExpression': return ['op', 'Call', ''];
            case 'IndexExpression': return ['op', 'Index', '[ ]'];
            default: return ['', n.type.replace('Statement', 'Stmt').replace('Expression', 'Expr').replace('Literal', ''), ''];
        }
    }
    function astToHtml(n) {
        if (!n) return '';
        const [cls, label, detail] = astLabel(n);
        const kids = astChildren(n).filter(Boolean);
        return `<li><span class="ast-node ${cls}"><b>${esc(label)}</b>${detail !== '' ? `<i>${esc(detail)}</i>` : ''}</span>` +
            (kids.length ? `<ul>${kids.map(astToHtml).join('')}</ul>` : '') + '</li>';
    }
    function treeHtml(node) { return `<ul>${astToHtml(node)}</ul>`; }

    global.Monkey = { goType, lex, parse, run, str, inspect, treeHtml, esc, Environment, P, PNAME, precedences };
})(window);
