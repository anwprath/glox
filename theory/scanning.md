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


### 4.5 Recognizing Lexemes

In each turn of the loop, we scan a single token. We’ll start simple. Imagine if every lexeme were only a single character long. All you would need to do is consume the next character and pick a token type for it. Several lexemes are only a single character in Lox, so let’s start with those.
```java
// lox/Scanner.java
// add after scanTokens()

  private void scanToken() {
    char c = advance();
    switch (c) {
      case '(': addToken(LEFT_PAREN); break;
      case ')': addToken(RIGHT_PAREN); break;
      case '{': addToken(LEFT_BRACE); break;
      case '}': addToken(RIGHT_BRACE); break;
      case ',': addToken(COMMA); break;
      case '.': addToken(DOT); break;
      case '-': addToken(MINUS); break;
      case '+': addToken(PLUS); break;
      case ';': addToken(SEMICOLON); break;
      case '*': addToken(STAR); break; 
    }
  }
```

Again, we need a couple of helper methods.
```java
// lox/Scanner.java
  private char advance() {
    return source.charAt(current++);
  }

  private void addToken(TokenType type) {
    addToken(type, null);
  }

  private void addToken(TokenType type, Object literal) {
    String text = source.substring(start, current);
    tokens.add(new Token(type, text, literal, line));
  }
```

The advance() method consumes the next character in the source file and returns it. Where advance() is for input, addToken() is for output. It grabs the text of the current lexeme and creates a new token for it. We’ll use the other overload to handle tokens with literal values soon.


#### 4.5.1 Lexical errors

Right now, erroneous characters get silently discarded. They aren’t used by the Lox language, but that doesn’t mean the interpreter can pretend they aren’t there. Instead, we report an error.


```java
// lox/Scanner.java
// in scanToken()
      default:
        Lox.error(line, "Unexpected character.");
        break;

    }
```

Note that the erroneous character is still consumed by the earlier call to advance(). That’s important so that we don’t get stuck in an infinite loop.

Note also that we keep scanning. There may be other errors later in the program. It gives our users a better experience if we detect as many of those as possible in one go.


#### 4.5.2 Operators

We have single-character lexemes working, but that doesn’t cover all of Lox’s operators. What about !? It’s a single character, right? Sometimes, yes, but if the very next character is an equals sign, then we should instead create a != lexeme. Note that the ! and = are not two independent operators. You can’t write ! = in Lox and have it behave like an inequality operator. That’s why we need to scan != as a single lexeme. Likewise, <, >, and = can all be followed by = to create the other equality and comparison operators.

For all of these, we need to look at the second character.
```java
      case '*': addToken(STAR); break; 

// lox/Scanner.java
// in scanToken()

      case '!':
        addToken(match('=') ? BANG_EQUAL : BANG);
        break;
      case '=':
        addToken(match('=') ? EQUAL_EQUAL : EQUAL);
        break;
      case '<':
        addToken(match('=') ? LESS_EQUAL : LESS);
        break;
      case '>':
        addToken(match('=') ? GREATER_EQUAL : GREATER);
        break;


      default:
```
Those cases use this new method:
```java
// lox/Scanner.java
// add after scanToken()

  private boolean match(char expected) {
    if (isAtEnd()) return false;
    if (source.charAt(current) != expected) return false;

    current++;
    return true;
  }
```

### 4.6 Longer Lexemes

We’re still missing one operator: / for division. That character needs a little special handling because comments begin with a slash too.

```java
// lox/Scanner.java
// in scanToken()

      case '/':
        if (match('/')) {
          // A comment goes until the end of the line.
          while (peek() != '\n' && !isAtEnd()) advance();
        } else {
          addToken(SLASH);
        }
        break;

```
This is similar to the other two-character operators, except that when we find a second /, we don’t end the token yet. Instead, we keep consuming characters until we reach the end of the line.

This is our general strategy for handling longer lexemes. After we detect the beginning of one, we shunt over to some lexeme-specific code that keeps eating characters until it sees the end.

We’ve got another helper:
```java
// lox/Scanner.java
// add after match()

  private char peek() {
    if (isAtEnd()) return '\0';
    return source.charAt(current);
  }
```

#### 4.6.1 String literals

We’ll do strings first, since they always begin with a specific character, `"`.

`      case '"': string(); break;`

That calls:
```java
// lox/Scanner.java

  private void string() {
    while (peek() != '"' && !isAtEnd()) {
      if (peek() == '\n') line++;
      advance();
    }

    if (isAtEnd()) {
      Lox.error(line, "Unterminated string.");
      return;
    }

    // The closing ".
    advance();

    // Trim the surrounding quotes.
    String value = source.substring(start + 1, current - 1);
    addToken(STRING, value);
  }
```

#### 4.6.2 Number literals

We don’t allow a leading or trailing decimal point, so these are both invalid:

```
.1234
1234.
```
We could easily support the former, but I left it out to keep things simple. The latter gets weird if we ever want to allow methods on numbers like 123.sqrt().


To recognize the beginning of a number lexeme, we look for any digit. It’s kind of tedious to add cases for every decimal digit, so we’ll stuff it in the default case instead.

      default:

```java
// lox/Scanner.java
// in scanToken()
// replace 1 line

        if (isDigit(c)) {
          number();
        } else {
          Lox.error(line, "Unexpected character.");
        }

        break;
```

This relies on this little utility:
```java
// lox/Scanner.java
// add after peek()

  private boolean isDigit(char c) {
    return c >= '0' && c <= '9';
  } 
```

Once we know we are in a number, we branch to a separate method to consume the rest of the literal, like we do with strings.

```java
// lox/Scanner.java
// add after scanToken()

  private void number() {
    while (isDigit(peek())) advance();

    // Look for a fractional part.
    if (peek() == '.' && isDigit(peekNext())) {
      // Consume the "."
      advance();

      while (isDigit(peek())) advance();
    }

    addToken(NUMBER,
        Double.parseDouble(source.substring(start, current)));
  }

```