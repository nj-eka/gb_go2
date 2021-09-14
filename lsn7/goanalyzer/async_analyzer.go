package goanalyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func GetGoFuncCallCount(fileName, funcName string) (count int, err error) {
	set := token.NewFileSet()
	astFile, err := parser.ParseFile(set, fileName, nil, 0)
	if err != nil { return }
	funcDecls := []*ast.FuncDecl{}
	for _, decl := range astFile.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcDecls = append(funcDecls, funcDecl)
		}
	}
	for _, fn := range funcDecls {
		ast.Inspect(fn, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.GoStmt:
				if n.Call.Fun.(*ast.Ident).Name == funcName{
					count++
				}
			}
			return true
		})
	}
	return
}
