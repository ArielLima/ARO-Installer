package encryption

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"context"
	"encoding/hex"
	"io"
	"testing"
)

func TestNewAES256SHA512(t *testing.T) {
	for _, tt := range []struct {
		name    string
		key     []byte
		wantErr string
	}{
		{
			name: "valid",
			key:  make([]byte, 64),
		},
		{
			name:    "key too short",
			key:     make([]byte, 63),
			wantErr: "etm: key must be 64 bytes long",
		},
		{
			name:    "key too long",
			key:     make([]byte, 65),
			wantErr: "etm: key must be 64 bytes long",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAES256SHA512(context.Background(), tt.key)
			if err != nil && err.Error() != tt.wantErr ||
				err == nil && tt.wantErr != "" {
				t.Fatal(err)
			}
		})
	}
}

func TestAES256SHA512Open(t *testing.T) {
	for _, tt := range []struct {
		name       string
		key        []byte
		input      []byte
		wantOpened []byte
		wantErr    string
	}{
		{
			name:       "valid",
			key:        []byte("\x6a\x98\x95\x6b\x2b\xb2\x7e\xfd\x1b\x68\xdf\x5c\x40\xc3\x4f\x8b\xcf\xff\xe8\x17\xc2\x2d\xf6\x40\x2e\x5a\xb0\x15\x63\x4a\x2d\x2e\xab\x79\x86\x50\xfb\xce\xdc\x9d\xdd\x1c\x01\x32\xd6\x03\x99\xe6\x59\x81\x37\xb3\xdb\x67\x6f\x12\x34\x1d\xb9\x58\x18\x31\x30\x57"),
			input:      []byte("\xd9\x1c\x3c\x05\xb2\xf3\xc5\x93\x20\x9f\x9b\x67\x43\x8c\x0c\x3d\xe0\x80\x26\x59\x2a\x20\xb2\xe5\x5e\x30\xd6\xd1\x24\x1e\x34\x36\xbe\xfb\x79\x8e\x46\xb5\x95\xce\xe0\x79\x9c\x44\x5c\xaa\x83\x26\x92\xdb\x76\x34\x33\xe0\x0e\x0e\x54\xb2\x0b\x2f\xde\x63\x53\xf6"),
			wantOpened: []byte("test"),
		},
		{
			name:    "invalid - encrypted value tampered with",
			key:     []byte("\x6a\x98\x95\x6b\x2b\xb2\x7e\xfd\x1b\x68\xdf\x5c\x40\xc3\x4f\x8b\xcf\xff\xe8\x17\xc2\x2d\xf6\x40\x2e\x5a\xb0\x15\x63\x4a\x2d\x2e\xab\x79\x86\x50\xfb\xce\xdc\x9d\xdd\x1c\x01\x32\xd6\x03\x99\xe6\x59\x81\x37\xb3\xdb\x67\x6f\x12\x34\x1d\xb9\x58\x18\x31\x30\x57"),
			input:   []byte("\xda\x1c\x3c\x05\xb2\xf3\xc5\x93\x20\x9f\x9b\x67\x43\x8c\x0c\x3d\xe0\x80\x26\x59\x2a\x20\xb2\xe5\x5e\x30\xd6\xd1\x24\x1e\x34\x36\xbe\xfb\x79\x8e\x46\xb5\x95\xce\xe0\x79\x9c\x44\x5c\xaa\x83\x26\x92\xdb\x76\x34\x33\xe0\x0e\x0e\x54\xb2\x0b\x2f\xde\x63\x53\xf6"),
			wantErr: "message authentication failed",
		},
		{
			name:    "invalid - too short",
			key:     []byte("\x6a\x98\x95\x6b\x2b\xb2\x7e\xfd\x1b\x68\xdf\x5c\x40\xc3\x4f\x8b\xcf\xff\xe8\x17\xc2\x2d\xf6\x40\x2e\x5a\xb0\x15\x63\x4a\x2d\x2e\xab\x79\x86\x50\xfb\xce\xdc\x9d\xdd\x1c\x01\x32\xd6\x03\x99\xe6\x59\x81\x37\xb3\xdb\x67\x6f\x12\x34\x1d\xb9\x58\x18\x31\x30\x57"),
			input:   make([]byte, 31),
			wantErr: "encrypted value too short",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cipher, err := NewAES256SHA512(context.Background(), tt.key)
			if err != nil {
				t.Fatal(err)
			}

			opened, err := cipher.Open(tt.input)
			if err != nil && err.Error() != tt.wantErr ||
				err == nil && tt.wantErr != "" {
				t.Fatal(err)
			}

			if !bytes.Equal(tt.wantOpened, opened) {
				t.Error(string(opened))
			}
		})
	}
}

func TestAES256SHA512Seal(t *testing.T) {
	for _, tt := range []struct {
		name       string
		key        []byte
		randReader io.Reader
		input      []byte
		wantSealed []byte
		wantErr    string
	}{
		{
			name:       "valid",
			key:        []byte("\x6a\x98\x95\x6b\x2b\xb2\x7e\xfd\x1b\x68\xdf\x5c\x40\xc3\x4f\x8b\xcf\xff\xe8\x17\xc2\x2d\xf6\x40\x2e\x5a\xb0\x15\x63\x4a\x2d\x2e\xab\x79\x86\x50\xfb\xce\xdc\x9d\xdd\x1c\x01\x32\xd6\x03\x99\xe6\x59\x81\x37\xb3\xdb\x67\x6f\x12\x34\x1d\xb9\x58\x18\x31\x30\x57"),
			randReader: bytes.NewBufferString("\xd9\x1c\x3c\x05\xb2\xf3\xc5\x93\x20\x9f\x9b\x67\x43\x8c\x0c\x3d\x9c\x33\x5b\x16\xd6\x9a\x9c\xf2"),
			input:      []byte("test"),
			wantSealed: []byte("\xd9\x1c\x3c\x05\xb2\xf3\xc5\x93\x20\x9f\x9b\x67\x43\x8c\x0c\x3d\xe0\x80\x26\x59\x2a\x20\xb2\xe5\x5e\x30\xd6\xd1\x24\x1e\x34\x36\xbe\xfb\x79\x8e\x46\xb5\x95\xce\xe0\x79\x9c\x44\x5c\xaa\x83\x26\x92\xdb\x76\x34\x33\xe0\x0e\x0e\x54\xb2\x0b\x2f\xde\x63\x53\xf6"),
		},
		{
			name:       "rand.Read EOF",
			key:        make([]byte, 64),
			randReader: &bytes.Buffer{},
			wantErr:    "EOF",
		},
		{
			name:       "rand.Read unexpected EOF",
			key:        make([]byte, 64),
			randReader: bytes.NewBufferString("X"),
			wantErr:    "unexpected EOF",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cipher, err := NewAES256SHA512(context.Background(), tt.key)
			if err != nil {
				t.Fatal(err)
			}

			cipher.(*aes256Sha512).randReader = tt.randReader

			sealed, err := cipher.Seal(tt.input)
			if err != nil && err.Error() != tt.wantErr ||
				err == nil && tt.wantErr != "" {
				t.Fatal(err)
			}

			if !bytes.Equal(tt.wantSealed, sealed) {
				t.Error(hex.EncodeToString(sealed))
			}
		})
	}
}
