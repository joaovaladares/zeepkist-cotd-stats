package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	csvDir = "csv"

	mapIndexUrl = "https://docs.google.com/spreadsheets/d/1rqtVPKeDxEaBfbNl7whL5HjmzeqIDk8mj3xOtfACyeE/export?format=csv&gid=1863499630"
	// map-index csv columns
	leftEvent  = 0
	leftMapper = 2
	leftName   = 3

	rightEvent  = 8
	rightMapper = 10
	rightName   = 11
)

type MapRecord struct {
	Event  string
	Mapper string
	Name   string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func downloadCsvAndSaveFile(url string, fileName string) {
	resp, err := http.Get(url)
	if err != nil {
		check(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s/%s.csv", csvDir, fileName))
	if err != nil {
		check(err)
	}
	defer file.Close()

	io.Copy(file, resp.Body)
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func parseRecords(records [][]string) []MapRecord {
	var out []MapRecord

	for _, row := range records {
		if isDataLeft(row) {
			out = append(out, MapRecord{
				Event:  valueAt(row, leftEvent),
				Mapper: valueAt(row, leftMapper),
				Name:   valueAt(row, leftName),
			})
		}
		if isDataRight(row) {
			out = append(out, MapRecord{
				Event:  valueAt(row, rightEvent),
				Mapper: valueAt(row, rightMapper),
				Name:   valueAt(row, rightName),
			})
		}
	}

	return out
}

func valueAt(row []string, idx int) string {
	if idx >= 0 && idx < len(row) {
		return row[idx]
	}
	return ""
}

func isHeaderRow(row []string) bool {
	return valueAt(row, leftEvent) == "COTD #" &&
		valueAt(row, leftMapper) == "Mapper" &&
		valueAt(row, leftName) == "Map Name"
}

func isDataLeft(row []string) bool {
	return valueAt(row, leftEvent) != "" &&
		valueAt(row, leftMapper) != "" &&
		valueAt(row, leftName) != "" &&
		!isHeaderRow(row)
}

func isDataRight(row []string) bool {
	return valueAt(row, rightEvent) != "" &&
		valueAt(row, rightMapper) != "" &&
		valueAt(row, rightName) != "" &&
		!isHeaderRow(row)
}

func main() {
	// Download mapIndex CSV
	downloadCsvAndSaveFile(mapIndexUrl, "map_index")

	records := readCsvFile("csv/map_index.csv")

	parsedRecords := parseRecords(records)
	for _, rec := range parsedRecords {
		fmt.Println(rec)
	}
}
