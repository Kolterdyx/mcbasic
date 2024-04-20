[![build](https://github.com/Kolterdyx/mcbasic/actions/workflows/go.yml/badge.svg)](https://github.com/Kolterdyx/mcbasic/actions/workflows/go.yml)

# MCBasic

MCBasic is a statically typed scripting language with a mixture of C syntax and Python features.
It compiles to a datapack for Minecraft 1.20.4 and above.

This project is still WIP and is not fit for actual use.

## Features

This is a list of features that are currently implemented or planned for the future.

- [x] **Variables**: Declare and assign variables.
- [x] **Functions**: Define and call functions.
- [ ] **Variadic functions**: Functions that can take a variable number of arguments.
- [ ] **Control Flow**:
  - [x] **If-else statements**.
  - [x] **Recursion**.
  - [ ] **Loops**.
- [x] **Operators**:
  - [x] **Arithmetic**.
  - [x] **Comparison**.
  - [x] **Logical operators**.
- [ ] **Data Types**:
  - [x] **Integer**.
  - [x] **Fixed point**.
  - [x] **String**.
  - [x] **Boolean**: Booleans are represented as numbers, with 0 being false and 1 being true. There isn't a separate boolean type.
  - [ ] **Array**.
  - [ ] **Struct**: User-defined data types.
- [x] **Comments**: Single-line comments with `#`.
- [x] **String operations**:
  - [x] **Concatenation**.
  - [x] **Length**.
  - [x] **Substring / slicing**.


## Quirks

These are some of the quirks of the language that you should be aware of. They may be fixed in the future.

- **No global scope**: Only functions can be declared at the top level. Variables must be declared inside functions.
- **No type inference**: The compiler does not infer the type of variable from its value.
- **No type casting**: You cannot cast a variable from one type to another.
- **Recursive functions**: Recursive functions have not been thoroughly tested and may not work as expected.


## Example

```
func add(a: fixed, b: fixed) fixed {
    return a + b;
}

# Recursion is possible! Careful not to hit the maxCommandChainLength or the maxCommandForkCount limit 
func fib(n: int) int {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
}

func fizbuzz(n: int) {
    # since for/while loops aren't supported yet, here is a cut down version of the fizzbuzz function
    # that evaluates a single value at a time instead of all the numbers up to n
    if (n % 3 == 0 and n % 5 == 0) {
        print("FizzBuzz");
    } else if (n % 3 == 0) {
        print("Fizz");
    } else if (n % 5 == 0) {
        print("Buzz");
    } else {
        print(n);
    }
}

func main() {
    let x: fixed = 5.51;
    let y: fixed = 10.73;
    let z: fixed = add(x, y);
    
    # string concatenation
    let hello: str = "Hello" + " world!";
    
    # print function can take any argument type
    print(hello);
    print(z);
    
    # Other types can be appended to strings. The string must be on the left side of the operation for this to work
    print("10th fib: " + fib(10));
    
    # You can also slice strings
    print(hello[0:5]);
    # Even with literals and negative indices
    print("Hello world!"[-6:-1]);
    # Expressions can be used as indices. The result must be an integer
    print(hello[0:5 + 1]);
    # Out of bounds exceptions are checked, and errors are logged to the chat
    # Negative indices are also checked
    # print(hello[0:100]);  # This will return out of the main function and print an error message
    # print(hello[0:-100]);  # This will return out of the main function and print an error message
}
```

## Contributing

Feel free to contribute to this project by opening an issue or a pull request.
Optimization suggestions are also welcome, as I'm still learning about compilers.

To avoid duplicated work, please submit an issue before starting a contribution, and explain what your contribution will implement.

## Development documentation

- [Implementation details](docs/implementation_details.md)

## License

This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). See the [LICENSE](LICENSE) file for more details.
