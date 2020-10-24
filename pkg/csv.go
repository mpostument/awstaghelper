package pkg

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

// WriteCsv writes the csv data to disk
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

// ReadCsv reads a csv file from disk
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

func addHeadersToCsv(tagsToRead string, resourceIDHeader string) [][]string {
	var rows [][]string
	headers := []string{resourceIDHeader}
	headers = append(headers, strings.Split(tagsToRead, ",")...)
	rows = append(rows, headers)
	return rows
}

func addTagsToCsv(tagsToRead string, tags map[string]string, rows [][]string, resourceID string) [][]string {
	var resultTags []string
	for _, key := range strings.Split(tagsToRead, ",") {
		resultTags = append(resultTags, tags[key])
	}
	rows = append(rows, append([]string{resourceID}, resultTags...))
	return rows
}
