package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const mapIndexUrl = "https://docs.google.com/spreadsheets/d/1rqtVPKeDxEaBfbNl7whL5HjmzeqIDk8mj3xOtfACyeE/export?format=csv&gid=1863499630"
const csvDir = "csv"

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
	var mapRecords []MapRecord
	// FIXME: Make it a bit more modular and avoid duplicated code
	for i, line := range records {
		if i > 4 && i < 15 { // Lines with a normal cup and troll cup at the same time
			var leftRec, rightRec MapRecord
			for j, field := range line {
				switch j {
				case 0:
					leftRec.Event = field
				case 2:
					leftRec.Mapper = field
				case 3:
					leftRec.Name = field
				case 8:
					rightRec.Event = field
				case 10:
					rightRec.Mapper = field
				case 11:
					rightRec.Name = field
				}
			}
			mapRecords = append(mapRecords, leftRec)
			mapRecords = append(mapRecords, rightRec)
		}
		if i > 15 { // Lines with no troll cups
			var rec MapRecord
			for j, field := range line {
				switch j {
				case 0:
					rec.Event = field
				case 2:
					rec.Mapper = field
				case 3:
					rec.Name = field
				}
			}
			mapRecords = append(mapRecords, rec)
		}
	}
	return mapRecords
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
