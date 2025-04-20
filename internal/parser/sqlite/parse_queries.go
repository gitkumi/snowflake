package parser

import (
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/gitkumi/snowflake/grammar/sqlite"
)

type QueryType string

const (
	Many QueryType = "many"
	One  QueryType = "one"
	Exec QueryType = "exec"
)

type Query struct {
	Name     string
	Type     QueryType
	SQL      string
	Table    string
	Params   []string
	Columns  []string
	IsSelect bool
	IsInsert bool
	IsUpdate bool
	IsDelete bool
}

type SQLQueryExtractor struct {
	sqlite.BaseSQLiteParserListener
	Query                 Query
	CurrentTable          string
	InSelectExpr          bool
	InTableExpr           bool
	InColumnsList         bool
	InValuesClause        bool
	ExtractedColumns      []string
	ExtractedTables       []string
	ExtractedPlaceholders int
}

func NewSQLQueryExtractor() *SQLQueryExtractor {
	return &SQLQueryExtractor{
		Query: Query{
			Params:  []string{},
			Columns: []string{},
		},
	}
}

func (s *SQLQueryExtractor) EnterSelect_stmt(ctx *sqlite.Select_stmtContext) {
	s.Query.IsSelect = true
}

func (s *SQLQueryExtractor) EnterInsert_stmt(ctx *sqlite.Insert_stmtContext) {
	s.Query.IsInsert = true
}

func (s *SQLQueryExtractor) EnterUpdate_stmt(ctx *sqlite.Update_stmtContext) {
	s.Query.IsUpdate = true
}

func (s *SQLQueryExtractor) EnterDelete_stmt(ctx *sqlite.Delete_stmtContext) {
	s.Query.IsDelete = true
}

func (s *SQLQueryExtractor) EnterTable_name(ctx *sqlite.Table_nameContext) {
	if s.Query.Table == "" && ctx.GetText() != "" {
		s.Query.Table = ctx.GetText()
	}
}

func (s *SQLQueryExtractor) EnterColumn_name(ctx *sqlite.Column_nameContext) {
	if s.Query.IsInsert && s.InColumnsList {
		s.Query.Columns = append(s.Query.Columns, ctx.GetText())
	}
}

func (s *SQLQueryExtractor) EnterColumn_name_list(ctx *sqlite.Column_name_listContext) {
	s.InColumnsList = true
}

func (s *SQLQueryExtractor) ExitColumn_name_list(ctx *sqlite.Column_name_listContext) {
	s.InColumnsList = false
}

func (s *SQLQueryExtractor) EnterLiteral_value(ctx *sqlite.Literal_valueContext) {
	// Check for bind parameter (?)
	if strings.Contains(ctx.GetText(), "?") {
		s.ExtractedPlaceholders++
		s.Query.Params = append(s.Query.Params, "?")
	}
}

func (s *SQLQueryExtractor) EnterExpr(ctx *sqlite.ExprContext) {
	if ctx.GetText() == "?" {
		s.Query.Params = append(s.Query.Params, "?")
	}
}

func (s *SQLQueryExtractor) FinalizeQuery(sql string) Query {
	s.Query.SQL = sql

	if len(s.Query.Params) == 0 && strings.Contains(sql, "?") {
		paramCount := strings.Count(sql, "?")
		for i := 0; i < paramCount; i++ {
			s.Query.Params = append(s.Query.Params, "?")
		}
	}

	if s.Query.IsSelect && len(s.Query.Columns) == 0 {
		columnRegex := regexp.MustCompile(`(?i)SELECT\s+(.*?)\s+FROM`)
		columnMatch := columnRegex.FindStringSubmatch(sql)
		if len(columnMatch) > 1 {
			if columnMatch[1] == "*" {
				s.Query.Columns = append(s.Query.Columns, "*")
			} else {
				columnsStr := columnMatch[1]
				columnList := strings.Split(columnsStr, ",")
				for _, col := range columnList {
					col = strings.TrimSpace(col)
					if strings.Contains(col, " AS ") {
						parts := strings.Split(strings.ToUpper(col), " AS ")
						col = strings.TrimSpace(parts[1])
					} else if strings.Contains(col, ".") {
						parts := strings.Split(col, ".")
						col = strings.TrimSpace(parts[1])
					}
					s.Query.Columns = append(s.Query.Columns, col)
				}
			}
		}
	}

	if s.Query.IsInsert && len(s.Query.Columns) == 0 {
		columnRegex := regexp.MustCompile(`(?i)INSERT\s+INTO\s+\w+\s*\(([^)]+)\)`)
		columnMatch := columnRegex.FindStringSubmatch(sql)
		if len(columnMatch) > 1 {
			columnsStr := columnMatch[1]
			columnList := strings.Split(columnsStr, ",")
			for _, col := range columnList {
				col = strings.TrimSpace(col)
				s.Query.Columns = append(s.Query.Columns, col)
			}
		}
	}

	return s.Query
}

func ParseQueries(sqlContent string) ([]Query, error) {
	var queries []Query

	queryDefRegex := regexp.MustCompile(`(?m)^--\s+name:\s+(\w+)\s+:(\w+)[\r\n]+([^-]*)`)
	matches := queryDefRegex.FindAllStringSubmatch(sqlContent, -1)

	for _, match := range matches {
		queryName := match[1]
		queryTypeStr := match[2]
		querySql := strings.TrimSpace(match[3])

		var queryType QueryType
		switch strings.ToLower(queryTypeStr) {
		case "many":
			queryType = Many
		case "one":
			queryType = One
		case "exec":
			queryType = Exec
		default:
			queryType = Exec
		}

		input := antlr.NewInputStream(querySql)
		lexer := sqlite.NewSQLiteLexer(input)
		tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
		p := sqlite.NewSQLiteParser(tokens)

		extractor := NewSQLQueryExtractor()

		tree := p.Parse()
		antlr.ParseTreeWalkerDefault.Walk(extractor, tree)

		query := extractor.FinalizeQuery(querySql)
		query.Name = queryName
		query.Type = queryType

		queries = append(queries, query)
	}

	return queries, nil
}
