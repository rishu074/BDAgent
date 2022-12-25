
# The BDAgent

This is the project made to upload and preserve backups of a specified directory.


## Features

- Unlimited uploaders
- FTP Uploading
- Local uploaders
- Secure uploading with Bearer tokens
- Websockets support

## How it works?

You basically, create a `config.yml` file, in which you specify the `token`,`uploaders`,`data-folder` and `chunk-sizes`.

When you run the bdclient and hit its api, the bdclient starts uploading the data to bdagent, and bdagent sends those chunks to configured places.

### For example
We have a uploader, lets call it `uploader-1`\
We have a agent configured, lets call it `agent-a`

When the api of `uploader-1`, gets hitted by `crontab` or by manually, its starts sending the chunks as defined.

The `uploader-1` will log every possible debugging output to console.


## Installation

This agent is only built for ubuntu/linux based distributions.

```bash
# Create the directory for agent
$ mkdir /var/apps
$ mkdir /var/apps/bdagent
$ cd /var/apps/bdagent

# Install the latest binaries
$ wget https://github.com/NotRoyadma/BDAgent/releases/latest/download/agent
$ wget https://github.com/NotRoyadma/BDAgent/releases/latest/download/bdagent.service
```
    


