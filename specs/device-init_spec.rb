require 'serverspec'
set :backend, :exec

describe "device-init" do
  context "binary" do
    it "exists" do
      file = file('/usr/local/bin/device-init')

      expect(file).to exist
    end

    it "is x86-64" do
      result = command('file /usr/local/bin/device-init').stdout

      expect(result).to contain("x86-64")
    end
  end

  context "hostname" do
    before(:each) do
      system('device-init hostname set device-tester >> /dev/null')
      system('rm -f /boot/device-init.yaml')
    end

    context "with command hostname" do
      it "shows hostname" do
        hostname_cmd_result = command('hostname').stdout
        device_init_cmd_result = command('device-init hostname show').stdout

        expect(device_init_cmd_result).to contain(hostname_cmd_result)
      end

      it "sets hostname" do
        old_hostname_cmd_result = command('hostname').stdout
        device_init_cmd_result = command('device-init hostname set black-pearl').stdout
        new_hostname_cmd_result = command('hostname').stdout

        expect(old_hostname_cmd_result).to contain("device-tester")
        expect(device_init_cmd_result).to contain("Set")
        expect(new_hostname_cmd_result).to contain("black-pearl")
      end
    end

    context "with config-file" do
      it "sets hostname" do
        status = command(%q(echo -n 'hostname: "black-pearl"\n' > /boot/device-init.yaml)).exit_status
        expect(status).to be(0)

        old_hostname_cmd_result = command('hostname').stdout
        device_init_cmd_result = command('device-init --config').stdout
        new_hostname_cmd_result = command('hostname').stdout

        expect(old_hostname_cmd_result).to contain("device-tester")
        expect(device_init_cmd_result).to contain("Set")
        expect(new_hostname_cmd_result).to contain("black-pearl")
      end
    end
  end

  context "wifi" do
    context "with command wifi" do
      before(:each) do
        expect(command('rm -Rf /etc/network/interfaces.d').exit_status).to be(0)
      end

      it "shows config" do
        network_interface_dir = file('/etc/network/interfaces.d/')
        expect(network_interface_dir.exists?).to be(false)

        wlan0 = "allow-hotplug wlan1\n\nauto wlan0\niface wlan1 inet dhcp\nwpa-ssid MySecondNetwork\nwpa-psk 32919a0369631b758391d00e2aaaf66e6ab61b61949cc853c45410fbf4910442"
        status = command(%Q(mkdir -p /etc/network/interfaces.d && echo -n '#{wlan0}' > /etc/network/interfaces.d/wlan0)).exit_status
        expect(status).to be(0)

        device_init_cmd_result = command('device-init wifi show -i wlan0')
        expect(device_init_cmd_result.exit_status).to be(0)
        expect(device_init_cmd_result.stdout).to contain('wlan0 in /etc/network/interfaces.d/wlan0')
        expect(device_init_cmd_result.stdout).to contain('------------------------------------------')
        expect(device_init_cmd_result.stdout).to contain('wpa-psk 32919a0369631b758391d00e2aaaf66e6ab61b61949cc853c45410fbf4910442')
      end

      it "fails for non-existent interface" do
        network_interface_dir = file('/etc/network/interfaces.d/')
        expect(network_interface_dir.exists?).to be(false)

        device_init_cmd_result = command('device-init wifi show -i wlan3')
        expect(device_init_cmd_result.exit_status).to be(0)
        expect(device_init_cmd_result.stdout).to contain("Could not open wifi configuration for interface 'wlan3'")
      end

      it "sets config" do
        network_interface_dir = file('/etc/network/interfaces.d/')
        expect(network_interface_dir.exists?).to be(false)

        device_init_cmd_result = command('device-init wifi set -i wlan0 -s MyNetwork -p my_secret_password')
        expect(device_init_cmd_result.exit_status).to be(0)

        network_interface_dir = file('/etc/network/interfaces.d/')
        expect(network_interface_dir.exists?).to be(true)
        expect(network_interface_dir.directory?).to be(true)

        wlan_config_file = file('/etc/network/interfaces.d/wlan0')
        expect(wlan_config_file.exists?).to be(true)
        expect(wlan_config_file).to contain('allow-hotplug wlan0')
        expect(wlan_config_file).to contain('auto wlan0')
        expect(wlan_config_file).to contain('iface wlan0 inet dhcp')
        expect(wlan_config_file).to contain('wpa-ssid MyNetwork')
        expect(wlan_config_file).to contain('wpa-psk a53576661368249ebfa26f8828669ad0e6f0523154b55404b33a21ca1242b845')
      end
    end

    context "with config-file" do
      let(:config_one_wifi_interface)  { File.read(File.join(File.dirname(__FILE__), 'testdata', 'one_wifi_interface.yaml')) }
      let(:config_two_wifi_interfaces) { File.read(File.join(File.dirname(__FILE__), 'testdata', 'two_wifi_interfaces.yaml')) }
      let(:config_no_wifi)             { File.read(File.join(File.dirname(__FILE__), 'testdata', 'no_wifi_interface.yaml')) }

      context "sets config" do
        before(:each) do
          expect(command('rm -Rf /etc/network/interfaces.d').exit_status).to be(0)
        end

        it "creates configuration for one interface entry" do
          status = command(%Q(echo -n '#{config_one_wifi_interface}' > /boot/device-init.yaml)).exit_status
          expect(status).to be(0)

          network_interface_dir = file('/etc/network/interfaces.d/')
          expect(network_interface_dir.exists?).to be(false)

          device_init_cmd_result = command('device-init --config')
          expect(device_init_cmd_result.exit_status).to be(0)

          network_interface_dir = file('/etc/network/interfaces.d/')
          expect(network_interface_dir.exists?).to be(true)
          expect(network_interface_dir.directory?).to be(true)

          wlan_config_file = file('/etc/network/interfaces.d/wlan0')
          expect(wlan_config_file.exists?).to be(true)
          expect(wlan_config_file).to contain('allow-hotplug wlan0')
          expect(wlan_config_file).to contain('auto wlan0')
          expect(wlan_config_file).to contain('iface wlan0 inet dhcp')
          expect(wlan_config_file).to contain('wpa-ssid MyNetwork')
          expect(wlan_config_file).to contain('wpa-psk a53576661368249ebfa26f8828669ad0e6f0523154b55404b33a21ca1242b845')
        end

        it "creates configuration for two interface entries" do
          status = command(%Q(echo -n '#{config_two_wifi_interfaces}' > /boot/device-init.yaml)).exit_status
          expect(status).to be(0)

          network_interface_dir = file('/etc/network/interfaces.d/')
          expect(network_interface_dir.exists?).to be(false)

          device_init_cmd_result = command('device-init --config')
          expect(device_init_cmd_result.exit_status).to be(0)

          network_interface_dir = file('/etc/network/interfaces.d/')
          expect(network_interface_dir.exists?).to be(true)
          expect(network_interface_dir.directory?).to be(true)

          wlan_config_file = file('/etc/network/interfaces.d/wlan1')
          expect(wlan_config_file.exists?).to be(true)
          expect(wlan_config_file).to contain('allow-hotplug wlan1')
          expect(wlan_config_file).to contain('auto wlan1')
          expect(wlan_config_file).to contain('iface wlan1 inet dhcp')
          expect(wlan_config_file).to contain('wpa-ssid MySecondNetwork')
          expect(wlan_config_file).to contain('wpa-psk 32919a0369631b758391d00e2aaaf66e6ab61b61949cc853c45410fbf4910442')
        end

        it "creates no configuration if there is no 'wifi' key in device-init.yaml" do
          status = command(%Q(echo -n '#{config_no_wifi}' > /boot/device-init.yaml)).exit_status
          expect(status).to be(0)

          network_interface_dir = file('/etc/network/interfaces.d/')
          expect(network_interface_dir.exists?).to be(false)

          device_init_cmd_result = command('device-init --config')
          expect(device_init_cmd_result.exit_status).to be(0)

          network_interface_dir = file('/etc/network/interfaces.d/')
          expect(network_interface_dir.exists?).to be(false)
        end
      end

      context "show config" do
        it "for interfaces in device-init.yaml" do
          status = command(%Q(echo -n '#{config_one_wifi_interface}' > /boot/device-init.yaml)).exit_status
          expect(status).to be(0)

          wlan0 = "allow-hotplug wlan1\n\nauto wlan0\niface wlan1 inet dhcp\nwpa-ssid MySecondNetwork\nwpa-psk 32919a0369631b758391d00e2aaaf66e6ab61b61949cc853c45410fbf4910442"
          status = command(%Q(mkdir -p /etc/network/interfaces.d && echo -n '#{wlan0}' > /etc/network/interfaces.d/wlan0)).exit_status
          expect(status).to be(0)

          device_init_cmd_result = command('device-init wifi show -c')
          expect(device_init_cmd_result.exit_status).to be(0)
          expect(device_init_cmd_result.stdout).to contain('wlan0 in /etc/network/interfaces.d/wlan0')
          expect(device_init_cmd_result.stdout).to contain('------------------------------------------')
          expect(device_init_cmd_result.stdout).to contain('wpa-psk 32919a0369631b758391d00e2aaaf66e6ab61b61949cc853c45410fbf4910442')
        end
      end
    end
  end

  context "docker" do
    context "preload-images" do
      let(:preload_docker_images_tar_gz)  { File.read(File.join(File.dirname(__FILE__), 'testdata', 'preload_docker_images_tar_gz.yaml')) }
      let(:preload_docker_images_tar)  { File.read(File.join(File.dirname(__FILE__), 'testdata', 'preload_docker_images_tar.yaml')) }
      let(:preload_docker_images_non_existant)  { File.read(File.join(File.dirname(__FILE__), 'testdata', 'preload_docker_images_non_existant.yaml')) }

      before(:each) do
        cmd_search_image = command('docker images | grep -q busybox')
        if cmd_search_image.exit_status == 0
          rmi_image_cmd = command('docker rmi -f busybox')
          expect(rmi_image_cmd.exit_status).to be(0)
        end

        echo_config_cmd = command(%Q(rm -f /boot/device-init.yaml))
        expect(echo_config_cmd.exit_status).to be(0)
      end

      it "preloads local tar file images" do
        echo_config_cmd = command(%Q(echo -n '#{preload_docker_images_tar}' > /boot/device-init.yaml))
        expect(echo_config_cmd.exit_status).to be(0)

        device_init_cmd = command('device-init --config')
        expect(device_init_cmd.exit_status).to be(0)
        expect(device_init_cmd.stdout).to contain('Imported image: /specs/testdata/busybox.tar')

        docker_images_cmd = command('docker images')
        expect(docker_images_cmd.stdout).to contain('busybox')
      end

      it "preloads local tar.gz file images" do
        echo_config_cmd = command(%Q(echo -n '#{preload_docker_images_tar_gz}' > /boot/device-init.yaml))
        expect(echo_config_cmd.exit_status).to be(0)

        device_init_cmd = command('device-init --config')
        expect(device_init_cmd.exit_status).to be(0)
        expect(device_init_cmd.stdout).to contain('Imported image: /specs/testdata/busybox.tar.gz')

        docker_images_cmd = command('docker images')
        expect(docker_images_cmd.stdout).to contain('busybox')
      end

      it "doesn't preload images that were already imported" do
        echo_config_cmd = command(%Q(echo -n '#{preload_docker_images_tar_gz}' > /boot/device-init.yaml))
        expect(echo_config_cmd.exit_status).to be(0)

        mkdir_cmd = command(%Q(mkdir -p /var/log/device-init))
        expect(mkdir_cmd.exit_status).to be(0)

        echo_logfile_cmd = command(%Q(echo -n '/specs/testdata/busybox.tar.gz' > /var/log/device-init/preloaded_images.log))
        expect(echo_logfile_cmd.exit_status).to be(0)

        device_init_cmd = command('device-init --config')
        expect(device_init_cmd.exit_status).to be(0)

        expect(device_init_cmd.stdout).to contain('Already imported image: /specs/testdata/busybox.tar.gz')

        docker_images_cmd = command('docker images')
        expect(docker_images_cmd.stdout).to_not contain('busybox')
      end

      it "doesn't preload images that do not exist" do
        echo_config_cmd = command(%Q(echo -n '#{preload_docker_images_non_existant}' > /boot/device-init.yaml))
        expect(echo_config_cmd.exit_status).to be(0)

        device_init_cmd = command('device-init --config')
        expect(device_init_cmd.exit_status).to be(0)
        expect(device_init_cmd.stdout).to contain('Image file does not exist: /var/not/exists/busybox.tar.gz')
      end

      it "doesn't choke on a config file without docker key" do
        echo_config_cmd = command(%Q(echo -n '' > /boot/device-init.yaml))
        expect(echo_config_cmd.exit_status).to be(0)

        device_init_cmd = command('device-init --config')
        expect(device_init_cmd.exit_status).to be(0)
      end
    end
  end

  context "cluster-lab" do
    let(:cluster_lab_enabled)  { File.read(File.join(File.dirname(__FILE__), 'testdata', 'cluster_lab_enabled.yaml')) }
    let(:cluster_lab_disabled)  { File.read(File.join(File.dirname(__FILE__), 'testdata', 'cluster_lab_disabled.yaml')) }

    it "calls systemctl to enable cluster-lab service" do
      echo_config_cmd = command(%Q(echo -n '#{cluster_lab_enabled}' > /boot/device-init.yaml))
      expect(echo_config_cmd.exit_status).to be(0)

      device_init_cmd = command('device-init --config')
      expect(device_init_cmd.exit_status).to be(0)

      expect(device_init_cmd.stdout).to contain('Unable to enable cluster-lab')
    end

    it "calls systemctl to disable cluster-lab service" do
      echo_config_cmd = command(%Q(echo -n '#{cluster_lab_disabled}' > /boot/device-init.yaml))
      expect(echo_config_cmd.exit_status).to be(0)

      device_init_cmd = command('device-init --config')
      expect(device_init_cmd.exit_status).to be(0)

      expect(device_init_cmd.stdout).to contain('Unable to disable cluster-lab')
    end
  end
end

