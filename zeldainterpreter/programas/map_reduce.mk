let map = fn(arr, f) {
  let iter = fn(arr, acumulado) {
    if (len(arr) == 0) {
      acumulado
    } else {
      iter(rest(arr), push(acumulado, f(first(arr))));
    }
  };
  iter(arr, []);
};

let reduce = fn(arr, inicial, f) {
  let iter = fn(arr, resultado) {
    if (len(arr) == 0) {
      resultado
    } else {
      iter(rest(arr), f(resultado, first(arr)));
    }
  };
  iter(arr, inicial);
};

let numeros = [1, 2, 3, 4, 5];

puts("map(numeros, cuadrado) =");
puts(map(numeros, fn(x) { x * x }));

puts("suma con reduce =");
puts(reduce(numeros, 0, fn(acc, x) { acc + x }));
