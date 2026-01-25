#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum TokenKind {
    EndOfInput,
    OpMultiply,
    OpDivide,
    OpMod,
    OpAdd,
    OpSubtract,
    OpNegate,
    OpNot,
    OpLess,
    OpLessEqual,
    OpGreater,
    OpGreaterEqual,
    OpEqual,
    OpNotEqual,
    OpAssign,
    OpAnd,
    OpOr,
    KeywordIf,
    KeywordElse,
    KeywordWhile,
    KeywordPrint,
    KeywordPutc,
    OpenParen,
    CloseParen,
    OpenBrace,
    CloseBrace,
    Semicolon,
    Comma,
    Indentifier,
    Integer,
    String
}

#[derive(Debug, Clone)]
pub struct Token {
    pub kind: TokenKind,
    pub value: String
}

impl Token { // why value into?
    pub fn new(kind: TokenKind, value: impl Into<String>) -> Self {
        Token { kind, value: value.into() }
    }

    pub fn is_one_of_many(&self, expected_tokens: &[TokenKind]) -> bool {
        expected_tokens.contains(&self.kind)
    }

    pub fn debug(&self) {
        if self.is_one_of_many(&[TokenKind::Indentifier, TokenKind::Integer, TokenKind::String]) {
            println!("{:?} ({})", self.kind, self.value)
        } else {
            println!("{:?} ()", self.kind)
        }
    }
}
