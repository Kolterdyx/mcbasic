
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
- [x] **Operators**:
  - [x] **Arithmetic**.
  - [ ] **Comparison**.
  - [ ] **Logical operators**.

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