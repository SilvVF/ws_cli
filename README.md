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
### flags
- -addr "host:port"
- -path "/route"
- -ticker 1000  (time in milliseconds)
