#!/bin/bash
#
# ZEUS Error Dump
# Timestamp: [Sat Feb 6 19:55:07 2021]
# Error: exit status 2
# StdErr: 
# # github.com/dreadl0ck/maltego
# ./utils.go:202:10: undefined: netcapPrefix
# ./utils.go:360:5: undefined: netcapMachinePrefix
# # github.com/dreadl0ck/maltego [github.com/dreadl0ck/maltego.test]
# ./utils.go:202:10: undefined: netcapPrefix
# ./utils.go:360:5: undefined: netcapMachinePrefix
# 


#!/bin/bash
UBUNTU_IMAGE="dreadl0ck/trx-ubuntu"
ALPINE_IMAGE="dreadl0ck/trx-alpine"
VERSION="1.0"



go test ./...