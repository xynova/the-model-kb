package main

import (
	"os"

	"github.com/xynova/roster-mill/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
