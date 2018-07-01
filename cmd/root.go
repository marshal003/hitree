// Copyright © 2018 Vinit Kumar Rai <vinitrai.marshal@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"

	tree "github.com/marshal003/hitree/core"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var opt tree.Options
var cfgFile string

//BuildInfo populated from goreleaser
type buildInfo struct {
	version string
	commit  string
	date    string
}

func (b buildInfo) Print() {
	fmt.Printf("\n\nVersion: %v, Build: %v, Build Date: %v\n\n", b.version, b.commit, b.date)
}

var release buildInfo

func releaseDetails(version, date, commit interface{}) {
	release.version = version.(string)
	release.commit = commit.(string)
	release.date = date.(string)
}

func setColorOption(cmd *cobra.Command, colorHolder *tree.Colorize, flag string) {
	color := viper.GetString(flag)
	var ok bool
	if *colorHolder, ok = tree.ColorMap[strings.ToLower(color)]; !ok {
		fmt.Printf("Invalid color specified for %s. Possible Colors are: %v", flag, reflect.ValueOf(tree.ColorMap).MapKeys())
		*colorHolder = tree.ColorMap["blue"]
	}
}

func initOptions() {
	opt.DirOnly = viper.GetBool("dironly")
	opt.IncludeHidden = viper.GetBool("all")
	opt.ShowFullPath = viper.GetBool("fullpath")
	opt.NoReport = viper.GetBool("noreport")
	opt.FollowLink = viper.GetBool("followlink")
	opt.Prune = viper.GetBool("prune")
	opt.MaxLevel = int16(viper.GetInt("level"))
	opt.ExcludePattern = viper.GetString("excludepattern")
	opt.IncludePattern = viper.GetString("includepattern")
	opt.Indent = viper.GetInt("jsonindent")
	opt.OutputPath = viper.GetString("output")
	opt.FileLimit = viper.GetInt("filelimit")
	opt.PrintGID = (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && viper.GetBool("group")
	opt.PrintUID = (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && viper.GetBool("user")
	opt.PrintSize = viper.GetBool("size")
	opt.PrintProtection = viper.GetBool("protection")
	opt.PrintModTime = viper.GetBool("modtime")
	opt.SortReverse = viper.GetBool("reverse")
	opt.SortByModTime = viper.GetBool("sortbymodtime")
	opt.TimeFormat = viper.GetString("timefmt")
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hitree",
	Short: "Print tree structure of the directory",
	Long: `Golang implementation of popular tree command from linux.q
Note: windows 10 has issue with ansi color, so for this release we have
disable color outputs on windows platform.
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		nocolor := (viper.GetBool("nocolor") || (runtime.GOOS == "windows"))
		tree.InitAurora(!nocolor)
		initOptions()
		setColorOption(cmd, &opt.DirColor, "dircolor")
		setColorOption(cmd, &opt.FileColor, "filecolor")
		setColorOption(cmd, &opt.SymLinkColor, "symlinkcolor")
		setColorOption(cmd, &opt.PipeColor, "pipecolor")
		setColorOption(cmd, &opt.TLinkColor, "tlinkcolor")
		setColorOption(cmd, &opt.LLinkColor, "llinkcolor")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			release.Print()
			return nil
		}
		path := "."
		if len(args) >= 1 {
			path = args[0]
		}

		root, err := tree.TraverseDir(path, opt, 0)
		if err != nil {
			return err
		}
		return sendOutput(cmd, root)
	},
}

func sendOutput(cmd *cobra.Command, root tree.Tree) error {
	asJSON, _ := cmd.Flags().GetBool("json")
	var w io.WriteCloser
	if opt.OutputPath != "stdout" {
		f, err := os.Create(opt.OutputPath)
		if err != nil {
			panic(fmt.Sprintf("Unable to open output file %v", err))
		}
		w = f
	} else {
		w = os.Stdout
	}
	defer w.Close()

	if !asJSON {
		root.Print(w, opt)
		return nil
	}
	res, err := root.AsJSONString(opt)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\n%s\n", res)
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, commit, date interface{}) {
	releaseDetails(version, date, commit)
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.Flags().SortFlags = false
	RootCmd.PersistentFlags().SortFlags = false
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hitree.yaml)")
	RootCmd.Flags().BoolP("version", "v", false, "Version of hitree command")
	RootCmd.Flags().BoolP("json", "j", false, "Print Tree structure as JSON")
	RootCmd.Flags().Int("jsonindent", 2, "JSON Indentation")
	RootCmd.Flags().StringP("output", "o", "stdout", "Put result in the output file")
	RootCmd.Flags().BoolP("dironly", "d", false, "List only directories")
	RootCmd.Flags().BoolP("all", "a", false, "List all files & directories including hidden ones")
	RootCmd.Flags().BoolP("fullpath", "f", false, "Print full path prefix for all files")
	RootCmd.Flags().BoolP("noreport", "", false, "Omits printing of the file and directory report at the end of the tree listing.")
	RootCmd.Flags().BoolP("followlink", "l", false, "Follow link and list files in the link is for a directory")
	RootCmd.Flags().BoolP("prune", "", false, "Makes tree prune empty directories from the output")
	RootCmd.Flags().BoolP("nocolor", "n", false, "Turn colorization off always")
	RootCmd.Flags().Int16P("level", "L", -1, "Max display depth of the directory tree")

	//New
	RootCmd.Flags().Int("filelimit", -1, "Do not descend directories that contain more than # entries.")
	RootCmd.Flags().String("timefmt", "Jan 2 13:10", "Prints (implies -D) and formats the date according to the format string")
	RootCmd.Flags().BoolP("protection", "p", false, "Print Protection on file")
	RootCmd.Flags().BoolP("size", "s", false, "Print Size on file")
	RootCmd.Flags().BoolP("user", "u", false, "Print the username, or UID")
	RootCmd.Flags().BoolP("group", "g", false, "Print the group name, or GID")
	RootCmd.Flags().BoolP("modtime", "D", false, "Print the date of the last modification time for the file listed")
	RootCmd.Flags().BoolP("reverse", "r", false, "Sort the output in reverse alphabetic order")
	RootCmd.Flags().BoolP("sortbymodtime", "t", false, "Sort the output by last modification time instead of alphabetically")

	//Pattern flags
	RootCmd.Flags().StringP("includepattern", "P", "", "List only those files which matches to wild-card pattern")
	RootCmd.Flags().StringP("excludepattern", "I", "", "Do not list those files that match the wild-card pattern.")

	//Color flag
	RootCmd.Flags().String("dircolor", "gray", "Directory Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	RootCmd.Flags().String("filecolor", "green", "File Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	RootCmd.Flags().String("symlinkcolor", "blue", "SymLink Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	RootCmd.Flags().String("tlinkcolor", "brown", "TLink Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	RootCmd.Flags().String("llinkcolor", "brown", "Pipe Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	RootCmd.Flags().String("pipecolor", "brown", "Pipe Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")

	//Bind viper
	viper.BindPFlag("filelimit", RootCmd.Flags().Lookup("filelimit"))
	viper.BindPFlag("timefmt", RootCmd.Flags().Lookup("timefmt"))
	viper.BindPFlag("protection", RootCmd.Flags().Lookup("protection"))
	viper.BindPFlag("size", RootCmd.Flags().Lookup("size"))
	viper.BindPFlag("user", RootCmd.Flags().Lookup("user"))
	viper.BindPFlag("group", RootCmd.Flags().Lookup("group"))
	viper.BindPFlag("modtime", RootCmd.Flags().Lookup("modtime"))
	viper.BindPFlag("reverse", RootCmd.Flags().Lookup("reverse"))
	viper.BindPFlag("sortbymodtime", RootCmd.Flags().Lookup("sortbymodtime"))

	viper.BindPFlag("dironly", RootCmd.Flags().Lookup("dironly"))
	viper.BindPFlag("output", RootCmd.Flags().Lookup("output"))
	viper.BindPFlag("nocolor", RootCmd.Flags().Lookup("nocolor"))
	viper.BindPFlag("all", RootCmd.Flags().Lookup("all"))
	viper.BindPFlag("fullpath", RootCmd.Flags().Lookup("fullpath"))
	viper.BindPFlag("noreport", RootCmd.Flags().Lookup("noreport"))
	viper.BindPFlag("followlink", RootCmd.Flags().Lookup("followlink"))
	viper.BindPFlag("level", RootCmd.Flags().Lookup("level"))
	viper.BindPFlag("includepattern", RootCmd.Flags().Lookup("includepattern"))
	viper.BindPFlag("excludepattern", RootCmd.Flags().Lookup("excludepattern"))
	viper.BindPFlag("jsonindent", RootCmd.Flags().Lookup("jsonindent"))

	viper.BindPFlag("nocolor", RootCmd.Flags().Lookup("nocolor"))
	viper.BindPFlag("dircolor", RootCmd.Flags().Lookup("dircolor"))
	viper.BindPFlag("filecolor", RootCmd.Flags().Lookup("filecolor"))
	viper.BindPFlag("symlinkcolor", RootCmd.Flags().Lookup("symlinkcolor"))
	viper.BindPFlag("tlinkcolor", RootCmd.Flags().Lookup("tlinkcolor"))
	viper.BindPFlag("llinkcolor", RootCmd.Flags().Lookup("llinkcolor"))
	viper.BindPFlag("pipecolor", RootCmd.Flags().Lookup("pipecolor"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hitree" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hitree")
	}
	viper.AutomaticEnv() // read in environment variables that match
}
