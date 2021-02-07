#!/bin/bash
#
# ZEUS Error Dump
# Timestamp: [Sun Feb 7 17:13:15 2021]
# Error: exit status 2
# StdErr: 
# # github.com/dreadl0ck/maltego/cmd/maltego-gen
# cmd/maltego-gen/main.go:71:14: cannot use maltego.NewStringField(f.Name, f.Description) (type *maltego.PropertyField) as type maltego.PropertyField in assignment
# cmd/maltego-gen/main.go:73:14: cannot use maltego.NewRequiredStringField(f.Name, f.Description) (type *maltego.PropertyField) as type maltego.PropertyField in assignment
# 


#!/bin/bash
UBUNTU_IMAGE="dreadl0ck/trx-ubuntu"
ALPINE_IMAGE="dreadl0ck/trx-alpine"
VERSION="1.0"



go install github.com/dreadl0ck/maltego/cmd/maltego-gen