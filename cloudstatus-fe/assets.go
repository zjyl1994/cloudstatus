package cloudstatusfe

import (
	"embed"
)

//go:embed build/client/*
var FrontendAssets embed.FS
