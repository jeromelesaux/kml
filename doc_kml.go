package kml

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type DocKml struct {
	Kml       *Kml
	IconsPath []string
	Kmlname   string
}

type zipArchive struct {
	Name string
	Body string
}

func NewDocKml(kmlname string) *DocKml {
	return &DocKml{Kml: NewKML("", 0),
		Kmlname: kmlname}
}

func (k *DocKml) Save(filename string) error {
	var err error

	defer func() {
		if err != nil {
			fmt.Fprint(os.Stderr, "Error with message "+err.Error())
			return
		}
	}()
	content, err := xml.Marshal(k.Kml)
	if err != nil {
		return err
	}

	archive, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer archive.Close()
	zip := zip.NewWriter(archive)

	kmlContent := string(content)
	var files = []zipArchive{{"doc.xml", kmlContent}}

	for _, icon := range k.IconsPath {
		f, err := os.Open(icon)
		if err != nil {
			return err
		}
		stats, err := f.Stat()
		fBytes := make([]byte, stats.Size())
		buffIcon := bufio.NewReader(f)
		buffIcon.Read(fBytes)
		filename := filepath.Base(icon)
		files = append(files, zipArchive{Name: filename, Body: string(fBytes)})
	}

	for _, file := range files {
		f, err := zip.Create(file.Name)
		if err != nil {
			return err
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			return err
		}
	}
	zip.Close()

	return nil
}
