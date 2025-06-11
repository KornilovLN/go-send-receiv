// generator.go
package generator

import (
	"encoding/binary"
	"math"
	"time"

	//"gendata-project/shared/sensorvalue"
	"gendata-project/shared/types"
)

func GenAnalog(value float32, shift byte) []byte {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], math.Float32bits(value))
	return append([]byte{types.TypeAnalog | shift}, bytes[:]...)
}

func GenFloat(value float32, shift byte) []byte {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], math.Float32bits(value))
	return append([]byte{types.TypeFloat | shift}, bytes[:]...)
}

func GenInt(value int32, shift byte) []byte {
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], uint32(value))
	return append([]byte{types.TypeInt | shift}, bytes[:]...)
}

func GenFixed(value float32, shift byte) []byte {
	scaled := int32(math.Round(float64(value) * math.Pow10(int(shift))))
	var bytes [4]byte
	binary.LittleEndian.PutUint32(bytes[:], uint32(scaled))
	return append([]byte{types.TypeFixed | shift}, bytes[:]...)
}

func GenDiscrete(b1, b2, b3, b4 byte) []byte {
	return []byte{
		types.TypeDiscrete | 0x00,
		b1, b2, b3, b4,
	}
}

func GIDCreater(Region byte, Nnpp byte, Eblock byte) byte {
	var GID byte = 0x00
	GID = (Region & types.RegionCheck) << types.RegionShift
	GID |= (Nnpp & types.NppCheck) << types.NppShift
	GID |= (Eblock & types.EblockCheck) << types.EblockShift
	return GID
}

func HeadBlockCreater(
	GID byte,
	DataType byte,
	ListNum byte,
	ListVer byte,
	Index int16,
	NumParams int16,
	Timestamp time.Time,
) []byte {
	header := make([]byte, types.HeaderLength)

	header[0] = GID
	header[1] = DataType
	header[2] = ListNum
	header[3] = ListVer

	binary.LittleEndian.PutUint16(header[4:6], uint16(Index))
	binary.LittleEndian.PutUint16(header[6:8], uint16(NumParams))

	unixTime := uint32(Timestamp.Unix())
	binary.LittleEndian.PutUint32(header[8:12], unixTime)

	return header
}

func GenBlock(header []byte) []byte {
	if len(header) != types.HeaderLength {
		panic("Header должен быть 12 байт")
	}

	numParams := int16(binary.LittleEndian.Uint16(header[6:8]))
	packet := make([]byte, 0, types.HeaderLength+int(numParams)*types.PayloadSize)
	packet = append(packet, header...)

	for i := int16(0); i < numParams; i++ {
		var param []byte

		switch header[1] {
		case types.TypeAnalog:
			param = GenAnalog(float32(i)-486.57, 0x05)
		case types.TypeFloat:
			param = GenFloat(float32(i)+3.1415927, 0x00)
		case types.TypeFixed:
			param = GenFixed(float32(i)+59.27, 0x02)
		case types.TypeInt:
			param = GenInt(int32(i)+54321, 0x00)
		case types.TypeDiscrete:
			val := int16(10)
			b1 := byte((val + i) % 256)
			b2 := byte(((val + i) + 1) % 256)
			b3 := byte(((val + i) + 2) % 256)
			b4 := byte(((val + i) + 3) % 256)
			param = GenDiscrete(b1, b2, b3, b4)
		default:
			param = make([]byte, types.PayloadSize)
		}

		packet = append(packet, param...)
	}

	return packet
}
