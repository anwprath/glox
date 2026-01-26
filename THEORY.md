# Theory

### Just snippets I found relevant

### The Parts of a Language
##### _http://localhost:8000/a-map-of-the-territory.html_

#### 2.1.1 Scanning 
A scanner (or lexer) takes in the linear stream of characters and chunks them together into a series of something more akin to “words”. In programming languages, each of these words is called a token. Some tokens are single characters, like ( and ,. Others may be several characters long, like numbers (123), string literals ("hi!"), and identifiers (min).

#### 2.1.2 Parsing

A parser takes the flat sequence of tokens and builds a tree structure that mirrors the nested nature of the grammar. These trees have a couple of different names—parse tree or abstract syntax tree—depending on how close to the bare syntactic structure of the source language they are. In practice, language hackers usually call them syntax trees, ASTs, or often just trees.

#### 2.1.3 Static analysis

In an expression like a + b, we know we are adding a and b, but we don’t know what those names refer to. Are they local variables? Global? Where are they defined?

The first bit of analysis that most languages do is called binding or resolution. For each identifier, we find out where that name is defined and wire the two together. This is where scope comes into play—the region of source code where a certain name can be used to refer to a certain declaration.

If the language is statically typed, this is when we type check. Once we know where a and b are declared, we can also figure out their types. Then if those types don’t support being added to each other, we report a type error.

All this semantic insight that is visible to us from analysis needs to be stored somewhere:

- Often, it gets stored right back as attributes on the syntax tree itself—extra fields in the nodes that aren’t initialized during parsing but get filled in later.

- Other times, we may store data in a lookup table off to the side. Typically, the keys to this table are identifiers—names of variables and declarations. In that case, we call it a symbol table and the values it associates with each key tell us what that identifier refers to.

- The most powerful bookkeeping tool is to transform the tree into an entirely new data structure that more directly expresses the semantics of the code. 

Everything up to this point is considered the front end of the implementation. 

#### 2.1.4 Intermediate representations

The front end of the pipeline is specific to the source language the program is written in. The back end is concerned with the final architecture where the program will run.

In the middle, the code may be stored in some intermediate representation (IR) that isn’t tightly tied to either the source or destination forms (hence “intermediate”). Instead, the IR acts as an interface between these two languages.

This lets you support multiple source languages and target platforms with less effort. Say you want to implement Pascal, C, and Fortran compilers, and you want to target x86, ARM, and, I dunno, SPARC. Normally, that means you’re signing up to write nine full compilers: Pascal→x86, C→ARM, and every other combination.

A shared intermediate representation reduces that dramatically. You write one front end for each source language that produces the IR. Then one back end for each target architecture. Now you can mix and match those to get every combination.


#### 2.1.5 Optimization


A simple example is constant folding: if some expression always evaluates to the exact same value, we can do the evaluation at compile time and replace the code for the expression with its result. If the user typed in this:

`pennyArea = 3.14159 * (0.75 / 2) * (0.75 / 2);`

we could do all of that arithmetic in the compiler and change the code to:

`pennyArea = 0.4417860938;`

Some keywords to get you started are:

- constant propagation
- common subexpression elimination
- loop invariant code motion
- global value numbering
- strength reduction
- scalar replacement of aggregates
- dead code elimination
- loop unrolling

#### 2.1.6 Code generation

The last step is converting it to a form the machine can actually run. In other words, generating code (or code gen), where “code” here usually refers to the kind of primitive assembly-like instructions a CPU runs and not the kind of “source code” a human might want to read.

Finally, we are in the back end, descending the other side of the mountain. From here on out, our representation of the code becomes more and more primitive, like evolution run in reverse, as we get closer to something our simple-minded machine can understand.

Hackers like Martin Richards and Niklaus Wirth, of BCPL and Pascal fame, respectively, made their compilers produce virtual machine code. Instead of instructions for some real chip, they produced code for a hypothetical, idealized machine. Wirth called this p-code for portable, but today, we generally call it bytecode because each instruction is often a single byte long.


#### 2.1.7 (Language) Virtual Machine

Virtual machine (VM), a program that emulates a hypothetical chip supporting your virtual architecture at runtime. Running bytecode in a VM is slower than translating it to native code ahead of time because every instruction must be simulated at runtime each time it executes


#### 2.1.8 Runtime 

We usually need some services that our language provides while the program is running. For example, if the language automatically manages memory, we need a garbage collector going in order to reclaim unused bits. 

### 2.2 Shortcuts and Alternate Routes

#### 2.2.1 Single-pass compilers

Some simple compilers interleave parsing, analysis, and code generation so that they produce output code directly in the parser, without ever allocating any syntax trees or other IRs. 


#### 2.2.2 Tree-walk interpreters

Some programming languages begin executing code right after parsing it to an AST (with maybe a bit of static analysis applied). To run the program, the interpreter traverses the syntax tree one branch and leaf at a time, evaluating each node as it goes.


#### 2.2.3 Transpilers

You write a front end for your language. Then, in the back end, instead of doing all the work to lower the semantics to some primitive target language, you produce a string of valid source code for some other language that’s about as high level as yours. Then, you use the existing compilation tools for that language as your escape route off the mountain and down to something you can execute.

They used to call this a source-to-source compiler or a transcompiler.

#### 2.2.4 Just-in-time compilation

On the end user’s machine, when the program is loaded—either from source in the case of JS, or platform-independent bytecode for the JVM and CLR—you compile it to native code for the architecture their computer supports. Naturally enough, this is called just-in-time compilation.

The most sophisticated JITs insert profiling hooks into the generated code to see which regions are most performance critical and what kind of data is flowing through them. Then, over time, they will automatically recompile those hot spots with more advanced optimizations.


### 2.3 Compilers and Interpreters 



- Compiling is an implementation technique that involves translating a source language to some other—usually lower-level—form. When you generate bytecode or machine code, you are compiling. When you transpile to another high-level language, you are compiling too.

- When we say a language implementation “is a compiler”, we mean it translates source code to some other form but doesn’t execute it. The user has to take the resulting output and run it themselves.

- Conversely, when we say an implementation “is an interpreter”, we mean it takes in source code and executes it immediately. It runs programs “from source”.


### 3.0 The Lox Language
###### http://localhost:8000/the-lox-language.html

#### 3.1 Hello, Lox

Here’s your very first taste of Lox:

```C
// Your first Lox program!
print "Hello, world!";
```

#### 3.2.1 Dynamic typing

Lox is dynamically typed. Variables can store values of any type, and a single variable can even store values of different types at different times. If you try to perform an operation on values of the wrong type—say, dividing a number by a string—then the error is detected and reported at runtime.


#### 3.2.2 Automatic memory management

High-level languages exist to eliminate error-prone, low-level drudgery, and what could be more tedious than manually managing the allocation and freeing of storage? 

There are two main techniques for managing memory: reference counting and tracing garbage collection (usually just called garbage collection or GC). Ref counters are much simpler to implement—I think that’s why Perl, PHP, and Python all started out using them. But, over time, the limitations of ref counting become too troublesome. All of those languages eventually ended up adding a full tracing GC, or at least enough of one to clean up object cycles.

### 3.3 Data Types

1. Booleans: Lox has a dedicated Boolean type. 

    ```C
    true;  // Not false.
    false; // Not *not* false.
    ```

2. Numbers: Lox has only one kind of number: double-precision floating point. Since floating-point numbers can also represent a wide range of integers, that covers a lot of territory, while keeping things simple.

    ```C
    1234;  // An integer.
    12.34; // A decimal number.C
    ```

3. Strings: Like most languages, they are enclosed in double quotes.
    ```C
    "I am a string";
    "";    // The empty string.
    "123"; // This is a string, not a number.
    ```

4. Nil: It represents “no value”. 


### 3.4 Expressions

#### 3.4.1 Arithmetic

```C
add + me;
subtract - me;
multiply * me;
divide / me;
```

The subexpressions on either side of the operator are operands. Because there are two of them, these are called binary operators. 

Because the operator is fixed in the middle of the operands, these are also called infix operators (as opposed to prefix operators where the operator comes before the operands, and postfix where it comes after).

One arithmetic operator is actually both an infix and a prefix one. The - operator can also be used to negate a number.

```C
-negateMe;
```

#### 3.4.2 Comparison and Equality


less < than;
lessThan <= orEqual;
greater > than;
greaterThan >= orEqual;

We can test two values of any kind for equality or inequality.

```C
1 == 2;         // false.
"cat" != "dog"; // true.
```

Even different types.

```C
314 == "pi"; // false.
```

Values of different types are never equivalent.

```C
123 == "123"; // false.
```

#### 3.4.3 Logical operators

The not operator, a prefix !, returns false if its operand is true, and vice versa.

```C
!true;  // false.
!false; // true.
```
The other two logical operators really are control flow constructs in the guise of expressions. An and expression determines if two values are both true. It returns the left operand if it’s false, or the right operand otherwise.

```C
true and false; // false.
true and true;  // true.
```

And an or expression determines if either of two values (or both) are true. It returns the left operand if it is true and the right operand otherwise.

```C
false or false; // false.
true or false;  // true.
```

The reason and and or are like control flow structures is that they short-circuit. Not only does and return the left operand if it is false, it doesn’t even evaluate the right one in that case. Conversely, if the left operand of an or is true, the right is skipped.


#### 3.4.4 Precedence and grouping

All of these operators have the same precedence and associativity that you’d expect coming from C. In cases where the precedence isn’t what you want, you can use () to group stuff.

```js
var average = (min + max) / 2;
```

#### 3.5 Statements

Where an expression’s main job is to produce a value, a statement’s job is to produce an effect. Since, by definition, statements don’t evaluate to a value, to be useful they have to otherwise change the world in some way—usually modifying some state, reading input, or producing output.

#### 3.6 Variables

You declare variables using `var` statements. If you omit the initializer, the variable’s value defaults to nil.

```C
var imAVariable = "here is my value";
var iAmNil;
```

Once declared, you can, naturally, access and assign a variable using its name.

```C
var breakfast = "bagels";
print breakfast; // "bagels".
breakfast = "beignets";
print breakfast; // "beignets".
```

#### 3.7 Control Flow

An if statement executes one of two statements based on some condition.

```C
if (condition) {
  print "yes";
} else {
  print "no";
}
```

Finally, we have for loops.

```C
for (var a = 1; a < 10; a = a + 1) {
  print a;
}
```

### 3.8 Functions

In Lox, you do that with `fun`.
```C
fun printSum(a, b) {
  print a + b;
}
```

From here on out:

 - An argument is an actual value you pass to a function when you call it. So a function call has an argument list. Sometimes you hear actual parameter used for these.

 - A parameter is a variable that holds the value of the argument inside the body of the function. Thus, a function declaration has a parameter list. Others call these formal parameters or simply formals.


### 3.9 Closures


Functions are first class in Lox, which just means they are real values that you can get a reference to, store in variables, pass around, etc. This works:

```C
fun addPair(a, b) {
  return a + b;
}

fun identity(a) {
  return a;
}

print identity(addPair)(1, 2); // Prints "3".
```
Since function declarations are statements, you can declare local functions inside another function.

```C
fun outerFunction() {
  fun localFunction() {
    print "I'm local!";
  }

  localFunction();
}
```

If you combine local functions, first-class functions, and block scope, you run into this interesting situation:
```C
fun returnFunction() {
  var outside = "outside";

  fun inner() {
    print outside;
  }

  return inner;
}

var fn = returnFunction();
fn();
```
Here, inner() accesses a local variable declared outside of its body in the surrounding function. Is this kosher? Now that lots of languages have borrowed this feature from Lisp, you probably know the answer is yes.

For that to work, inner() has to “hold on” to references to any surrounding variables that it uses so that they stay around even after the outer function has returned. We call functions that do this closures. These days, the term is often used for any first-class function, though it’s sort of a misnomer if the function doesn’t happen to close over any variables.


### 3.9 Classes

For a dynamically typed language, objects are pretty handy. We need some way of defining compound data types to bundle blobs of stuff together.

#### 3.9.3 Classes or Prototypes

In class-based languages, there are two core concepts: instances and classes. Instances store the state for each object and have a reference to the instance’s class. Classes contain the methods and inheritance chain. To call a method on an instance, there is always a level of indirection. You look up the instance’s class and then you find the method _there_

Prototype-based languages merge these two concepts. There are only objects—no classes—and each individual object may contain state and methods. Objects can directly inherit from each other (or “delegate to” in prototypal lingo):

#### 3.9.4 Classes in Lox

You declare a class and its methods like so:

```C
class Breakfast {
  cook() {
    print "Eggs a-fryin'!";
  }

  serve(who) {
    print "Enjoy your breakfast, " + who + ".";
  }
}
```

The body of a class contains its methods. They look like function declarations but without the fun keyword. When the class declaration is executed, Lox creates a class object and stores that in a variable named after the class. Just like functions, classes are first class in Lox.


```C
// Store it in variables.
var someVariable = Breakfast;

// Pass it to functions.
someFunction(Breakfast);
```
Call a class like a function, and it produces a new instance of itself.
```C
var breakfast = Breakfast();
print breakfast; // "Breakfast instance".
```

#### 3.9.5 Instantiation and Initialization

The idea behind object-oriented programming is encapsulating behavior and state together. To do that, you need fields. Lox, like other dynamically typed languages, lets you freely add properties onto objects.
```C
breakfast.meat = "sausage";
breakfast.bread = "sourdough";
```

Assigning to a field creates it if it doesn’t already exist.

If you want to access a field or method on the current object from within a method, you use good old this.

```C
class Breakfast {
  serve(who) {
    print "Enjoy your " + this.meat + " and " +
        this.bread + ", " + who + ".";
  }

  // ...
}
```
Part of encapsulating data within an object is ensuring the object is in a valid state when it’s created. To do that, you can define an initializer. If your class has a method named init(), it is called automatically when the object is constructed. Any parameters passed to the class are forwarded to its initializer.

```C
class Breakfast {
  init(meat, bread) {
    this.meat = meat;
    this.bread = bread;
  }

  // ...
}

var baconAndToast = Breakfast("bacon", "toast");
baconAndToast.serve("Dear Reader");
// "Enjoy your bacon and toast, Dear Reader."

```


#### 3.9.6 Inheritance

Every object-oriented language lets you not only define methods, but reuse them across multiple classes or objects. For that, Lox supports single inheritance. When you declare a class, you can specify a class that it inherits from using a less-than (<) operator.

```C
class Brunch < Breakfast {
  drink() {
    print "How about a Bloody Mary?";
  }
}
```

Why the < operator? I didn’t feel like introducing a new keyword like extends. Lox doesn’t use : for anything else so I didn’t want to reserve that either. Instead, I took a page from Ruby and used <.

Here, Brunch is the derived class or subclass, and Breakfast is the base class or superclass.

Even the init() method gets inherited. In practice, the subclass usually wants to define its own init() method too. But the original one also needs to be called so that the superclass can maintain its state. We need some way to call a method on our own instance without hitting our own methods.

```C
class Brunch < Breakfast {
  init(meat, bread, drink) {
    super.init(meat, bread);
    this.drink = drink;
  }
}
```

### 3.10 The Standard Library

- `print` statement
- clock() - epoch