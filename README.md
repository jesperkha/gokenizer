<br>

<div align="center">

<img src=".github/assets/logo.png" width="50%"></img>

*A simple text tokenizer*

</div>

<br>

## How it works

Gokenizer uses pattern strings composed of *static* words and *classes*.

```go
pattern := "var {word} = {number}"
// "var " and " = " are static patterns. Note that whitespace is included.
// {word} and {number} are classes. They are special pre-defined patterns.
```

The following classes are defined by default:

- `lbrace` and `rbrace`: for a static `{` and `}` respectively
- `word`: any alphabetical string
- `char`: a single letter
- `number`: any numerical string
- `float`: any numerical string including a period `.`
- `symbol`: any printable ascii character that is not a number or letter
- `line`: gets all characters before a newline character `\n`
- `base64`: any base64 string, does not check length
- `hex`: hexadecimal string, including `#`

Note that patterns are checked in the order they are defined, therefore it is usually preferred to define specific patterns first, and more general ones last. `ClassX` functions have no immediate effect, but must be run before using the defined class.

## Basic example

```go
func main() {
    tokr := gokenizer.New()

    // Create a pattern that looks for any word ending with a !
    tokr.Pattern("{word}!", func(tok gokenizer.Token) error {
        fmt.Println(tok.Lexeme)
        return nil
    })

    // Run the tokenizer on a given string
    tokr.Run("foo! 123 + bar")
}
```

As you can see, the tokenizer only matched the first word "foo!".

```sh
$ go run .
foo!
```

> Notice that the callback given to `Pattern()` returns an error. This error is returned by `Run()`.

## Getting the string from a class

You can get the parsed token from a class by using `Token.Get()`:

> In cases where you use the same class more than once, use `Token.GetAt()`

```go
func main() {
    tokr := gokenizer.New()

    tokr.Pattern("{word}{number}", func (tok gokenizer.Token) error {
        word := tok.Get("word").Lexeme
        number := tok.Get("number").Lexeme

        fmt.Println(word, number)
        return nil
    })


    tokr.Run("username123")
}
```

```sh
$ go run .
username 123
```

## Creating your own class

You can create a new class a few different ways:

- `.Class()`: takes the class name and a function that returns true as long as the given character is part of your class.
- `.ClassPattern()`: takes a class name and a pattern it uses.
- `.ClassAny()`: takes a class name and a list of patterns, of which only one has to match.
- `.ClassOptional()`: takes a class name and a list of patterns, of which one *or* none have to match.

```go
tokr.Class("notA", func (b byte) bool {
    return b != 'A'
})

tokr.ClassPattern("username", "username: {word}")

tokr.ClassAny("games", "Elden Ring", "The Sims {number}")

tokr.ClassOptional("whitespace", " ", "\t")
```

When you have nested classes as in the `.ClassPattern()` and `.ClassAny()` examples above, you can still access the value of each previous class:

```go
tokr.Pattern("{username}", func (tok gokenizer.Token) error {
    fmt.Println(tok.Get("username").Get("word").Lexeme)
    return nil
})

tokr.Run("username: John")
```

```sh
$ go run .
John
```

## Example: Parsing a .env file

Here is an example for a .env file parser that prints out each key-value pair and comment:

```go
func main() {
    tokr := gokenizer.New()

    tokr.ClassPattern("comment", "#{line}")
    tokr.ClassOptional("space?", " ")

    // Syntactic sugar
    tokr.ClassPattern("key", "{word}")
    tokr.ClassPattern("value", "{line}")

    tokr.ClassPattern("keyValue", "{key}{space?}={space?}{value}")

    // Prints all key value pairs
    tokr.Pattern("{keyValue}", func(t gokenizer.Token) error {
        key := t.Get("keyValue").Get("key")
        value := t.Get("keyValue").Get("value")

        fmt.Printf("Key: %s, Value: %s", key.Lexeme, value.Lexeme)
        return nil
    })

    // Prints all comments
    tokr.Pattern("{comment}", func(t gokenizer.Token) error {
        line := t.Lexeme[:len(t.Lexeme)-1] // Remove newline

        fmt.Println("Comment: " + line)
        return nil
    })

    b, err := os.ReadFile(".env")
    if err != nil {
        return
    }

    tokr.Run(string(b))
}
```

## Limitations

Gokenizer is designed to be as minimal and straight forward as possible, and therefore comes with a few limitations:

- **No look-ahead parsing:** The tokenizer does not look ahead when parsing patterns, therefore patterns with overlapping match requirements will not work. Example: `{word}bar` will never be matched as any word, including one ending in "bar" will be part of the `word` class. This may be fixed in a later update with a new class/pattern type.
- **More complex user classes:** Currently you cannot easily define a "complex" class that interacts with the internal string iterator. As of now the idea is for the user to create their own checks in the callback functions to more general patterns. However, more control may be given to the user when creating classes in a later update.
