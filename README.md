# go-netcollect
Simple PoC to collect metrics for a specific network.

## Building images
```
$ docker-compose build
docker-compose build
Building srv1
Step 1/3 : FROM alpine
 ---> 11cd0b38bc3c
Step 2/3 : COPY entry.sh /
 ---> Using cache
 ---> e7bc7f937b1c
Step 3/3 : CMD ["/entry.sh"]
 ---> Using cache
 ---> b3a520015724

Successfully built b3a520015724
Successfully tagged pinger:latest
Building srv2
Step 1/3 : FROM alpine
 ---> 11cd0b38bc3c
Step 2/3 : COPY entry.sh /
 ---> Using cache
 ---> e7bc7f937b1c
Step 3/3 : CMD ["/entry.sh"]
 ---> Using cache
 ---> b3a520015724

Successfully built b3a520015724
Successfully tagged pinger:latest
Building srv3
Step 1/3 : FROM alpine
 ---> 11cd0b38bc3c
Step 2/3 : COPY entry.sh /
 ---> Using cache
 ---> e7bc7f937b1c
Step 3/3 : CMD ["/entry.sh"]
 ---> Using cache
 ---> b3a520015724

Successfully built b3a520015724
Successfully tagged pinger:latest
```
## Fire up testbed
```
$ docker stack deploy -c docker-compose.yml test
Ignoring unsupported options: build

Creating network test_testnet
Creating service test_srv1
Creating service test_srv2
Creating service test_srv3
```

## Hook into `test_testnet`

```
$ go run main.go "test_.*"
### test_testnet
## Start Collector for: 04e56092a46d6fa5c6ada408d8318c654534b086a663c6af21ac5a41e0dafd53
## Start Collector for: 42df758522df90df7c7e73d1eb564487f77c26aec1db23462cbfb72efbc2e8b8
## Start Collector for: ce9baff33adbae0c4e4fd5f8051994e223eb1eeb812bbd44b0a3ac3a5f7459fa
cntId:04e56092a46d6fa5c6ada408d8318c654534b086a663c6af21ac5a41e0dafd53 time:2018-07-23T09:23:08.548533 eth0.RxBytes:21602 eth0.TxBytes:21560 eth1.RxBytes:1292 eth1.TxBytes:504
cntId:42df758522df90df7c7e73d1eb564487f77c26aec1db23462cbfb72efbc2e8b8 time:2018-07-23T09:23:08.551612 eth0.RxBytes:21322 eth0.TxBytes:21322 eth1.RxBytes:788 eth1.TxBytes:0
cntId:ce9baff33adbae0c4e4fd5f8051994e223eb1eeb812bbd44b0a3ac3a5f7459fa time:2018-07-23T09:23:08.556277 eth0.RxBytes:21560 eth0.TxBytes:21518 eth1.RxBytes:1054 eth1.TxBytes:224
cntId:ce9baff33adbae0c4e4fd5f8051994e223eb1eeb812bbd44b0a3ac3a5f7459fa time:2018-07-23T09:23:09.55956 eth0.RxBytes:21756 eth0.TxBytes:21714 eth1.RxBytes:1054 eth1.TxBytes:224
cntId:04e56092a46d6fa5c6ada408d8318c654534b086a663c6af21ac5a41e0dafd53 time:2018-07-23T09:23:09.563722 eth0.RxBytes:21798 eth0.TxBytes:21756 eth1.RxBytes:1292 eth1.TxBytes:504
cntId:42df758522df90df7c7e73d1eb564487f77c26aec1db23462cbfb72efbc2e8b8 time:2018-07-23T09:23:09.567108 eth0.RxBytes:21518 eth0.TxBytes:21518 eth1.RxBytes:788 eth1.TxBytes:0
*snip*
```

## ToDo

 - [ ] Only output device connected to network