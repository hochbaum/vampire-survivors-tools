package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

const (
	pattern     = `(.*const _0x[a-f0-9]+=)(!0x1|true|false)(,(?:_0x[a-f0-9]+=.*?(?:,|;)){10,}.*)`
	defaultPath = `C:\Program Files (x86)\Steam\steamapps\common\Vampire Survivors\resources\app\.webpack\renderer\main.bundle.js`
)

func main() {
	path := flag.String("path", defaultPath, "Specifies the path to the game code.")
	debug := flag.Bool("debug", false, "Changes the debug mode value.")
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Could not find your game at %s. Please check the location and specify it using `--path`.", *path)
			os.Exit(1)
		}
		panic(err)
	}

	defer file.Close()

	code, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	result := regex.FindAllSubmatch(code, -1)
	if len(result) == 0 {
		panic("could not find debug constant")
	}

	replaced := regex.ReplaceAll(code, []byte(fmt.Sprintf("${1}%t${3}", *debug)))
	if err := os.WriteFile(*path, replaced, 0644); err != nil {
		panic(err)
	}
}
