<h3 align="center">Bd Agent</h3>

---

<p align="center"> Helps in recieveing files from BdClient and serving it to users
    <br> 
</p>

### Installing

A step by step series of examples that tell you how to get app env running.

First Clone this git repo
Lets assume your installation directory as (DIR)

```
cd (DIR) && mkdir logs
```

Then Configure the config.yml and bdagent.service and there you go

```
mv (DIR)/bdagent.service /etc/systemd/system/ 
systemctl enable --now bdagent.service
```