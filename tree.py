import json
from typing import List
from treelib import Node, Tree


class ExprType:
    BinaryExprType = "Binary"
    GroupingExprType = "Grouping"
    LiteralExprType = "Literal"
    UnaryExprType = "Unary"
    VariableExprType = "Variable"
    FunctionCallExprType = "FunctionCall"


class StmtType:
    ExpressionStmtType = "Expression"
    PrintStmtType = "Print"
    VariableDeclarationStmtType = "VariableDeclaration"
    FunctionDeclarationStmtType = "FunctionDeclaration"
    VariableAssignmentStmtType = "VariableAssignment"
    BlockStmtType = "Block"


class Node:
    def __init__(self, node_type, children):
        self.node_type = node_type
        self.children = children

    def __str__(self):
        # return f"{self.node_type}({', '.join(str(child) for child in self.children)})"
        return f"{self.node_type}"

    def __repr__(self):
        return str(self)


def get_nodes(data: List):
    nodes = []
    for stmt_expr in data:
        if stmt_expr is None or not isinstance(stmt_expr, dict):
            continue
        match stmt_expr['type']:
            case StmtType.VariableDeclarationStmtType:
                nodes.append(Node(StmtType.VariableDeclarationStmtType, [stmt_expr['name'], Node(StmtType.ExpressionStmtType, get_nodes([stmt_expr.get('initializer')]))]))
            case StmtType.FunctionDeclarationStmtType:
                nodes.append(Node(StmtType.FunctionDeclarationStmtType, [stmt_expr['name'], list(map(lambda x: x['Literal'], stmt_expr['parameters'])), Node(StmtType.BlockStmtType, get_nodes(stmt_expr['body']['statements']))]))
            case StmtType.VariableAssignmentStmtType:
                nodes.append(Node(StmtType.VariableAssignmentStmtType, [stmt_expr['name'], get_nodes([stmt_expr['value']])]))
            case StmtType.BlockStmtType:
                nodes.append(Node(StmtType.BlockStmtType, get_nodes(stmt_expr['statements'])))
            case StmtType.ExpressionStmtType:
                nodes.append(Node(StmtType.ExpressionStmtType, get_nodes([stmt_expr['expression']])))
            case StmtType.PrintStmtType:
                nodes.append(Node(StmtType.PrintStmtType, get_nodes([stmt_expr['expression']])))
            case ExprType.LiteralExprType:
                nodes.append(Node(ExprType.LiteralExprType, [stmt_expr['value']]))
            case ExprType.VariableExprType:
                nodes.append(Node(ExprType.VariableExprType, [stmt_expr['name']]))
            case ExprType.BinaryExprType:
                nodes.append(Node(ExprType.BinaryExprType, [Node(StmtType.ExpressionStmtType, get_nodes([stmt_expr['left']])), stmt_expr['operator'], Node(StmtType.ExpressionStmtType, get_nodes([stmt_expr['right']]))]))
            case ExprType.UnaryExprType:
                nodes.append(Node(ExprType.UnaryExprType, [stmt_expr['operator'], Node(StmtType.ExpressionStmtType, get_nodes([stmt_expr['expression']]))]))
            case ExprType.FunctionCallExprType:
                nodes.append(Node(ExprType.FunctionCallExprType, [stmt_expr['callee'], get_nodes(stmt_expr['arguments'])]))
            case _:
                raise ValueError(f"Unknown statement type: {stmt_expr['type']}")
    return nodes




def main():
    with open('output.json') as f:
        data = json.load(f)

    assert data['type'] == StmtType.BlockStmtType
    root = Node(StmtType.BlockStmtType, get_nodes(data['statements']))

    tree = Tree()

    def add_nodes(node: Node, parent=None, n=1):
        if not isinstance(node, Node):
            return
        if parent is None:
            parent = tree.create_node(str(node), identifier=str(0))
        for i, child in enumerate(node.children):
            if isinstance(child, Node):
                child_node = tree.create_node(str(child), parent=parent, identifier=f"{parent.identifier}{i + 1}_{n}")
                add_nodes(child, child_node, n + 1)
            elif isinstance(child, list) and all(isinstance(x, Node) for x in child):
                for j, c in enumerate(child):
                    child_node = tree.create_node(str(c), parent=parent, identifier=f"{parent.identifier}{i + 1}_{j + 1}_{n}")
                    add_nodes(c, child_node, n + 1)
            else:
                if parent is not None and parent.tag == StmtType.FunctionDeclarationStmtType and isinstance(child, list):
                    tree.create_node(f"Args: {[x for x in child]}", parent=parent, identifier=parent.identifier + str(i + 1) + str(n))
                else:
                    tree.create_node(str(child), parent=parent, identifier=parent.identifier + str(i + 1) + str(n))

    add_nodes(root)
    print(tree.show(stdout=False))


if __name__ == '__main__':
    main()
