package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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
			Members: []string{"Left Expr", "Operator yaglox.Token", "Right Expr"},
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
			Members: []string{"Operator yaglox.Token", "Right Expr"},
		},
	})
}

func defineAST(outputDir, interfaceName string, defs []*Definition) {
	exePwd, err := os.Executable()
	tmplDir := filepath.Dir(exePwd)
	fmt.Println("template path: ", tmplDir)
	tmplFilePath := path.Join(tmplDir, "expression.go.tmpl")
	filePath := path.Join(outputDir, "expression.go")

	tmplStr, err := os.ReadFile(tmplFilePath)
	if err != nil {
		panic(err.Error())
	}

	tmpl := template.Must(template.New("expressions").Parse(string(tmplStr)))

	buf := &bytes.Buffer{}
	data := &TemplateData{
		BaseName:    interfaceName,
		Definitions: defs,
	}
	err = tmpl.Execute(buf, data)
	if err != nil {
		panic(fmt.Sprintln("template hydration failed: ", err.Error()))
	}

	err = os.WriteFile(filePath, buf.Bytes(), fs.ModePerm) // 0777
	if err != nil {
		panic(fmt.Sprintln("hydrated file creation failed: ", err.Error()))
	}

	//f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, fs.ModePerm) // 0777
}
