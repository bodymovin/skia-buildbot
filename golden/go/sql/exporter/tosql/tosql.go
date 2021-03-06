// This executable generates a go file that contains the SQL schema for Gold defined as a string.
// By doing this, we have the source of truth as a documented go struct, which can be used in a
// more flexible way than having the SQL as the source of truth.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"go.skia.org/infra/go/sklog"
	"go.skia.org/infra/golden/go/sql/schema"
)

func main() {
	outputFile := flag.String("output_file", "", "The name of the file to write to.")
	outputPkg := flag.String("output_pkg", "", "The name of the package to output to.")
	flag.Parse()
	cwd, err := os.Getwd()
	if err != nil {
		sklog.Fatalf("Could not get working dir")
	}

	generatedText := generateSQL(schema.Tables{}, *outputPkg)
	out := filepath.Join(cwd, *outputFile)
	err = ioutil.WriteFile(out, []byte(generatedText), 0666)
	if err != nil {
		sklog.Fatalf("Could not write SQL to %s: %s", out, err)
	}
	sklog.Infof("Tables schema written to %s.\n", out)
}

// generateSQL takes in a "table type", that is a table whose fields are slices. Each field
// will be interpreted as a table. The sql struct tags will be used to generate the SQL schema.
// A package name is taken in to be included in the returned string. If a malformed type is passed
// in, this function will panic.
func generateSQL(inputType interface{}, pkg string) string {
	header := fmt.Sprintf("package %s\n\n// Generated by //golden/go/sql/exporter/tosql\n// DO NOT EDIT\n\nconst Schema = `", pkg)

	body := strings.Builder{}
	t := reflect.TypeOf(inputType)
	for i := 0; i < t.NumField(); i++ {
		table := t.Field(i) // Fields of the outer type are expected to be tables.
		if table.Type.Kind() != reflect.Slice {
			panic(`Expected table should be a slice: ` + table.Name)
		}
		body.WriteString("CREATE TABLE IF NOT EXISTS ")
		body.WriteString(table.Name)
		body.WriteString(" (")
		row := table.Type.Elem()
		wasFirst := true
		for j := 0; j < row.NumField(); j++ {
			col := row.Field(j)
			sqlText, ok := col.Tag.Lookup("sql")
			if !ok {
				panic(`Field missing "sql" tag:` + table.Name + "." + row.Name())
			}
			if !wasFirst {
				body.WriteString(",")
			}
			wasFirst = false
			body.WriteString("\n  ")
			body.WriteString(strings.TrimSpace(sqlText))
		}
		body.WriteString("\n);\n")
	}
	const footer = "`\n"
	return header + body.String() + footer
}
