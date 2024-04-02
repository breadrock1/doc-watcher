package cmd

import (
	"doc-notifier/internal/pkg/options"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var serviceOptions *options.Options

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "internal",
	Short: "Launch internal service to load endpoints from watcher directory",
	Long: `
		Launch internal service to load endpoints from watcher directory.
	`,

	Run: func(cmd *cobra.Command, _ []string) {
		fromEnv, _ := cmd.Flags().GetBool("from-env")
		disabledDotenv, _ := cmd.Flags().GetBool("without-dotenv")

		var parseErr error
		if fromEnv {
			serviceOptions, parseErr = options.LoadFromEnv(disabledDotenv)
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
	flags.StringArrayP("watched-dirs", "w", []string{"./indexer"}, "A local directory path to watch fs-events")
	flags.StringP("service-address", "n", "0.0.0.0:2893", "Address of current watcher service")

	flags.StringP("ocr-address", "o", "http://localhost:1231", "Address of current watcher service")
	flags.StringP("ocr-mode", "a", "read-raw-file", "Address of current watcher service")

	flags.StringP("docsearch-address", "d", "http://localhost:2892", "An doc-seacher address with port")

	flags.StringP("tokenizer-address", "t", "http://localhost:8001", "fs-notifier service host address")
	flags.StringP("tokenizer-mode", "b", "assistant", "An llm address with port")
	flags.IntP("chunk-size", "l", 800, "An llm address with port")
	flags.IntP("size-overlap", "p", 100, "An llm address with port")
	flags.UintP("tokenizer-timeout", "x", 300, "Tokenizer timeout seconds")
	flags.BoolP("return-chunks", "r", true, "Load config from env.")
	flags.BoolP("chunk-by-self", "c", false, "Store document as doc-chunks.")

	flags.BoolP("from-env", "e", false, "Parse options from env.")
	flags.BoolP("without-dotenv", "z", false, "Parse options from native env.")
}

func LoadFromCli(cmd *cobra.Command) (*options.Options, error) {
	var parseOptionErr error

	var watchedDirectories []string
	var tokenizerTimeout uint
	var chunkSize, chunkOverlap int
	var returnChunksFlag, chunkBySelfFlag bool
	var tokenizerServiceAddr, tokenizerServiceMode string
	var notifierAddr, docSearchAddr, ocrServiceAddr, ocrServiceMode string

	flags := cmd.Flags()

	if notifierAddr, parseOptionErr = flags.GetString("service-address"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if watchedDirectories, parseOptionErr = flags.GetStringArray("watched-dirs"); parseOptionErr != nil {
		return nil, parseOptionErr
	}

	if ocrServiceAddr, parseOptionErr = flags.GetString("ocr-address"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if ocrServiceMode, parseOptionErr = flags.GetString("ocr-mode"); parseOptionErr != nil {
		return nil, parseOptionErr
	}

	if docSearchAddr, parseOptionErr = flags.GetString("docsearch-address"); parseOptionErr != nil {
		return nil, parseOptionErr
	}

	if tokenizerServiceAddr, parseOptionErr = flags.GetString("tokenizer-address"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if tokenizerServiceMode, parseOptionErr = flags.GetString("tokenizer-mode"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if chunkSize, parseOptionErr = flags.GetInt("chunk-size"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if chunkOverlap, parseOptionErr = flags.GetInt("size-overlap"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if tokenizerTimeout, parseOptionErr = flags.GetUint("tokenizer-timeout"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if returnChunksFlag, parseOptionErr = flags.GetBool("return-chunks"); parseOptionErr != nil {
		return nil, parseOptionErr
	}
	if chunkBySelfFlag, parseOptionErr = flags.GetBool("chunk-by-self"); parseOptionErr != nil {
		return nil, parseOptionErr
	}

	return &options.Options{
		WatcherServiceAddress: notifierAddr,
		WatchedDirectories:    watchedDirectories,

		OcrServiceAddress: ocrServiceAddr,
		OcrServiceMode:    ocrServiceMode,

		DocSearchAddress: docSearchAddr,

		TokenizerServiceAddress: tokenizerServiceAddr,
		TokenizerServiceMode:    tokenizerServiceMode,
		TokenizerChunkSize:      chunkSize,
		TokenizerChunkOverlap:   chunkOverlap,
		TokenizerReturnChunks:   returnChunksFlag,
		TokenizerChunkBySelf:    chunkBySelfFlag,
		TokenizerTimeout:        tokenizerTimeout,
	}, nil
}
