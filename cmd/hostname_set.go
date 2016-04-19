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
	"os/exec"
	"regexp"
	"strings"
)

// setCmd represents the set command
var setHostnameCmd = &cobra.Command{
	Use:   "set [hostname]",
	Short: "Set a hostname",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			set_hostname(args[0])
		} else {
			cmd.Help()
		}
	},
}

func init() {
	hostnameCmd.AddCommand(setHostnameCmd)
}

func set_hostname(args ...string) {
	hostname := ""

	// if we have hostname in config file use that
	if config.IsSet("hostname") {
		hostname = config.GetString("hostname")
	}

	// if we have a hostname as command line arg use that
	if len(args) > 0 {
		hostname = args[0]
	}

	if hostname == "" && cfgFile == "" {
		fmt.Println("missing hostname argument")
		return
	}

	if hostname != "" {
		err := ioutil.WriteFile("/etc/hostname", []byte(hostname), 0644)
		if err != nil {
			panic(err)
		}

		input, err := ioutil.ReadFile("/etc/hosts")
		if err != nil {
			panic(err)
		}

		lines := strings.Split(string(input), "\n")
		lines_new := []string{}

		for i, line := range lines {
			if strings.Contains(line, "127.0.0.1	localhost") {
				lines_new = append(lines_new, lines[0:i+1]...)
				lines_new = append(lines_new, fmt.Sprintf("127.0.0.1	%s", hostname))
				lines_new = append(lines_new, lines[i+1:]...)
			}
		}
		output := strings.Join(lines_new, "\n")
		err = ioutil.WriteFile("/etc/hosts", []byte(output), 0644)
		if err != nil {
			panic(err)
		}

		err = exec.Command("hostname", hostname).Run()
		if err != nil {
			fmt.Println("Unable to set hostname: ", err)
		}

		// ensure that dhcp server and avahi daemon are aware of new hostname
		for _, interfaceName := range activeInterfaces() {
			err = exec.Command("/sbin/ifdown", interfaceName).Run()
			if err != nil {
				fmt.Println("Unable to bring interface down: ", interfaceName, err)
			}

			err = exec.Command("/sbin/ifup", interfaceName).Run()
			if err != nil {
				fmt.Println("Unable to bring interface up: ", interfaceName, err)
			}
		}

		err = exec.Command("/bin/systemctl", "restart", "avahi-daemon.service").Run()
		if err != nil {
			fmt.Println("Unable to restart avahi-daemon: ", err)
		}

		err = exec.Command("/bin/systemctl", "restart", "rsyslog.service").Run()
		if err != nil {
			fmt.Println("Unable to restart rsyslog: ", err)
		}

		fmt.Printf("Set hostname: %s\n", hostname)
	}
}

func activeInterfaces() []string {
	var result []string
	output, err := exec.Command("ip", "link").Output()
	if err != nil {
		fmt.Println("Could not run 'ip link'", err)
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		interfaceIsUp, _ := regexp.MatchString("state UP", line)
		if interfaceIsUp {
			re := regexp.MustCompile(`^\d*:\s([a-z0-9@]*):`)
			result = append(result, re.FindStringSubmatch(line)[1])
		}
	}
	return result
}
