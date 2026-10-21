# Code Generator Task (Rosetta Code "Compiler" series, C-like subset)

A summary of the task in my own words, written because rosettacode.org is down. It is not the original text, so check details against the real page once the site is back up.

## Where this fits

The task is the third stage of a pipeline for a small C-like language (keywords `if`, `else`, `while`, `print`, `putc`; integers, strings, identifiers, and the usual operators):

```
lex < prog.t | parse | gen | vm
```

1. **Lexer** turns source into tokens.
2. **Parser** turns tokens into a flattened AST.
3. **Code generator (this task)** turns the AST into text-form virtual assembly.
4. **VM** runs the assembly.

Read from a file and/or stdin, write to a file and/or stdout.

## Task

Read the parser's flattened AST and emit assembly for the stack-based VM.

## Input: flattened AST

One node per line, in pre-order. The first token on a line is the node type.

- `;` is a null node (no child).
- `Identifier name`, `Integer n` and `String "text"` are leaves, with the value as the second field.
- Any other node has no value on its line. Its left child follows, then its right child.

Loader sketch:

```
load_ast():
    line = readline(); tokens = split(line, respecting double quotes)
    if tokens[0] == ";": return null
    if len(tokens) > 1: return leaf(tokens[0], tokens[1])
    left = load_ast(); right = load_ast()
    return node(tokens[0], left, right)
```

Node types: Identifier, String, Integer, Sequence, If, Prtc, Prts, Prti, While, Assign, Negate, Not, Multiply, Divide, Mod, Add, Subtract, Less, LessEqual, Greater, GreaterEqual, Equal, NotEqual, And, Or.

Shapes worth knowing:
- `If`: left is the condition, right is a second `If` node whose left is the then-branch and whose right is the else-branch (or `;`).
- `While`: left is the condition, right is the body.
- `Assign`: left is an Identifier, right is the expression.
- `Sequence`: left is a statement, right is the rest.
- `Prtc` / `Prts` / `Prti` and `Negate` / `Not` are unary, with the child on the left and `;` on the right.

## Output format

1. Header: `Datasize: <N> Strings: <M>`
   - N is the number of distinct variables (each a 32-bit slot).
   - M is the number of distinct string constants.
2. The M string constants, one per line, in first-use order, with their quotes.
3. The code listing: one instruction per line, each starting with its byte address.

Variables are numbered 0, 1, 2… in first-use order, and strings the same way, with duplicates reused.

## VM model

- A 32-bit integer stack, with `sp` pointing at the next free slot.
- Bytecode in a byte array, with `pc` pointing at the current instruction.
- Each instruction is 1 byte, except these, which carry a 4-byte (32-bit) operand and so take 5 bytes: `fetch`, `store`, `push`, `jmp`, `jz`.

| Instruction | Operand | Effect |
|---|---|---|
| `fetch [i]` | variable index | push `data[i]` |
| `store [i]` | variable index | pop into `data[i]` |
| `push n` | 32-bit int | push n (also used for string indices) |
| `add sub mul div mod` | none | pop two, push result |
| `lt gt le ge eq ne` | none | pop two, push 1 or 0 |
| `and or` | none | pop two, push logical result |
| `neg` `not` | none | operate on top of stack |
| `jmp (n) addr` | relative offset | unconditional jump |
| `jz (n) addr` | relative offset | pop; jump if zero |
| `prtc` | none | pop, print as a character |
| `prti` | none | pop, print as an integer |
| `prts` | none | pop a string-pool index, print that string |
| `halt` | none | stop |

Opcode numbering is an internal matter. The reference solutions list them in the order above (fetch, store, push, add, sub, mul, div, mod, lt, gt, le, ge, eq, ne, and, or, neg, not, jmp, jz, prtc, prts, prti, halt).

### Jump display

`jmp (n) addr` and `jz (n) addr` show the signed relative offset `n` and the resolved absolute address `addr`. The offset is measured from the address of the operand, which is the instruction's address + 1, so `addr = (instruction address + 1) + n`.

Example from the reference output: `jz (43) 65` at address 21 gives 22 + 43 = 65, and `jmp (-51) 10` at address 60 gives 61 − 51 = 10.

## Generation rules

- Identifier: `fetch [var]`
- Integer: `push n`
- String: `push strIndex`
- Assign: generate the right side, then `store [var]`
- Sequence: left, then right
- Binary operators: left, then right, then the opcode
- Unary operators / Prtc / Prti / Prts: child, then the opcode
- If: condition, `jz` (patched to the end of the then-branch, or to the else label), then-branch. With an else, add a `jmp` over the else-branch and patch both jumps.
- While: remember the loop-top address, generate the condition, `jz` (patched to the loop exit), the body, then `jmp` back to the top. Patch the `jz` to the address after the `jmp`.
- End of program: `halt`

## Worked example

Source `while.t`:

```
count = 1;
while (count < 10) {
    print("count is: ", count, "\n");
    count = count + 1;
}
```

Expected `gen` output:

```
Datasize: 1 Strings: 2
"count is: "
"\n"
   0 push  1
   5 store [0]
  10 fetch [0]
  15 push  10
  20 lt
  21 jz     (43) 65
  26 push  0
  31 prts
  32 fetch [0]
  37 prti
  38 push  1
  43 prts
  44 fetch [0]
  49 push  1
  54 add
  55 store [0]
  60 jmp    (-51) 10
  65 halt
```

The listing I pulled from the page printed only `"count is: "` in the string section, even though the header says 2 strings. The `"\n"` line above is what a correct generator should print, and it is consistent with `push 1` before the second `prts`. This looks like a copy/format glitch in the page, so verify it when the site is back.

## Testing

Pass the sample programs from the task's companion page, "Compiler/Sample programs" (compare against the reference C and Python implementations).