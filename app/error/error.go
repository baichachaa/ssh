package ek

import (
	"fmt"
	"os"
)

func CheckError(err error, info string) {
	if err != nil {
		fmt.Printf("%s. error: %s\n", info, err)
		os.Exit(1)
	}
}
