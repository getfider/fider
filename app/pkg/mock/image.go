package mock

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"sync"
)

var uniformPNGCache sync.Map

// UniformPNG returns a valid 1-bit grayscale PNG of the given dimensions where every pixel is black.
// It is streamed row by row, so generating it never allocates the full image in memory.
// Large dimensions compress extremely well, which makes this useful for building
// "decompression bomb" images in tests (e.g. 12000x12000 is only a few KB).
// Results are cached, callers must not modify the returned slice.
func UniformPNG(width, height int) []byte {
	key := [2]int{width, height}
	if cached, ok := uniformPNGCache.Load(key); ok {
		return cached.([]byte)
	}
	content := uniformPNG(width, height)
	uniformPNGCache.Store(key, content)
	return content
}

func uniformPNG(width, height int) []byte {
	var idat bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&idat, zlib.BestCompression)
	row := make([]byte, 1+(width+7)/8) // filter byte (0 = none) + packed 1-bit pixels
	for i := 0; i < height; i++ {
		_, _ = zw.Write(row)
	}
	_ = zw.Close()

	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:4], uint32(width))
	binary.BigEndian.PutUint32(ihdr[4:8], uint32(height))
	ihdr[8] = 1 // bit depth
	ihdr[9] = 0 // color type: grayscale

	var out bytes.Buffer
	out.Write([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'})
	writePNGChunk(&out, "IHDR", ihdr)
	writePNGChunk(&out, "IDAT", idat.Bytes())
	writePNGChunk(&out, "IEND", nil)
	return out.Bytes()
}

func writePNGChunk(out *bytes.Buffer, name string, data []byte) {
	_ = binary.Write(out, binary.BigEndian, uint32(len(data)))
	crc := crc32.NewIEEE()
	_, _ = crc.Write([]byte(name))
	_, _ = crc.Write(data)
	out.WriteString(name)
	out.Write(data)
	_ = binary.Write(out, binary.BigEndian, crc.Sum32())
}

// GIFHeader returns the header of a GIF with the given logical screen dimensions (max 65535).
// It contains no image data: enough for image.DecodeConfig, but not for image.Decode.
func GIFHeader(width, height int) []byte {
	b := []byte("GIF89a")
	b = binary.LittleEndian.AppendUint16(b, uint16(width))
	b = binary.LittleEndian.AppendUint16(b, uint16(height))
	b = append(b, 0x80, 0, 0)             // global color table present (2 entries), bg color, aspect ratio
	b = append(b, 0, 0, 0, 255, 255, 255) // global color table
	return b
}
