package db

import (
	"strings"
)

const (
	argSchema = "{schema}"
	argTable  = "{table}"
)

func WithSchema(query, schema string) string {
	return strings.ReplaceAll(query, argSchema, schema)
}

func WithSchemaAndTable(query, schema, table string) string {
	return strings.NewReplacer(argSchema, schema, argTable, table).Replace(query)
}
