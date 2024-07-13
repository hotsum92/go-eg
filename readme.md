[【Go】公式ツール eg を使って効率的にGoのコードをリファクタリングする](https://qiita.com/_ken_/items/0ee4e9a98b923a79418f)

[Apply transformations to Go code with eg rakyll.org](https://rakyll.org/eg/)

[sample template](https://github.com/golang/tools/tree/master/refactor/eg/testdata)

```
eg -t ./template.go -w ./main.go
```

```
eg --help
This tool implements example-based refactoring of expressions.

The transformation is specified as a Go file defining two functions,
'before' and 'after', of identical types.  Each function body consists
of a single statement: either a return statement with a single
(possibly multi-valued) expression, or an expression statement.  The
'before' expression specifies a pattern and the 'after' expression its
replacement.

        package P
        import ( "errors"; "fmt" )
        func before(s string) error { return fmt.Errorf("%s", s) }
        func after(s string)  error { return errors.New(s) }

The expression statement form is useful when the expression has no
result, for example:

        func before(msg string) { log.Fatalf("%s", msg) }
        func after(msg string)  { log.Fatal(msg) }

The parameters of both functions are wildcards that may match any
expression assignable to that type.  If the pattern contains multiple
occurrences of the same parameter, each must match the same expression
in the input for the pattern to match.  If the replacement contains
multiple occurrences of the same parameter, the expression will be
duplicated, possibly changing the side-effects.

The tool analyses all Go code in the packages specified by the
arguments, replacing all occurrences of the pattern with the
substitution.

So, the transform above would change this input:
        err := fmt.Errorf("%s", "error: " + msg)
to this output:
        err := errors.New("error: " + msg)

Identifiers, including qualified identifiers (p.X) are considered to
match only if they denote the same object.  This allows correct
matching even in the presence of dot imports, named imports and
locally shadowed package names in the input program.

Matching of type syntax is semantic, not syntactic: type syntax in the
pattern matches type syntax in the input if the types are identical.
Thus, func(x int) matches func(y int).

This tool was inspired by other example-based refactoring tools,
'gofmt -r' for Go and Refaster for Java.


LIMITATIONS
===========

EXPRESSIVENESS

Only refactorings that replace one expression with another, regardless
of the expression's context, may be expressed.  Refactoring arbitrary
statements (or sequences of statements) is a less well-defined problem
and is less amenable to this approach.

A pattern that contains a function literal (and hence statements)
never matches.

There is no way to generalize over related types, e.g. to express that
a wildcard may have any integer type, for example.

It is not possible to replace an expression by one of a different
type, even in contexts where this is legal, such as x in fmt.Print(x).

The struct literals T{x} and T{K: x} cannot both be matched by a single
template.


SAFETY

Verifying that a transformation does not introduce type errors is very
complex in the general case.  An innocuous-looking replacement of one
constant by another (e.g. 1 to 2) may cause type errors relating to
array types and indices, for example.  The tool performs only very
superficial checks of type preservation.


IMPORTS

Although the matching algorithm is fully aware of scoping rules, the
replacement algorithm is not, so the replacement code may contain
incorrect identifier syntax for imported objects if there are dot
imports, named imports or locally shadowed package names in the input
program.

Imports are added as needed, but they are not removed as needed.
Run 'goimports' on the modified file for now.

Dot imports are forbidden in the template.


TIPS
====

Sometimes a little creativity is required to implement the desired
migration.  This section lists a few tips and tricks.

To remove the final parameter from a function, temporarily change the
function signature so that the final parameter is variadic, as this
allows legal calls both with and without the argument.  Then use eg to
remove the final argument from all callers, and remove the variadic
parameter by hand.  The reverse process can be used to add a final
parameter.

To add or remove parameters other than the final one, you must do it in
stages: (1) declare a variant function f' with a different name and the
desired parameters; (2) use eg to transform calls to f into calls to f',
changing the arguments as needed; (3) change the declaration of f to
match f'; (4) use eg to rename f' to f in all calls; (5) delete f'.
```


