let newAdder = fn(x) {
  fn(y) { x + y }
};

let sumarDos = newAdder(2);
let sumarDiez = newAdder(10);

puts("sumarDos(40) =");
puts(sumarDos(40));
puts("sumarDiez(5) =");
puts(sumarDiez(5));

let compose = fn(f, g) {
  fn(x) { f(g(x)) }
};
let doble = fn(x) { x * 2 };
let dobleYMasDos = compose(sumarDos, doble);

puts("compose(sumarDos, doble)(20) =");
puts(dobleYMasDos(20));
