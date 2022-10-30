// Code generated by @dung.ht and @bao.nq1 - DO NOT EDIT.
package store

import (
	"context"
	"database/sql"
	"regexp"
	"strconv"
	"strings"
)

var omittedSign = regexp.MustCompile(`(?m)--\*\$(\d+)$`)

func (q *Queries) filterOmitted(ctx context.Context, query string, omitted []bool, args []interface{}) (*sql.Stmt, string, []interface{}) {
	var filterQuery string
	for _, line := range strings.Split(query, "\n") {
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
			argRe := regexp.MustCompile(`(?m)([^*]\$)` + s)
			filterQuery = argRe.ReplaceAllString(filterQuery, "${1}"+r)
		}
	}

	filterStmt, _ := q.db.PrepareContext(ctx, filterQuery)

	return filterStmt, filterQuery, filterArgs
}
