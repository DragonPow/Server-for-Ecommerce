package cmd

import (
	"log"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) (err error) {
	if app.Run(os.Args) != nil {
		panic(err)
	}
	return err
}
