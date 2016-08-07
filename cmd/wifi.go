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
)

var cmdInterfaceName string
var cmdSsid string
var cmdPassword string

type Credentials struct {
	Ssid, Password string
}

type WifiConfig struct {
	Interfaces map[string]Credentials
}

var myWifiConfig WifiConfig

var networkInterfacesPath = "/etc/network/interfaces.d"

// wifiCmd represents the wifi command
var wifiCmd = &cobra.Command{
	Use:   "wifi",
	Short: "set WiFi settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(wifiCmd)

	wifiCmd.PersistentFlags().StringVarP(&cmdInterfaceName, "interface_name", "i", "", "Name of your wireless network interface e.g. wlan0")
}

func readWifiConfig() {
	err := config.UnmarshalKey("wifi", &myWifiConfig)
	if err != nil {
		fmt.Println("Could not unmarshal WiFi configuration")
	}
}
