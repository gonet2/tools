#!/bin/sh

##################################################
###   client proto & api
##################################################
Head="\
using UnityEngine;\n\
using System;\n\
using System.Collections.Generic;\n"

printf "${Head}namespace NetProto.Proto {\n" > proto.cs
gawk -f client/proto.awk proto.txt >> proto.cs 
#gawk -f client/proto_func.awk proto.txt >> proto.cs 
printf "}" >> proto.cs


printf "${Head}namespace NetProto.Api {\n" > api.cs
gawk -f client/api.awk api.txt >> api.cs 
gawk -f client/api_rcode.awk api.txt >> api.cs
printf "public Dictionary<ushort, Func<ByteArray, object>> Handler = new Dictionary<ushort, Func<ByteArray, object>>(){\n" >> api.cs
gawk -v from=0 -v to=65535 -f client/api_bind_req.awk api.txt >> api.cs
printf "}\n" >> api.cs
printf "}" >> api.cs

mv -f proto.cs ./proto_code/NetProto.cs
mv -f api.cs ./proto_code/NetApi.cs
