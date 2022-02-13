package main

import (
	"fmt"
	"vampire-survivors-tools/vampires"
)

func main() {
	const path = "<path to levelDB>"

	save, err := vampires.ParseSave(path)
	if err != nil {
		panic(err)
	}

	if save.CheatCodeUsed {
		fmt.Println("You cheater!")
	} else {
		fmt.Println("Honest player.")
	}
}
