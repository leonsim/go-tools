package main

import (
	"flag"
	"fmt"
	"regexp"

	"os"

	"bytes"

	"github.com/corsc/go-tools/package-coverage/generator"
	"github.com/corsc/go-tools/package-coverage/parser"
	"github.com/corsc/go-tools/package-coverage/utils"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s\n", r)
		}
	}()

	verbose := false
	coverage := false
	singleDir := false
	clean := false
	print := false
	slack := false
	ignoreDirs := ""
	webhook := ""
	prefix := ""
	depth := 0
	var matcher *regexp.Regexp

	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&coverage, "c", false, "generate coverage")
	flag.BoolVar(&singleDir, "s", false, "only generate for the supplied directory (no recursion)")
	flag.BoolVar(&clean, "d", false, "clean")
	flag.BoolVar(&print, "p", false, "print coverage to stdout")
	flag.BoolVar(&slack, "slack", false, "output coverage to slack")
	flag.StringVar(&ignoreDirs, "i", `./\.git.*|./_.*`, "ignore regex specified directory")
	flag.StringVar(&webhook, "webhook", "", "Slack webhook URL")
	flag.StringVar(&prefix, "prefix", "", "Prefix to be removed from the output (currently only supported by Slack output)")
	flag.IntVar(&depth, "depth", 0, "How many levels of coverage to output (default is 0 = all) (currently only supported by Slack output)")
	flag.Parse()

	if !verbose {
		utils.VerboseOff()
	}

	startDir := utils.GetCurrentDir()
	path := getPath()

	if ignoreDirs != "" {
		matcher = regexp.MustCompile(ignoreDirs)
	}

	if coverage {
		if singleDir {
			generator.CoverageSingle(path, matcher)
		} else {
			generator.Coverage(path, matcher)
		}
	}

	// switch back to start dir
	err := os.Chdir(startDir)
	if err != nil {
		panic(err)
	}

	if print {
		buffer := bytes.Buffer{}

		if singleDir {
			parser.PrintCoverageSingle(&buffer, path, matcher)
		} else {
			parser.PrintCoverage(&buffer, path, matcher)
		}

		fmt.Print(buffer.String())
	}

	if slack {
		if singleDir {
			parser.SlackCoverageSingle(path, matcher, webhook, prefix, depth)
		} else {
			parser.SlackCoverage(path, matcher, webhook, prefix, depth)
		}
	}

	if clean {
		generator.Clean(path, matcher)
	}
}

func getPath() string {
	path := flag.Arg(0)
	if path == "" {
		panic("Please include a directory as the last argument")
	}
	return path
}
