# Maintainer: Andrey Kitsul <a.kitsul@zarplata.ru>

pkgname=zabbix-agent-extension-certchecker
pkgver=20170928.14_84a37b7
pkgrel=1
pkgdesc="Extension for zabbix-agentd for monitoring SSL certificate file"
arch=('any')
license=('GPL')
makedepends=('go')
depends=('zabbix-agent')
install="install.sh"
source=("git+https://github.com/zarplata/$pkgname.git#branch=master")
md5sums=(
    'SKIP'
    )

make_zabbix_config() {
    userparam_string_discovery="UserParameter=certificate.discovery[*], /usr/bin/zabbix-agent-extension-certchecker --discovery --path \$1"
    userparam_string="UserParameter=certificate.stats[*], /usr/bin/zabbix-agent-extension-certchecker -z \$1 -c \$2 -k \$3"

    echo "$userparam_string_discovery" > "$pkgname.conf"
    echo "$userparam_string" >> "$pkgname.conf"
}


pkgver() {
	cd "$pkgname"
    make ver
}

build() {
    make_zabbix_config

    cd "$pkgname"
    make 
}

package() {
	cd "$srcdir/$pkgname"
    ZBX_INC_DIR=/etc/zabbix/zabbix_agentd.conf.d/

    install -Dm 0755 .out/"${pkgname}" "${pkgdir}/usr/bin/${pkgname}"
    install -Dm 0644 ../"${pkgname}.conf" "${pkgdir}${ZBX_INC_DIR}${pkgname}.conf"
    
}
