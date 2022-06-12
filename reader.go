package vap

import (
	"bufio"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"io"
	"os"
	"unsafe"
)

// ReadFile reads a Vap file from the given path.
func ReadFile(path string) (*Vap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()
	return Read(bufio.NewReader(f))
}

// Read reads a Vap from the given byte reader.
func Read(b interface {
	io.Reader
	io.ByteReader
}) (*Vap, error) {
	r := protocol.NewReader(b, 0)

	var version uint16
	r.Uint16(&version)
	if version != currentVersion {
		return nil, fmt.Errorf("vap: unsupported version: %v", version)
	}

	v := &Vap{}

	var gamesLen uint8
	r.Uint8(&gamesLen)

	v.games = make([]uint8, gamesLen)
	for i := range v.games {
		r.Uint8(&v.games[i])
	}

	r.String(&v.name)

	for i := 0; i < 3; i++ {
		r.Uint64((*uint64)(unsafe.Pointer(&v.positions[i][0])))
		r.Uint64((*uint64)(unsafe.Pointer(&v.positions[i][1])))
		r.Uint64((*uint64)(unsafe.Pointer(&v.positions[i][2])))
	}

	r.Uint32(&v.bounds[0])
	r.Uint32(&v.bounds[1])
	r.Uint32(&v.bounds[2])

	var paletteLen byte
	r.Uint8(&paletteLen)

	v.palette = make([]state, paletteLen)
	for i := range v.palette {
		r.String(&v.palette[i].name)
		r.NBT(&v.palette[i].properties, nbt.LittleEndian)
	}

	r.Bytes(&v.data)
	return v, nil
}
