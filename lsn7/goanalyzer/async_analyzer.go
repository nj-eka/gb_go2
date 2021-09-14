package goanalyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func GetGoFuncCallCount(fileName, funcName string) (count int, err error) {
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, fileName, nil, 0)
	if err != nil { return }
	for _, decl := range astFile.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			ast.Inspect(funcDecl, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.GoStmt:
					if n.Call.Fun.(*ast.Ident).Name == funcName{
						count++
					}
				}
				return true
			})
		}
	}
	return
}
