let fibonacci = fn(n) {
  if (n < 2) {
    n
  } else {
    fibonacci(n - 1) + fibonacci(n - 2)
  }
};

puts("fibonacci(10) =");
puts(fibonacci(10));

puts("fibonacci(15) =");
puts(fibonacci(15));
