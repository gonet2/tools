#!/bin/sh

##################################################
###   client proto & api
##################################################
#modify the path to your self path.
export PATH_AGENT=/go/gonet2/agent
export PATH_GAME=/go/gonet2/game

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
	printf "}" >> api.go
	printf "}" >> api.go
	go fmt ./
	mv -f proto.go $PATH_AGENT/src/client_handler
	mv -f api.go $PATH_AGENT/src/client_handler

else
	gawk -v from=1001 -v to=65535 -f server/api_bind_req.awk api.txt >> api.go 
	printf "}" >> api.go
	printf "}" >> api.go
	go fmt ./
	mv -f proto.go $PATH_GAME/src/client_handler
	mv -f api.go $PATH_GAME/src/client_handler
fi

