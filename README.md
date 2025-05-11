[![build](https://github.com/Kolterdyx/mcbasic/actions/workflows/build.yml/badge.svg)](https://github.com/Kolterdyx/mcbasic/actions/workflows/build.yml)

# MCBasic

MCBasic is a statically typed scripting language with a mixture of C syntax and Python features.
It compiles to a datapack for Minecraft 1.21.5. Other versions might or might not work, idk, Mojang keeps introducing breaking changes ¯\\\_(ツ)_/¯.

This project is still WIP and is not fit for actual use.

## Features

This is a list of features that are currently implemented or planned for the future.

- [x] **Variables**: Declare and assign variables.
- [x] **Functions**: Define and call functions.
- [ ] **Variadic functions**: Functions that can take a variable number of arguments.
- [ ] **Control Flow**:
  - [x] **If-else statements**.
  - [x] **Recursion**.
  - [ ] **Loops**. Classic for/while loops will not be implemented, but alternative methods of repeating execution are planned.
- [x] **Operators**:
  - [x] **Arithmetic**.
  - [x] **Comparison**.
  - [x] **Logical operators**.
- [ ] **Data Types**:
  - [x] **Integer**.
  - [x] **Double**.
  - [x] **String**.
  - [x] **Boolean**: Booleans do not have a type. Other types are evaluated to be truthy or falsy. Falsy values are `0`, `0.0`, `""`.
  - [x] **List**: Statically typed lists.
  - [ ] **Struct**: User-defined data types.
- [x] **Comments**: Single-line comments with `#`.
- [x] **String operations**:
  - [x] **Concatenation**.
  - [x] **Length**.
  - [ ] **Substring / slicing**.


## Quirks

These are some of the quirks of the language that you should be aware of. They may be fixed in the future.

- **No global scope**: Only functions can be declared at the top level. Variables must be declared inside functions. Variables are scoped to the function they are declared in.
- **No type inference**: The compiler does not infer the type of variable from its value.
- **No type casting**: You cannot cast a variable from one type to another.
- **Recursive functions**: Recursive functions have not been thoroughly tested and may not work as expected.


## Usage

First you need to install the compiler. It's a standalone binary, so you can download it and run it from anywhere.
You can download the [latest version](https://github.com/Kolterdyx/mcbasic/releases/latest) from the releases page.

### Syntax highlighting
If you want syntax highlighting for the language, you can use the [VSCode extension](https://marketplace.visualstudio.com/items?itemName=mcbasic-lang.mcbasic-lang) I made.
It's pretty basic, but it's better than nothing.

### Project setup

There are two ways to set up a project: using the `mcbasic` command line tool or manually creating the files.

#### Using the command line tool

You can use the `mcbasic` command line tool to create a new project. Run the following command:

```sh
mcbasic init
```
This will prompt you for some information about your project, such as the name, namespace, and if you want to initialize a git repository.
It will create a new directory with the name of your project and create the necessary files for you.

#### Manually creating the files
If you prefer to set up the project manually, you can do so by creating the necessary files yourself.

Create a new directory for your project and create a file called `project.toml` with the following content:
```toml
[Project]
name = "Example project"
namespace = "example"
description = "An example project"
version = "0.1.0"
entrypoint = "main.mcb"
```

This file is used to store metadata about your project, and it tells the compiler where to start compiling your code and
how to generate the pack.mcmeta file.

Then, create a new file called `main.mcb` with the following content:

```go # Go is used for syntax highlighting, but the language is not Go
func load() {
    mcb:print("Hello, world!");
}
```

This is the entry point of your program. The `load` function is called when the datapack is loaded.

You can also create a function called `tick` that will be called every tick.

### Compiling

Run the following command inside the project folder to compile your code to a datapack:

```sh
mcbasic build
```

This will generate a new directory called `build` with the compiled datapack inside.
You can then move this directory to your Minecraft world's datapacks folder.

> **Note**: The compiler will overwrite the `build` directory every time you compile your code.
> Make sure to back up any changes you make manually to the compiled datapack.
> You can also specify a different output directory by passing the `--output` flag to the compiler.
> This can be useful if you want to compile directly to the datapacks folder of your world.

### Running

Load the datapack in your Minecraft world and join or reload the world to see the output of your program.
if you are joining a world, the message may appear before you join the world. You can always run the command `/function <your_namespace>:main` to see the output again.


## Hello world example

```go # Go is used for syntax highlighting, but the language is not Go
func load() {
    mcb:print("Hello, world!");
}
```

You can find more examples in the [examples](examples) folder.

## Contributing

Feel free to contribute to this project by opening an issue or a pull request.
Optimization suggestions are also welcome, as I'm still learning about compilers.

To avoid duplicated work, please submit an issue before starting a contribution, and explain what your contribution will implement.

## Development documentation

- [Implementation details](docs/implementation_details.md)

## License

This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html). See the [LICENSE](LICENSE) file for more details.
