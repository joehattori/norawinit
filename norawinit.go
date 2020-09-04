package norawinit

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
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
	FactTypes: []analysis.Fact{new(initFact)},
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

type initFact map[string]string

func (*initFact) AFact() {}
func (f *initFact) String() string {
	var arr []string
	for ty, fn := range *f {
		arr = append(arr, fmt.Sprintf("%s: %s", ty, fn))
	}
	return "wrappers(" + strings.Join(arr, ", ") + ")"
}

func run(pass *analysis.Pass) (interface{}, error) {
	initWrappers := make(initFact)
	funcScopes := map[string]*posRange{}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	file := pass.Files[0]
	cmap := ast.NewCommentMap(pass.Fset, file, file.Comments)

	doImport := func(spec *ast.ImportSpec) {
		pkg := imported(pass.TypesInfo, spec)
		var fact initFact
		if pass.ImportPackageFact(pkg, &fact) {
			for ty, wrapper := range fact {
				initWrappers[ty] = wrapper
			}
		}
	}

	inspect.Preorder([]ast.Node{(*ast.GenDecl)(nil)}, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)
		switch genDecl.Tok {
		case token.IMPORT:
			for _, spec := range genDecl.Specs {
				doImport(spec.(*ast.ImportSpec))
			}
		case token.TYPE:
			for _, spec := range genDecl.Specs {
				spec := spec.(*ast.TypeSpec)
				var comments []*ast.CommentGroup
				if genDecl.Lparen == 0 && genDecl.Rparen == 0 {
					comments = cmap.Filter(genDecl).Comments()
				} else {
					comments = cmap.Filter(spec).Comments()
				}
				if len(comments) == 0 {
					continue
				}
				comment := comments[0].Text()
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
		case *ast.SelectorExpr:
			if funcName, ok := initWrappers[t.Sel.Name]; ok {
				if rng, ok := funcScopes[funcName]; ok && !rng.contains(n.Pos()) || !ok {
					pass.Reportf(n.Pos(), fmt.Sprintf("%s should be initialized in %s.", t.Sel.Name, funcName))
				}
			}
		}
	})

	fact := make(initFact)
	for ty, wrapper := range initWrappers {
		fact[ty] = wrapper
	}

	if len(fact) > 0 {
		pass.ExportPackageFact(&fact)
	}

	return nil, nil
}

func imported(info *types.Info, spec *ast.ImportSpec) *types.Package {
	obj, ok := info.Implicits[spec]
	if !ok {
		obj = info.Defs[spec.Name]
	}
	return obj.(*types.PkgName).Imported()
}
