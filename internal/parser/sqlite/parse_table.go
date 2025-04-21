package parser

import (
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/gitkumi/snowflake/grammar/sqlite"
)

type TableSchema struct {
	TableName string
	Columns   []ColumnDef
}

type ColumnDef struct {
	Name            string
	Type            string
	Constraints     []string
	IsPrimaryKey    bool
	IsAutoIncrement bool
	IsNotNull       bool
}

type SQLTableExtractor struct {
	sqlite.BaseSQLiteParserListener
	CurrentTable             string
	CurrentColumnName        string
	CurrentColumnType        string
	CurrentColumnConstraints []string
	Tables                   map[string][]ColumnDef
	IsInColumnDef            bool
	IsInConstraint           bool
	IsPrimaryKey             bool
	IsAutoIncrement          bool
	IsNotNull                bool
	RawSQL                   string
}

func NewSQLTableExtractor() *SQLTableExtractor {
	return &SQLTableExtractor{
		Tables: make(map[string][]ColumnDef),
	}
}

func (s *SQLTableExtractor) EnterCreate_table_stmt(ctx *sqlite.Create_table_stmtContext) {
	if ctx.Table_name() != nil {
		s.CurrentTable = strings.TrimSpace(ctx.Table_name().GetText())
		if _, exists := s.Tables[s.CurrentTable]; !exists {
			s.Tables[s.CurrentTable] = []ColumnDef{}
		}
	}
}

func (s *SQLTableExtractor) EnterColumn_def(ctx *sqlite.Column_defContext) {
	s.IsInColumnDef = true
	s.CurrentColumnConstraints = []string{}
	s.IsPrimaryKey = false
	s.IsAutoIncrement = false
	s.IsNotNull = false

	if ctx.Column_name() != nil {
		s.CurrentColumnName = strings.TrimSpace(ctx.Column_name().GetText())
	}

	if ctx.Type_name() != nil {
		s.CurrentColumnType = strings.TrimSpace(ctx.Type_name().GetText())
	}
}

func (s *SQLTableExtractor) ExitColumn_def(ctx *sqlite.Column_defContext) {
	if s.CurrentTable != "" && s.CurrentColumnName != "" {
		if s.RawSQL != "" {
			s.extractDefaultConstraints()
		}

		column := ColumnDef{
			Name:            s.CurrentColumnName,
			Type:            s.CurrentColumnType,
			Constraints:     s.CurrentColumnConstraints,
			IsPrimaryKey:    s.IsPrimaryKey,
			IsAutoIncrement: s.IsAutoIncrement,
			IsNotNull:       s.IsNotNull,
		}

		s.Tables[s.CurrentTable] = append(s.Tables[s.CurrentTable], column)
	}

	s.IsInColumnDef = false
	s.CurrentColumnName = ""
	s.CurrentColumnType = ""
}

func (s *SQLTableExtractor) extractDefaultConstraints() {
	columnRegex := regexp.MustCompile(
		`(?i)` + regexp.QuoteMeta(s.CurrentColumnName) + `\s+` +
			regexp.QuoteMeta(s.CurrentColumnType) +
			`(?:\s+NOT\s+NULL)?(?:\s+UNIQUE)?(?:\s+DEFAULT\s+([^\s,)]+))?`)

	matches := columnRegex.FindStringSubmatch(s.RawSQL)
	if len(matches) > 1 && matches[1] != "" {
		defaultValue := strings.TrimSpace(matches[1])
		hasDefault := false
		for _, constraint := range s.CurrentColumnConstraints {
			if strings.HasPrefix(constraint, "DEFAULT ") {
				hasDefault = true
				break
			}
		}

		if !hasDefault {
			s.CurrentColumnConstraints = append(s.CurrentColumnConstraints, "DEFAULT "+defaultValue)
		}
	}
}

func (s *SQLTableExtractor) EnterColumn_constraint(ctx *sqlite.Column_constraintContext) {
	if !s.IsInColumnDef {
		return
	}

	s.IsInConstraint = true

	if ctx.PRIMARY_() != nil && ctx.KEY_() != nil {
		s.IsPrimaryKey = true
		s.CurrentColumnConstraints = append(s.CurrentColumnConstraints, "PRIMARY KEY")
	}

	if ctx.AUTOINCREMENT_() != nil {
		s.IsAutoIncrement = true
		s.CurrentColumnConstraints = append(s.CurrentColumnConstraints, "AUTOINCREMENT")
	}

	if ctx.NOT_() != nil && ctx.NULL_() != nil {
		s.IsNotNull = true
		s.CurrentColumnConstraints = append(s.CurrentColumnConstraints, "NOT NULL")
	}

	if ctx.UNIQUE_() != nil {
		s.CurrentColumnConstraints = append(s.CurrentColumnConstraints, "UNIQUE")
	}

	if ctx.DEFAULT_() != nil && ctx.Expr() != nil {
		defaultValue := strings.TrimSpace(ctx.Expr().GetText())
		s.CurrentColumnConstraints = append(s.CurrentColumnConstraints, "DEFAULT "+defaultValue)
	}
}

func (s *SQLTableExtractor) ExitColumn_constraint(ctx *sqlite.Column_constraintContext) {
	s.IsInConstraint = false
}

func ParseTable(sqlContent string) ([]TableSchema, error) {
	input := antlr.NewInputStream(sqlContent)

	lexer := sqlite.NewSQLiteLexer(input)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := sqlite.NewSQLiteParser(tokens)

	extractor := NewSQLTableExtractor()
	extractor.RawSQL = sqlContent

	tree := p.Parse()
	antlr.ParseTreeWalkerDefault.Walk(extractor, tree)

	var schemas []TableSchema
	for tableName, columns := range extractor.Tables {
		schema := TableSchema{
			TableName: tableName,
			Columns:   columns,
		}
		schemas = append(schemas, schema)
	}

	return schemas, nil
}
