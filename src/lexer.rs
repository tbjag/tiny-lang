use regex::Regex;
use crate::token::{Token};

#[derive(Clone)]
enum Handler {
    Default(Token, String),
    Skip,
    String,
    Character,
    Identifier,
    Integer
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
        let inner = &s[1..s.len()-1];
        
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
    
    ch as u32
}

    fn handle_pattern(&mut self, handler: &Handler, regex: &Regex) {
        match handler {
            Handler::Default(token, value) => {
                self.advance_n(value.len());
                self.push(token.clone());
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
                    self.push(Token::String(match_str.clone()));
                    self.advance_n(len);
                }
            }
            Handler::Character => {
                if let Some(mat) = regex.find(self.remainder()) {
                    let match_str = mat.as_str();
                    let char = self.char_literal_to_number(match_str);
                    let len = match_str.len() + 2;
                    self.push(Token::Integer(char.to_string()));
                    self.advance_n(len);
                }
            }
            Handler::Identifier => {
                if let Some(mat) = regex.find(self.remainder()) {
                    let match_str = mat.as_str().to_string();
                    let len = match_str.len();
                    self.push(Token::Indentifier(match_str));
                    self.advance_n(len);
                }
            }
            Handler::Integer => {
                if let Some(mat) = regex.find(self.remainder()) {
                    let match_str = mat.as_str().to_string();
                    let len = match_str.len();
                    self.push(Token::Integer(match_str));
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

    lexer.push(Token::EndOfInput);

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
                handler: Handler::Default(Token::KeywordPrint, "print".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"putc").unwrap(),
                handler: Handler::Default(Token::KeywordPutc, "putc".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"while").unwrap(),
                handler: Handler::Default(Token::KeywordWhile, "while".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"if").unwrap(),
                handler: Handler::Default(Token::KeywordIf, "if".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"else").unwrap(),
                handler: Handler::Default(Token::KeywordElse, "else".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"[_a-zA-Z][_a-zA-Z0-9]*").unwrap(),
                handler: Handler::Identifier,
            },
            RegexPattern {
                regex: Regex::new(r"[0-9]+").unwrap(),
                handler: Handler::Integer,
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
                handler: Handler::Default(Token::OpenParen, "(".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\)").unwrap(),
                handler: Handler::Default(Token::CloseParen, ")".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\{").unwrap(),
                handler: Handler::Default(Token::OpenBrace, "{".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\}").unwrap(),
                handler: Handler::Default(Token::CloseBrace, "}".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"==").unwrap(),
                handler: Handler::Default(Token::OpEqual, "==".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"!=").unwrap(),
                handler: Handler::Default(Token::OpNotEqual, "!=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"=").unwrap(),
                handler: Handler::Default(Token::OpAssign, "=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"!").unwrap(),
                handler: Handler::Default(Token::OpNot, "!".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"<=").unwrap(),
                handler: Handler::Default(Token::OpLessEqual, "<=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"<").unwrap(),
                handler: Handler::Default(Token::OpLess, "<".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r">=").unwrap(),
                handler: Handler::Default(Token::OpGreaterEqual, ">=".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r">").unwrap(),
                handler: Handler::Default(Token::OpGreater, ">".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"&&").unwrap(),
                handler: Handler::Default(Token::OpAnd, "&&".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\|\|").unwrap(),
                handler: Handler::Default(Token::OpOr, "||".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r";").unwrap(),
                handler: Handler::Default(Token::Semicolon, ";".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r",").unwrap(),
                handler: Handler::Default(Token::Comma, ",".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\+").unwrap(),
                handler: Handler::Default(Token::OpAdd, "+".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"-").unwrap(),
                handler: Handler::Default(Token::OpSubtract, "-".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"/").unwrap(),
                handler: Handler::Default(Token::OpDivide, "/".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"\*").unwrap(),
                handler: Handler::Default(Token::OpMultiply, "*".to_string()),
            },
            RegexPattern {
                regex: Regex::new(r"%").unwrap(),
                handler: Handler::Default(Token::OpMod, "%".to_string()),
            },
        ],
    }
}