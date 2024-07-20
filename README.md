# cuttle
Server macros for quick actions.

### Enable Ping on Linux
You need to run the following in the container build or on the server to enable UDP ping. I'm avoiding ICMP ping to eliminate the need to run this application with privilege.

```
sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"
```

On WSL2 add this to %HOMEPATH%/.wslconfig
```
[wsl2]
kernelCommandLine=sysctl.net.ipv4.ping_group_range=\"0 2147483647\"
```
## Config
cli flags > env vars > config file > defaults

### Config Environment Variables
All cuttle environment variables should be the all uppercase version of the yaml config names and should contain the prefix "CUTTLE_".
Example:
```
export CUTTLE_API_HOST=192.169.1.50
```

### Config File
If a config file is specified through the cli flag or environment variable, cuttle will consider it mission critical and will panic if the file is missing or cannot be loaded.

CLI Arg:
```
-c, --config Specify the full path and file name for the config file.
```
Example:
```
cuttle-server -c /tmp/cuttle.yaml
```

Env Var:
```
export CUTTLE_CONFIG_FILE=/path/to/config.yaml
```

If no cli flag or environment varaible is found, cuttle will check two default locations for a config file. If the file is found in these locations, cuttle will treat it as mission critical and panic if it fails to load the config.

Config file will be checked for in this order:
```
~/cuttle.yaml
~/.config/cuttle/config.yaml
`pwd`/config.yaml
`pwd`/cuttle.yaml
```

If no config file is provided, cuttle will move on to checking for env variables and cli flags.

The way servers, connectors, and tests like ssh work may need to change later. Different tests might need different connectors to be used against the same server (one username needed for a simple echo while another needed to test reading a protected file or starting a service). For simplicity, maybe allowing server+connector to be defined in the group is best? Then group+tile selection matter and are controlled by the profile (a profile only allows for a selected set of groups and tiles). You would need to change profiles to access the privileged group and tile set.