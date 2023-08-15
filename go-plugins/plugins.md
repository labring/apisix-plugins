# go plugins

## try_path

### introduction
`try_path` is a plugin that references the `try_files` implementation in nginx. It can check if a path exists before accessing it, and if it doesn't exist, it will attempt the next path in the sequence until the last one is reached.

### how to use

name: plugin name, current use `try-path`

value: value is a string (serialized JSON object) that contains `paths` and `host`. `paths` are the paths to be attempted, and `$uri` can be used to replace the current requested path. `host` is the address used to determine if the path exists.

```json
{
  "name": "try-path",
  "value": "{\"paths\":[\"$uri\", \"$uri/\", \"$uriindex.html\", \"$uri/index.html\"], \"host\":\"http://hostname\"}"
}
```