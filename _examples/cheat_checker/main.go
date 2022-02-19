package main

import (
	"fmt"
	"github.com/hochbaum/vampire-survivors-tools/vampires"
)

func main() {
	const path = "path/to/your/levelDB"
	save, db, err := vampires.OpenSaveFile(path)
	if err != nil {
		panic(err)
	}

	// The DB must be closed by you when you are done with it.
	defer db.Close()

	if save.CheatCodeUsed {
		fmt.Println("Hackerman!")
		fmt.Println("Covering tracks.")
		save.CheatCodeUsed = false
	} else {
		fmt.Println("Have some coins, real gamer!")
		save.Coins += 1337
	}

	if err := vampires.StoreSaveFile(save, db); err != nil {
		panic(err)
	}
}
