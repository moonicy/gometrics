package main

import (
	_ "embed"
	"encoding/json"
	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

//go:embed config.json
var fileString string

// ConfigData описывает структуру файла конфигурации.
type ConfigData struct {
	StaticCheck []string
}

// NoExitInMainAnalyzer checks for os.Exit call in main function of main package
var NoExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noExitInMain",
	Doc:  "checks for os.Exit call in main function of main package",
	Run:  run,
}

// run реализует поиск и анализ main пакетов на os.Exit
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if fn.Name.Name != "main" || fn.Recv != nil {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				callExpr, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					if pkg, ok := fun.X.(*ast.Ident); ok && pkg.Name == "os" && fun.Sel.Name == "Exit" {
						pass.Reportf(n.Pos(), "calling os.Exit in the main function of main package is prohibited")
					}
				}
				return true
			})
		}
	}

	return nil, nil
}

func main() {
	var cfg ConfigData
	if err := json.Unmarshal([]byte(fileString), &cfg); err != nil {
		panic(err)
	}
	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		NoExitInMainAnalyzer,
		errcheck.Analyzer,
		ineffassign.Analyzer,
	}
	checks := make(map[string]bool)
	for _, v := range cfg.StaticCheck {
		checks[v] = true
	}
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}
