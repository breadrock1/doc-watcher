package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"doc-notifier/internal/pkg/watcher"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "internal",
	Short: "Launch internal service to load files from watcher directory",
	Long: `
		Launch internal service to load files from watcher directory.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		assistantAddr, _ := cmd.Flags().GetString("assistant-address")
		searcherAddr, _ := cmd.Flags().GetString("searcher-address")
		watcherPath, _ := cmd.Flags().GetStringArray("watcher-path")
		fmt.Println(assistantAddr, searcherAddr, watcherPath)

		service := watcher.New(&watcher.Options{
			SearcherAddress:  searcherAddr,
			AssistantAddress: assistantAddr,
			WatchDirectories: watcherPath,
		})

		service.RunWatcher()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("searcher-address", "s", "localhost:2892", "An docseacher address with port")
	rootCmd.Flags().StringP("assistant-address", "a", "localhost:8000", "An assistant address with port")
	rootCmd.Flags().StringArrayP("watcher-path", "p", []string{"/archiver"}, "A local directory path to watch fs-events")
}
