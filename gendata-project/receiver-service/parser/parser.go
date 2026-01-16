package parser

import (
	"encoding/binary"
	"fmt"
	"time"

	"gendata-project/shared/sensorvalue"
	"gendata-project/shared/types"
)

func ParseHeader(data []byte) (types.BlockHeader, error) {
	if len(data) < types.HeaderLength {
		return types.BlockHeader{}, fmt.Errorf("данные слишком короткие: %d байт", len(data))
	}

	return types.BlockHeader{
		GID:       data[0],
		DType:     data[1],
		Nlist:     data[2],
		Vlist:     data[3],
		Index:     int16(binary.LittleEndian.Uint16(data[4:6])),
		NumParams: int16(binary.LittleEndian.Uint16(data[6:8])),
		Timestamp: time.Unix(int64(binary.LittleEndian.Uint32(data[8:12])), 0),
	}, nil
}

func ParseParameters(data []byte, numParams int16) ([]sensorvalue.SensorValue, error) {
	params := make([]sensorvalue.SensorValue, 0, numParams)

	for i := 0; i < int(numParams); i++ {
		offset := types.HeaderLength + i*types.PayloadSize
		if offset+types.PayloadSize > len(data) {
			return nil, fmt.Errorf("недостаточно данных для параметра %d", i)
		}

		paramData := data[offset : offset+types.PayloadSize]
		params = append(params, sensorvalue.SensorValue{
			Pasport:  paramData[0],
			RawValue: [4]byte{paramData[1], paramData[2], paramData[3], paramData[4]},
		})
	}

	return params, nil
}

func GetRegionGID(GID byte) byte {
	return (GID & types.RegionMask) >> types.RegionShift
}

func GetNPPGID(GID byte) byte {
	return (GID & types.NppMask) >> types.NppShift
}

func GetEblockGID(GID byte) byte {
	return (GID & types.EblockMask) >> types.EblockShift
}

func GetGIDString(GID byte) string {
	region := GetRegionGID(GID)
	npp := GetNPPGID(GID)
	eblock := GetEblockGID(GID)

	return fmt.Sprintf("{%s, NPP: %d, Eblock: %d}",
		types.RegionString[region], npp, eblock)
}
