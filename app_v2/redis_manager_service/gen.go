//go:generate sqlc generate
//go:generate go run gen.go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

const path = "internal/database/store/"

const dynamicGo = `// Code generated by @dung.ht and @bao.nq1 - DO NOT EDIT.
package store

import (
	"context"
	"database/sql"
	"regexp"
	"strconv"
	"strings"
)

var omittedSign = regexp.MustCompile(` + "`" + `(?m)--\*\$(\d+)$` + "`" + `)

func (q *Queries) filterOmitted(ctx context.Context, queries string, omitted []bool, args []interface{}) (*sql.Stmt, string, []interface{}) {
	var filterQuery string
	for _, line := range strings.Split(queries, "\n") {
		matches := omittedSign.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			argIndex, _ := strconv.Atoi(matches[0][1])
			if omitted[argIndex-1] {
				continue
			}
		}
		filterQuery += line + "\n"
	}

	var (
		filterArgs = make([]interface{}, 0, len(args))
		argMap     = make(map[string]string)
		argIdx     = 1
	)
	for i, o := range omitted {
		if o {
			argMap[strconv.Itoa(i+1)] = ""
		} else {
			argMap[strconv.Itoa(i+1)] = strconv.Itoa(argIdx)
			argIdx++
			filterArgs = append(filterArgs, args[i])
		}
	}
	for i, _ := range omitted {
		s := strconv.Itoa(i + 1)
		r := argMap[s]
		if r != "" && r != s {
			argRe := regexp.MustCompile(` + "`" + `(?m)([^*]\$)` + "`" + ` + s)
			filterQuery = argRe.ReplaceAllString(filterQuery, "${1}"+r)
		}
	}

	filterStmt, _ := q.db.PrepareContext(ctx, filterQuery)

	return filterStmt, filterQuery, filterArgs
}
`

func main() {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`(?m)// :dynamic\nfunc \(q \*Queries\) (\w+)\(ctx context\.Context, arg (\w+)\)([^{]+){\n(\s*)(_|rows|result), err := q.(exec|queries)\(ctx, (q.\w+), (\w+),((.|\n)*?)\)\n`)
	substitution := []byte("// :dynamic\nfunc (q *Queries) $1(ctx context.Context, arg $2)$3{\n${4}stmt, queries, args := q.filterOmitted(ctx, $8, arg.Omitted, []interface{}{$9})\n$4$5, err := q.$6(ctx, stmt, queries, args...,)\n")

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql.go") {
			continue
		}

		fileName := path + f.Name()

		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal(fmt.Errorf("error read file %v", err))
		}

		// patch methods
		parsedData := re.ReplaceAll(data, substitution)

		// patch struct definitions
		matches := re.FindAllSubmatch(data, -1)
		for _, m := range matches {
			structName := string(m[2])
			re2 := regexp.MustCompile(`(?m)type ` + structName + ` struct {\n(\s+)`)
			substitution2 := "type " + structName + " struct {\n${1}Omitted []bool\n$1"
			parsedData = re2.ReplaceAll(parsedData, []byte(substitution2))
		}

		err = ioutil.WriteFile(fileName, parsedData, 0o600)
		if err != nil {
			log.Fatal(fmt.Errorf("error write file %v", err))
		}

		fmt.Printf("generated %s\n", fileName)
	}

	// Gen
	dynamicFileName := path + "dynamic_gen.go"
	err = ioutil.WriteFile(dynamicFileName, []byte(dynamicGo), 0o600)
	if err != nil {
		log.Fatal(fmt.Errorf("error write file %v", err))
	}

	fmt.Printf("generated %s\n", dynamicFileName)
}
