package core

//Options Data model to hold command line options
type Options struct {
	IncludeHidden   bool
	DirOnly         bool
	ShowFullPath    bool
	NoReport        bool
	FollowLink      bool
	Prune           bool
	PrintProtection bool
	PrintSize       bool
	PrintUID        bool
	PrintGID        bool
	PrintModTime    bool
	SortReverse     bool
	SortByModTime   bool
	FileLimit       int
	MaxLevel        int16
	Indent          int
	TimeFormat      string
	IncludePattern  string
	ExcludePattern  string
	OutputPath      string
	DirColor        Colorize
	FileColor       Colorize
	SymLinkColor    Colorize
	TLinkColor      Colorize
	LLinkColor      Colorize
	PipeColor       Colorize
}

// DefaultOptions A utility method to create default Options for hitree command
// This is intensionally created for test cases
func DefaultOptions() Options {
	opt := Options{
		IncludeHidden:  false,
		DirOnly:        false,
		ShowFullPath:   false,
		NoReport:       false,
		FollowLink:     false,
		Prune:          false,
		MaxLevel:       -1,
		IncludePattern: "",
		ExcludePattern: "",
		DirColor:       ColorMap["gray"],
		FileColor:      ColorMap["gray"],
		SymLinkColor:   ColorMap["gray"],
		TLinkColor:     ColorMap["gray"],
		LLinkColor:     ColorMap["gray"],
		PipeColor:      ColorMap["gray"],
	}
	return opt
}
