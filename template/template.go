package template

import "embed"

//go:embed base/*
var BaseFiles embed.FS

//go:embed fragments/database/*
var DatabaseFragments embed.FS
