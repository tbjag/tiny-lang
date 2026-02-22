use crate::lexer::Token;

#[derive(Debug, Clone)]
pub enum Expression {
    Number(i64),
    String(String),
    Symbol(String),
    Binary {
        left: Box<Expression>,
        operator: Token,
        right: Box<Expression>,
    }
}

#[derive(Debug, Clone)]
pub enum Statement {
    Block {body: Vec<Statement>},
    Expression {expression: Expression}
}