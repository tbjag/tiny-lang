use std::{env, fs, io};

use tiny_lang::lexer::tokenize;


fn main() -> io::Result<()> {
    let args: Vec<String> = env::args().collect();
    let input = fs::read_to_string(args[1].as_str())?;
    let source = input.clone();
    let tokens = tokenize(source);
       
    println!("Program:\n{}", input);
    println!("-------------------------");
    println!("Tokens:");
    for _token in tokens {
        println!("{:?}", _token);
    }

    Ok(())
}