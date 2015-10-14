#!/bin/sh

##################################################
###   proto & api
##################################################
#modify the path to your self path.
#export PATH_AGENT=/go/gonet2/agent/src/client_handler
#export PATH_GAME=/go/gonet2/game/src/client_handler
export PATH_AGENT=./proto_code/agent/
export PATH_GAME=./proto_code/game/
export PATH_CLIENT=./proto_code/client/

go run api.go ./templates/agent/api.tmpl < api.txt > $PATH_AGENT/api.go
go fmt $PATH_AGENT/api.go

go run api.go ./templates/game/api.tmpl < api.txt > $PATH_GAME/api.go
go fmt $PATH_GAME/api.go

go run api.go ./templates/client/api.tmpl < api.txt > $PATH_CLIENT/NetApi.cs

go run proto.go ./templates/agent/proto.tmpl < proto.txt > $PATH_AGENT/proto.go
go fmt $PATH_AGENT/proto.go

go run proto.go ./templates/game/proto.tmpl < proto.txt > $PATH_GAME/proto.go
go fmt $PATH_GAME/proto.go

go run proto.go ./templates/client/proto.tmpl < proto.txt > $PATH_CLIENT/NetProto.cs
