package main

import (
	"os"

	"github.com/xynova/roster-pick/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
