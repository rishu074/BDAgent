#!/bin/bash

# The Node Script 
# - Runs on the pterodactyl's Nodes

# Check if its the root
if [ "$(id -u)" != "0" ]; then
   echo "You should mind running the script as a root user. :)" 1>&2
   exit 1
fi

echo "Starting Backup Script."
# define args 
args=("$@")

echo $args[0]

# Check if pterodactyl directory exists 
if [ ! -d "/var/lib/pterodactyl/voluems" ]; then
    echo "No pterodactyl/voluems directory."
    exit 1
fi
