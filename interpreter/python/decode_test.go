/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Apache License 2.0.
 * See the file "LICENSE" for details.
 */

package python

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-ebpf-profiler/libpf"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeArm64Stubs(t *testing.T) {
	val := decodeStubArgumentWrapperARM64(
		[]byte{
			0x40, 0x0a, 0x00, 0x90, 0x01, 0xd4, 0x43, 0xf9,
			0x22, 0x60, 0x17, 0x91, 0x40, 0x00, 0x40, 0xf9,
			0xa2, 0xff, 0xff, 0x17},
		0, 0, 0)
	assert.Equal(t, libpf.SymbolValue(1496), val, "PyEval_ReleaseLock stub test")

	val = decodeStubArgumentWrapperARM64(
		[]byte{
			0x80, 0x12, 0x00, 0xb0, 0x02, 0xd4, 0x43, 0xf9,
			0x41, 0xf4, 0x42, 0xf9, 0x61, 0x00, 0x00, 0xb4,
			0x40, 0xc0, 0x17, 0x91, 0xad, 0xe4, 0xfe, 0x17},
		0, 0, 0)
	assert.Equal(t, libpf.SymbolValue(1520), val, "PyGILState_GetThisThreadState test")

	// Python 3.10.12 on ARM64 Nix
	val = decodeStubArgumentWrapperARM64(
		[]byte{
			0x40, 0x1a, 0x00, 0xd0, // adrp	x0, 0xffffa0eff000 <mknodat@got.plt>
			0x00, 0xa0, 0x46, 0xf9, // ldr	x0, [x0, #3392]
			0x01, 0x28, 0x41, 0xf9, // ldr	x1, [x0, #592]
			0x61, 0x00, 0x00, 0xb4, // cbz	x1, 0xffffa0bb53c8 <PyGILState_GetThisThreadState+24>
			0x00, 0x5c, 0x42, 0xb9, // ldr	w0, [x0, #604]
			0x93, 0x20, 0xff, 0x17, // b	0xffffa0b7d610 <pthread_getspecific@plt>
			0x00, 0x00, 0x80, 0xd2, // mov	x0, #0x0
			0xc0, 0x03, 0x5f, 0xd6, // ret
		},
		0, 0, 0)
	assert.Equal(t, libpf.SymbolValue(604), val, "PyGILState_GetThisThreadState test")
}
