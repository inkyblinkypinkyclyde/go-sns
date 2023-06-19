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
gateway=$(route -n get default | grep 'gateway' | awk '{print $2}')
interface=$(route -n get default | grep 'interface' | awk '{print $2}')
mac_address=$(ifconfig $interface | awk '/ether/ {gsub(/:/,""); print $2}')
ipSuffix=$(ifconfig | awk '/inet / {print $2}' | awk -F. '{print $NF}' | tail -n 1)
curl 192.168.1.27:8080/http/$ipSuffix/$mac_address/New%20Event%20From%20Bash%20Script/$url_encoded_string