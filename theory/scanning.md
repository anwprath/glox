# Scanning

The scanner takes in raw source code as a series of characters and groups it into a series of chunks we call tokens. These are the meaningful “words” and “punctuation” that make up the language’s grammar.


## Error reporting

```java
  static void error(int line, String message) {
    report(line, "", message);
  }

  private static void report(int line, String where,
                             String message) {
    System.err.println(
        "[line " + line + "] Error" + where + ": " + message);
    hadError = true;
  }
```
This error() function and its report() helper tells the user some syntax error occurred on a given line. That is really the bare minimum to be able to claim you even have error reporting.

```java
// lox/Lox.java
// in class Lox

public class Lox {

  static boolean hadError = false;
  
  ```

We’ll use this to ensure we don’t try to execute code that has a known error. Also, it lets us exit with a non-zero exit code like a good command line citizen should.

```java
// lox/Lox.java
// in runFile()
run(new String(bytes, Charset.defaultCharset()));


// Indicate an error in the exit code.
if (hadError) System.exit(65);


  ```

We need to reset this flag in the interactive loop. If the user makes a mistake, it shouldn’t kill their entire session.



```java
// lox/Lox.java
// in runPrompt()
run(line);
hadError = false;

   
```

### 4.2 Lexemes and Tokens

```var language = "lox";```

Here, var is the keyword for declaring a variable. That three-character sequence “v-a-r” means something. But if we yank three letters out of the middle of language, like “g-u-a”, those don’t mean anything on their own.

That’s what lexical analysis is about. Our job is to scan through the list of characters and group them together into the smallest sequences that still represent something. Each of these blobs of characters is called a lexeme. The lexemes are only the raw substrings of the source code. However, in the process of grouping character sequences into lexemes, we also stumble upon some other useful information.  When we take the lexeme and bundle it together with that other data, the result is a token.

#### 4.2.1 Token Type

 The parser wants to know not just that it has a lexeme for some identifier, but that it has a reserved word, and which keyword it is.

The parser could categorize tokens from the raw lexeme by comparing the strings, but that’s slow and kind of ugly. Instead, at the point that we recognize a lexeme, we also remember which kind of lexeme it represents. We have a different type for each keyword, operator, bit of punctuation, and literal type:

```java
// lox/TokenType.java
package com.craftinginterpreters.lox;

enum TokenType {
  // Single-character tokens.
  LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE,
  COMMA, DOT, MINUS, PLUS, SEMICOLON, SLASH, STAR,

  // One or two character tokens.
  BANG, BANG_EQUAL,
  EQUAL, EQUAL_EQUAL,
  GREATER, GREATER_EQUAL,
  LESS, LESS_EQUAL,

  // Literals.
  IDENTIFIER, STRING, NUMBER,

  // Keywords.
  AND, CLASS, ELSE, FALSE, FUN, FOR, IF, NIL, OR,
  PRINT, RETURN, SUPER, THIS, TRUE, VAR, WHILE,

  EOF
}
```

#### 4.2.2 Literal value

There are lexemes for literal values—numbers and strings and the like. Since the scanner has to walk each character in the literal to correctly identify it, it can also convert that textual representation of a value to the living runtime object that will be used by the interpreter later.

#### 4.2.3 Location information

Back when I was preaching the gospel about error handling, we saw that we need to tell users where errors occurred. Tracking that starts here. In our simple interpreter, we note only which line the token appears on, but more sophisticated implementations include the column and length too.

We take all of this data and wrap it in a class.
```java
//lox/Token.java
package com.craftinginterpreters.lox;

class Token {
  final TokenType type;
  final String lexeme;
  final Object literal;
  final int line; 

  Token(TokenType type, String lexeme, Object literal, int line) {
    this.type = type;
    this.lexeme = lexeme;
    this.literal = literal;
    this.line = line;
  }

  public String toString() {
    return type + " " + lexeme + " " + literal;
  }
}
```
Now we have an object with enough structure to be useful for all of the later phases of the interpreter.

### 4.3 Regular Languages and Expressions

Starting at the first character of the source code, the scanner figures out what lexeme the character belongs to, and consumes it and any following characters that are part of that lexeme. When it reaches the end of that lexeme, it emits a token.

Then it loops back and does it again, starting from the very next character in the source code. 


### 4.4 The Scanner Class

Without further ado, let’s make ourselves a scanner.
```java
// lox/Scanner.java
// create new file
...
import static com.craftinginterpreters.lox.TokenType.*; 

class Scanner {
  private final String source;
  private final List<Token> tokens = new ArrayList<>();

  Scanner(String source) {
    this.source = source;
  }
}
```

We store the raw source code as a simple string, and we have a list ready to fill with tokens we’re going to generate. The aforementioned loop that does that looks like this:
```java
// lox/Scanner.java
// add after Scanner()

  List<Token> scanTokens() {
    while (!isAtEnd()) {
      // We are at the beginning of the next lexeme.
      start = current;
      scanToken();
    }

    tokens.add(new Token(EOF, "", null, line));
    return tokens;
  }
```
The scanner works its way through the source code, adding tokens until it runs out of characters. Then it appends one final “end of file” token. That isn’t strictly needed, but it makes our parser a little cleaner.

This loop depends on a couple of fields to keep track of where the scanner is in the source code.

  private final List<Token> tokens = new ArrayList<>();
```java
// lox/Scanner.java
// in class Scanner

  private int start = 0;
  private int current = 0;
  private int line = 1;
```

The start and current fields are offsets that index into the string. The start field points to the first character in the lexeme being scanned, and current points at the character currently being considered. The line field tracks what source line current is on so we can produce tokens that know their location.

Little helper function that tells us if we’ve consumed all the characters:
```java
// lox/Scanner.java
// add after scanTokens()

  private boolean isAtEnd() {
    return current >= source.length();
  }
```