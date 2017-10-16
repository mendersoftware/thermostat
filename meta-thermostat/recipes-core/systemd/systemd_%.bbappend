FILESEXTRAPATHS_prepend := "${THISDIR}/files:"

PACKAGECONFIG_append = " networkd"

SRC_URI += "file://eth.network"

FILES_${PN} += "{sysconfdir}/systemd/network/"

do_install_append() {
    install -d ${D}${sysconfdir}/systemd/network/
    install -m 0644 ${WORKDIR}/eth.network ${D}${prefix}/lib/systemd/network/ 
}

