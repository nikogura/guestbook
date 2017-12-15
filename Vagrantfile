# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"
default_box = "ubuntu/trusty64"

virtual_machines = [
    # db needs to come up first
    {
        :name => "db",
        :provision => [
            "installers/db.sh"
        ]
    },
    {
        :name => "app",
        :provision => [
            "installers/app.sh"
        ],
        :forwarded_ports => [
            host: 8080,
            guest: 8080,
        ],

    },
]

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  virtual_machines.each do |machine|
    config.vm.define machine[:name] do |box|
      box.vm.box = machine[:box] ? machine[:box] : default_box
      box.vm.hostname = machine[:name]

      box.vm.network :private_network, :auto_network => true
      box.vm.provision :hosts, :sync_hosts => true

      box.vm.provider 'virtualbox' do |p|
        p.customize ["modifyvm", :id, "--natdnshostresolver1", "on"]
        p.customize ["modifyvm", :id, "--natdnsproxy1", "on"]
      end

      if machine[:forwarded_ports]
        machine[:forwarded_ports].each do |pfwd|
          box.vm.network "forwarded_port", guest: pfwd[:guest], host: pfwd[:host]
        end
      end

      if machine[:files]
        machine[:files].each do |file|
          box.vm.provision "file", source: file[:src], destination: file[:dst]
        end
      end

      if machine[:provision]
        machine[:provision].each do |script|
          box.vm.provision "shell", path: script, privileged: true
        end
      end
    end
  end
end

