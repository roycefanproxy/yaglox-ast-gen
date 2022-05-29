package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Definition struct {
	Name    string
	Members []string
}

type TemplateData struct {
	BaseName     string
	Definitions  []*Definition
	VisitorTypes []*VisitorType
}

type VisitorType struct {
	Type string
	Name string
}

var funcMap = template.FuncMap{
	"add":             add,
	"typeParam":       typeParam,
	"needVoidVisitor": needVoidVisitor,
}

func add(a, b int) int {
	return a + b
}

func typeParam(typeStr string) string {
	if typeStr == "" {
		return "Void"
	}

	return fmt.Sprintf("[%s]", typeStr)
}

func needVoidVisitor(visitorTypes []*VisitorType) bool {
	for _, vt := range visitorTypes {
		if vt.Name == "" {
			return true
		}
	}

	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: yaglox-ast-gen <output directory with filename>")
		os.Exit(64)
	}

	outputDir := os.Args[1]
	visitorTypes := []*VisitorType{
		{
			Name: "String",
			Type: "string",
		},
		{
			Name: "Interface",
			Type: "interface{}",
		},
		{
			Name: "",
			Type: "",
		},
	}

	defineAST(outputDir, "Expr", []*Definition{
		{
			Name:    "Assign",
			Members: []string{"Name Token", "Value Expr"},
		},
		{
			Name:    "Binary",
			Members: []string{"Left Expr", "Operator Token", "Right Expr"},
		},
		{
			Name:    "Grouping",
			Members: []string{"Expression Expr"},
		},
		{
			Name:    "Literal",
			Members: []string{"Value interface{}"},
		},
		{
			Name:    "Unary",
			Members: []string{"Operator Token", "Right Expr"},
		},
		{
			Name:    "Variable",
			Members: []string{"Name Token"},
		},
	}, visitorTypes)

	defineAST(outputDir, "Stmt", []*Definition{
		{
			Name:    "ExprStmt",
			Members: []string{"Expression Expr"},
		},
		{
			Name:    "PrintStmt",
			Members: []string{"Expression Expr"},
		},
		{
			Name:    "VarDeclStmt",
			Members: []string{"Name Token", "Initializer Expr"},
		},
		{
			Name:    "BlockStmt",
			Members: []string{"Statements []Stmt"},
		},
	}, visitorTypes)

}

func defineAST(outputDir, interfaceName string, defs []*Definition, visitorTypes []*VisitorType) {
	filename := fmt.Sprintf("%s.go", strings.ToLower(interfaceName))
	filePath := filepath.Join(outputDir, filename)
	tmpl := template.Must(template.New("expressions").Funcs(funcMap).Parse(TemplateSource))

	buf := &bytes.Buffer{}
	data := &TemplateData{
		BaseName:     interfaceName,
		Definitions:  defs,
		VisitorTypes: visitorTypes,
	}
	err := tmpl.Execute(buf, data)
	if err != nil {
		panic(fmt.Sprintln("template hydration failed: ", err.Error()))
	}

	err = os.WriteFile(filePath, buf.Bytes(), fs.ModePerm) // 0777
	if err != nil {
		panic(fmt.Sprintln("hydrated file creation failed: ", err.Error()))
	}

	//f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, fs.ModePerm) // 0777
}
