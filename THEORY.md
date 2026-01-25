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
