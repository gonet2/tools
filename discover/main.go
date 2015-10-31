package main

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
)

const (
	raw = "https://raw.githubusercontent.com/gonet2/libs/master/services/services.go"
)

func main() {
	resp, err := http.Get(raw)
	if err != nil {
		log.Fatal(err)
	}
	// parser
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", resp.Body, 0)
	if err != nil {
		log.Fatal(err)
	}

	// remove Init()
LOOP:
	for k := range f.Decls {
		switch f.Decls[k].(type) {
		case *ast.FuncDecl:
			decl := f.Decls[k].(*ast.FuncDecl)
			if decl.Name.Name == "Init" {
				f.Decls = append(f.Decls[:k], f.Decls[k+1:]...)
				break LOOP
			}
		}
	}

	// rewrite
	format.Node(os.Stdout, fset, f)
}
