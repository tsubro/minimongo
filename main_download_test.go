package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const byteSize = 2048 //change the byte size accorning to the chunk size needed
var wg sync.WaitGroup

type CustomWriter struct {
	Total uint64
}

func (cw *CustomWriter) Write(p []byte) (int, error) {
	n := len(p)
	cw.Total += uint64(n)
	//log.Println("Total Bytes Read : ", cw.Total) //Logging to check if chunk is working fine or not
	return n, nil
}

func mainTest() {
	log.Println("Download Started !!!")
	links := []string{
		// "http://www.africau.edu/images/default/sample.pdf",
		// "http://www.pdf995.com/samples/pdf.pdf",
		"https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		"https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		"https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		"https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",

		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",

		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",

		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",

		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
		// "https://test-subro.s3.us-east-2.amazonaws.com/chinesefont.doc",
	}

	for i, link := range links {
		wg.Add(1)
		go downloader(i, link)
	}
	wg.Wait()
	log.Println("Download Completed !!!")
}

func downloader(fileNum int, URL string) {

	client := &http.Client{Transport: &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}}

	req, _ := http.NewRequest("GET", URL, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create("File_" + strconv.Itoa(fileNum) + ".doc")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		resp.Body.Close()
		out.Close()
		wg.Done()
	}()

	cw := &CustomWriter{}
	b := make([]byte, byteSize)
	if _, err = io.CopyBuffer(out, io.TeeReader(resp.Body, cw), b); err != nil {
		log.Fatal(err)
	}
}
