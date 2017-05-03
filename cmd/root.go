// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type coloredPattern struct {
	color   string
	pattern string
}

type coloredPatterns []coloredPattern

func (i *coloredPatterns) String() string {
	return "red,(ERROR)"
}

func (i *coloredPatterns) Set(value string) error {
	parts := strings.SplitN(value, ",", 2)
	if len(parts) != 2 {
		return errors.New(fmt.Sprintf("value %s isn't correctly formatted", value))
	}

	*i = append(*i, coloredPattern{color: parts[0], pattern: parts[1]})
	return nil
}

func (i *coloredPatterns) Type() string {
	return "color pattern combo"
}

var cfgFile string
var useStdout bool
var stdoutPrefix string
var colors coloredPatterns = make(coloredPatterns, 0, 0)

type colorFunc func(...interface{}) string
type lookup struct {
	re *regexp.Regexp
	co colorFunc
}

func Paint(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)
	var printerLookup []lookup = make([]lookup, 0, len(colors))

	for _, cp := range colors {
		var cf colorFunc
		switch cp.color {
		case "black":
			cf = color.New(color.FgBlack).SprintFunc()
		case "red":
			cf = color.New(color.FgRed).SprintFunc()
		case "green":
			cf = color.New(color.FgGreen).SprintFunc()
		case "yellow":
			cf = color.New(color.FgYellow).SprintFunc()
		case "blue":
			cf = color.New(color.FgBlue).SprintFunc()
		case "magenta":
			cf = color.New(color.FgMagenta).SprintFunc()
		case "cyan":
			cf = color.New(color.FgCyan).SprintFunc()
		case "white":
			cf = color.New(color.FgWhite).SprintFunc()
		default:
			cf = color.New(color.FgRed).SprintFunc()
		}

		re := regexp.MustCompile(cp.pattern)
		printerLookup = append(printerLookup, lookup{re, cf})
	}

	for scanner.Scan() {
		line := scanner.Text()
		for _, l := range printerLookup {
			line = l.re.ReplaceAllString(line, l.co("$1"))
		}
		if len(stdoutPrefix) != 0 {
			line = stdoutPrefix + line
		}
		fmt.Println(line)
	}
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "paint",
	Short: "Paint output based on regex",
	Long:  `Paint stdin and/or stderr based on simple regexs`,
	Run:   Paint,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.paint.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolVarP(&useStdout, "stdout", "o", true, "Paint stdout")
	RootCmd.Flags().StringVarP(&stdoutPrefix, "stdout-prefix", "S", "", "Prefix for stdout")
	RootCmd.Flags().VarP(&colors, "colors", "c", "color/pattern combination. Pattern must have a capture group to paint: 'red,(ERROR)' would print all ERROR in red")
	// RootCmd.Flags().BoolP("stderr", "e", false, "Paint stdout")
	// RootCmd.Flags().StringP("stderr-prefix", "E", "", "Prefix for stderr")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".paint") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
