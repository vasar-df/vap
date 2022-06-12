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

// WriteFile writes the given Vap to the given path.
func WriteFile(path string, v *Vap) error {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	w := bufio.NewWriter(f)
	defer func() {
		_ = w.Flush()
		_ = f.Close()
	}()
	return Write(w, v)
}

// Write writes the given Vap to the given byte buffer.
func Write(b interface {
	io.Writer
	io.ByteWriter
}, v *Vap) error {
	w := protocol.NewWriter(b, 0)

	w.Uint16(&currentVersion)

	gamesLen := uint8(len(v.games))
	w.Uint8(&gamesLen)
	for i := range v.games {
		w.Uint8(&v.games[i])
	}

	w.String(&v.name)

	for i := 0; i < 3; i++ {
		w.Uint64((*uint64)(unsafe.Pointer(&v.positions[i][0])))
		w.Uint64((*uint64)(unsafe.Pointer(&v.positions[i][1])))
		w.Uint64((*uint64)(unsafe.Pointer(&v.positions[i][2])))
	}

	w.Uint32(&v.bounds[0])
	w.Uint32(&v.bounds[1])
	w.Uint32(&v.bounds[2])

	paletteLen := byte(len(v.palette))
	w.Uint8(&paletteLen)
	for _, s := range v.palette {
		w.String(&s.name)
		w.NBT(&s.properties, nbt.LittleEndian)
	}

	w.Bytes(&v.data)
	return nil
}
