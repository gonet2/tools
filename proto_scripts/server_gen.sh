#!/bin/sh

##################################################
###   client proto & api
##################################################
printf "package client_handler\n" > proto.go
gawk -f server/proto.awk proto.txt >> proto.go 
gawk -f server/proto_func.awk proto.txt >> proto.go 

printf "package client_handler\n" > api.go
printf "\n" >> api.go
printf "import \"misc/packet\"\n" >> api.go
printf "import . \"types\"\n" >> api.go
printf "\n" >> api.go

gawk -f server/api.awk api.txt >> api.go 
gawk -f server/api_rcode.awk api.txt >> api.go

printf "var Handlers map[int16]func(*Session, *packet.Packet) []byte\n" >> api.go
printf "func init() {\n" >> api.go
printf "Handlers = map[int16]func(*Session, *packet.Packet) []byte {\n" >> api.go
if [ "$1" = "agent" ]; then
	gawk -v from=0 -v to=1000 -f server/api_bind_req.awk api.txt >> api.go
else
	gawk -v from=1001 -v to=65535 -f server/api_bind_req.awk api.txt >> api.go 
fi
printf "}" >> api.go
printf "}" >> api.go

mv -f proto.go ./proto_code/
mv -f api.go ./proto_code/
go fmt ./proto_code/
