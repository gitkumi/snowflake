package generate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Field struct {
	Name string
	Type string
}

type Method struct {
	Name       string
	Params     []Field
	ReturnType []string
}

type ParamStruct struct {
	Name   string
	Fields []Field
}

type TemplateData struct {
	Name         string
	Database     string
	Methods      []Method
	ParamStructs map[string]ParamStruct
	Resource     string
	ModuleName   string
}

func Generate() {
	projectName := "acme"
	repoPath := "testdata/books.sql.go"
	database := "postgres"
	moduleName := "repo"

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, repoPath, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("parse error: %v", err)
	}

	crudPrefixes := []string{"List", "Get", "Create", "Update", "Delete"}
	methods, resource := extractCRUDMethods(node, crudPrefixes, moduleName)

	if resource == "" {
		log.Fatalf("could not determine resource name from CRUD methods")
	}

	paramStructNames := []string{
		fmt.Sprintf("Create%sParams", resource),
		fmt.Sprintf("Update%sParams", resource),
	}
	paramStructs := extractStructFields(node, paramStructNames)

	data := TemplateData{
		Name:         projectName,
		Database:     database,
		Methods:      methods,
		ParamStructs: paramStructs,
		Resource:     resource,
		ModuleName:   moduleName,
	}

	fmt.Printf("Generating service for resource: %s\n", resource)

	if err := generateService(data); err != nil {
		log.Fatalf("template error: %v", err)
	}
}

func extractCRUDMethods(file *ast.File, crudPrefixes []string, moduleName string) ([]Method, string) {
	var methods []Method
	resourceCandidates := make(map[string]int)

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}

		if sel, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "Queries" {
				methodName := fn.Name.Name

				for _, prefix := range crudPrefixes {
					if strings.HasPrefix(methodName, prefix) {
						var params []Field
						if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
							for _, param := range fn.Type.Params.List {
								if isContextType(param.Type) {
									continue
								}

								var paramType bytes.Buffer
								printer.Fprint(&paramType, token.NewFileSet(), param.Type)
								paramTypeStr := paramType.String()

								if strings.HasSuffix(paramTypeStr, "Params") {
									paramTypeStr = fmt.Sprintf("%s.%s", moduleName, paramTypeStr)
								}

								for _, name := range param.Names {
									params = append(params, Field{
										Name: name.Name,
										Type: paramTypeStr,
									})
								}
							}
						}

						var returnTypes []string
						if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
							for _, result := range fn.Type.Results.List {
								var typeBuf bytes.Buffer
								printer.Fprint(&typeBuf, token.NewFileSet(), result.Type)
								typeStr := typeBuf.String()

								if !isBuiltinType(typeStr) {
									if strings.HasPrefix(typeStr, "[]") {
										innerType := strings.TrimPrefix(typeStr, "[]")
										if !isBuiltinType(innerType) {
											typeStr = fmt.Sprintf("[]%s.%s", moduleName, innerType)
										}
									} else {
										typeStr = fmt.Sprintf("%s.%s", moduleName, typeStr)
									}
								}

								returnTypes = append(returnTypes, typeStr)
							}
						}

						methods = append(methods, Method{
							Name:       methodName,
							Params:     params,
							ReturnType: returnTypes,
						})

						resourceName := strings.TrimPrefix(methodName, prefix)
						if resourceName != "" {
							resourceCandidates[resourceName]++
						}

						break
					}
				}
			}
		}
	}

	var resource string
	var highestCount int
	for candidate, count := range resourceCandidates {
		if count > highestCount {
			resource = candidate
			highestCount = count
		}
	}

	return methods, resource
}

func isBuiltinType(typeName string) bool {
	builtinTypes := map[string]bool{
		"bool":       true,
		"byte":       true,
		"complex128": true,
		"complex64":  true,
		"error":      true,
		"float32":    true,
		"float64":    true,
		"int":        true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"int8":       true,
		"rune":       true,
		"string":     true,
		"uint":       true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uint8":      true,
		"uintptr":    true,
	}

	return builtinTypes[typeName]
}

func isContextType(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "context" && sel.Sel.Name == "Context"
		}
	}
	return false
}

func extractStructFields(file *ast.File, targetStructs []string) map[string]ParamStruct {
	result := make(map[string]ParamStruct)
	structSet := make(map[string]bool)
	for _, name := range targetStructs {
		structSet[name] = true
	}

	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			name := typeSpec.Name.Name
			if !structSet[name] {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			var fields []Field
			for _, f := range structType.Fields.List {
				if len(f.Names) == 0 {
					continue
				}
				var buf bytes.Buffer
				printer.Fprint(&buf, token.NewFileSet(), f.Type)
				fields = append(fields, Field{
					Name: f.Names[0].Name,
					Type: buf.String(),
				})
			}
			result[name] = ParamStruct{
				Name:   name,
				Fields: fields,
			}
		}
	}
	return result
}

func generateService(data TemplateData) error {
	tmplPath := filepath.Join("internal", "commands", "generate", "templates", "service_test.go.templ")

	content, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New("service_test.go.templ").
		Funcs(template.FuncMap{
			"lowerFirst": func(s string) string {
				return strings.ToLower(s[:1]) + s[1:]
			},
			"title": strings.Title,
			"lower": strings.ToLower,
		}).
		Parse(string(content))

	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	outputFileName := fmt.Sprintf("testdata/%s_service.go", strings.ToLower(data.Resource))
	out, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	return tmpl.Execute(out, data)
}
