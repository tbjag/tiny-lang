use regex::bytes::Regex;
use std::{fs, io};

#[derive(Debug, PartialEq)]
pub enum Token {
    OpMultiply,
    OpDivide,
    OpMod,
    OpAdd,
    OpSubtract,
    OpNegate,
    OpLess,
    OpLessEqual,
    OpGreater,
    OpGreaterEqual,
    OpEqual,
    OpNotEqual,
    OpNot,
    OpAssign,
    OpAnd,
    OpOr,

    LeftParen,
    RightParen,
    LeftBrace,
    RightBrace,
    Semicolon,
    Comma,

    KeywordIf,
    KeywordElse,
    KeywordWhile,
    KeywordPrint,
    KeywordPutc,

    Identifier(String),
    Integer(i64),
    String(String),

    EOF,
    INVALID,
}

pub struct Lexer {
    input: Vec<u8>,
    position: usize,
    digit_re: Regex,
    ident_re: Regex,
    first_ident_re: Regex,
}

const OPS: &[(&[u8], fn() -> Token)] = &[
    (b"==", || Token::OpEqual),
    (b"!=", || Token::OpNotEqual),
    (b"<=", || Token::OpLessEqual),
    (b">=", || Token::OpGreaterEqual),
    (b"&&", || Token::OpAnd),
    (b"||", || Token::OpOr),
];

impl Lexer {
    pub fn new(input: Vec<u8>) -> Self {
        Self {
            input: input,
            position: 0,
            digit_re: Regex::new(r"[0-9]").unwrap(),
            ident_re: Regex::new(r"[_a-zA-Z0-9]").unwrap(),
            first_ident_re: Regex::new(r"[_a-zA-Z]").unwrap(),
        }
    }

    /// Return the next token from the input.
    /// Should return Token::EOF when input is exhausted.
    pub fn next_token(&mut self) -> Token {
        self.skip_trivia();

        let b = match self.peek() {
            Some(b) => b,
            None => return Token::EOF,
        };

        for (pat, mk) in OPS {
            if self.match_bytes(pat) {
                return mk();
            }
        }

        if self.match_bytes(b"\"") {
            let s = self.read_string();
            return Token::String(s);
        }


        if self.first_ident_re.is_match(&[b]) {
            let ident_val = self.read_identifier();
            let x = match ident_val.as_str() {
                "print" => Token::KeywordPrint,
                "putc" => Token::KeywordPutc,
                "if" => Token::KeywordIf,
                "else" => Token::KeywordElse,
                "while" => Token::KeywordWhile,
                _ => Token::Identifier(ident_val),
            };
            return x;
        }

        if self.digit_re.is_match(&[b]) {
            let digit_val = self.read_number();
            return Token::Integer(digit_val);
        }

        if self.match_bytes(b"\'") {
            if let Some(c) = self.read_literal() {
                return Token::Integer(c)
            }
            return  Token::INVALID;
        }


        match b {
            b',' => {
                self.position += 1;
                Token::Comma
            }
            b'(' => {
                self.position += 1;
                Token::LeftParen
            }
            b')' => {
                self.position += 1;
                Token::RightParen
            }
            b';' => {
                self.position += 1;
                Token::Semicolon
            }
            b'=' => {
                self.position += 1;
                Token::OpAssign
            }
            b'-' => {
                self.position += 1;
                Token::OpSubtract
            }
            b'+' => {
                self.position += 1;
                Token::OpAdd
            }
            b'>' => {
                self.position += 1;
                Token::OpGreater
            }
            b'<' => {
                self.position += 1;
                Token::OpLess
            }
            b'*' => {
                self.position += 1;
                Token::OpMultiply
            }
            b'/' => {
                self.position += 1;
                Token::OpDivide
            }
            b'{' => {
                self.position += 1;
                Token::LeftBrace
            }
            b'}' => {
                self.position += 1;
                Token::RightBrace
            }
            b'!' => {
                self.position += 1;
                Token::OpNot
            }
            b'%' => {
                self.position += 1;
                Token::OpMod
            }
            _ => {
                self.position += 1;
                Token::INVALID
            }
        }
    }

    /// Advance the cursor by one byte and return that byte.
    /// If at end of input, return None.
    fn advance(&mut self) -> Option<u8> {
        if self.position < self.input.len() {
            let byte = self.input[self.position];
            self.position += 1;
            Some(byte)
        } else {
            None
        }
    }

    /// Peek at the current byte without consuming it.
    /// If at end of input, return None.
    fn peek(&self) -> Option<u8> {
        if self.position < self.input.len() {
            let byte = self.input[self.position];
            Some(byte)
        } else {
            None
        }
    }

    /// If the next byte matches `expected`, consume it and return true.
    /// Otherwise, do nothing and return false.
    fn match_bytes(&mut self, bytes: &[u8]) -> bool {
        if self.input.get(self.position..self.position + bytes.len()) == Some(bytes) {
            self.position += bytes.len();
            true
        } else {
            false
        }
    }

    /// Consume whitespace characters (space, newline, tab, carriage return).
    /// Should stop at the first non-whitespace byte.
    fn skip_whitespace(&mut self) {
        while self.position < self.input.len() && self.input[self.position].is_ascii_whitespace() {
            self.position += 1;
        }
    }

    /// See if there is a comment, skip until end of comment
    fn skip_comments(&mut self) {
        // Keep skipping comments as long as we see them (handles /*...*/ /*...*/)
        while self.match_bytes(b"/*") {
            // Scan until closing */
            loop {
                // If we reached end-of-input, stop (or you could mark INVALID elsewhere)
                if self.peek().is_none() {
                    return;
                }
                if self.match_bytes(b"*/") {
                    break;
                }
                self.advance();
            }
        }
    }

    fn skip_trivia(&mut self) {
        loop {
            let before = self.position;
            self.skip_whitespace();
            self.skip_comments();
            self.skip_whitespace();
            if self.position == before {
                break; // no progress => no more trivia
            }
        }
    }

    fn read_string(&mut self) -> String {
        let start = self.position;

        while let Some(b) = self.peek() {
            if b == b'"' {
                break;
            }
            self.advance();
        }
        self.advance();

        let v = self.input[start..self.position-1].to_vec();
        let s = String::from_utf8(v).unwrap();
        s
    }

    fn read_literal(&mut self) -> Option<i64> {
        let byte: u8 = match self.peek()? {
            b'\\' => {
                self.advance(); // consume '\'
                match self.peek()? {
                    b'n'  => { self.advance(); b'\n' } // 10
                    b't'  => { self.advance(); b'\t' } // 9
                    b'r'  => { self.advance(); b'\r' } // 13
                    b'\\' => { self.advance(); b'\\' } // 92
                    b'\'' => { self.advance(); b'\'' } // 39
                    _ => return None,
                }
            }
            b'\n' | b'\r' | b'\'' => return None,
            b => {
                self.advance();
                b
            }
        };

        // must end with '
        if self.peek() != Some(b'\'') {
            return None;
        }
        self.advance(); // consume closing '

        Some(byte as i64)
    }

    /// Read an identifier starting with the already-consumed first byte.
    /// Should consume [a-zA-Z0-9_]* and return the full identifier.
    fn read_identifier(&mut self) -> String {
        let start = self.position;

        while let Some(b) = self.peek() {
            if !self.ident_re.is_match(&[b]) {
                break;
            }
            self.advance();
        }

        let v = self.input[start..self.position].to_vec();
        let s = String::from_utf8(v).unwrap();
        s
    }

    /// Read a number starting with the already-consumed first digit.
    /// Should consume consecutive digits and return the parsed value.
    fn read_number(&mut self) -> i64 {
        let start = self.position;

        while let Some(b) = self.peek() {
            if !self.digit_re.is_match(&[b]) {
                break;
            }
            self.advance();
        }

        let v = self.input[start..self.position].to_vec();
        let num: i64 = std::str::from_utf8(&v).unwrap().parse().unwrap();
        num
    }
}

fn main() -> io::Result<()> {
    let input = fs::read("tests/7.t")?;
    let mut lexer = Lexer::new(input);

    let mut token = lexer.next_token();

    while token != Token::EOF {
        println!("{:?}", token);
        token = lexer.next_token();
    }

    Ok(())
}
