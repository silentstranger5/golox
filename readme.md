# Lox, dynamic interpreted programming language

This Abstract Syntax Tree Interpreter is carefully crafted using Go programming language.

Sources:

- [Crafting Interpreters](https://craftinginterpreters.com) book, Robert Nystrom.
- [jlox](https://github.com/munificent/craftinginterpreters/tree/master/java/com/craftinginterpreters), Java Lox implementation.

## How to use

```
git clone https://github.com/silentstranger5/golox
cd golox
go run lox # run lox interpreter on the fly
go build lox # or compile it
./lox
# you can also execute files
go run lox test.lox
```

```
> print "Hello, world!";
Hello, world!
> ...
```

Familiarize yourself with syntax and capabilities of the Lox Programming Language [here](https://craftinginterpreters.com/the-lox-language.html)
