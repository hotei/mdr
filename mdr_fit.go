// mdr_fit.go

package mdr

import (
	"time"
)

func Fit2Time(t uint32) time.Time {
	return time.Unix(int64(t)+631065600, 0)
}

var Fit_crc_table [16]uint16 = [16]uint16{
	0x0000, 0xCC01, 0xD801, 0x1400, 0xF001, 0x3C00, 0x2800, 0xE401,
	0xA001, 0x6C00, 0x7800, 0xB401, 0x5000, 0x9C01, 0x8801, 0x4400,
}

func FitCRC_Get16(crc uint16, bite byte) uint16 {

	var tmp uint16

	// compute checksum of lower four bits of byte
	tmp = Fit_crc_table[crc&0xF]
	crc = (crc >> 4) & 0x0FFF
	crc = crc ^ tmp ^ Fit_crc_table[bite&0xF]

	// now compute checksum of upper four bits of byte
	tmp = Fit_crc_table[crc&0xF]
	crc = (crc >> 4) & 0x0FFF
	crc = crc ^ tmp ^ Fit_crc_table[(bite>>4)&0xF]

	return crc
}

func FitCRC_Calc16(data []byte, size uint32) uint16 {
	return FitCRC_Update16(0, data, size)
}

func FitCRC_Update16(crc uint16, data_ptr []byte, size uint32) uint16 {

	for i := 0; i < int(size); i++ {
		crc = FitCRC_Get16(crc, data_ptr[i])
	}

	return crc
}

func Float16To64(v int16) float64 {
	return float64(v) * (256.0 / 65536.0)
}

// ============================================================================  SemiT
type SemiT int32

func (semi SemiT) Degrees() float64 {
	rv := (float64(semi) / float64(0x7fffffff)) * 180
	if rv > 180.0 {
		rv -= 360
	}
	return rv
}

//
func DegreeFromSemicircle(semi int32) float64 {
	rv := (float64(semi) / float64(0x7fffffff)) * 180
	if rv > 180.0 {
		rv -= 360
	}
	return rv
}
