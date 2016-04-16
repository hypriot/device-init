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
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/pbkdf2"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"
)

// setCmd represents the set command
var setWifiCmd = &cobra.Command{
	Use:   "set",
	Short: "Set WiFi settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		set_wifi()
	},
}

func set_wifi() {
	var err error

	readWifiConfig()

	const interfaceTemlate = `allow-hotplug {{.Name}}

auto {{.Name}}
iface {{.Name}} inet dhcp
wpa-ssid {{.Ssid}}
wpa-psk {{.Psk}}
`

	// if we have command line parameters only add those to our wifi configuration
	if cmdInterfaceName != "" && cmdSsid != " " && cmdPassword != "" {
		for key := range myWifiConfig.Interfaces {
			delete(myWifiConfig.Interfaces, key)
		}
		myWifiConfig.Interfaces = make(map[string]Credentials)
		myWifiConfig.Interfaces[cmdInterfaceName] = Credentials{Ssid: cmdSsid, Password: cmdPassword}
	}

	for interfaceName, interfaceCredentials := range myWifiConfig.Interfaces {

		type templateVariables struct {
			Name, Ssid, Psk string
		}

		variables := templateVariables{
			Name: interfaceName,
			Ssid: interfaceCredentials.Ssid,
			Psk:  create_encrypted_psk([]byte(interfaceCredentials.Password), []byte(interfaceCredentials.Ssid)),
		}

		err = os.MkdirAll(networkInterfacesPath, 0755)
		if err != nil {
			fmt.Println("Could not create path: ", networkInterfacesPath)
		}

		configFilePath := path.Join(networkInterfacesPath, interfaceName)
		if _, err := os.Stat(configFilePath); err == nil {
			filepath, filename := path.Dir(configFilePath), path.Base(configFilePath)
			backupFile := "." + filename + ".backup"
			backupPath := path.Join(filepath, backupFile)
			err = os.Rename(configFilePath, backupPath)
			if err != nil {
				fmt.Println("Could not backup file ", backupPath, ": ", err)
			}
		}

		f, err := os.Create(configFilePath)
		defer f.Close()
		if err != nil {
			fmt.Println("Could not create file: "+configFilePath+": ", err)
		}

		t := template.Must(template.New("config").Parse(interfaceTemlate))

		err = t.Execute(f, variables)
		if err != nil {
			fmt.Println("Error writing configuration:", err)
		}

		if interfaceExistsAndIsDown(interfaceName) {
			err := exec.Command("ifup", interfaceName).Run()
			if err != nil {
				fmt.Println("Could not bring up interface :", interfaceName)
			}
		}

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
func create_encrypted_psk(password, salt []byte) string {
	defer clear(password)
	result := pbkdf2.Key(password, salt, 4096, 32, sha1.New)
	return hex.EncodeToString(result)
}

func interfaceExistsAndIsDown(interfaceName string) bool {
	output, err := exec.Command("ip", "link").Output()
	if err != nil {
		fmt.Println("Could not run 'ip link'", err)
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		interfaceExists, _ := regexp.MatchString(interfaceName, line)
		interfaceIsDown, _ := regexp.MatchString("state DOWN", line)
		if interfaceExists && interfaceIsDown {
			return true
		}
	}
	return false
}
