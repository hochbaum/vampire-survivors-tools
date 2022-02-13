# vampire-survivors-tools
A small set of tools for the Vampire Survivors Game written in Go.

This repository contains an unmarshaler for Vampire Survivors' LevelDB save files and a CLI tool which allows you to toggle debug mode. More tools will be added when I feel like it.

## Using the CLI tool
```
$ go build
$ ./vampire-survivors-tools.exe --debug # Enables debug mode.
$ ./vampire-survivors-tools.exe         # Disables debug mode.
```

## Using the unmarshaler library
Run `go get github.com/hochbaum/vampire-survivors-tools`

```go
saveFile, err := vampires.ParseSave("path/to/your/leveldb/")
if err != nil {
  panic(err)
}
...
```
