#!/usr/bin/env bash
export GO111MODULE="auto" #"" #go env -w  GO111MODULE=auto
cd /Users/daparici/code/github.com/davidaparicio/ambari-to-opsgenie
/Users/daparici/homebrew/opt/go/libexec/bin/go run /Users/daparici/code/github.com/davidaparicio/ambari-to-opsgenie/cmd/xbar/main.go

#  <xbar.title>Check Ambari</xbar.title>
#  <xbar.version>v1.0</xbar.version>
#  <xbar.author>David Aparicio</xbar.author>
#  <xbar.author.github>davidaparicio</xbar.author.github>
#  <xbar.desc>The script tries to get CRITICAL and WARNING alert numbers from Ambari Dashboard, through REST API service.</xbar.desc>
#  <xbar.image>https://raw.githubusercontent.com/christophschlosser/bitbar-plugins/checkhosts/Network/checkhosts.png</xbar.image>
#  <xbar.dependencies>go</xbar.dependencies>
#  <xbar.abouturl>https://davidaparicio.gitlab.io</xbar.abouturl>