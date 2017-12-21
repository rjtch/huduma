package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

//Version represents which version of hudumu is used
var Version string

var versionCmd = cobra.Command{
	Run: version,
	Use: "version",
}

func version(cmd *cobra.Command, args []string) {
	fmt.Println("1.0.0", Version)
}
