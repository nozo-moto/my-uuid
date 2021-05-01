package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type UUID [16]byte // 128bit

func main() {
	rander := rand.Reader
	var uuid UUID
	_, err := io.ReadFull(rander, uuid[:])
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // 46~50bit にバージョン情報を挿入
	// input  | x x x x x x x x |
	// & 0x0f | 0 0 0 0 1 1 1 1 |
	// | 0x40 | 0 1 0 0 0 0 0 0 |
	// output | 0 1 0 0 x x x x |

	uuid[8] = (uuid[8] & 0x3f) | 0x80 // 65~66bitにvariantを設定
	// input  | x x x x x x x x |
	// & 0x0f | 0 0 1 1 1 1 1 1 |
	// | 0x40 | 1 0 0 0 0 0 0 0 |
	// output | 1 0 x x x x x x |

	fmt.Print(uuid)
}

func (uuid UUID) encodeHex(dst []byte) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

func (uuid UUID) String() string {
	var dst [36]byte
	uuid.encodeHex(dst[:])
	return string(dst[:])
}
