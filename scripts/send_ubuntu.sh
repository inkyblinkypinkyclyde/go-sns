#!/bin/bash

if [ -z "$1" ]; then
  echo "Please provide a string as an argument."
  exit 1
fi

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

string="$1"
url_encoded_string=$(urlencode "$string")
gateway=$(ip route | awk '/default/ {print $3}')
interface=$(ip route get $gateway | awk '/dev/ {print $3}')
mac_address=$(ip link show dev $interface | awk '/ether/ {gsub(/:/,""); print $2}')
ip_address=$(hostname -I | awk '{print $1}' | awk -F. '{print $NF}')
curl 192.168.1.17:8080/http/$ip_address/$mac_address/Task%20Completed/$url_encoded_string