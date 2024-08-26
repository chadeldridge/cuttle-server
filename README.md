# cuttle

Server macros for quick actions.

# Planned v1 Features

- Select a profile then a server group. Be presented with a list of actions. Clicking on an action will perform the action against all servers in the group. Servers will be listed with a status icon on the right. Logs from the actions will be shown at the bottom.
- Profiles dictate which actions can be ran against which groups.
- Only Admins can create actions, add servers and connections, create groups, and profiles.
- User Groups must be used to provide users with access to profiles.

## Post v1 Planned Features

- Satellite servers. Main server can send action requests to satellite servers in other vlans or locations. This would allow for easier testing in environments where special vlan access is needed, dns spoofing is setup, etc.
- Allow external auth and group queries. Active Directory, LDAP, Oauth2, or allow users to provide their own script the server calls out to (performance concerns?).

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
`pwd`/cuttle.yaml
`pwd`/config.yaml
```

If no config file is provided, cuttle will move on to checking for env variables and cli flags.

The way servers, connectors, and tests like ssh work may need to change later. Different tests might need different connectors to be used against the same server (one username needed for a simple echo while another needed to test reading a protected file or starting a service). For simplicity, maybe allowing server+connector to be defined in the group is best? Then group+tile selection matter and are controlled by the profile (a profile only allows for a selected set of groups and tiles). You would need to change profiles to access the privileged group and tile set.