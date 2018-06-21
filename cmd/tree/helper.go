package tree

import "github.com/logrusorgru/aurora"

//Colorize A function type which is being used in Options
type Colorize func(interface{}) aurora.Value

//Options holds command line options
type Options struct {
	IncludeHidden  bool
	DirOnly        bool
	ShowFullPath   bool
	NoReport       bool
	FollowLink     bool
	Prune          bool
	MaxLevel       int16
	IncludePattern string
	ExcludePattern string
	DirColor       Colorize
	FileColor      Colorize
	SymLinkColor   Colorize
	TLinkColor     Colorize
	LLinkColor     Colorize
	PipeColor      Colorize
}

var au aurora.Aurora

//InitAurora Initialize Aurora, which is being used to print colorize output
func InitAurora(enableColorOutput bool) {
	au = aurora.NewAurora(enableColorOutput)
}

// Create a partial function based on color and bold, which then
// uses aurora to colorize the string provided to the message.
// This way, we are restricting the uses of aurora instance only at one place and
// are free to use any other library for coloring
var colorCurry = func(color aurora.Color, bold bool) Colorize {
	return func(message interface{}) aurora.Value {
		if bold {
			return au.Colorize(message, color).Bold()
		}
		return au.Colorize(message, color)
	}
}

// ColorMap A map of color to its colorize function, which will be obtained by
// calling abaove colorCurry method with required color
var ColorMap = map[string]Colorize{
	//gray
	"gray":  colorCurry(aurora.GrayFg, false),
	"grayb": colorCurry(aurora.GrayFg, true),

	//gree
	"green":  colorCurry(aurora.GreenFg, false),
	"greenb": colorCurry(aurora.GreenFg, true),

	//blue
	"blue":  colorCurry(aurora.BlueFg, false),
	"blueb": colorCurry(aurora.BlueFg, true),

	//brown
	"brown":  colorCurry(aurora.BrownFg, false),
	"brownb": colorCurry(aurora.BrownFg, true),

	//red
	"red":  colorCurry(aurora.RedFg, false),
	"redb": colorCurry(aurora.RedFg, true),

	//black
	"black":  colorCurry(aurora.BlackFg, false),
	"blackb": colorCurry(aurora.BlackFg, true),

	//magenta
	"magenta":  colorCurry(aurora.MagentaFg, false),
	"magentab": colorCurry(aurora.MagentaFg, true),

	//cyan
	"cyan":  colorCurry(aurora.CyanFg, false),
	"cyanb": colorCurry(aurora.CyanFg, true),
}
