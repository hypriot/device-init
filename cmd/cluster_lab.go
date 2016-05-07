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
	"os/exec"
	"strconv"
)

type ClusterLabConfig struct {
	Service map[string]string
}

var myClusterLabConfig ClusterLabConfig

// cluster-labCmd represents the cluster-lab command
var clusterlabCmd = &cobra.Command{
	Use:   "cluster-lab",
	Short: "manage cluster-lab settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		manageClusterLab()
	},
}

func init() {
	RootCmd.AddCommand(clusterlabCmd)
}

func manageClusterLab() {
	readClusterLabConfig()

	runOnBoot, err := strconv.ParseBool(myClusterLabConfig.Service["run_on_boot"])
	if err != nil {
		fmt.Println("Could not parse string to boolean value", err)
	}

	if runOnBoot {
		err = exec.Command("/usr/local/bin/cluster-lab", "start").Run()
		if err != nil {
			fmt.Println("Unable to start cluster-lab:", err)
		}
		fmt.Println("Running")
	}
}

func readClusterLabConfig() {
	err := config.UnmarshalKey("clusterlab", &myClusterLabConfig)
	if err != nil {
		fmt.Println("Could not unmarshal Cluster-Lab configuration")
	}
}
