package align

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "align",
	Short: "Face-align images and generate video output",
}

func Execute() error {
	return rootCmd.Execute()
}
