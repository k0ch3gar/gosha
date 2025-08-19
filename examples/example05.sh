#!/usr/bin/gosha

var fileName = "info.log"

if ! -f fileName {
  $(touch $fileName)
}

var logPath = "/var/log/anaconda/syslog"
if ! -f logPath {
  logPath = "/var/log/installer/syslog"
}

var lines = $(sudo cat $logPath | awk '$5 == "INFO:" { print }')

return lines