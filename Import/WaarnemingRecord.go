package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// WaarnemingRecord ...
type WaarnemingRecord struct {
	ID              int
	NederlandseNaam string
	OfficieleNaam   string
}

// GetPhotos ...
func (w *WaarnemingRecord) GetPhotos() {
	types := []string{
		"onbekend",

		"adult",
		"subadult",
		"juveniel",
		"pullus",
		"onvolwassen",
		"afwijkend",

		"adult zomerkleed",
		"adult winterkleed",
		"winterkleed",
		"zomerkleed",
		"eclips",

		"eerste jaar",
		"eerste kalenderjaar",
		"tweede kalenderjaar",
		"derde kalenderjaar",
		"vierde kalenderjaar",
		"eerste najaar",
		"eerste winter",
		"tweede winter",
		"derde winter",
		"eerste zomer",
		"tweede zomer",
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i] <= types[j]
	})

	url := fmt.Sprintf("https://waarneming.nl/species/%d/", w.ID)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("%s - %s, bestaat niet", w.OfficieleNaam, w.NederlandseNaam)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div[class=\"app-content-body\"]>h3").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		selector := fmt.Sprintf("div[class=\"app-content-body\"] > h3:nth-of-type(%d)+div[class=\"row\"] figure > a", i+1)
		doc.Find(selector).Each(func(i int, s2 *goquery.Selection) {
			path, exists := s2.Attr("href")
			if !exists {
				log.Fatal("href does not exist")
			}
			url := fmt.Sprintf("https://waarneming.nl%s", path)
			baseFolder := "C:\\Users\\jonat\\Documents\\Passerine\\Import\\Fotos"
			folderPath := fmt.Sprintf("%s\\%s\\%s", baseFolder, w.NederlandseNaam, title)
			filename := filepath.Base(path)
			err := DownloadFile(folderPath, filename, url)
			if err != nil {
				log.Fatal(err)
			}
		})
	})
}

// GetWaarnemingpuntnlRecord ...
func (w *WaarnemingRecord) GetWaarnemingpuntnlRecord() {
	url := fmt.Sprintf("https://waarneming.nl/species/search/?q=%s&species_group=1", strings.Replace(w.NederlandseNaam, " ", "+", -1))
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("div[class=\"app-content-body\"] table a[href^=\"/species/\"]").Each(func(i int, s *goquery.Selection) {
		title := s.Text()

		if !strings.HasPrefix(title, fmt.Sprintf("%s - ", w.NederlandseNaam)) {
			return
		}

		if i > 0 {
			fmt.Printf("%s heeft %de species: '%s", w.NederlandseNaam, i+1, title)
			return
		}

		href, exists := s.Attr("href")
		if !exists {
			log.Fatal("href bestaat niet")
		}

		_, err := fmt.Sscanf(href, "/species/%d/", &w.ID)
		if err != nil {
			log.Fatal("parse ging fout")
		}

		fmt.Println(w)
	})
}

// DownloadFile ...
func DownloadFile(folderPath, filename, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = os.Stat(folderPath)
	if os.IsNotExist(err) {
		os.MkdirAll(folderPath, os.ModePerm)
	}

	out, err := os.Create(filepath.Join(folderPath, filename))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
