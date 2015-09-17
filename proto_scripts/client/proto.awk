###########################################################
## generate proto payload struct 
##
@include "client/header.awk"
BEGIN { RS = "==="; FS ="\n" 
print "\npublic abstract class NetBase{"
print "\tpublic ushort NetMsgId;"
print "\tpublic abstract void Pack(ByteArray w);"
print "}\n"
}
{

	pack_code = ""
	unpack_code = ""
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /[A-Za-z_]+=/) {
			name = substr($i,1, match($i,/=/)-1)
			print "\npublic class " name " : NetBase {"
			typeok = "true"
		} else {
			if ($i!="" && typeok) {	
				print(_field($i))
				pack_code = pack_code "\t\t" _writer($i)
				unpack_code = unpack_code  _reader($i)
			}
		}

	}

	if (typeok) {
		#writer
		print "\n\tpublic override void Pack(ByteArray w) {"
		print pack_code
		print "\t}"
		#reader
		print "\tpublic  static "name" UnPack(ByteArray reader){"
		print "\t\tUInt16 narr = 0;"
		print "\t\t"name " tbl = new " name "();"
		print unpack_code
		print "\t\treturn tbl;\n\t}"

		print "}"
	}

	typeok=false
}
END { }

function _field(line) {
	split(line, a, " ")

	if (a[2] in TYPES) {
		return "\tpublic " TYPES[a[2]] " " a[1] ";"
	} else if (a[2] == "array") {
		if (a[3] in TYPES) {
			return "\tpublic " TYPES[a[3]] "[] " a[1] ";"
		} else {
			return "\tpublic " a[3] "[] " a[1] ";"
		}
	} else {
		return "\tpublic " a[2]" " a[1] ";"
	}
}

function _writer(line) {
	split(line, a, " ")

	if (a[2] in WRITERS) {
		return "w." WRITERS[a[2]] "(this." a[1] ");\n"
	} else if (a[2] == "array") {
		ret = "w.WriteUnsignedInt16((UInt16)this." a[1] ".Length);\n"
		if (a[3] == "byte") {
			ret = ret "w.WriteRawBytes(this."a[1]");\n"
			return ret
		} else if (a[3] in TYPES) {
			ret = ret "\tforeach (int k in this." a[1] ") {\n"
				ret = ret "\t\tw." WRITERS[a[3]] "(this." a[1] "[k]);\n"
			return ret "\t}\n"
		} else {
			ret = ret "\tforeach ("a[3]" k in this." a[1] ") {\n"
				ret = ret "\t\tk.Pack(w);\n"
			return ret "\t}\n"
		}
	} else {
		return "this." a[1] ".Pack(w);\n"
	}
}

function _reader(line){
	if (line ~ /^#.*/ || line ~ /^===/) {
		return ""
	}
	ret =  "\t"
	split(line, a, " ")
	if (a[2] ==  "array") {
		if (a[3] == "byte") { 		## bytes
			ret = ret "\ttbl."a[1]" = reader.ReadBytes();"
			
		} else if (a[3] in READERS) {	## primitives
			
			ret = ret "\tnarr = reader.ReadUnsignedInt16();\n"
			ret = ret "\t\ttbl."a[1]" = new "TYPES[a[3]]"[narr];\n"
			ret = ret "\t\tfor (int i = 0; i < narr; i++) {\n"
			ret = ret "\t\t\ttbl."a[1]"[i] =reader."READERS[a[3]]"();\n"
			#ret = ret "\t\t\ttbl."a[1]"[i] =  v;\n"
			ret = ret "\t\t}"
			
		} else {	## struct
			
			ret = ret "\tnarr = reader.ReadUnsignedInt16();\n"
			ret = ret "\t\ttbl."a[1]" = new "a[3]"[narr];\n"
			ret = ret "\t\tfor (int i = 0; i < narr; i++){\n"
			ret = ret "\t\t\ttbl."a[1]"[i] = "a[3]".UnPack(reader);\n"
			ret = ret "\t\t}"
		
		}
	} else if (!(a[2] in READERS)) {
		ret = ret "\ttbl."a[1]"  = "a[2]".UnPack(reader);"
	} else {
		ret = ret "\ttbl."a[1]"  = reader." READERS[a[2]] "();"
	}
	return ret "\n";
}



