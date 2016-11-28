package parser

import (
	"io/ioutil"
	"log"
	"regexp"
	"sort"
)

// get coverage using the paths and exclusions supplied
func getCoverageData(paths []string, dirMatcher, fileMatcher *regexp.Regexp) ([]string, coverageByPackage) {
	var contents string
	for _, path := range paths {
		if dirMatcher.FindString(path) != "" {
			log.Printf("Printing of coverage for path '%s' skipped due to skipDir regex '%s'", path, dirMatcher.String())
			continue
		}

		contents += getFileContents(path)
	}

	return getCoverageByContents(contents, fileMatcher)
}

// get coverage from supplied string (used after concatenating all the individual coverage files together
func getCoverageByContents(contents string, fileMatcher *regexp.Regexp) ([]string, coverageByPackage) {
	coverageData := getCoverageByPackage(contents, fileMatcher)
	pkgs := getSortedPackages(coverageData)

	return pkgs, coverageData
}

// will calculate and return the coverage for a package or packages from the supplied coverage file contents
func getCoverageByPackage(contents string, fileMatcher *regexp.Regexp) coverageByPackage {
	coverageData := parseLines(contents, fileMatcher)
	return coverageData
}

func getFileContents(filename string) string {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(contents)
}

func getSortedPackages(coverageData coverageByPackage) []string {
	output := []string{}

	for pkg := range coverageData {
		output = append(output, pkg)
	}

	sort.Strings(output)

	return output
}

// calculate the coverage and statement counts from the supplied data
func getStats(cover *coverage) (float64, float64) {
	stmts := float64(cover.selfStatements + cover.childStatements)
	stmtsCovered := float64(cover.selfCovered + cover.childCovered)

	covered := (stmtsCovered / stmts) * 100

	return covered, stmts
}
