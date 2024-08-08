/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Apache License 2.0.
 * See the file "LICENSE" for details.
 */

package elfunwindinfo

import (
	"encoding/base64"
	"os"
	"testing"

	sdtypes "github.com/open-telemetry/opentelemetry-ebpf-profiler/nativeunwind/stackdeltatypes"

	"github.com/stretchr/testify/require"
)

// Base64-encoded data from /usr/bin/volname on a stock debian box, the smallest
// 64-bit executable on my system (about 6k).
var usrBinVolname = `f0VMRgIBAQAAAAAAAAAAAAMAPgABAAAA8AkAAAAAAABAAAAAAAAAADgRAAAAAAAAAAAAAEAAOAAJ
AEAAGwAaAAYAAAAFAAAAQAAAAAAAAABAAAAAAAAAAEAAAAAAAAAA+AEAAAAAAAD4AQAAAAAAAAgA
AAAAAAAAAwAAAAQAAAA4AgAAAAAAADgCAAAAAAAAOAIAAAAAAAAcAAAAAAAAABwAAAAAAAAAAQAA
AAAAAAABAAAABQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFQNAAAAAAAAVA0AAAAAAAAAACAA
AAAAAAEAAAAGAAAAgA0AAAAAAACADSAAAAAAAIANIAAAAAAAkAIAAAAAAACwAgAAAAAAAAAAIAAA
AAAAAgAAAAYAAACYDQAAAAAAAJgNIAAAAAAAmA0gAAAAAADAAQAAAAAAAMABAAAAAAAACAAAAAAA
AAAEAAAABAAAAFQCAAAAAAAAVAIAAAAAAABUAgAAAAAAAEQAAAAAAAAARAAAAAAAAAAEAAAAAAAA
AFDldGQEAAAA+AsAAAAAAAD4CwAAAAAAAPgLAAAAAAAAPAAAAAAAAAA8AAAAAAAAAAQAAAAAAAAA
UeV0ZAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAABS
5XRkBAAAAIANAAAAAAAAgA0gAAAAAACADSAAAAAAAIACAAAAAAAAgAIAAAAAAAABAAAAAAAAAC9s
aWI2NC9sZC1saW51eC14ODYtNjQuc28uMgAEAAAAEAAAAAEAAABHTlUAAAAAAAIAAAAGAAAAIAAA
AAQAAAAUAAAAAwAAAEdOVQCSX5P2bs4LXU0AZhU77QH4cIow5gIAAAATAAAAAQAAAAYAAAAAAQAA
AAAAAgAAAAATAAAAOfKLHAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACeAAAAIAAAAAAAAAAA
AAAAAAAAAAAAAACBAAAAEgAAAAAAAAAAAAAAAAAAAAAAAAB9AAAAEgAAAAAAAAAAAAAAAAAAAAAA
AAAuAAAAEgAAAAAAAAAAAAAAAAAAAAAAAAA4AAAAEgAAAAAAAAAAAAAAAAAAAAAAAABcAAAAEgAA
AAAAAAAAAAAAAAAAAAAAAABJAAAAEgAAAAAAAAAAAAAAAAAAAAAAAACMAAAAEgAAAAAAAAAAAAAA
AAAAAAAAAAC6AAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAdAAAAEgAAAAAAAAAAAAAAAAAAAAAAAAAL
AAAAEgAAAAAAAAAAAAAAAAAAAAAAAABpAAAAEgAAAAAAAAAAAAAAAAAAAAAAAAAnAAAAEgAAAAAA
AAAAAAAAAAAAAAAAAADJAAAAIAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAAEgAAAAAAAAAAAAAAAAAA
AAAAAABOAAAAEgAAAAAAAAAAAAAAAAAAAAAAAADdAAAAIAAAAAAAAAAAAAAAAAAAAAAAAABuAAAA
IgAAAAAAAAAAAAAAAAAAAAAAAABiAAAAEQAYACAQIAAAAAAACAAAAAAAAAAAbGliYy5zby42AF9f
cHJpbnRmX2NoawBleGl0AHNldGxvY2FsZQBwZXJyb3IAZGNnZXR0ZXh0AF9fc3RhY2tfY2hrX2Zh
aWwAcmVhZABfX2ZwcmludGZfY2hrAGxzZWVrAHN0ZGVycgBvcGVuAF9fY3hhX2ZpbmFsaXplAGJp
bmR0ZXh0ZG9tYWluAF9fbGliY19zdGFydF9tYWluAF9JVE1fZGVyZWdpc3RlclRNQ2xvbmVUYWJs
ZQBfX2dtb25fc3RhcnRfXwBfSnZfUmVnaXN0ZXJDbGFzc2VzAF9JVE1fcmVnaXN0ZXJUTUNsb25l
VGFibGUAR0xJQkNfMi4zLjQAR0xJQkNfMi40AEdMSUJDXzIuMi41AAAAAAAAAgACAAIAAwACAAIA
AgAAAAIABAACAAIAAAACAAQAAAACAAIAAAAAAAAAAQADAAEAAAAQAAAAAAAAAHQZaQkAAAQA9wAA
ABAAAAAUaWkNAAADAAMBAAAQAAAAdRppCQAAAgANAQAAAAAAAIANIAAAAAAACAAAAAAAAADwCgAA
AAAAAIgNIAAAAAAACAAAAAAAAACwCgAAAAAAAAgQIAAAAAAACAAAAAAAAAAIECAAAAAAAHAPIAAA
AAAABgAAAAEAAAAAAAAAAAAAAHgPIAAAAAAABgAAAAIAAAAAAAAAAAAAAIAPIAAAAAAABgAAAAMA
AAAAAAAAAAAAAIgPIAAAAAAABgAAAAQAAAAAAAAAAAAAAJAPIAAAAAAABgAAAAUAAAAAAAAAAAAA
AJgPIAAAAAAABgAAAAYAAAAAAAAAAAAAAKAPIAAAAAAABgAAAAcAAAAAAAAAAAAAAKgPIAAAAAAA
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
GQEAAAAAAAALAAAAAAAAABgAAAAAAAAAFQAAAAAAAAAAAAAAAAAAAAMAAAAAAAAAWA8gAAAAAAAH
AAAAAAAAACgGAAAAAAAACAAAAAAAAAAQAgAAAAAAAAkAAAAAAAAAGAAAAAAAAAAeAAAAAAAAAAgA
AAAAAAAA+///bwAAAAABAAAIAAAAAP7//28AAAAA6AUAAAAAAAD///9vAAAAAAEAAAAAAAAA8P//
bwAAAAC6BQAAAAAAAPn//28AAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJgNIAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
CBAgAAAAAAA1ZjkzZjY2ZWNlMGI1ZDRkMDA2NjE1M2JlZDAxZjg3MDhhMzBlNi5kZWJ1ZwAAAAAF
PtqRAC5zaHN0cnRhYgAuaW50ZXJwAC5ub3RlLkFCSS10YWcALm5vdGUuZ251LmJ1aWxkLWlkAC5n
bnUuaGFzaAAuZHluc3ltAC5keW5zdHIALmdudS52ZXJzaW9uAC5nbnUudmVyc2lvbl9yAC5yZWxh
LmR5bgAuaW5pdAAucGx0AC5wbHQuZ290AC50ZXh0AC5maW5pAC5yb2RhdGEALmVoX2ZyYW1lX2hk
cgAuZWhfZnJhbWUALmluaXRfYXJyYXkALmZpbmlfYXJyYXkALmpjcgAuZHluYW1pYwAuZGF0YQAu
YnNzAC5nbnVfZGVidWdsaW5rAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALAAAAAQAAAAIAAAAAAAAAOAIAAAAAAAA4AgAAAAAA
ABwAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAEwAAAAcAAAACAAAAAAAAAFQCAAAAAAAA
VAIAAAAAAAAgAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAACEAAAAHAAAAAgAAAAAAAAB0
AgAAAAAAAHQCAAAAAAAAJAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA0AAAA9v//bwIA
AAAAAAAAmAIAAAAAAACYAgAAAAAAACQAAAAAAAAABQAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAPgAA
AAsAAAACAAAAAAAAAMACAAAAAAAAwAIAAAAAAADgAQAAAAAAAAYAAAABAAAACAAAAAAAAAAYAAAA
AAAAAEYAAAADAAAAAgAAAAAAAACgBAAAAAAAAKAEAAAAAAAAGQEAAAAAAAAAAAAAAAAAAAEAAAAA
AAAAAAAAAAAAAABOAAAA////bwIAAAAAAAAAugUAAAAAAAC6BQAAAAAAACgAAAAAAAAABQAAAAAA
AAACAAAAAAAAAAIAAAAAAAAAWwAAAP7//28CAAAAAAAAAOgFAAAAAAAA6AUAAAAAAABAAAAAAAAA
AAYAAAABAAAACAAAAAAAAAAAAAAAAAAAAGoAAAAEAAAAAgAAAAAAAAAoBgAAAAAAACgGAAAAAAAA
EAIAAAAAAAAFAAAAAAAAAAgAAAAAAAAAGAAAAAAAAAB0AAAAAQAAAAYAAAAAAAAAOAgAAAAAAAA4
CAAAAAAAABcAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAegAAAAEAAAAGAAAAAAAAAFAI
AAAAAAAAUAgAAAAAAAAQAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAQAAAAAAAAAH8AAAABAAAABgAA
AAAAAABgCAAAAAAAAGAIAAAAAAAAaAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAACIAAAA
AQAAAAYAAAAAAAAA0AgAAAAAAADQCAAAAAAAAMICAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAA
AAAAjgAAAAEAAAAGAAAAAAAAAJQLAAAAAAAAlAsAAAAAAAAJAAAAAAAAAAAAAAAAAAAABAAAAAAA
AAAAAAAAAAAAAJQAAAABAAAAAgAAAAAAAACgCwAAAAAAAKALAAAAAAAAWAAAAAAAAAAAAAAAAAAA
AAQAAAAAAAAAAAAAAAAAAACcAAAAAQAAAAIAAAAAAAAA+AsAAAAAAAD4CwAAAAAAADwAAAAAAAAA
AAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAqgAAAAEAAAACAAAAAAAAADgMAAAAAAAAOAwAAAAAAAAc
AQAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAALQAAAAOAAAAAwAAAAAAAACADSAAAAAAAIAN
AAAAAAAACAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAACAAAAAAAAADAAAAADwAAAAMAAAAAAAAAiA0g
AAAAAACIDQAAAAAAAAgAAAAAAAAAAAAAAAAAAAAIAAAAAAAAAAgAAAAAAAAAzAAAAAEAAAADAAAA
AAAAAJANIAAAAAAAkA0AAAAAAAAIAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAANEAAAAG
AAAAAwAAAAAAAACYDSAAAAAAAJgNAAAAAAAAwAEAAAAAAAAGAAAAAAAAAAgAAAAAAAAAEAAAAAAA
AACDAAAAAQAAAAMAAAAAAAAAWA8gAAAAAABYDwAAAAAAAKgAAAAAAAAAAAAAAAAAAAAIAAAAAAAA
AAgAAAAAAAAA2gAAAAEAAAADAAAAAAAAAAAQIAAAAAAAABAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA
CAAAAAAAAAAAAAAAAAAAAOAAAAAIAAAAAwAAAAAAAAAgECAAAAAAABAQAAAAAAAAEAAAAAAAAAAA
AAAAAAAAACAAAAAAAAAAAAAAAAAAAADlAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAQEAAAAAAAADQA
AAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAQAAAAMAAAAAAAAAAAAAAAAAAAAAAAAARBAA
AAAAAAD0AAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAA==`

var firstDeltas = sdtypes.StackDeltaArray{
	{Address: 0x850, Hints: sdtypes.UnwindHintKeep,
		Info: deltaRSP(16, 0)},
	{Address: 0x856, Info: deltaRSP(24, 0)},
	{Address: 0x860, Hints: sdtypes.UnwindHintKeep, Info: deltaRSP(8, 0)},
	{Address: 0x8d1, Info: deltaRSP(16, 16)},
	{Address: 0x8d2, Info: deltaRSP(24, 16)},
	{Address: 0x8e7, Info: deltaRSP(80, 16)},
}

func TestExtractStackDeltasFromFilename(t *testing.T) {
	buffer, err := base64.StdEncoding.DecodeString(usrBinVolname)
	require.NoError(t, err)
	// Write the executable file to a temporary file, and the symbol
	// file, too.
	exeFile, err := os.CreateTemp("/tmp", "dwarf_extract_elf_")
	require.NoError(t, err)
	defer exeFile.Close()
	_, err = exeFile.Write(buffer)
	require.NoError(t, err)
	err = exeFile.Sync()
	require.NoError(t, err)
	defer os.Remove(exeFile.Name())
	filename := exeFile.Name()

	var data sdtypes.IntervalData
	err = Extract(filename, &data)
	require.NoError(t, err)
	for _, delta := range data.Deltas {
		t.Logf("%#v", delta)
	}
	require.Equal(t, data.Deltas[:len(firstDeltas)], firstDeltas)
}
