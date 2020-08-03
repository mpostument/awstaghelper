package commonLib

import (
	"encoding/csv"
	"log"
	"os"
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
