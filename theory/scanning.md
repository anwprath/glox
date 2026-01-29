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