# vampire-survivors-tools
A small set of tools for the Vampire Survivors game written in Go.

As of right now, this repository only contains a library for writing and saving save files and a command line tool which
scans your game for the debug mode constant and enables or disables it.

## Using the CLI tool
```
$ go build
$ ./vampire-survivors-tools.exe --debug # Enables debug mode.
$ ./vampire-survivors-tools.exe         # Disables debug mode.
```

## Using the unmarshaler library
Run `go get github.com/hochbaum/vampire-survivors-tools`

Take a look into the [the examples directory](https://github.com/hochbaum/vampire-survivors-tools/tree/master/_examples)!