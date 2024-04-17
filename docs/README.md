
# This is the development documentation for MCBasic.

## Implementation details

### General structure of the compiler

The compiler takes care of converting the generated AST from the parser into a Minecraft datapack.
It is implemented with a visitor pattern, where each node in the AST has a corresponding visitor method.
The visitor methods are responsible for generating the appropriate commands for the datapack.

[Here](https://craftinginterpreters.com/representing-code.html) is a good resource for understanding the visitor pattern.

[Here](implementation_details.md) is how each node in the AST is handled by the visitor.