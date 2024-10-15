/*
Package main реализует анализатор кода Go, использующий несколько инструментов статического анализа для проверки кода.

Используемые анализаторы:
 1. **printf.Analyzer** - Этот анализатор проверяет корректность использования функций, подобных `fmt.Printf`, и следит за тем, чтобы спецификаторы формата соответствовали переданным аргументам.
    Источник: `golang.org/x/tools/go/analysis/passes/printf`

 2. **shadow.Analyzer** - Этот анализатор обнаруживает затененные переменные, которые возникают, когда новое объявление переменной случайно использует уже существующее имя переменной в той же области видимости.
    Источник: `golang.org/x/tools/go/analysis/passes/shadow`

 3. **structtag.Analyzer** - Этот анализатор проверяет корректность тегов структур, следя за тем, чтобы они были правильно оформлены и соответствовали стандартам Go.
    Источник: `golang.org/x/tools/go/analysis/passes/structtag`

 4. **NoExitInMainAnalyzer** - Пользовательский анализатор, который проверяет наличие вызовов `os.Exit` в функции `main` пакета `main`, что не рекомендуется в определённых сценариях (например, для тестирования или плавного завершения программы).
    Определён в данном пакете.

 5. **errcheck.Analyzer** - Этот анализатор проверяет, что возвращаемые ошибки из функций правильно обрабатываются и не игнорируются.
    Источник: `github.com/kisielk/errcheck/errcheck`

 6. **ineffassign.Analyzer** - Этот анализатор обнаруживает и сообщает о присваиваниях переменным, которые никогда не используются до того, как будут перезаписаны, что неэффективно и может указывать на ошибку.
    Источник: `github.com/gordonklaus/ineffassign/pkg/ineffassign`

7. **Анализаторы Staticcheck** - Набор анализаторов из Staticcheck, используемый для выявления различных ошибок и улучшений в коде Go. Этот анализатор настраивается через файл конфигурации `config.json`. В данном примере анализируются следующие статические проверки:
  - `SA`: Обнаружение потенциальных ошибок безопасности.
  - `S1005`: Замена вызовов функций с неиспользуемыми результатами.
  - `ST1008`: Проверка корректности структур.
  - `QF1012`: Оптимизация и устранение неэффективных операций.
    Источник: `honnef.co/go/tools/staticcheck`

### Как запускать

	Запустите этот код как Go-программу. Она использует пакет multichecker, который позволяет запустить несколько анализаторов одновременно.
	Пример команды для запуска: go run main.go
	Это выполнит указанные в коде и конфигурации анализаторы.
*/
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

// NoExitInMainAnalyzer checks for os.Exit call in main function of main package.
var NoExitInMainAnalyzer = &analysis.Analyzer{
	Name: "noExitInMain",
	Doc:  "checks for os.Exit call in main function of main package",
	Run:  run,
}

// run реализует поиск и анализ main пакетов на os.Exit.
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
