package cli

import (
	"fmt"
	"os"

	derrors "github.com/xynova/roster-pick/internal/errors"
)

// Execute runs the Cobra root command.
func Execute() int {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return mapPickExit(err)
	}
	return 0
}

func mapPickExit(err error) int {
	if derrors.IsCode(err, derrors.CodeInvalidArgument) {
		return 2
	}
	return 1
}
