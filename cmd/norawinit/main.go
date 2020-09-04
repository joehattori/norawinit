package main

import (
	"github.com/joehattori/norawinit"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(norawinit.Analyzer) }
