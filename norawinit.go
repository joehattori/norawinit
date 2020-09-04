package norawinit

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "norawinit is a tool to limit initialization of structs to designated functions."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "norawinit",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const mark = "initWrapper"

type posRange struct {
	from token.Pos
	to   token.Pos
}

func (r *posRange) contains(n token.Pos) bool {
	return r.from <= n && n <= r.to
}

var idMatcher = regexp.MustCompile(`^[a-zA-Z_]+\w*`)
var initWrappers = map[string]string{}
var funcScopes = map[string]*posRange{}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	file := pass.Files[0]
	cmap := ast.NewCommentMap(pass.Fset, file, file.Comments)

	inspect.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)
		comments := cmap.Filter(genDecl).Comments()
		if len(comments) == 0 {
			return
		}
		for i, spec := range genDecl.Specs {
			switch spec := spec.(type) {
			case *ast.TypeSpec:
				comment := comments[i].Text()
				if idx := strings.Index(comment, mark); idx >= 0 {
					restStr := strings.TrimSpace(comment[idx+len(mark):])
					if restStr[0] != ':' {
						continue
					}
					restStr = strings.TrimSpace(restStr[1:])
					initWrappers[spec.Name.Name] = idMatcher.FindString(restStr)
				}
			}
		}
	})

	inspect.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		fn := n.(*ast.FuncDecl)
		funcScopes[fn.Name.Name] = &posRange{fn.Pos(), fn.End()}
	})

	inspect.Preorder([]ast.Node{(*ast.CompositeLit)(nil)}, func(n ast.Node) {
		switch t := n.(*ast.CompositeLit).Type.(type) {
		case *ast.Ident:
			if funcName, ok := initWrappers[t.Name]; ok {
				if rng, ok := funcScopes[funcName]; ok && !rng.contains(n.Pos()) {
					pass.Reportf(n.Pos(), fmt.Sprintf("%s should be initialized in %s.", t.Name, funcName))
				}
			}
		}
	})

	return nil, nil
}
