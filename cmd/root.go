package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	use     = "category-service"
	version = "1.0.0-beta"
)

var rootCmd = &cobra.Command{
	Use:   use,
	Short: fmt.Sprintf("USAGE %s [OPTIONS]", os.Args[0]),
	Long:  fmt.Sprintf(`USAGE %s [OPTIONS] : see --help for details`, os.Args[0]),
	RunE:  executeRootCommand,
}

func init() {
	rootCmd.Flags().BoolVarP(&rootFlags.Version, "version", "V", false, "show version information.")
	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigFile, "config", "c", "config.yaml", "application config path")
}

// Execute initializes whole application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func executeRootCommand(cmd *cobra.Command, _ []string) error {
	if rootFlags.Version {
		_, err := fmt.Printf("%s v. %s\n", use, version)
		return err
	} else {
		return cmd.Help()
	}
}
