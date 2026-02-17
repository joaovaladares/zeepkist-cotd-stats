package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joaovaladares/cotd-tracker/internal/core"
	"github.com/joho/godotenv"
)

const (
	csvDir = "csv"

	mapIndexUrl = "https://docs.google.com/spreadsheets/d/1rqtVPKeDxEaBfbNl7whL5HjmzeqIDk8mj3xOtfACyeE/export?format=csv&gid=1863499630"

	workshopBrowseURL = "https://api.steampowered.com/IPublishedFileService/QueryFiles/v1/"
	zeepkistAppId     = "1440670"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func searchWorkshop(appID, apiKey, text string) {
	u, err := url.Parse(workshopBrowseURL)
	check(err)
	q := u.Query()
	q.Set("key", apiKey)
	q.Set("querytype", "12")
	q.Set("cursor", "*")
	q.Set("appid", appID)
	q.Set("search_text", text)
	q.Set("filetype", "0")
	q.Set("numberperpage", "10")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	check(err)
	defer resp.Body.Close()

	var data any
	err = json.NewDecoder(resp.Body).Decode(&data)
	check(err)

	fmt.Printf("%+v\n", data)
}

func downloadCsvAndSaveFile(url string, fileName string) {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	file, err := os.Create(fmt.Sprintf("%s/%s.csv", csvDir, fileName))
	check(err)
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	steamApiKey := os.Getenv("STEAM_API_KEY")

	downloadCsvAndSaveFile(mapIndexUrl, "map_index")

	records := readCsvFile("csv/map_index.csv")

	parsedRecords := core.ParseMapIndexRecs(records)
	fmt.Println(parsedRecords[100])
	searchWorkshop(zeepkistAppId, steamApiKey, parsedRecords[100].Name)
	// for _, rec := range parsedRecords {
	// 	fmt.Println(rec)
	// }
}
