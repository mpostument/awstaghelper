package pkg

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func WriteCsv(data [][]string, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}
}

func ReadCsv(filename string) [][]string {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Cannot open csv file", err)
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Fatal("Cannot read csv file", err)
	}

	return csvLines
}

func addHeaders(tagsToRead string, resourceIdHeader string) [][]string {
	var rows [][]string
	headers := []string{resourceIdHeader}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	return rows
}
