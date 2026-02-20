package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"loglinter/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}