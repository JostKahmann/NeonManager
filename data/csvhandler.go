package data

import (
	"encoding/csv"
	"fmt"
	"os"
)

func parseCsv[T any](file *os.File, parse func([]string) T) []T {
	reader := csv.NewReader(file)
	items := make([]T, 0)
	for record, err := reader.Read(); err == nil && record != nil; record, err = reader.Read() {
		items = append(items, parse(record))
	}
	return items
}

func writeCsv[T any](file *os.File, stringify func(T) []string, items []T) error {
	writer := csv.NewWriter(file)
	for _, item := range items {
		itemStr := stringify(item)
		err := writer.Write(itemStr)
		if err != nil {
			return fmt.Errorf("error writing item %v csv: %w", itemStr, err)
		}
	}
	return nil
}
