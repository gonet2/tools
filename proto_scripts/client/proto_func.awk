###########################################################
## generate protocol packet reader
##
@include "client/header.awk"
BEGIN { RS = ""; FS ="\n"}
{
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/ || $i ~ /^===/) {
			continue
		}
		split($i, a, " ")
		if (a[1] ~ /[A-Za-z_]+=/) {
			name = substr(a[1],1, match(a[1],/=/)-1)
			print "public "name" PKT_"name"(ByteArray reader){"
			print  "\t"name " tbl = new " name "();"
			typeok = "true"
		} else if (a[2] ==  "array") {
			if (a[3] == "byte") { 		## bytes
				print "\ttbl."a[1]" = reader.ReadBytes()"
				
			} else if (a[3] in READERS) {	## primitives
				print "\t{"
				print "\tnarr := reader.ReadU16()"
				print "\tfor (i:=0;i<int(narr);i++ ) {"
				print "\t\tv := reader."READERS[a[3]]"()"
				print "\t\ttbl."a[1]" = append(tbl."a[1]", v)"
				print "\t}\n"
				print "\t}\n"
			} else {	## struct
				print "\t{"
				print "\tnarr := reader.ReadU16()"
				print "\ttbl."a[1]"=make([]"a[3]",narr)"
				print "\tfor i:=0;i<int(narr);i++ {"
				print "\t\ttbl."a[1]"[i], err = PKT_"a[3]"(reader)"
				print "\t}\n"
				print "\t}\n"
			}
		}
		else if (!(a[2] in READERS)) {
			print "\t\ttbl."a[1]"  = PKT_"a[2]"(reader)"
		}
		else {
			print "\ttbl."a[1]"  = reader." READERS[a[2]] "()"
		}
	}

	if (typeok) {
		print "\treturn tbl"
		print "}\n"
	}

	typeok=false
}
END { }
