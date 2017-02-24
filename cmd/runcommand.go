package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

// runCmd represents the run command to execute shell commands
var runCmd = &cobra.Command{
	Use:   "runcommand",
	Short: "Run shell commands",
	Long:  `Run shell commands that are defined in device-init.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		runCommands()
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func runCommands() {
	commandList := config.GetStringSlice("runcmd")
	for _, command := range commandList {
		err := exec.Command("sh", "-c", command).Run()
		if err != nil {
			fmt.Printf("Unable to run command %v. Reason: %v\n", command, err)
		}
	}
}
