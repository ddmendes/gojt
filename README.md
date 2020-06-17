# GOJT The Go JSON Tool

Utility tool for parsing and reading json documents from terminal

# Usage

```sh
cat doc.json | gojt [path|keys] <json_path>
```

Gojt always reads piped data being useful to get info from inline JSON
documents in your disk or in your clipboard

```sh
pbpaste | gojt path '.foo.bar'
```

## Commands

### path

Print object in a given path

```sh
➜  gojt git:(master) cat doc.json
{"foo":{"bar":"baz"},"numbers":[1,2,3,4]}
➜  gojt git:(master) cat doc.json | gojt path .foo.bar
"baz"
```

Print arrays using both map or array indexing

```sh
➜  gojt git:(master) ✗ cat doc.json
{"foo":{"bar":"baz"},"numbers":[1,2,3,4]}
➜  gojt git:(master) ✗ cat doc.json | gojt path .numbers.2
3
➜  gojt git:(master) ✗ cat doc.json | gojt path '.numbers[2]'
3
```

### keys

Print the available keys on a given path

```sh
➜  gojt git:(master) ✗ echo '{"foo":{"bar":"baz"},"numbers":[1,2,3,4]}' | pbcopy
➜  gojt git:(master) ✗ pbpaste | ./gojt keys .
[
  "foo",
  "numbers"
]
➜  gojt git:(master) ✗ pbpaste | ./gojt keys .foo
[
  "bar"
]
```

# References Used for Development

- [encoding/json reference](https://golang.org/pkg/encoding/json/)
- [Read piped data](https://flaviocopes.com/go-shell-pipes/)
- [CLI: Cobra](https://github.com/spf13/cobra)
- [Fake data into os.Stdin during tests](https://stackoverflow.com/questions/46365221/fill-os-stdin-for-function-that-reads-from-it)
