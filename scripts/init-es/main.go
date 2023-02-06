package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("ERROR: schema directory not supplied")
		os.Exit(1)
	}
	dir := os.Args[1]
	elasticHost := os.Getenv("Elastic_HOST")
	if elasticHost == "" {
		elasticHost = "http://localhost:9200/"
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("ERROR: unable to walk dir:", err)
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		value, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("ERROR: read file:", err)
			return nil
		}

		filename := strings.TrimSuffix(info.Name(), ".json")
		req, _ := http.NewRequest(http.MethodPut, elasticHost+filename+"?pretty", bytes.NewBuffer(value))
		req.Header.Add("content-type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("ERROR: fail in storing index based on schema %s: %v", path, err)
			return nil
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(res.Body)
			log.Printf("ERROR: fail in storing index based on schema %s: code=%d resp=%s", path, res.StatusCode, string(b))
			return nil
		}
		log.Printf("INFO: index for file %s is created!", filename)

		return nil
	})

	if err != nil {
		fmt.Println("ERROR: while walking directory:", err)
	}
}
