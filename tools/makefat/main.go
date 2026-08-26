// Command makefat fuses per-architecture macOS Mach-O binaries into a single
// universal ("fat") binary, so one artifact runs on both Intel and Apple
// Silicon. It reads each thin binary's Mach-O header for its CPU type, then
// writes the fat wrapper by hand — pure Go, so it runs on any build host with no
// macOS and no `lipo`.
//
// Usage: makefat <output> <thin-macho> <thin-macho> [...]
//
// [impl->component~build-tooling~1]
package main

import (
	"bytes"
	"debug/macho"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintln(os.Stderr, "usage: makefat <output> <thin-macho> <thin-macho> [...]")
		os.Exit(2)
	}
	if err := fuse(os.Args[1], os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "makefat: %v\n", err)
		os.Exit(1)
	}
}

// Fat-binary on-disk layout: a big-endian fat_header (magic + count) followed by
// one fat_arch record per slice, then the slices themselves at aligned offsets.
const (
	fatMagic  = 0xCAFEBABE
	fatHeader = 8  // sizeof(fat_header)
	fatArch   = 20 // sizeof(fat_arch)
	alignBits = 14 // 2^14 = 16 KiB — satisfies Apple Silicon's page size
)

type slice struct {
	data   []byte
	cpu    macho.Cpu
	subCpu uint32
	offset uint32
}

func fuse(out string, inputs []string) error {
	var slices []slice
	for _, in := range inputs {
		data, err := os.ReadFile(in)
		if err != nil {
			return err
		}
		f, err := macho.NewFile(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("%s: not a Mach-O binary: %w", in, err)
		}
		slices = append(slices, slice{data: data, cpu: f.Cpu, subCpu: f.SubCpu})
		f.Close()
	}

	align := uint32(1) << alignBits
	pos := uint32(fatHeader + fatArch*len(slices))
	for i := range slices {
		pos = (pos + align - 1) &^ (align - 1)
		slices[i].offset = pos
		pos += uint32(len(slices[i].data))
	}

	buf := make([]byte, pos)
	be := binary.BigEndian
	be.PutUint32(buf[0:], fatMagic)
	be.PutUint32(buf[4:], uint32(len(slices)))
	for i, s := range slices {
		o := fatHeader + fatArch*i
		be.PutUint32(buf[o+0:], uint32(s.cpu))
		be.PutUint32(buf[o+4:], s.subCpu)
		be.PutUint32(buf[o+8:], s.offset)
		be.PutUint32(buf[o+12:], uint32(len(s.data)))
		be.PutUint32(buf[o+16:], alignBits)
		copy(buf[s.offset:], s.data)
	}
	return os.WriteFile(out, buf, 0o755)
}
