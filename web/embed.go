// Package web embeds the static web assets so they can be served by the webapp binary.
package web

import "embed"

//go:embed static/*
var StaticFS embed.FS
