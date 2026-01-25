use regex::Regex;
use crate::token::{Token, TokenKind};

#[derive(Clone)]
enum Handler {
    Default(TokenKind, String),
    Skip,
    String,
    Character,
    Symbol(TokenKind),
}

struct RegexPattern {
    regex: Regex,
    handler: Handler
}

pub struct Lexer {
    patterns: Vec<RegexPattern>,
    pub tokens: Vec<Token>,
    source: String,
    pos: usize
}

impl Lexer {
    fn advance_n(&mut self, n: usize) {
        self.pos += n;
    }

    fn push(&mut self, token: Token) {
        self.tokens.push(token);
    }

    fn remainder(&self) -> &str {
        &self.source[self.pos..]
    }

    fn at_eof(&self) -> bool {
        self.pos >= self.source.len()
    }

    fn char_literal_to_number(&self, s: &str) -> u32 {
        // Remove the surrounding quotes
        let inner = &s[1..s.len()-1];
        
        // Handle escape sequences
        let ch = if inner.starts_with('\\') {
            match &inner[1..] {
                "n" => '\n',
                "t" => '\t',
                "r" => '\r',
                "\\" => '\\',
                "'" => '\'',
                "0" => '\0',
                _ => panic!("Unknown escape sequence"),
            }
        } else {
            inner.chars().next().unwrap()
    };
    
    ch as u32  // or ch as i32 if you want signed
}

    fn handle_pattern(&mut self, handler: &Handler, regex: &Regex) {
        match handler {
            Handler::Default(kind, value) => {
                self.advance_n(value.len());
                self.push(Token::new(*kind, value.clone()));
            }
            Handler::Skip => {
                if let Some(mat) = regex.find(self.remainder()) {
                    self.advance_n(mat.end());
                }
            }
            Handler::String => {
                if let Some(mat) = regex.find(self.remainder()) {
                    let match_str = mat.as_str();
                    let match_str = &match_str[1..match_str.len()-1].to_string();
                    let len = match_str.len() + 2;
                    self.push(Token::new(TokenKind::String, match_str));
                    self.advance_n(len);
                }
            }
            Handler::Character => {
                if let Some(mat) = regex.find(self.remainder()) {
                    let match_str = mat.as_str();
                    let char = self.char_literal_to_number(match_str);
                    let len = match_str.len() + 2;
                    self.push(Token::new(TokenKind::Integer, char.to_string()));
                    self.advance_n(len);
                }
            }
            Handler::Symbol(kind) => {
                if let Some(mat) = regex.find(self.remainder()) {
                    let match_str = mat.as_str().to_string();
                    let len = match_str.len();
                    self.push(Token::new(*kind, match_str));
                    self.advance_n(len);
                }
            }
        }
    }
}

pub fn tokenize(source: impl Into<String>) -> Vec<Token> {
    let mut lexer = create_lexer(source);

    while !lexer.at_eof() {
        let mut matched = false;
        let mut match_info = None;
        for i in 0..lexer.patterns.len() {
            let remainder = lexer.remainder();
            if let Some(mat) = lexer.patterns[i].regex.find(remainder) {
                if mat.start() == 0 {
                    match_info = Some(i);
                    break;
                }
            }
        }

        if let Some(i) = match_info {
            let handler = lexer.patterns[i].handler.clone();
            let regex = lexer.patterns[i].regex.clone();
            lexer.handle_pattern(&handler, &regex);
            matched = true;
        }

        if !matched {
            panic!("Lexer::Error -> unrecognized token near {}", lexer.remainder());
        }
    }

    lexer.push(Token::new(TokenKind::EndOfInput, "EOF"));

    lexer.tokens
}

fn create_lexer(source: impl Into<String>) -> Lexer {
    Lexer {
        pos: 0,
        source: source.into(),
        tokens: Vec::new(),
        patterns: vec![
            RegexPattern {
                regex: Regex::new(r"print").unwrap(),
                handler: Handler::Default(TokenKind::KeywordPrint, "print".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"putc").unwrap(),
                handler: Handler::Default(TokenKind::KeywordPutc, "putc".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"while").unwrap(),
                handler: Handler::Default(TokenKind::KeywordWhile, "while".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"if").unwrap(),
                handler: Handler::Default(TokenKind::KeywordIf, "if".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"else").unwrap(),
                handler: Handler::Default(TokenKind::KeywordElse, "else".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"[_a-zA-Z][_a-zA-Z0-9]*").unwrap(),
                handler: Handler::Symbol(TokenKind::Indentifier),
            },
            RegexPattern {
                regex: Regex::new(r"[0-9]+").unwrap(),
                handler: Handler::Symbol(TokenKind::Integer),
            },
            RegexPattern {
                regex: Regex::new(r#""[^"]*""#).unwrap(),
                handler: Handler::String,
            },
            RegexPattern {
                regex: Regex::new(r"'([^'\n]|\\n|\\\\)'").unwrap(),
                handler: Handler::Character,
            },
            RegexPattern {
                regex: Regex::new(r"(?s)/\*.*?\*/").unwrap(),
                handler: Handler::Skip,
            },
            RegexPattern {
                regex: Regex::new(r"\s+").unwrap(),
                handler: Handler::Skip,
            },
            RegexPattern {
                regex: Regex::new(r"\(").unwrap(),
                handler: Handler::Default(TokenKind::OpenParen, "(".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\)").unwrap(),
                handler: Handler::Default(TokenKind::CloseParen, ")".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\{").unwrap(),
                handler: Handler::Default(TokenKind::OpenBrace, "{".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\}").unwrap(),
                handler: Handler::Default(TokenKind::CloseBrace, "}".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"==").unwrap(),
                handler: Handler::Default(TokenKind::OpEqual, "==".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"!=").unwrap(),
                handler: Handler::Default(TokenKind::OpNotEqual, "!=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"=").unwrap(),
                handler: Handler::Default(TokenKind::OpAssign, "=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"!").unwrap(),
                handler: Handler::Default(TokenKind::OpNot, "!".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"<=").unwrap(),
                handler: Handler::Default(TokenKind::OpLessEqual, "<=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"<").unwrap(),
                handler: Handler::Default(TokenKind::OpLess, "<".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r">=").unwrap(),
                handler: Handler::Default(TokenKind::OpGreaterEqual, ">=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r">").unwrap(),
                handler: Handler::Default(TokenKind::OpGreater, ">".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"&&").unwrap(),
                handler: Handler::Default(TokenKind::OpAnd, "&&".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\|\|").unwrap(),
                handler: Handler::Default(TokenKind::OpOr, "||".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r";").unwrap(),
                handler: Handler::Default(TokenKind::Semicolon, ";".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r",").unwrap(),
                handler: Handler::Default(TokenKind::Comma, ",".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\+").unwrap(),
                handler: Handler::Default(TokenKind::OpAdd, "+".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"-").unwrap(),
                handler: Handler::Default(TokenKind::OpSubtract, "-".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"/").unwrap(),
                handler: Handler::Default(TokenKind::OpDivide, "/".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\*").unwrap(),
                handler: Handler::Default(TokenKind::OpMultiply, "*".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"%").unwrap(),
                handler: Handler::Default(TokenKind::OpMod, "%".to_string()),
            },
        ],
    }
}