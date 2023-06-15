#!/bin/bash
urlencode() {
  local string="$1"
  local length="${#string}"
  local encoded_string=""

  for (( i = 0; i < length; i++ )); do
    local char="${string:i:1}"
    case $char in
      [a-zA-Z0-9.~_-])
        encoded_string+="$char" ;;
      *)
        encoded_string+=$(printf '%%%02X' "'$char") ;;
    esac
  done

  echo "$encoded_string"
}
if [ -z "$1" ]; then
  echo "Please provide a string as an argument."
  exit 1
fi
string="$1"
gateway=$(ip route | awk '/default/ {print $3}')
interface=$(ip route get $gateway | awk '/dev/ {print $3}')
mac_address=$(ip link show dev $interface | awk '/ether/ {gsub(/:/,""); print $2}')
ipLine=$(ifconfig | grep 192)
ipSuffix=${ipLine:16:2}
echo "IP address of $interface: $ipSuffix"
echo "MAC address of $interface: $mac_address"

echo "192.168.1.17:8080/$ipSuffix/$mac_address/Task%20Completed/$string"
