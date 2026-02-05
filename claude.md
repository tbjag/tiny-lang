# Claude Code Context

## Project Overview

This is a compiler for "Tiny Language" written in Rust. It's a personal learning project reimplementing compiler concepts (lexing, parsing, IR, optimization) based on a university compiler course.

The language specification is based on the [Rosetta Code compiler series](https://rosettacode.org/wiki/Compiler/lexical_analyzer).

## Current Status

**Completed:**
- Lexer (tokenizer) - fully functional

**In Progress:**
- Lexer refactoring (move keywords to symbol table, add unit tests)

**Planned:**
- Parser (syntax analyzer, AST generation)
- Code generator (WASM or FASM targets)
- SSA-based IR with .dot graph output
- Optimizations

## Architecture

```
src/
├── main.rs      # CLI entry point - reads file, tokenizes, prints tokens
├── lib.rs       # Library exports (Token, tokenize)
├── lexer.rs     # Lexer implementation with regex-based pattern matching
└── token.rs     # Token enum definition (32 variants)

tests/
└── *.lang       # Sample Tiny Language programs for testing
```

## Key Components

### Token Types (`src/token.rs`)
- Keywords: `if`, `else`, `while`, `print`, `putc`
- Operators: `+`, `-`, `*`, `/`, `%`, `<`, `>`, `<=`, `>=`, `==`, `!=`, `&&`, `||`, `!`, `=`
- Delimiters: `(`, `)`, `{`, `}`, `;`, `,`
- Literals: `Identifier(String)`, `Integer(String)`, `String(String)`
- Special: `EndOfInput`

### Lexer (`src/lexer.rs`)
- Uses `regex` crate for pattern matching
- Patterns are matched in priority order (keywords before identifiers)
- `Handler` enum determines how each match is processed:
  - `Default(Token, usize)` - fixed token with length
  - `Skip` - whitespace/comments
  - `String`, `Character`, `Identifier`, `Integer` - extract content

### Public API
```rust
use tiny_lang::{Token, tokenize};
let tokens: Vec<Token> = tokenize("print(42);");
```

## Running

```bash
# Build
cargo build

# Run on a test file
cargo run -- tests/06.lang

# Run tests (when added)
cargo test
```

## Language Features

- Integer literals (including negative: `-42`)
- String literals (`"hello"`)
- Character literals converted to integers (`'a'` -> `97`)
- Escape sequences: `\n`, `\t`, `\r`, `\\`, `\'`, `\0`
- Multi-line comments (`/* ... */`)
- C-style control flow (`if`, `else`, `while`)
- Print statements (`print`, `putc`)

## Development Notes

- Keywords must be matched before the identifier pattern to avoid being captured as identifiers
- Character literals are converted to their numeric value and stored as `Integer` tokens
- The lexer panics on unrecognized input with position info
- Note: `Token::Indentifier` has a typo (should be `Identifier`)
