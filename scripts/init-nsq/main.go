package main

import (
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
		log.Println("ERROR: nsq schema directory not supplied")
		os.Exit(1)
	}
	dir := os.Args[1]
	nsqHost := os.Getenv("NSQ_HOST")
	if nsqHost == "" {
		nsqHost = "http://localhost:4151"
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("ERROR: unable to walk dir:", err)
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		topicName := strings.TrimSuffix(info.Name(), ".json")
		req, _ := http.NewRequest(http.MethodPost, nsqHost+"/topic/create?topic="+topicName, nil)
		req.Header.Set("content-type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("ERROR: fail in creating topic based on schema %s: %v", path, err)
			return nil
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			b, _ := ioutil.ReadAll(res.Body)
			log.Printf("ERROR: fail in creating topic based on schema %s: code=%d resp=%s", path, res.StatusCode, string(b))
			return nil
		}
		log.Printf("INFO: topic %s created!", topicName)

		return nil
	})

	if err != nil {
		fmt.Println("ERROR: while walking directory:", err)
	}

}
