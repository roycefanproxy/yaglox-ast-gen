package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"
)

type Definition struct {
	Name    string
	Members []string
}

type TemplateData struct {
	BaseName    string
	Definitions []*Definition
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: yaglox-ast-gen <output directory>")
		os.Exit(64)
	}

	outputDir := os.Args[1]

	defineAST(outputDir, "Expr", []*Definition{
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
	})
}

func defineAST(outputDir, interfaceName string, defs []*Definition) {
	filePath := path.Join(outputDir, "expression.go")

	tmpl := template.Must(template.New("expressions").Parse(TemplateSource))

	buf := &bytes.Buffer{}
	data := &TemplateData{
		BaseName:    interfaceName,
		Definitions: defs,
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
