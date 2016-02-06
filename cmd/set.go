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
	"strings"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
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
	hostnameCmd.AddCommand(setCmd)
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

	if hostname == "" {
		fmt.Println("missing hostname argument")
		return
	}

	err := ioutil.WriteFile("/etc/hostname", []byte(hostname), 0644)
	if err != nil {
		panic(err)
	}

	input, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "127.0.0.1	localhost") {
			lines[i] = fmt.Sprintf("127.0.0.1	localhost	%s", hostname)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("/etc/hosts", []byte(output), 0644)
	if err != nil {
		panic(err)
	}

	set_hostname_cmd := exec.Command("hostname", hostname)
	err = set_hostname_cmd.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Set hostname: %s\n", hostname)
}
