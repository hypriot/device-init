# device-init
[![Build Status](https://travis-ci.org/hypriot/device-init.svg?branch=master)](https://travis-ci.org/hypriot/device-init)

**Disclaimer:** This is still work in progress and not suited for production use. We are currently adding `device-init` to the image-builder-rpi repo. Please give us feedback and file issues if you have some wishes, enhancements or bugs found. Pull requests are also welcome!

The `device-init` runs while booting your device and can customize some settings of the device.

To replace the `/boot/occidentalis.txt` (as we want this feature for more boards than the RPi) and to enhance the possibilities to customize the SD card image at first boot, we want to introduce this tool `device-init` that read customizations from the FAT partition `/boot/device-init.yaml` so it is easier to write it onto the SD card image.

* ☑ Set hostname
* ☐ Set WiFi SSID/PSK
* ☐ Set timezone
* ☐ Set locale
* ☐ Set static IP (see hypriot/flash#25)
* ☐ Pull some docker images
* ☐ Install a list of DEB packages
* ☐ Run a custom script from /boot
* ...

## The /boot/device-init.yaml

The `device-init` tool reads the file `/boot/device-init.yaml` to initialize several settings while booting your device.

### hostname

```yaml
hostname: "black-pearl"
```
