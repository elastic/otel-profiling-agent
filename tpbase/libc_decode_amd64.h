/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. Licensed under the Apache License 2.0.
 * See the file "LICENSE" for details.
 */

//go:build amd64

#ifndef LIBC_DECODE_X86_64
#define LIBC_DECODE_X86_64

#include <stdint.h>

uint32_t decode_pthread_getspecific(const uint8_t* code, size_t codesz);

#endif
