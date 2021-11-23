Download file at URL resuming when interrupted. Adapted from Go in Practice, technique 50.

```
go build
./download https://go.dev/dl/go1.17.3.darwin-amd64.pkg
# Hit Ctrl-C and try again. Check the file size.
./download https://go.dev/dl/go1.17.3.darwin-amd64.pkg
```
