
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
$ mkdir logs
$ cd /var/apps/bdagent

# Install the latest binaries
$ wget https://github.com/NotRoyadma/BDAgent/releases/latest/download/agent
$ wget https://github.com/NotRoyadma/BDAgent/releases/latest/download/bdagent.service
```

### Setup configuration file

```bash
# Create a config file
$ touch config.yml
```

#### Data format of config file

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`      | `string` | **Required**. Name of the application |
| `version`      | `string` | **Don't change it**. |
| `port`      | `integer` | **Required**, Port of the application |
| `nodes`      | `yaml-array` | **Required**, Names of the uploaders  |
| `dataDirectory`      | `string` | **Required**, The name of the data directory, *without any slash*  |
| `data_file`      | `string` | **Required**, File which is set in client.  |
| `token`      | `string` | **Required**, The authorization token for BDclient  |
| `BashFile`      | `string` | **Required**, Let it be as it is  |
| `IP_HEADER`      | `string` | **Required**, leave it `default` or if using cloudflare change it  |
| `ftp`      | `yaml-object` | **Required**, The ftp configuration  |
| `chunk_size`      | `integer` | **Required**, Dont change unless you don't know about it  |

An example demostration to config file
```
name: "Auto backup dnxrg"
version: "1.0.1"
port: 1337
nodes: 
  - game1
  - game2
  - game3
dataDirectory: "data"
data_file: "data.zip"
token: "SomerandomToken"
BashFile: "./avails/download.sh"
IP_HEADER: "default"
ftp: 
  enabled: false
  uri: "172.105.33.245:21"
  user: "username@ftp"
  password: "somepassword"
chunk_size: 4000000
```
### Setting up Systemd service
```bash
# Copy the service to systemd directory
$ cd /var/apps/bdagent
$ mv bdagent.service /etc/systemd/system/
$ systemctl enable --now bdagent.service
```

### There are prebuilt loggers for http and application
```bash
# To view the app logs
$ cat /var/apps/bdagent/logs/app.log

# To view the http logs
$ cat /var/apps/bdagent/logs/http.log

# To view the error logs (if any)
$ cat /var/apps/bdagent/logs/app.error.log

# You can also view the systemctl service status by doing
$ systemctl status bdagent.service

# You can view live http,app and error logs by doing
$ journalctl -u bdagent.service -e --follow
```


## Authors

- [@NotRoyadma](https://www.github.com/NotRoyadma)

