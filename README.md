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

Input `1` to list all online users:
```shell
Online Users:
test
alice
```

Input `2` to send file to other online user. You need to input target username and file path in your local machine:
```shell
Input target username:
alice
Input file path:
/path/to/your/file
```
Then just wait for the file to be sent!

If you're the receiving client, you'll get notifications in the terminal:
```shell
[INFO] Start receiving file xxx
[INFO] File xxx received successfully!
```

Last, input `3` for logout:
```
[INFO] Logged out
```

## Enjoy your file sharing experience!