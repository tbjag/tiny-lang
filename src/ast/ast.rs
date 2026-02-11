

trait Statement {
    fn statement();
}

trait Expression {
    fn expression();
}

// LITERALS

struct NumberExpression {
    value: f64
}

impl NumberExpression {
    fn expression(self) {}
}

struct StringExpression {
    value: String
}

impl StringExpression {
    fn expression(self) {}
}

struct SymbolExpression {
    value: String
}

impl SymbolExpression {
    fn expression(self) {}
}

// COMPLEX

// struct BinaryExpression {
//     left 
// }
