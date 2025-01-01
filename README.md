### building
download wsclient.exe from releases or build from source
 ```shell
 go build main.go
 ```
### run 
add wsclient.exe to path
```shell
 wsclient -addr "localhost:8080" -path "/echo"
```
connects to the ws and reads from stdin
escape seq is os default

### flags
- -addr "host:port"
- -path "/route"
- -ticker 1000  (time in milliseconds)
