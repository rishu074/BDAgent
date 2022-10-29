#!/bin/bash

nodeURl=$1
nodeDir=$2
output=$3

cd $nodeDir
rm -rf *
wget -O "$output" "$nodeURl"
