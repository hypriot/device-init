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
        device_init_cmd_result = command('device-init').stdout
        new_hostname_cmd_result = command('hostname').stdout

        expect(old_hostname_cmd_result).to contain("device-tester")
        expect(device_init_cmd_result).to contain("Set")
        expect(new_hostname_cmd_result).to contain("black-pearl")
      end
    end
  end

end
