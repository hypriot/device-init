# device-init
[![Build Status](https://travis-ci.org/hypriot/device-init.svg?branch=master)](https://travis-ci.org/hypriot/device-init)

The `device-init` runs while booting your device and can customize some settings of the device.

To replace the `/boot/occidentalis.txt` (as we want this feature for more boards than the RPi) and to enhance the possibilities to customize the SD card image at first boot, we want to introduce this tool `device-init` that read customizations from the FAT partition `/boot/device-init.yaml` so it is easier to write it onto the SD card image.

* ☑ Set hostname
* ☑ Set WiFi SSID/PSK
* ☐ Set timezone
* ☐ Set locale
* ☑ Set static IP
* ☑ Configure network interfaces
* ☑ Pull some docker images
* ☐ Install a list of DEB packages
* ☑ Run a custom script from /boot
* ...

## The /boot/device-init.yaml

The `device-init` tool reads the file `/boot/device-init.yaml` to initialize several settings while booting your device.

### Hostname

```yaml
hostname: "black-pearl"
```

### Network Settings

```yaml
network:
  interfaces:
    eth0:
      # sets a static IP
      address: 192.168.13.37
      netmask: 255.255.255.0
      gateway: 192.168.13.1
      dnsnameservers:
        - 8.8.8.8
        - 8.8.4.4
      dnssearch:
        example.com
    eth1:
      # uses dhcp
    wlan0:
      ssid: "MyNetwork"
      password: "secret_password"
```

### Wifi Settings

This option only allows creating wlan intefaces and setting SSID/PSK. For advanced configuration see [network configuration](#network-settings).

```yaml
wifi:
  interfaces:
    wlan0:
      ssid: "MyNetwork"
      password: "secret_password"
```

### Docker Preload Images Settings
device-init can preload local image files into the Docker engine on boot.
Those images have to be exported via 'docker save image-name > image-name.tar'.
It is recommended to compress the output of 'docker save' with 'gzip image-name.tar' which results in a image-name.tar.gz file.

```yaml
docker:
  images:
    - "/path/to/image-name.tar.gz"
    - "/path/to/another-image-name.tar"
```

### Run a Command
device-init can execute a list of commands on boot. e.g.:
```yaml
runcmd:
  - "apt-get install package-name"
  - "curl https://raw.githubusercontent.com/.../myscript.sh | sh"
```
Make sure that the command lines do not produce YAML syntax errors. You can check [here](http://www.yamllint.com/)  

### Hypriot Cluster-Lab
device-init can start the Hypriot Cluster-Lab on start up by setting the 'run_on_boot' option to 'true'.

```
clusterlab:
  service:
    run_on_boot: "false"
```


## Buy us a beer!

This FLOSS software is funded by donations only. Please support us to maintain and further improve it!

<a href="https://liberapay.com/Hypriot/donate"><img alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg"></a>

