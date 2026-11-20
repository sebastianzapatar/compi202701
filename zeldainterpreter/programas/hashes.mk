let personas = [
  {"nombre": "Ada", "lenguaje": "Go", "edad": 36},
  {"nombre": "Alan", "lenguaje": "Monkey", "edad": 41}
];

let saludar = fn(p) {
  puts("Hola, " + p["nombre"] + "! Te gusta " + p["lenguaje"]);
};

saludar(personas[0]);
saludar(personas[1]);

puts("Edad de Alan:");
puts(personas[1]["edad"]);

puts("Llave que no existe:");
puts(personas[0]["apellido"]);
