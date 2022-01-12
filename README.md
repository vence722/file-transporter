# file-transporter
A very light-weight file sharing platform, including server and client

## Installation
```shell
git clone https://github.com/vence722/file-transporter.git
cd file-transporter
go build 
```

## Usage
Running server on port `8123`:
```shell
./file-transporter server :8123
```

For client logins to server `10.0.0.1:8123` as user `test`:
```shell
./file-transporter client 10.0.0.1:8123 test
```

In client's terminal, a menu is prompted:
```shell
Choose Action:
(1) List online users
(2) Transfer file
(3) Logout
Please input your action >>>
```
Please follow the tips to send file to other users that logged in to this server.