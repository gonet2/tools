#!/bin/sh
#上传numbers.xlsx到etcd
FILE=numbers.xlsx
base64 $FILE > $FILE.base64
curl http://192.168.99.100:2379/v2/keys/numbers -XPUT --data-urlencode value@$FILE.base64
rm $FILE.base64
