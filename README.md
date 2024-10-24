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

- `word`: any alphabetical string
- `number`: any numerical string
- `symbol`: any printable ascii character that is not a number or letter
- `line`: gets all characters before a newline character `\n`
- `lbrace` and `rbrace`: for a static `{` and `}` respectively

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
	fmt.Println(tokr.Run("foo! 123 + bar"))
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

You can create a new class two different ways. With `.Class()` you provide the class name and a function that returns true as long as the given character is part of your class. With `.ClassFromPattern()` you just provide a class name and a pattern it uses:

```go
tokr.Class("notA", func (b byte) bool {
    return b != 'A'
})

tokr.ClassFromPattern("username", "username: {word}")
```

When you have nested classes as in the `.ClassFromPattern()` example above, you can still access the value of each previous class:

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
