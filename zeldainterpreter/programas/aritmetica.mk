puts("== Precedencia de operadores ==");
puts(1 + 2 * 3);
puts((1 + 2) * 3);
puts(-5 + 10 / 2);
puts(7 / 2);

puts("== Comparaciones y booleanos ==");
puts(10 > 5 == true);
puts(!true);
puts(!!0);

puts("== if es una expresion ==");
let edad = 20;
let mensaje = if (edad > 17) { "mayor de edad" } else { "menor de edad" };
puts(mensaje);
