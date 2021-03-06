package core

import "github.com/logrusorgru/aurora"

//Colorize A function type which is being used in Options
type Colorize func(interface{}) aurora.Value

var au aurora.Aurora

// InitAurora Utility method to initialize Aurora(ANSI color library),
// being used to print colorize output.
// Aurora Provides flag to allow enabling or disabling output
func InitAurora(enableColorOutput bool) {
	au = aurora.NewAurora(enableColorOutput)
}

// colorCurry: A partial function based on color and bold, which then
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

// ColorMap A map of supported colors, wich uses internal partial function and returns
// a function of type Colorize. Uses of partial function allows us to limit the uses of
// Aurora's way of colorizing the string at one place and gives us the flexibility to
// easily switch to any other lib in future.
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
