use std::{env, fs, io};

use tiny_lang::{Lexer, Token};

fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let input = fs::read(args[1].as_str())?;
    let mut lexer = Lexer::new(input);

    let mut token = lexer.next_token();

    while token != Token::EOF {
        println!("{:?}", token);
        token = lexer.next_token();
    }

    Ok(())
}
