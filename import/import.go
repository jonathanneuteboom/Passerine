package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadVogelFile(filename string) {
	fileBytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sliceData := strings.Split(string(fileBytes), "\n")

	for i, species := range sliceData {
		names := strings.Split(species, "\t")

	}

}

func main() {
	vogels := ReadVogelFile("nederlandse vogels.txt")
}
