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