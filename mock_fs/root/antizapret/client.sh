#!/bin/bash
#
# MOCK SCRIPT for adding/removing clients.
# This script is a simplified version of the original client.sh
# It does not generate any real keys or configs.
# It only creates and deletes dummy files to simulate client management.
#
# To be executed from the project root directory.
# Example: ./mock_fs/root/antizapret/client.sh 1 my-test-client
#
set -e

# The root of the mock file system.
MOCK_ROOT="mock_fs"

# Helper function to resolve paths
abs_path() {
    echo "$MOCK_ROOT/$1"
}

ROOT_PATH=$(abs_path "root/antizapret")
ETC_PATH=$(abs_path "etc")


handle_error() {
	echo "Error at line $1: $2"
	exit 1
}
trap 'handle_error $LINENO "$BASH_COMMAND"' ERR

export LC_ALL=C

askClientName(){
	if ! [[ "$CLIENT_NAME" =~ ^[a-zA-Z0-9_-]{1,32}$ ]]; then
		echo
		echo 'Enter client name: 1–32 alphanumeric characters (a-z, A-Z, 0-9) with underscore (_) or dash (-)'
		until [[ "$CLIENT_NAME" =~ ^[a-zA-Z0-9_-]{1,32}$ ]]; do
			read -rp 'Client name: ' -e CLIENT_NAME
		done
	fi
}

# This function is not used in the mock script for client cert expiration
# but kept for compatibility with the original script's call signature.
askClientCertExpire(){
    :
}

setServerHost_FileName(){
	if [[ -z "$1" ]]; then
		SERVER_HOST="$SERVER_IP"
	else
		SERVER_HOST="$1"
	fi

	FILE_NAME="${CLIENT_NAME#antizapret-}"
	FILE_NAME="${FILE_NAME#vpn-}"
	FILE_NAME="${FILE_NAME}-(${SERVER_HOST})"
}

setServerIP(){
    # Mocked IP
	SERVER_IP="1.2.3.4"
}

# No-op render function
render() {
    cat "$1"
}

# No-op init function
initOpenVPN(){
	:
}

addOpenVPN(){
	setServerHost_FileName "$OPENVPN_HOST"
	
    # Create a dummy file to be listed by listOpenVPN
    touch "$ETC_PATH/openvpn/easyrsa3/pki/issued/$CLIENT_NAME.crt"

    # Create dummy config files
	echo "Mock OpenVPN config for ${CLIENT_NAME}" > "$ROOT_PATH/client/openvpn/antizapret/antizapret-$FILE_NAME.ovpn"
	echo "Mock OpenVPN config for ${CLIENT_NAME}" > "$ROOT_PATH/client/openvpn/antizapret-udp/antizapret-$FILE_NAME-udp.ovpn"
	echo "Mock OpenVPN config for ${CLIENT_NAME}" > "$ROOT_PATH/client/openvpn/antizapret-tcp/antizapret-$FILE_NAME-tcp.ovpn"
	echo "Mock OpenVPN config for ${CLIENT_NAME}" > "$ROOT_PATH/client/openvpn/vpn/vpn-$FILE_NAME.ovpn"
	echo "Mock OpenVPN config for ${CLIENT_NAME}" > "$ROOT_PATH/client/openvpn/vpn-udp/vpn-$FILE_NAME-udp.ovpn"
	echo "Mock OpenVPN config for ${CLIENT_NAME}" > "$ROOT_PATH/client/openvpn/vpn-tcp/vpn-$FILE_NAME-tcp.ovpn"

	echo "Mock OpenVPN profile files created for client '$CLIENT_NAME'"
}

deleteOpenVPN(){
	setServerHost_FileName "$OPENVPN_HOST"
	echo
	
    # Delete the dummy file
    rm -f "$ETC_PATH/openvpn/easyrsa3/pki/issued/$CLIENT_NAME.crt"

    # Delete dummy config files
	rm -f $ROOT_PATH/client/openvpn/antizapret/antizapret-$FILE_NAME.ovpn
	rm -f $ROOT_PATH/client/openvpn/antizapret-udp/antizapret-$FILE_NAME-udp.ovpn
	rm -f $ROOT_PATH/client/openvpn/antizapret-tcp/antizapret-$FILE_NAME-tcp.ovpn
	rm -f $ROOT_PATH/client/openvpn/vpn/vpn-$FILE_NAME.ovpn
	rm -f $ROOT_PATH/client/openvpn/vpn-udp/vpn-$FILE_NAME-udp.ovpn
	rm -f $ROOT_PATH/client/openvpn/vpn-tcp/vpn-$FILE_NAME-tcp.ovpn

	echo "Mock OpenVPN client '$CLIENT_NAME' successfully deleted"
}

listOpenVPN(){
	[[ -n "$CLIENT_NAME" ]] && return
	echo
	echo 'OpenVPN client names:'
	ls "$ETC_PATH/openvpn/easyrsa3/pki/issued" | sed 's/\.crt$//' | grep -v "^antizapret-server$" | sort
}

# No-op init function
initWireGuard(){
	:
}

addWireGuard(){
	setServerHost_FileName "$WIREGUARD_HOST"
	echo

    # Add a dummy client entry for listWireGuard to work
    echo "# Client = ${CLIENT_NAME}" >> "$ETC_PATH/wireguard/antizapret.conf"
    echo "# Client = ${CLIENT_NAME}" >> "$ETC_PATH/wireguard/vpn.conf"

	echo "Mock WG config for ${CLIENT_NAME}" > "$ROOT_PATH/client/wireguard/antizapret/antizapret-$FILE_NAME-wg.conf"
	echo "Mock AWG config for ${CLIENT_NAME}" > "$ROOT_PATH/client/amneziawg/antizapret/antizapret-$FILE_NAME-am.conf"
	echo "Mock WG config for ${CLIENT_NAME}" > "$ROOT_PATH/client/wireguard/vpn/vpn-$FILE_NAME-wg.conf"
	echo "Mock AWG config for ${CLIENT_NAME}" > "$ROOT_PATH/client/amneziawg/vpn/vpn-$FILE_NAME-am.conf"

	echo "Mock WireGuard/AmneziaWG profile files created for client '$CLIENT_NAME'"
}

deleteWireGuard(){
	setServerHost_FileName "$WIREGUARD_HOST"
	echo

	if ! grep -q "# Client = ${CLIENT_NAME}" "$ETC_PATH/wireguard/antizapret.conf" && ! grep -q "# Client = ${CLIENT_NAME}" "$ETC_PATH/wireguard/vpn.conf"; then
		echo "Failed to delete client '$CLIENT_NAME'! Please check if client exists"
		exit 6
	fi

	sed -i.bak "/^# Client = ${CLIENT_NAME}$/d" "$ETC_PATH/wireguard/antizapret.conf"
	sed -i.bak "/^# Client = ${CLIENT_NAME}$/d" "$ETC_PATH/wireguard/vpn.conf"
    rm -f "$ETC_PATH/wireguard/antizapret.conf.bak"
    rm -f "$ETC_PATH/wireguard/vpn.conf.bak"

	rm -f $ROOT_PATH/client/{wireguard,amneziawg}/antizapret/antizapret-$FILE_NAME-*.conf
	rm -f $ROOT_PATH/client/{wireguard,amneziawg}/vpn/vpn-$FILE_NAME-*.conf

	echo "Mock WireGuard/AmneziaWG client '$CLIENT_NAME' successfully deleted"
}

listWireGuard(){
	[[ -n "$CLIENT_NAME" ]] && return
	echo
	echo 'WireGuard/AmneziaWG client names:'
	cat "$ETC_PATH/wireguard/antizapret.conf" "$ETC_PATH/wireguard/vpn.conf" | grep -E "^# Client" | cut -d '=' -f 2 | sed 's/ //g' | sort -u
}

recreate(){
	echo
    
    # This is a mock implementation. We just list existing clients and "recreate" their files.
    
    echo "Recreating OpenVPN profiles..."
    listOpenVPN | while read -r client; do
        CLIENT_NAME=$client addOpenVPN
    done

    echo "Recreating WireGuard/AmneziaWG profiles..."
    listWireGuard | while read -r client; do
        CLIENT_NAME=$client addWireGuard
    done
}

backup(){
	echo
    echo "This is a mock backup. No files will be created."
}

source "$ROOT_PATH/setup"
umask 022
setServerIP

OPTION=$1
CLIENT_NAME=$2
CLIENT_CERT_EXPIRE=$3

if ! [[ "$OPTION" =~ ^[1-8]$ ]]; then
	echo
	echo 'Please choose option:'
	echo '    1) OpenVPN - Add client'
	echo '    2) OpenVPN - Delete client'
	echo '    3) OpenVPN - List clients'
	echo '    4) WireGuard/AmneziaWG - Add client'
	echo '    5) WireGuard/AmneziaWG - Delete client'
	echo '    6) WireGuard/AmneziaWG - List clients'
	echo '    7) (Re)create client profile files'
	echo '    8) Backup configuration and clients'
	until [[ "$OPTION" =~ ^[1-8]$ ]]; do
		read -rp 'Option choice [1-8]: ' -e OPTION
	done
fi

case "$OPTION" in
	1)
		echo "OpenVPN - Add client $CLIENT_NAME"
		askClientName
		addOpenVPN
		;;
	2)
		echo "OpenVPN - Delete client $CLIENT_NAME"
		listOpenVPN
		askClientName
		deleteOpenVPN
		;;
	3)
		echo 'OpenVPN - List clients'
		listOpenVPN
		;;
	4)
		echo "WireGuard/AmneziaWG - Add client $CLIENT_NAME"
		askClientName
		addWireGuard
		;;
	5)
		echo "WireGuard/AmneziaWG - Delete client $CLIENT_NAME"
		listWireGuard
		askClientName
		deleteWireGuard
		;;
	6)
		echo 'WireGuard/AmneziaWG - List clients'
		listWireGuard
		;;
	7)
		echo '(Re)create client profile files'
		recreate
		;;
	8)
		echo 'Backup configuration and clients'
		backup
		;;
esac
exit 0
