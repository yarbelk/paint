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
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var useStdout bool
var stdoutPrefix, colorStr, pattern string

func Paint(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)
	var colorFunc func(...interface{}) string
	switch colorStr {
	case "black":
		colorFunc = color.New(color.FgBlack).SprintFunc()
	case "red":
		colorFunc = color.New(color.FgRed).SprintFunc()
	case "green":
		colorFunc = color.New(color.FgGreen).SprintFunc()
	case "yellow":
		colorFunc = color.New(color.FgYellow).SprintFunc()
	case "blue":
		colorFunc = color.New(color.FgBlue).SprintFunc()
	case "magenta":
		colorFunc = color.New(color.FgMagenta).SprintFunc()
	case "cyan":
		colorFunc = color.New(color.FgCyan).SprintFunc()
	case "white":
		colorFunc = color.New(color.FgWhite).SprintFunc()
	default:
		colorFunc = color.New(color.FgRed).SprintFunc()
	}

	re := regexp.MustCompile(pattern)

	for scanner.Scan() {
		line := scanner.Text()
		painted := re.ReplaceAllString(line, colorFunc("$1"))
		if len(stdoutPrefix) != 0 {
			painted = colorFunc(stdoutPrefix) + painted
		}
		fmt.Println(painted)
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
	RootCmd.Flags().StringVarP(&colorStr, "color", "c", "red", "color to paint")
	RootCmd.Flags().StringVarP(&pattern, "pattern", "p", "", "Pattern, like (DEBUG|INFO).  It must have a capture group.")
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
