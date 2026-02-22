use std::{env, fs, io};

use tiny_lang::lexer::tokenize;
use tiny_lang::{Parser, parse_statement};

fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let input = fs::read_to_string(args[1].as_str())?;
    let source = input.clone();
    let tokens = tokenize(source);
       
    println!("Program:\n{}", input);
    println!("-------------------------");
    // println!("Tokens:");
    // for token in tokens {
    //     println!("{:?}", token);
    // }

    let mut parser = Parser::new(tokens);
    let stmt = parse_statement(&mut parser);

    println!("{:#?}", stmt);


    Ok(())
}