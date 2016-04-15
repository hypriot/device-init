// Copyright © 2016 Govinda Fichtner <govinda.fichtner@googlemail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config = viper.New()

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "device-init",
	Short: "Initialize your device",
	Long: `device-init allows you to configure various aspect of your devices.
This ranges from configuration as simple as setting a hostname to
more complex stuff as configuring WiFi access.`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" {
			set_all_commands()
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func set_all_commands() {
	// If a config file is found, do stuff for all settings that are present
	if err := config.ReadInConfig(); err == nil {
		set_hostname()
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /boot/device-init.yaml)")
	RootCmd.PersistentFlags().Lookup("config").NoOptDefVal = "/boot/device-init.yaml"
}

// initConfig reads in config file
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		config.SetConfigFile(cfgFile)

		// If a config file is found, read it in.
		if err := config.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", config.ConfigFileUsed())
		}
	}

}
