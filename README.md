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

doc.json
```json
{"foo": {"bar": "baz"}, "numbers": [1, 2, 3, 4]}
```

Command:
```
cat doc.json | gojt path  .foo.bar
```

Output:
```
baz
```

### keys

Print the available keys on a given path

Clipboard:
```json
{"foo": {"bar": "baz"}, "numbers": [1, 2, 3, 4]}
```

Command:
```
pbpaste | gojt keys .
```

Output:
```
foo
numbers
```

Command:
```
pbpaste | gojt keys .foo
```

Output:
```
bar
```

# References Used for Development

- [Read piped data](https://flaviocopes.com/go-shell-pipes/)
- [CLI: Cobra](https://github.com/spf13/cobra)
