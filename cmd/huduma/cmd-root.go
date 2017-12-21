package main

import (
	"github.com/spf13/cobra"
)

//RootCommand configur our cobra command
func RootCommand() *cobra.Command {
	rootCmd := cobra.Command{
		Use:  "huduma",
		Long: `Huduma is a leighweight rest api that was implemented to learn how Cobra, Viper in rest api can interact with each others.`,
		Run:  run,
	}

	rootCmd.PersistentFlags().StringP("config", "C", "", "An explicit config file to use")
	rootCmd.Flags().IntP("port", "p", 0, "the port to do things on")

	rootCmd.AddCommand(&versionCmd)

	return &rootCmd
}

func run(cmd *cobra.Command, args []string) {
	initConf()
	//initAPIConfig()
}
