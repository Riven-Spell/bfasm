# BFasm
### Simple language that compiles to brainfuck.

- Memory is not used in a tape, but there are variables.
- Variables are a single cell, unless written as var[number], in which case multiple ones are reserved.
- Variables are defined at the start of the program.
- Note: they might not be in the same order as defined at the start, and/or they might have data inbetween.
- Comments are defined as // until end of line.
- Every operation is on its own line.
- Operations are case-insensitive. Variable names are case-sensitive.

### Supported operations
```
WHILE   var
IF 	var
UNTIL   var
END

SET     var      val
CPY     srcvar   destvar
ADD     var      (var/val)
SUB     var      (var/val)
MUL	var	 (var/val)
DIV	var	 (var/val)

//both implemented by means of a simple . or ,
READ    var      (number of chars to read)
PRINT   var      (number of chars to print)

//assumes terminates with memory in the same position and that it doesnt mess with code outside the var.
BF      memaddr  code
```

### Variable definition
At the start of the program variables are defined as follows:
```
varname
arrname[arrsize]
```
Variable/program are seperated by a line that contains nothing but `!`.

### Input methods
#### val
- Decimal numbers: `215`
- Hexadecimal numbers: `0x36`
- Chars/strings (' and " are equivalent): `'string'` `"c"`
#### var
- Variable: `varname`
- Array: `arrname[val]`

### Example program
```
helloworld[12]
!
SET helloworld "Hello World!"
PRINT helloworld 12
```
