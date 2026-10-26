# Tiny Lang
Practice with Go as well as writing my own compiler. Inspired by my advanced compiler class by Michael Franz.


## References:
- https://rosettacode.org/wiki/Compiler/lexical_analyzer
- https://rosettacode.org/wiki/Compiler/syntax_analyzer
- https://rosettacode.org/wiki/Compiler/code_generator
- https://rosettacode.org/wiki/Compiler/Sample_programs

## Notes:
- Re-used the same tokens for the lexer as well as the ast since there was not much a difference with them.
- Instead of parsing from scratch between tasks, send the bare instructions to save on work.
- Definitely not tested and robust enough for things like variable assignment.
- At some point I want to implement SSA which isn't described in the original spec.

## Future todos:
- Need to clean/refactor some of the variables and debug print statements.
- Dot notation for clean IR graphs.
- implement basic ssa on variable assignment