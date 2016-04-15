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
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
	"strings"
	"unicode/utf8"
)

// showCmd represents the show command
var showWifiCmd = &cobra.Command{
	Use:   "show",
	Short: "Show WiFi settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		readWifiConfig()

		// if we have command line parameters only add those to our wifi configuration
		if cmdInterfaceName != "" {
			for key := range myWifiConfig.Interfaces {
				delete(myWifiConfig.Interfaces, key)
			}
			fmt.Println(cmdInterfaceName)
			myWifiConfig.Interfaces = make(map[string]Credentials)
			myWifiConfig.Interfaces[cmdInterfaceName] = Credentials{Ssid: "", Password: ""}
		}

		if len(myWifiConfig.Interfaces) > 0 {
			for interfaceName := range myWifiConfig.Interfaces {
				configFilePath := path.Join(networkInterfacesPath, interfaceName)

				input, err := ioutil.ReadFile(configFilePath)
				if err == nil {
					lines := strings.Split(string(input), "\n")

					header := fmt.Sprintf("\n%s in %s\n", interfaceName, configFilePath)
					fmt.Printf(header)
					fmt.Println(strings.Repeat("-", utf8.RuneCountInString(header)))
					for _, line := range lines {
						fmt.Println(line)
					}
				} else {
					panic(err)
				}
			}
		} else {
			cmd.Help()
		}
	},
}

func init() {
	wifiCmd.AddCommand(showWifiCmd)
}
