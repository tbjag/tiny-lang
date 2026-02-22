use crate::{lexer::Token, parser::ast::{Expression, Statement}};

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
pub enum  BindingPower {
    Default,
    Comma,
    Assignment,
    Logical,
    Relational,
    Additive,
    Multiplicative,
    Unary,
    Call,
    Member,
    Primary
}

impl BindingPower {
    fn next(self) -> Self {
        match self {
            Self::Default        => Self::Comma,
            Self::Comma          => Self::Assignment,
            Self::Assignment     => Self::Logical,
            Self::Logical        => Self::Relational,
            Self::Relational     => Self::Additive,
            Self::Additive       => Self::Multiplicative,
            Self::Multiplicative => Self::Unary,
            Self::Unary          => Self::Call,
            Self::Call           => Self::Member,
            Self::Member         => Self::Primary,
            Self::Primary        => Self::Primary,
        }
    }
}

pub struct Parser {
    tokens: Vec<Token>,
    pos: usize
}

impl Parser {
    pub fn new(tokens: Vec<Token>) -> Self {
        Self {tokens, pos: 0}
    }

    fn peek(&self) -> &Token {
        self.tokens.get(self.pos).unwrap_or(&Token::EndOfInput)
    }

    fn next(&mut self) -> Token {
        let t = self.tokens[self.pos].clone();
        self.pos += 1;
        t
    }

    fn expect(&mut self, expected: &Token) -> Token {
        let t = self.next();
        assert_eq!(&t, expected, "Expected {:?}", expected);
        t
    }

}

fn infix_bp(op: &Token) -> Option<(BindingPower, BindingPower)> {
    match op {
        Token::OpAdd | Token::OpSubtract => Some((BindingPower::Additive, BindingPower::Additive.next())),
        Token::OpMultiply | Token::OpDivide => Some((BindingPower::Multiplicative, BindingPower::Multiplicative.next())),
        _ => None,
    }
}

pub fn parse_expression(parser: &mut Parser, min_bp: BindingPower) -> Expression {
    let mut lhs = match parser.next() {
        Token::Integer(n) => Expression::Number(n),
        Token::String(s) => Expression::String(s),
        Token::Indentifier(s) => Expression::Symbol(s),
        Token::OpenParen => {
            let inner = parse_expression(parser, BindingPower::Default);
            parser.expect(&Token::CloseParen);
            inner
        }
        t => panic!("Unexpected token in prefix position: {:?}", t),
    };

    loop {
        let op = parser.peek().clone();

        let (left_bp, right_bp) = match infix_bp(&op) {
            Some(bp) => bp,
            None => break,
        };

        if left_bp <= min_bp {
            break;
        }

        parser.next();

        let rhs = parse_expression(parser, right_bp);

        lhs = Expression::Binary { left: Box::new(lhs), operator: op, right: Box::new(rhs) }
    }

    lhs
}

pub fn parse_statement(parser: &mut Parser) -> Statement {
    let expr = parse_expression(parser, BindingPower::Default);
    parser.expect(&Token::Semicolon);
    Statement::Expression { expression: expr }
}