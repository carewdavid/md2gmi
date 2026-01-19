package main

import (
	"fmt"
	"os"

	"github.com/carewdavid/md2gmi/src/convert"
)

func main() {
	for _, filename := range os.Args[1:] {
		text, err := os.ReadFile(filename)
		if err != nil {
			_ = fmt.Errorf("Could not open %q\n", filename)
			continue
		}
		formatted := convert.Convert(string(text))
		fmt.Print(formatted)
	}
}
