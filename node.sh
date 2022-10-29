#!/bin/bash

# The Node Script 
# - Runs on the pterodactyl's Nodes

# Check if its the root
if [ "$(id -u)" != "0" ]; then
   echo "You should mind running the script as a root user. :)" 1>&2
   exit 1
fi

