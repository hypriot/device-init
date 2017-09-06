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
	"crypto/sha1"
	"encoding/hex"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/pbkdf2"
)

// setCmd represents the set command
var setWifiCmd = &cobra.Command{
	Use:   "set",
	Short: "Set WiFi settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		setWifi()
	},
}

func setWifi() {
	readWifiConfig()

	// if we have command line parameters only add those to our wifi configuration
	if cmdInterfaceName != "" && cmdSsid != " " && cmdPassword != "" {
		for key := range myWifiConfig.Interfaces {
			delete(myWifiConfig.Interfaces, key)
		}
		myWifiConfig.Interfaces = make(map[string]Credentials)
		myWifiConfig.Interfaces[cmdInterfaceName] = Credentials{Ssid: cmdSsid, Password: cmdPassword}
	}

	for interfaceName, interfaceCredentials := range myWifiConfig.Interfaces {
		variables := Interface{
			Name: interfaceName,
			SSID: interfaceCredentials.Ssid,
			Password:  createEncryptedPsk([]byte(interfaceCredentials.Password), []byte(interfaceCredentials.Ssid)),
		}
		interfaceString := generateInterfaceConfig(variables)
		applyInterfaceConfig(interfaceName, interfaceString)
	}
}

func init() {
	wifiCmd.AddCommand(setWifiCmd)
	setWifiCmd.Flags().StringVarP(&cmdSsid, "ssid", "s", "", "SSID is the name of your wireless network")
	setWifiCmd.Flags().StringVarP(&cmdPassword, "password", "p", "", "Password for your wireless network")
}

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

// https://godoc.org/golang.org/x/crypto
// http://docs.ros.org/diamondback/api/wpa_supplicant/html/wpa__passphrase_8c_source.html
func createEncryptedPsk(password, salt []byte) string {
	defer clear(password)
	result := pbkdf2.Key(password, salt, 4096, 32, sha1.New)
	return hex.EncodeToString(result)
}
