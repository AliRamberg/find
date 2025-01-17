package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/AliRamberg/find/pkg/search"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	filename   string
	searchTerm string
	err        error
)

var rootCmd = &cobra.Command{
	Use: "find",
	Run: func(cmd *cobra.Command, args []string) {
		var input string
		s, err := search.NewSearcher(filename)
		if err != nil {
			log.Fatalf("failed to create searcher: %v", err)
		}

		if searchTerm == "-" {
			buff := bufio.NewReaderSize(os.Stdin, search.BufferLimit)
			input, err = buff.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatalf("failed to read from stdin: %v", err)
			}
		} else {
			input = searchTerm
		}

		line, err := s.FindLine(input)
		if err != nil {
			log.Fatalf("failed to find line: %v", err)
		}
		fmt.Printf("Result: %s\n", formatLine(line))
	},
}

func formatLine(line string) string {
	if len(line) <= 20 {
		return line
	}
	return line[:10] + "..." + line[len(line)-10:]
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&filename, "file", "f", "", "File to search")
	rootCmd.Flags().StringVarP(&searchTerm, "term", "t", "", "Search term")

	rootCmd.MarkFlagRequired("file")
	rootCmd.MarkFlagRequired("term")
}
