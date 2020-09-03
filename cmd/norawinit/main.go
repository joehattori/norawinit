package main

import (
	"norawinit"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(norawinit.Analyzer) }

