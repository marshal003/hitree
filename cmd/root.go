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
	"os"
	"reflect"
	"strings"

	"github.com/marshal003/hitree/cmd/tree"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var opt tree.Options
var cfgFile string

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
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hitree",
	Short: "Print tree of the directory",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		nocolor := viper.GetBool("nocolor")
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
		root := "."
		if len(args) >= 1 {
			root = args[0]
		}

		tree, err := tree.TraverseDir(root, opt, 0)
		if err != nil {
			return err
		}
		tree.Print(opt)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hitree.yaml)")
	rootCmd.Flags().BoolP("dironly", "d", false, "List only directories")
	rootCmd.Flags().BoolP("all", "a", false, "List all files & directories including hidden ones")
	rootCmd.Flags().BoolP("fullpath", "f", false, "Print full path prefix for all files")
	rootCmd.Flags().BoolP("noreport", "", false, "Omits printing of the file and directory report at the end of the tree listing.")
	rootCmd.Flags().BoolP("followlink", "l", false, "Follow link and list files in the link is for a directory")
	rootCmd.Flags().BoolP("prune", "", false, "Makes tree prune empty directories from the output")
	rootCmd.Flags().Bool("nocolor", false, "Directory Color")
	rootCmd.Flags().Int16P("level", "L", -1, "Max display depth of the directory tree")
	rootCmd.Flags().StringP("includepattern", "P", "", "List only those files which matches to wild-card pattern")
	rootCmd.Flags().StringP("excludepattern", "I", "", "Do not list those files that match the wild-card pattern.")
	rootCmd.Flags().String("dircolor", "gray", "Directory Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	rootCmd.Flags().String("filecolor", "green", "File Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	rootCmd.Flags().String("symlinkcolor", "blue", "SymLink Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	rootCmd.Flags().String("tlinkcolor", "brown", "TLink Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	rootCmd.Flags().String("llinkcolor", "brown", "Pipe Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")
	rootCmd.Flags().String("pipecolor", "brown", "Pipe Color(gray/b, green/b, blue/b, brown/b, red/b, black/b, magenta/b, cyan/b)")

	//Bind viper
	viper.BindPFlag("dironly", rootCmd.Flags().Lookup("dironly"))
	viper.BindPFlag("nocolor", rootCmd.Flags().Lookup("nocolor"))
	viper.BindPFlag("all", rootCmd.Flags().Lookup("all"))
	viper.BindPFlag("fullpath", rootCmd.Flags().Lookup("fullpath"))
	viper.BindPFlag("noreport", rootCmd.Flags().Lookup("noreport"))
	viper.BindPFlag("followlink", rootCmd.Flags().Lookup("followlink"))
	viper.BindPFlag("level", rootCmd.Flags().Lookup("level"))
	viper.BindPFlag("includepattern", rootCmd.Flags().Lookup("includepattern"))
	viper.BindPFlag("excludepattern", rootCmd.Flags().Lookup("excludepattern"))

	viper.BindPFlag("nocolor", rootCmd.Flags().Lookup("nocolor"))
	viper.BindPFlag("dircolor", rootCmd.Flags().Lookup("dircolor"))
	viper.BindPFlag("filecolor", rootCmd.Flags().Lookup("filecolor"))
	viper.BindPFlag("symlinkcolor", rootCmd.Flags().Lookup("symlinkcolor"))
	viper.BindPFlag("tlinkcolor", rootCmd.Flags().Lookup("tlinkcolor"))
	viper.BindPFlag("llinkcolor", rootCmd.Flags().Lookup("llinkcolor"))
	viper.BindPFlag("pipecolor", rootCmd.Flags().Lookup("pipecolor"))
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

	// 	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
