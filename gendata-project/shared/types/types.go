// types.go
package types

import "time"

const (
	HeaderLength = 12
	PayloadSize  = 5

	// Типы данных
	TypeAnalog   = 0x10
	TypeFloat    = 0x30
	TypeFixed    = 0x20
	TypeInt      = 0x40
	TypeDiscrete = 0x50

	// Маски и сдвиги для GID
	RegionMask  = 0xE0
	RegionShift = 5
	RegionCheck = 0x07
	NppMask     = 0x18
	NppShift    = 3
	NppCheck    = 0x03
	EblockMask  = 0x07
	EblockShift = 0
	EblockCheck = 0x07
)

// Region один из 5 возможных для размещения станции
const (
	RegionWest  = 0x00 // Западный регион
	RegionCentr = 0x01 // Центральный регион
	RegionEast  = 0x02 // Восточный регион
	RegionSouth = 0x03 // Южный регион
	RegionNorth = 0x04 // Северный регион
)

var DataTypeSymbols = map[byte]rune{
	TypeAnalog:   'A', // Аналоговый тип данных
	TypeFloat:    'F', // Тип данных с плавающей точкой
	TypeFixed:    'I', // Фиксированный тип данных
	TypeInt:      'B', // Целочисленный тип данных
	TypeDiscrete: 'D', // Дискретный тип данных
}

// RegionString - строковые представления регионов.
var RegionString = map[byte]string{
	RegionWest:  "Region-West",
	RegionCentr: "Region-Centr",
	RegionEast:  "Region-East",
	RegionSouth: "Region-South",
	RegionNorth: "Region-North",
}

// BlockData представляет block of data в СПД.
// type BlockData struct {
type BlockHeader struct {
	GID byte // Group ID (мл.3 бита - номер энергоблока, ст.5 бит - группа)
	// [7-6 bits - регион, 5-3 bits - номер АЭС, 2-0 bits - Энергоблок]
	DType     byte      // Тип данных (0x00-аналог,0x10-float,0x20-fixed,0x30-int,0x40-discrete)
	Nlist     byte      // Номер списка (0-255)  		(но меньше 255, а 255 - это error)
	Vlist     byte      // Версия списка (0-255) 		(но гораздо меньше 255, а 255 - это error)
	Index     int16     // Индекс в списке (0-65535) 	(но меньше 65535, а 65535 - это error)
	NumParams int16     // Колич.дан. в блоке (0-65535) (но меньше 65535, а 65535 - это error)
	Timestamp time.Time // Время в формате Unix 		(секунды с 1970-01-01T00:00:00Z)
}

// DataPacket представляет блок данных передаваемый по сети
type DataPacket struct {
	Header BlockHeader
	Data   []byte
}
