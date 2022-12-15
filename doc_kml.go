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
	var files = []zipArchive{{"doc.kml", kmlContent}}

	if len(k.IconsPath) > 0 {
		for _, icon := range k.IconsPath {
			f, err := os.Open(icon)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while opening file [%s]  error : %v\n", icon, err)
				continue
			}
			stats, err := f.Stat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error while getting stats on file [%s]  error : %v\n", icon, err)
				continue
			}
			fBytes := make([]byte, stats.Size())
			buffIcon := bufio.NewReader(f)
			_, err = buffIcon.Read(fBytes)
			if err != nil {
				return err
			}
			filename := filepath.Base(icon)
			files = append(files, zipArchive{Name: "files/" + filename, Body: string(fBytes)})
		}
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
	fmt.Fprintf(os.Stdout, "Zipping kml document [%s]\n", filename)
	zip.Close()
	fmt.Fprintf(os.Stdout, "Kml document zip ended [%s]\n", filename)

	return nil
}
