#!/bin/bash
echo "Getting latest version..."

#echo "Stopping service"
systemctl stop watgbridge

#echo "go things"
sudo -u watgbridge go clean 2> /dev/null
sudo -u watgbridge git pull origin 2> /dev/null
sudo -u watgbridge go build 2> /dev/null

systemctl start watgbridge
echo "Everything is up to date now :)"