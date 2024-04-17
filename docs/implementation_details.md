
# Implementation Details

Here is a brief overview of how the AST is compiled into a Minecraft datapack.

## Steps in the compilation process

The compiler takes in the first node of the AST, which is the `Program` node, and starts the compilation process.
The program node, in turn, contains a list of function nodes.
You can think of the program node as the top level of the source code.

The compiler then visits each function node and generates the appropriate commands for the datapack.
Each function node is compiled into its own function in the datapack.
To separate the scope of each function, a macro is appended to each variable declared in each function.
All functions have at least that macro as a parameter, which is just an integer that increments with each function call.
The body of the function is compiled, and the function is created with "__wrapped" appended to its name.
This is to prevent the function from being called directly by the user. Then, a wrapper function
with the original name is created that takes in the parameters, sets the macro, and calls the wrapped function.
This allows the user to call the function without worrying about the scope macro.

## Expression evaluation

Expressions are evaluated by visiting each node in the AST and generating the appropriate commands.
The last command in the expression moves the result to the register variable RX, which is used by later expressions.

## Binary nodes

Binary nodes are nodes that have two children, such as an addition or subtraction operation.
These nodes are compiled by visiting the left and right children and then generating the appropriate commands.
This causes a collision in the register variable RX, which would be overwritten by the right child before the left child is used.
To prevent this, a custom register is used for each binary node, which is then moved to RX at the end of the evaluation.
This is done by creating a copy of RX with "_N" appended to the name, where N is an integer that increments with each new register.