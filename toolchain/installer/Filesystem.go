package installer

import "embed"

//go:embed programs/*
var Filesystem embed.FS
