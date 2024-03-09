package cmd

import (
	"doc-notifier/internal/pkg/options"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var serviceOptions *options.Options

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "internal",
	Short: "Launch internal service to load endpoints from watcher directory",
	Long: `
		Launch internal service to load endpoints from watcher directory.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		fromEnv, _ := cmd.Flags().GetBool("from-env")

		var parseErr error
		if fromEnv {
			serviceOptions, parseErr = options.LoadFromEnv()
		} else {
			serviceOptions, parseErr = LoadFromCli(cmd)
		}

		if parseErr != nil {
			log.Fatal(parseErr)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() *options.Options {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	return serviceOptions
}

func init() {
	flags := rootCmd.Flags()
	flags.BoolP("from-env", "e", false, "Load config from env.")
	flags.BoolP("load-chunks", "c", false, "Store document as doc-chunks.")
	flags.BoolP("read-raw-file", "m", false, "Read raw file data and load.")
	flags.StringArrayP("watcher-dir-path", "w", []string{"/archiver"}, "A local directory path to watch fs-events")
	flags.StringP("host-address", "s", "0.0.0.0:2893", "fs-notifier service host address")
	flags.StringP("llm-address", "l", "localhost:8000", "An llm address with port")
	flags.StringP("ocr-address", "r", "localhost:8000", "An ocr address with port")
	flags.StringP("doc-address", "d", "localhost:2892", "An docseacher address with port")
}

func LoadFromCli(cmd *cobra.Command) (*options.Options, error) {
	var parseErr error
	var watcherPath []string
	var storeChunksFlag, readRawFileFlag bool
	var serverAddr, llmAddr, ocrAddr, docSearchAddr string

	flags := cmd.Flags()
	if storeChunksFlag, parseErr = flags.GetBool("load-chunks"); parseErr != nil {
		return nil, parseErr
	}

	if readRawFileFlag, parseErr = flags.GetBool("read-raw-file"); parseErr != nil {
		return nil, parseErr
	}

	if serverAddr, parseErr = flags.GetString("host-address"); parseErr != nil {
		return nil, parseErr
	}

	if llmAddr, parseErr = flags.GetString("llm-address"); parseErr != nil {
		return nil, parseErr
	}

	if ocrAddr, parseErr = flags.GetString("ocr-address"); parseErr != nil {
		return nil, parseErr
	}

	if docSearchAddr, parseErr = flags.GetString("doc-address"); parseErr != nil {
		return nil, parseErr
	}

	if watcherPath, parseErr = flags.GetStringArray("watcher-dir-path"); parseErr != nil {
		return nil, parseErr
	}

	return &options.Options{
		ServerAddress:     serverAddr,
		LlmServiceAddress: llmAddr,
		DocSearchAddress:  docSearchAddr,
		OcrServiceAddress: ocrAddr,
		WatchDirectories:  watcherPath,
		StoreChunksFlag:   storeChunksFlag,
		ReadRawFileFlag:   readRawFileFlag,
	}, nil
}
