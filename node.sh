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
token=$1
port=$2

# Check if pterodactyl directory exists 
if [ ! -d "/var/lib/pterodactyl/volumes" ]; then
    echo "No pterodactyl/volumes directory."
    exit 1
fi

echo ""
echo "Token $token"
echo "Port $port"
echo ""

if [[ $# -eq 0 ]]; then
    echo "No args specified."
    exit 1
fi

# Check if nginx is install 
if [ -z "$(command -v nginx)" ]; then
    echo "Nginx should be installed to run this script."
    exit 1
elif [ -z "$(command -v zip)" ]; then
    echo "zip should be installed to run this script."
    exit 1
elif [ -z "$(command -v unzip)" ]; then
    echo "unzip should be installed to run this script."
    exit 1
elif [ -z "$(command -v wget)" ]; then
    echo "wget should be installed to run this script."
    exit 1
elif [ -z "$(command -v curl)" ]; then
    echo "curl should be installed to run this script."
    exit 1
fi

# We have nginx + pterodactyl

# check if nginx conf already exists
if [ -f "/etc/nginx/sites-enables/auto.conf" ]; then
    rm /etc/nginx/sites-enables/auto.conf
fi

# Clear uploads folder
mkdir /uploads
cd /uploads
rm -rf *

# Copy all server's data 
echo "Copying data from pterodactyl directory..."
cp -r /var/lib/pterodactyl/volumes/* /uploads/


# Loop through all the folders in /uploads and zip every type of data to data.zip in each folder 
cd /uploads

echo "Looping through servers"
for dir in */
do
        cd $dir
        echo "zipping $dir"
        zip -r data.zip . -i ./*
        shopt -s extglob
        shopt -s dotglob
        for name in *
        do
                if [ "$name" != "data.zip" ] ;
                then
                        echo $name
                        rm -rf $name
                fi
        done
        cd ../;
done

# Now we have exactly what we want in /uploads folder
# Send the request to client
curl -X POST "localhost:$port/upload" -H "token: $token"