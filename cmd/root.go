package cmd

import (
	"doc-notifier/internal/pkg/server"
	"doc-notifier/internal/pkg/watcher"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "internal",
	Short: "Launch internal service to load endpoints from watcher directory",
	Long: `
		Launch internal service to load endpoints from watcher directory.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		serviceAddress, _ := cmd.Flags().GetString("watcherService-address")
		llmServiceAddr, _ := cmd.Flags().GetString("llm-watcherService-address")
		ocrServiceAddr, _ := cmd.Flags().GetString("recognition-watcherService-address")
		docSearchServiceAddr, _ := cmd.Flags().GetString("docsearch-watcherService-address")
		watcherPath, _ := cmd.Flags().GetStringArray("watcher-directory-path")
		fmt.Println(ocrServiceAddr, docSearchServiceAddr, watcherPath)

		watcherService := watcher.New(&watcher.Options{
			LlmServiceAddress: llmServiceAddr,
			DocSearchAddress:  docSearchServiceAddr,
			OcrServiceAddress: ocrServiceAddr,
			WatchDirectories:  watcherPath,
		})

		tmp := strings.Split(serviceAddress, ":")
		servicePort, _ := strconv.Atoi(tmp[1])
		serverOptions := server.BuildOptions(tmp[0], servicePort)
		httpServer := server.New(serverOptions, watcherService)
		httpServer.RunServer()

		watcherService.RunWatcher()
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
	rootCmd.Flags().StringP("service-address", "s", "0.0.0.0:2893", "fs-notifier service host address")
	rootCmd.Flags().StringP("llm-service-address", "l", "localhost:8000", "An llm address with port")
	rootCmd.Flags().StringP("recognition-service-address", "r", "localhost:8000", "An ocr address with port")
	rootCmd.Flags().StringP("docsearch-service-address", "d", "localhost:2892", "An docseacher address with port")
	rootCmd.Flags().StringArrayP("watcher-directory-path", "w", []string{"/archiver"}, "A local directory path to watch fs-events")
}
