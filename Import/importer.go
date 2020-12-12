package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ReadVogelsFile ...
func ReadVogelsFile(filename string) []WaarnemingRecord {
	fileBytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sliceData := strings.Split(string(fileBytes), "\r\n")

	vogels := make([]WaarnemingRecord, len(sliceData))
	for i, species := range sliceData {
		names := strings.Split(species, "\t")
		vogels[i] = WaarnemingRecord{NederlandseNaam: names[0], OfficieleNaam: names[1]}
	}

	return vogels
}

func main() {
	waarnemingRecords := ReadVogelsFile("C:\\Users\\jonat\\Documents\\Passerine\\Import\\Nederlandse Vogelsoorten.txt")
	for i := range waarnemingRecords {
		waarnemingRecords[i].GetWaarnemingpuntnlRecord()
		waarnemingRecords[i].GetPhotos()
	}
	fmt.Println(waarnemingRecords)
}
