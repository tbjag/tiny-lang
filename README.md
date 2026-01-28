# Tiny Compiler Project - Rust

I really liked my compiler class. Specifically all of the CS problems: with
tokenizing, holding IR, and finally optimizing.This project is me giving all
of the concepts another try now that I know the end goal - and I am a better
programmer than before.

Based off of the tiny language because I threw away the spec of class project.

## Next Steps

Refactor Lexer:

- Add enums to type
- Remove debug from tokens
- move keywords into symbol?
- Write tests

Rest of Project:

- Build parser
- Generate instructions for WASM or FASM
- Build out SSA, generate .dot graphs
- Generate instructions

## Links

Specifications:

- [Lexer](https://rosettacode.org/wiki/Compiler/lexical_analyzer)
- [Parser](https://rosettacode.org/wiki/Compiler/syntax_analyzer)
- [Code Generator](https://rosettacode.org/wiki/Compiler/code_generator)
- [Interpreter](https://rosettacode.org/wiki/Compiler/AST_interpreter)
- [VM Interpreter](https://rosettacode.org/wiki/Compiler/virtual_machine_interpreter)
