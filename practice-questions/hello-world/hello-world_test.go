package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

const filePath = "main.go"

// TestPackageName tests that the main package is defined
func TestPackageName(t *testing.T) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)

	if err != nil {
		t.Fatalf("error parsing file. Please fix any errors in the file: %v", err)
	}

	if node.Name.Name != "main" {
		t.Errorf("package = %v, want 'main'", node.Name.Name)
	}
}

// TestImports tests if the 'fmt' package is imported
func TestImports(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		t.Fatalf("error parsing file. Please fix any errors in the file: %v", err)
	}
	if len(node.Imports) == 0 {
		t.Errorf("no imports found. Please import the 'fmt' package")
	}
	for _, importSpec := range node.Imports {
		if importSpec.Path.Value == `"fmt"` {
			return
		}
	}
	t.Errorf("no imports found. Please import the 'fmt' package")
}

// TestMain tests that the main function is defined and that a call to fmt.Println is called with "Hello World!"
func TestMain(t *testing.T) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, 0)

	if err != nil {
		t.Fatal("error parsing file. Please fix any errors in the file")
	}

	found := false

	for _, declaration := range node.Decls {
		if fn, isFn := declaration.(*ast.FuncDecl); isFn && fn.Name.Name == "main" {
			found = true
			if !printsHelloWorld(fn) {
				t.Errorf("main function does not print 'Hello World!'. Please fix this")
			}
			break
		}
	}

	if !found {
		t.Errorf("expected main function to be defined")
	}
}

// PrintsHelloWorld a helper function to check if the main function prints "Hello World!"
func printsHelloWorld(fn *ast.FuncDecl) bool {
	// iterate over each statement in the function body
	for _, stmt := range fn.Body.List {
		// check if the statement is an expression statement
		if exprStmt, isExpr := stmt.(*ast.ExprStmt); isExpr {
			// check if the expression is a call expression
			if callExpr, isCall := exprStmt.X.(*ast.CallExpr); isCall {
				// check if the function being called is a selector expression (e.g fmt.Println)
				if selExpr, isSel := callExpr.Fun.(*ast.SelectorExpr); isSel {
					// check if the selector expression is calling a function from the "fmt" package and the method being called is the Println method
					if ident, isIdent := selExpr.X.(*ast.Ident); isIdent && ident.Name == "fmt" && selExpr.Sel.Name == "Println" {
						// check if the call expression has exactly one arg
						if len(callExpr.Args) == 1 {
							// check if the argument is a basic literal and the value is "Hello World!"
							basicLit, _ := callExpr.Args[0].(*ast.BasicLit)
							if basicLit.Value == `"Hello World!"` {
								return true
							}
						}
					}
				}
			}
		}
	}
	return false
}
