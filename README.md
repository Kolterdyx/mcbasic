
# MCBasic

MCBasic is a statically typed scripting language with a mixture of C syntax and Python features.
It compiles to a datapack for Minecraft 1.20.4 and above.

This project is still WIP and is not fit for actual use.

## Features

- [x] **Variables**: Declare and assign variables.
- [x] **Functions**: Define and call functions.
- [ ] **Control Flow**:
  - [ ] **If-else statements**.
  - [ ] **Loops**.
- [ ] **Operators**:
  - [x] **Arithmetic**.
  - [ ] **Comparison**.
  - [ ] **Logical operators**.
- [ ] **Data Types**:
  - [x] **Number**.
  - [x] **String**.
  - [ ] **Array**.
  - [ ] **Boolean**.
  - [ ] **Structs**.
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
func add(a: number, b: number) {
    return a + b;
}

func main() {
    let x: number = 5;
    let y: number = 10;
    let z: number = add(x, y);
    print("Hello world!");
    print(z);
}
```

## Contributing

Feel free to contribute to this project by opening an issue or a pull request.
Optimization suggestions are also welcome, as I'm still learning about compilers.

To avoid duplicated work, please submit an issue before starting a contribution, and explain what your contribution will implement.
