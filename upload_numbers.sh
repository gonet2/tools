#!/bin/sh
#上传numbers.xlsx到etcd

d=$(base64 numbers.xlsx)
curl http://192.168.99.100:2379/v2/keys/numbers -XPUT -d value="$d"
