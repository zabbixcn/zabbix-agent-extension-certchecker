#zabbix-agent-extension-certchecker

zabbix-agent-extension-certchecker - this extension for monitoring SSL certifiicate file expire.

### Supported features

- Remaining time.
- DNS Names.
- Not Before time.
- Not After time.

### Monitoring valid only certificate and private key pair.

### Installation

#### Manual build

```sh
# Building
git clone https://github.com/zarplata/zabbix-agent-extension-certchecker.git
cd zabbix-agent-extension-certchecker
make

#Installing
make install

# By default, binary installs into /usr/bin/ and zabbix config in /etc/zabbix/zabbix_agentd.conf.d/ but,
# you may manually copy binary to your executable path and zabbix config to specific include directory
```

#### Arch Linux package
```sh
# Building
git clone https://github.com/zarplata/zabbix-agent-extension-certchecker.git
cd zabbix-agent-extension-certchecker
git checkout pkgbuild

makepkg

#Installing
pacman -U *.tar.xz
```

### Dependencies

zabbix-agent-extension-elasticsearch requires [zabbix-agent](http://www.zabbix.com/download) v2.4+ to run.

### Zabbix configuration
In order to start getting metrics, it is enough to import template and attach it to monitored node.

`WARNING:` You must define macro with name - `{$ZABBIX_SERVER_IP}` in global or local (template) scope with IP address of  zabbix server.
