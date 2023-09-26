package constants

import (
	"time"
)

const (
	// Version is the current version of picfit
	Version = "0.14.0"
)

var (
	// Branch is the compiled branch
	Branch string

	// Revision is the compiled revision
	Revision string

	// LatestCommitMessage is the latest commit message
	LatestCommitMessage string

	// BuildTime is the compiled build time
	BuildTime string

	// Compiler is the compiler used during build
	Compiler string
)

const (
	BottomLeft  = "bottom-left"
	BottomRight = "bottom-right"
	TopLeft     = "top-left"
	TopRight    = "top-right"
)

var StickPositions = []string{
	BottomLeft,
	BottomRight,
	TopLeft,
	TopRight,
}

const ModifiedTimeFormat = time.RFC1123

const RequestIDCtx = "request-id"
