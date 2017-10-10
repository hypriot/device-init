package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdEthernetInterfaceName string
var cmdAddress string
var cmdNetmask string
var cmdNetwork string
var cmdGateway string
var cmdBroadcast string

type InterfaceInfo struct {
	Address, Netmask, Network, Gateway, Broadcast string
}

type EthernetConfig struct {
	Interfaces map[string]InterfaceInfo
}

var myEthernetConfig EthernetConfig

// var networkInterfacesPath = "/etc/network/interfaces.d" // Already defined in wifi.go

// ethernetCmd represents the ethernet command
var ethernetCmd = &cobra.Command{
	Use:   "ethernet",
	Short: "set Ethernet settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(ethernetCmd)

	ethernetCmd.PersistentFlags().StringVarP(&cmdInterfaceName, "interface_name", "i", "", "Name of your network interface e.g. eth0")
}

func readEthernetConfig() {
	err := config.UnmarshalKey("ethernet", &myEthernetConfig)
	if err != nil {
		fmt.Println("Could not unmarshal Ethernet configuration")
	}
}
