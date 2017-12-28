# Go Interpreter 

# Parser 

### Expressions 
An expression is a combination of one or more of the following: 
* values 
* constants 
* variables 
* operators 
* functions 

that a programming language interprets (based on precendence and association) and computes another value.

In Monkey Langauge anything other than a let and return statement is an expression. There are different types of expressions.

#### prefix operators 
* -5 
* ++i 
* !true 

#### binary infix operators 
* a+5 
* 10-5 
* 10/2 

#### comparison operator 
* a == 3 
* foo < bar 
* foo >= 13 

#### Call and group expressions 
* 5 * (5-3) 
* 5/(3*9)
* add(foo, bar) 
* foo / (foo+bar) 

#### Function literals are also expressions
* let add = fn(x, y) { return x+y }; 
* if-else expressions as well 

`  let result = if (10>5) { 
                true } else {
                  false
                }`
                

### Pratt Parsing 
Also called top down precedence based parsing 
In the Monkey language there are two basic type of statements 
1. Let statements 
2. Return statements 

rest of the language is mainly expressions therefore the main parser program needs to have sections to handle: 
* statements 
* expression statements - these are statements that wrap expressions.  

There are two parsing functions related to each token type: 
1. When the token is used in a prefix expression 
2. when the token is used in an infix expression 

Following is how the pratt parsing function will implement parsing.

#### 1. Identifiers 
Identifiers are simplest expressions they can be used either standalone or as part of a context with other expressions
* add(foo, bar)  - function calls
* foobar + bar  - expression prefix 
* if (foo) { ... }  - conditionals.

in all these context the identifier will eventually evaluate to value. 

#### 2. Integer Literals 
Similar to the identifier the integer literals are also form part of the expressions. 
* add (5, 3)
* let x = 5; 
* 5 + 5;
* if (x > 5) {  ...} 

#### 3. Prefix Operators 
There are two prefix operators supported - and ! 
e.g. 
* 5 + -10 
* !true 
* -(5+3)
* !GreaterThan2(3)

There syntax of a prefix operator is  <prefix operator><expression> therefore the ast node for a prefix operator must be able to point to an expression. 


#### 4. Infix Operators - binary expression. 
These involve one operator and two literals or operands. 
<expression> <operator> <expression> 
e.g.  5 == 5 , 5 < 5, 5 <= 5 etc. 




