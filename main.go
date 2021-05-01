package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	mathRand "math/rand"
	"net"
	"time"
)

type UUID [16]byte // 128bit

func main() {
	uuidv1, err := NewV1()
	if err != nil {
		panic(err)
	}
	fmt.Println(uuidv1)

	uuidV4, err := NewV4()
	if err != nil {
		panic(err)
	}
	fmt.Print(uuidV4)
}

func NewV1() (uuid UUID, err error) {
	var gregorianBeginTime time.Time
	gregorianBeginTime, err = time.Parse(time.RFC3339, "1582-10-15T00:00:00Z")
	if err != nil {
		return
	}

	timestamp := time.Since(gregorianBeginTime).Nanoseconds() / 100
	// set time_low
	binary.BigEndian.PutUint32(uuid[0:], uint32(timestamp&0xffffffff))
	// set time_mid
	binary.BigEndian.PutUint16(uuid[4:], uint16((timestamp>>32)&0xffff))
	// set time_hi
	binary.BigEndian.PutUint16(uuid[6:], uint16((timestamp>>48)&0xffff))

	// set clk_seq_hi_res set clk_seq_low
	rander := rand.Reader
	var b [2]byte
	if _, err = io.ReadFull(rander, b[:]); err != nil {
		return
	}
	seq := int(b[0])<<8 | int(b[1])
	clockSeq := uint16(seq&0x3fff) | 0x8000
	binary.BigEndian.PutUint16(uuid[8:], clockSeq)

	// set macadd for 81~127
	macAddres, err := getMacAddr()
	if err != nil {
		return
	}
	node := randomChoiceByteArray(macAddres)[:]
	copy(uuid[10:], node[:])
	return
}

func NewV4() (uuid UUID, err error) {
	rander := rand.Reader
	_, err = io.ReadFull(rander, uuid[:])
	if err != nil {
		return
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

	return
}

func NewV6() (uuid UUID, err error) {
	var gregorianBeginTime time.Time
	gregorianBeginTime, err = time.Parse(time.RFC3339, "1582-10-15T00:00:00Z")
	if err != nil {
		return
	}

	timestamp := time.Since(gregorianBeginTime).Nanoseconds() / 100
	// set time_low
	binary.BigEndian.PutUint32(uuid[0:], uint32(timestamp&0xffffffff))
	// set time_mid
	binary.BigEndian.PutUint16(uuid[4:], uint16((timestamp>>32)&0xffff))
	// set time_hi_and_version
	timHighAndVersion := ((timestamp >> 48) & 0x0fff) | 0x6000
	binary.BigEndian.PutUint16(uuid[6:], uint16(timHighAndVersion))

	// set clk_seq_hi_res set clk_seq_low
	rander := rand.Reader
	var b [2]byte
	if _, err = io.ReadFull(rander, b[:]); err != nil {
		return
	}
	seq := int(b[0])<<8 | int(b[1])
	clockSeq := uint16(seq&0x3fff) | 0x8000
	binary.BigEndian.PutUint16(uuid[8:], clockSeq)

	// set macadd for 81~127
	macAddres, err := getMacAddr()
	if err != nil {
		return
	}
	node := randomChoiceByteArray(macAddres)[:]
	copy(uuid[10:], node[:])
	return
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

func getMacAddr() ([][]byte, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as [][]byte
	for _, ifa := range ifas {
		haddr := ifa.HardwareAddr
		if len(haddr) > 0 {
			as = append(as, haddr)
		}
	}
	return as, nil
}

func randomChoiceByteArray(arr [][]byte) []byte {
	mathRand.Seed(time.Now().Unix())
	return arr[mathRand.Intn(len(arr))]
}
