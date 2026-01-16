package display

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"gendata-project/receiver-service/internal/parser"
	"gendata-project/shared/sensorvalue"
	"gendata-project/shared/types"
)

const LINELEN = 80

func PrintHeaderLine(title string, char rune, length int) string {
	if len(title) >= length-4 {
		return strings.Repeat(string(char), length)
	}

	padding := (length - len(title) - 2) / 2
	leftPad := strings.Repeat(string(char), padding)
	rightPad := strings.Repeat(string(char), length-len(title)-2-padding)

	return leftPad + " " + title + " " + rightPad
}

func PrintHeader(h types.BlockHeader) {
	fmt.Println(PrintHeaderLine("Заголовок блока", '=', LINELEN))

	fmt.Printf("GID: 0x%02X => ", h.GID)
	fmt.Println(parser.GetGIDString(h.GID))

	fmt.Printf("Dtype: %c\n", types.DataTypeSymbols[h.DType])
	fmt.Printf("Nlist: %d\n", h.Nlist)
	fmt.Printf("Vlist: %d\n", h.Vlist)
	fmt.Printf("Index: %d\n", h.Index)
	fmt.Printf("N prm: %d\n", h.NumParams)
	fmt.Printf("%s\n", h.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

func PrintParameter(p sensorvalue.SensorValue, index int) {
	fmt.Printf("[%3d]  ", index+1)
	fmt.Printf("Psp 0x%02X = {", p.Pasport)
	fmt.Printf("%c ", types.DataTypeSymbols[p.DataType()])
	fmt.Printf("%v ", p.Validation())
	fmt.Printf("%d}  ", p.ShiftStatus())
	fmt.Printf("raw [% X]  ", p.RawValue)

	switch p.DataType() {
	case types.TypeAnalog, types.TypeFloat:
		value := binary.LittleEndian.Uint32(p.RawValue[:])
		floatVal := math.Float32frombits(value)
		fmt.Printf("float32:(0x%08X) %.4f\n", value, floatVal)

	case types.TypeFixed:
		shift := p.Pasport & 0x07
		intVal := int32(binary.LittleEndian.Uint32(p.RawValue[:]))
		fixedVal := float32(intVal) / float32(math.Pow10(int(shift)))
		fmt.Printf("fixed:(shift=%d) %.2f\n", shift, fixedVal)

	case types.TypeInt:
		intVal := int32(binary.LittleEndian.Uint32(p.RawValue[:]))
		fmt.Printf("int32:(0x%08X) %d\n", uint32(intVal), intVal)

	case types.TypeDiscrete:
		fmt.Printf("discrete:[% X]\n", p.RawValue)

	default:
		fmt.Println("  Неизвестный тип данных")
	}
}

func DisplayDataBlock(data []byte) error {
	header, err := parser.ParseHeader(data)
	if err != nil {
		return fmt.Errorf("ошибка парсинга заголовка: %v", err)
	}

	PrintHeader(header)

	params, err := parser.ParseParameters(data, header.NumParams)
	if err != nil {
		return fmt.Errorf("ошибка парсинга параметров: %v", err)
	}

	stroka := fmt.Sprintf("Параметры (%d)", len(params))
	fmt.Println(PrintHeaderLine(stroka, '=', LINELEN))
	for i, param := range params {
		PrintParameter(param, i)
	}

	fmt.Println(strings.Repeat("-", LINELEN))
	return nil
}
