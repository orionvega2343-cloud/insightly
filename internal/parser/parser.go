package parser

import (
	"encoding/csv"
	"insightly/internal/errs"
	"os"
	"strings"
)

func CsvParser(filepath string) (string, error) {
	//Открытие файла
	open, err := os.Open(filepath)
	if err != nil {
		return "", errs.ErrorOpenFile
	}

	//Ожидание закрытия
	defer open.Close()

	//Чтение файла
	read, err := csv.NewReader(open).ReadAll()
	if err != nil {
		return "", errs.ErrorReadFile
	}

	if len(read) <= 0 {
		return "", errs.ErrorEmptyFile
	}

	//Ответ
	var builder strings.Builder
	for _, v := range read {
		builder.WriteString(strings.Join(v, ","))
		builder.WriteString("\n")
	}
	return builder.String(), nil
}
