// sensor.go
package sensorvalue

import "gendata-project/shared/types"

const (
	DataTypeMask    byte = 0x70 // Маска для определения типа данных
	DataValidMask   byte = 0x01 // для получения бита валидности данных
	DataShiftMask   byte = 0x03 // Сдвиг бита валидности
	DataShiftStatus byte = 0x07 // Младшие 3 бита пасспорта - статус данных
)

// Структура SensorValue для хранения данных сенсоров
type SensorValue struct {
	Pasport  byte
	RawValue [4]byte
}

// Возвращает тип данных из Pasport
func (sv SensorValue) DataType() byte {
	return sv.Pasport & DataTypeMask
}

// Получение символа типа данных по значению кода типа из Pasport
func (sv *SensorValue) TypeSymbol() rune {
	return types.DataTypeSymbols[sv.DataType()]
}

func (sv SensorValue) Validation() string {
	switch (sv.Pasport >> DataShiftMask) & DataValidMask {
	case 0:
		return "Valid"
	case 1:
		return "Invalid"
	default:
		return "Unknown"
	}
}

func (sv SensorValue) ShiftStatus() byte {
	return sv.Pasport & DataShiftStatus
}
