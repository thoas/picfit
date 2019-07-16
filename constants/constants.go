package constants

const (
	// Version is the current version of picfit
	Version = "0.6.0"
)

var (
	// Branch is the compiled branch
	Branch string

	// Revision is the compiled revision
	Revision string

	// BuildTime is the compiled build time
	BuildTime string

	// Compiler is the compiler used during build
	Compiler string
)

const (
	TopRight    = "top-right"
	TopLeft     = "top-left"
	BottomRight = "bottom-right"
	BottomLeft  = "bottom-left"
)

var StickPositions = []string{
	TopRight,
	TopLeft,
	BottomRight,
	BottomLeft,
}
