# fp-restrictions

frp server plugin to support client restrictions for [frp](https://github.com/fatedier/frp).

fp-restrictions will run as one single process and accept HTTP requests from frps.

Based on [fp-multiuser](https://github.com/gofrp/fp-multiuser)

### Features

* Support user authentication and restrictions saved in file.

### Download

Download fp-restrictions binary file from [Release](https://github.com/fjfj/fp-restrictions/releases).

### Requirements

frp version >= v0.31.0

### Usage

1. Create file `restrictions` including all support usernames and restrictions.

    ```
    user1=123:{ "SubDomain":"test", "ProxyType":"http.?", "Locations":"[/path]", "UseEncryption":"true" }
    user2=abc
    ```

    One user each line. Username and token:restrictions are split by `=`. Restrictions are in JSON. All JSON values are regexp strings. 

2. Run fp-restrictions:

    `./fp-restrictions -l 127.0.0.1:7200 -f ./restrictions`

3. Register plugin in frps.

    ```
    # frps.ini
    [common]
    bind_port = 7000

    [plugin.restrictions]
    addr = 127.0.0.1:7700
    path = /handler
    ops = Login, NewProxy
    ```

4. Specify username and meta_token in frpc configure file.

    For user1:

    ```
    # frpc.ini
    [common]
    server_addr = x.x.x.x
    server_port = 7000
    user = user1
    meta_token = 123

    [http1]
    type = http
    local_port = 8080
    remote_port = 6000
    subdomain = test
    use_encryption = true
    locations = /path
    
    ```

    For user2:

    ```
    # frpc.ini
    [common]
    server_addr = x.x.x.x
    server_port = 7000
    user = user2
    meta_token = abc

    [ssh]
    type = tcp
    local_port = 22
    remote_port = 6000
    ```
