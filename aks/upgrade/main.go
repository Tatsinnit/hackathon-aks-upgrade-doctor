package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "aks-doctor",
	}

	cmd.AddCommand(createDemoCommand())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
