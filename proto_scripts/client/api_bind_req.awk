###########################################################
## Scripts for generate ProtoHandler map binding code
##
BEGIN { RS = ""; FS ="\n" }
{
	for (i=1;i<=NF;i++)
	{
		if ($i ~ /^#.*/) {
			continue
		}

		split($i, a, ":")
		if (a[1] == "packet_type") {
			array["packet_type"] = a[2]
		} else if (a[1] == "name") {
			if (a[2] ~ /.*_req$/) {
				break
			} else {
				array["name"] = a[2]
			}
		}
	}

	if ("packet_type" in array && "name" in array && array["packet_type"] >= from && array["packet_type"] <= to) {
		print "\t{"array["packet_type"]",\tHandle_"array["name"]"},"
	}

	delete array
}
END {
print "};\n"
}