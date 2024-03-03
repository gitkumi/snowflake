package template

import "embed"

//go:embed api/*
var ApiFiles embed.FS

//go:embed web/*
var WebFiles embed.FS
