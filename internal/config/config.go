package config

import (
	"encoding/json"
	//"fmt"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
)

type Target struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func ParseTargets(filepath string) ([]Target, error) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("error while opening %v: %v", filepath, err)
		return nil, fmt.Errorf("can't open file: %w", err)
	}

	defer file.Close()
	// var a []byte
	// _, err = file.Read(a)
	// if err != nil {
	// 	return nil
	// }
	// fmt.Println(a)
	decoder := json.NewDecoder(file)

	_, err = decoder.Token()

	if err == io.EOF {
		log.Println("empty file")
		return nil, fmt.Errorf("can't parse data: %w", err)
	}
	if err != nil {
		log.Printf("File: some problem with tokens: %v", err)
		return nil, fmt.Errorf("can't parse file: %w", err)
	}

	res := make([]Target, 0)

	for decoder.More() {
		var trgt Target
		err = decoder.Decode(&trgt)
		if err != nil {
			return nil, fmt.Errorf("file: error in decoding: %w", err)
		}
		URL, err := url.Parse(trgt.URL)
		if err != nil {
			log.Printf("Incorrect URL, skipping: %v\n", err)
			continue
		}
		if URL.Scheme == "" {
			log.Println("Incorrect URL: no schema, skipping")
			continue
		}

		res = append(res, trgt)

	}

	return res, nil

}
