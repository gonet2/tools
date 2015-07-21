###########################################################
## generate proto payload struct 
##
@include "client/header.awk"
BEGIN { RS = "==="; FS ="\n" 
print ""
}
{

	pack_code = ""
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /[A-Za-z_]+=/) {
			name = substr($i,1, match($i,/=/)-1)
			print "\npublic class " name " {"
			typeok = "true"
		} else {
			if ($i!="" && typeok) {	
				print(_field($i))
				pack_code = pack_code "\t" _writer($i)
			}
		}
	}

	if (typeok) {
		
		print "\n\tpublic void Pack(ByteArray w) {"
		print pack_code
		print "\t}"
		print "}"
	}

	typeok=false
}
END { }

function _field(line) {
	split(line, a, " ")

	if (a[2] in TYPES) {
		return "\tpublic " TYPES[a[2]] " " a[1]
	} else if (a[2] == "array") {
		if (a[3] in TYPES) {
			return "\tpublic " TYPES[a[3]] "[] " a[1]
		} else {
			return "\tpublic " a[3] "[] " a[1]
		}
	} else {
		return "\tpublic " a[2]" " a[1]
	}
}

function _writer(line) {
	split(line, a, " ")

	if (a[2] in WRITERS) {
		return "w." WRITERS[a[2]] "(this." a[1] ")\n"
	} else if (a[2] == "array") {
		ret = "w.WriteU16(uint16(len(this." a[1] ")))\n"
		if (a[3] == "byte") {
			ret = ret "w.WriteRawBytes(this."a[1]")\n"
			return ret
		} else if (a[3] in TYPES) {
			ret = ret "\tfor (int k in this." a[1] ") {\n"
				ret = ret "\t\tw." WRITERS[a[3]] "(this." a[1] "[k])\n"
			return ret "\t}\n"
		} else {
			ret = ret "\tfor (int k in this." a[1] ") {\n"
				ret = ret "\t\tthis."a[1]"[k].Pack(w)\n"
			return ret "\t}\n"
		}
	} else {
		return "this." a[1] ".Pack(w)\n"
	}
}
