package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setEthernetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set Ethernet settings",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		setEthernet()
	},
}

func setEthernet() {
	var err error

	readEthernetConfig()

	const interfaceTemplate = `allow-hotplug {{.Name}}

auto {{.Name}}
iface {{.Name}} inet static
  address {{.Address}}
  netmask {{.Netmask}}
  network {{.Network}}
  gateway {{.Gateway}}
  broadcast {{.Broadcast}}
`

	// if we have command line parameters only add those to our ethernet configuration
	if cmdAddress != "" &&
		cmdNetmask != "" &&
		cmdNetwork != "" &&
		cmdGateway != "" &&
		cmdBroadcast != "" {
		for key := range myEthernetConfig.Interfaces {
			delete(myEthernetConfig.Interfaces, key)
		}
		myEthernetConfig.Interfaces = make(map[string]InterfaceInfo)
		myEthernetConfig.Interfaces[cmdInterfaceName] = InterfaceInfo{
			Address:   cmdAddress,
			Netmask:   cmdNetmask,
			Network:   cmdNetwork,
			Gateway:   cmdGateway,
			Broadcast: cmdBroadcast,
		}
	}

	for interfaceName, interfaceInfo := range myEthernetConfig.Interfaces {

		type templateVariables struct {
			Name, Address, Netmask, Network, Gateway, Broadcast string
		}

		variables := templateVariables{
			Name:      interfaceName,
			Address:   interfaceInfo.Address,
			Netmask:   interfaceInfo.Netmask,
			Network:   interfaceInfo.Network,
			Gateway:   interfaceInfo.Gateway,
			Broadcast: interfaceInfo.Broadcast,
		}

		err = os.MkdirAll("/etc/network/interfaces.d/", 0755)
		if err != nil {
			fmt.Println("Could not create path: ", "/etc/network/interfaces.d/")
		}

		fmt.Println("TEST")
		configFilePath := path.Join("/etc/network/interfaces.d/", interfaceName)
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

		t := template.Must(template.New("config").Parse(interfaceTemplate))

		err = t.Execute(f, variables)
		if err != nil {
			fmt.Println("Error writing configuration:", err)
		}

		if interfaceExistsAndIsDown(interfaceName) {
			output, err := exec.Command("/sbin/ifup", interfaceName).CombinedOutput()
			if err != nil {
				message := fmt.Sprintf("Could not bring up interface %s: %s", interfaceName, err)
				fmt.Println(message)
			}
			fmt.Println(string(output)[:])
		}

		// try to bring the interface up once more but bring it down before
		if interfaceExistsAndIsDown(interfaceName) {
			output, err := exec.Command("/sbin/ifdown", interfaceName).CombinedOutput()
			if err != nil {
				message := fmt.Sprintf("Could not bring the interface down %s: %s ", interfaceName, err)
				fmt.Println(message)
			}
			fmt.Println(string(output)[:])
			output, err = exec.Command("/sbin/ifup", interfaceName).CombinedOutput()
			if err != nil {
				message := fmt.Sprintf("Could still not bring up interface %s: %s", interfaceName, err)
				fmt.Println(message)
			}
		}
	}
}

func init() {
	ethernetCmd.AddCommand(setEthernetCmd)
	setEthernetCmd.Flags().StringVarP(&cmdAddress, "address", "a", "", "The desired address")
	setEthernetCmd.Flags().StringVarP(&cmdNetmask, "netmask", "m", "", "The netmask of the network")
	setEthernetCmd.Flags().StringVarP(&cmdNetwork, "network", "n", "", "The network direction")
	setEthernetCmd.Flags().StringVarP(&cmdGateway, "gateway", "g", "", "The default gateway")
	setEthernetCmd.Flags().StringVarP(&cmdBroadcast, "broadcast", "b", "", "The broadcast address")
}
