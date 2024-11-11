package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: generateast <output_directory>")
		os.Exit(64)
	}
	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Assign		: name *Token, value Expr",
		"Binary		: left Expr, operator *Token, right Expr",
		"Call		: callee Expr, paren *Token, arguments []Expr",
		"Get		: object Expr, name *Token",
		"Grouping	: expression Expr",
		"Literal	: value any",
		"Logical	: left Expr, operator *Token, right Expr",
		"Set		: object Expr, name *Token, value Expr",
		"Super		: keyword *Token, method *Token",
		"This		: keyword *Token",
		"Unary		: operator *Token, right Expr",
		"Variable	: name *Token",
	})
	defineAst(outputDir, "Stmt", []string{
		"Block		: statements []Stmt",
		"Class		: name *Token, superclass *Variable," +
			" methods []*Function",
		"Expression	: expression Expr",
		"Function	: name *Token, params []*Token, body []Stmt",
		"If		: condition Expr, thenBranch Stmt, elseBranch Stmt",
		"Print		: expression Expr",
		"Return		: keyword *Token, value Expr",
		"Var		: name *Token, initializer Expr",
		"While		: condition Expr, body Stmt",
	})
}

func defineAst(outputDir, baseName string, types []string) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %s\n", err)
		os.Exit(64)
	}
	writer := bufio.NewWriter(f)
	writer.WriteString("package lox\n\n")
	defineVisitor(writer, baseName, types)
	writer.WriteString(fmt.Sprintf("type %s interface {\n", baseName))
	writer.WriteString(fmt.Sprintf("\tAccept(v %sVisitor) any\n", baseName))
	writer.WriteString("}\n\n")
	for _, etype := range types {
		className := strings.Split(etype, ":")[0]
		className = strings.Trim(className, " \t")
		fields := strings.Split(etype, ":")[1]
		fields = strings.Trim(fields, " \t")
		defineType(writer, baseName, className, fields)
	}
	writer.Flush()
}

func defineVisitor(writer *bufio.Writer, baseName string, types []string) {
	writer.WriteString(fmt.Sprintf("type %sVisitor interface {\n", baseName))
	for _, ftype := range types {
		typeName := strings.Split(ftype, ":")[0]
		typeName = strings.Trim(typeName, " \t")
		writer.WriteString(fmt.Sprintf(
			"\tVisit%s%s(%s *%s) any\n",
			typeName, baseName,
			strings.ToLower(baseName), typeName,
		),
		)
	}
	writer.WriteString("}\n\n")
}

func defineType(writer *bufio.Writer, baseName, className, fieldList string) {
	writer.WriteString(fmt.Sprintf("type %s struct {\n", className))
	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		ftype := strings.Split(field, " ")[1]
		name = strings.Title(name)
		writer.WriteString(fmt.Sprintf("\t%s %s\n", name, ftype))
	}
	writer.WriteString("}\n\n")
	writer.WriteString(fmt.Sprintf("func New%s(", className))
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		ftype := strings.Split(field, " ")[1]
		writer.WriteString(fmt.Sprintf("%s %s, ", name, ftype))
	}
	writer.WriteString(fmt.Sprintf(") %s {\n", baseName))
	writer.WriteString(fmt.Sprintf("\treturn &%s{ ", className))
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		writer.WriteString(fmt.Sprintf("%s, ", name))
	}
	writer.WriteString(" }\n}\n\n")
	writer.WriteString(fmt.Sprintf(
		"func (%c *%s) Accept(%cv %sVisitor) any {\n",
		unicode.ToLower(rune(className[0])),
		className,
		unicode.ToLower(rune(baseName[0])),
		baseName,
	),
	)
	writer.WriteString(fmt.Sprintf(
		"\treturn %cv.Visit%s%s(%c)\n",
		unicode.ToLower(rune(baseName[0])),
		className, baseName,
		unicode.ToLower(rune(className[0])),
	),
	)
	writer.WriteString("}\n\n")
}
