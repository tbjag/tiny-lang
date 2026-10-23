# Tiny Lang
Practice with Go as well as writing my own compiler. Inspired by my advanced compiler class by Michael Franz.


References:
- https://rosettacode.org/wiki/Compiler/lexical_analyzer
- https://rosettacode.org/wiki/Compiler/syntax_analyzer
- https://rosettacode.org/wiki/Compiler/code_generator


I re-used the same tokens for the lexer as well as the ast since there was not much a difference with them.


At some point I want to implement SSA which isn't described in the original spec.

## temp todos
- finish assign error handling during code gen
- rewrite emit output to match rosetta code version
- write tests agains all rosetta code
- onto code generation!