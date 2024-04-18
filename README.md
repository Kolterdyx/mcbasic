
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
- [ ] **Operators**:
  - [x] **Arithmetic**.
  - [x] **Comparison**.
  - [ ] **Logical operators**.
- [ ] **Data Types**:
  - [x] **Number**.
  - [x] **String**.
  - [ ] **Array**.
  - [x] **Boolean**: Booleans are represented as numbers, with 0 being false and 1 being true. There isn't a separate boolean type.
  - [ ] **Structs**: User-defined data types.
- [x] **Comments**: Single-line comments with `#`.
- [ ] **String operations**:
  - [ ] **Concatenation**.
  - [ ] **Length**.
  - [ ] **Substring**.


## Quirks

These are some of the quirks of the language that you should be aware of. They may be fixed in the future.

- **No global scope**: Only functions can be declared at the top level. Variables must be declared inside functions.
- **No type inference**: The compiler does not infer the type of variable from its value.
- **No type casting**: You cannot cast a variable from one type to another.
- **Recursive functions**: Recursive functions have not been thoroughly tested and may not work as expected.


## Example

```
func add(a: int, b: int) {
    return a + b;
}

func fib(n: int) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
}

func main() {
    let x: int = 5;
    let y: int = 10;
    let z: int = add(x, y);
    let hello: str = "Hello world!";
    print(hello);
    print(z);
    print(fib(10));
}
```

## Contributing

Feel free to contribute to this project by opening an issue or a pull request.
Optimization suggestions are also welcome, as I'm still learning about compilers.

To avoid duplicated work, please submit an issue before starting a contribution, and explain what your contribution will implement.

## Development documentation

- [Implementation details](implementation_details.md)

## License

This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). See the [LICENSE](LICENSE) file for more details.
