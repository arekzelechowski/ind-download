package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

const (
	wkhtmltopdf = "/usr/local/bin/wkhtmltopdf"

	urlFmt0 = "https://indd.adobe.com/view/publication/%s/1/publication.html"
	urlFmt1 = "https://indd.adobe.com/view/publication/%s/1/publication-%d.html"
)

func main() {
	err := os.Mkdir("out", 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(err)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4)
	for id := range readLines(f) {
		url := ""
		log.Printf("%s: Collecting HTML pages.\n", id)
		argList := []string{"-s", "A5", "-O", "landscape"}
		var page int
		for {
			if page == 0 {
				url = fmt.Sprintf(urlFmt0, id)
			} else {
				url = fmt.Sprintf(urlFmt1, id, page)
			}
			rsp, err := http.Get(url)
			if err != nil || rsp.StatusCode != 200 {
				log.Printf("%s: Last page found.\n", id)
				break
			}
			log.Printf("%s: Including page %d.\n", id, page)
			argList = append(argList, url)
			page++
		}
		if len(argList) == 2 {
			log.Printf("%s: Nothing to download.\n", id)
			continue
		}
		argList = append(argList, fmt.Sprintf(`out/%s.pdf`, id))
		sem <- struct{}{}
		wg.Add(1)
		log.Printf("%s: Downloading pages and converting to PDF.\n", id)
		go func(id string, argList []string) {
			err := exec.Command(wkhtmltopdf, argList...).Run()
			if err != nil {
				log.Printf("%s: PDF conversion failed. Reason: %s\n", id, err)
			}
			wg.Done()
			<-sem
		}(id, argList)
	}
	wg.Wait()
}

func readLines(f *os.File) <-chan string {
	ch := make(chan string, 16)
	go func() {
		r := bufio.NewReader(f)
		for {
			url, err := r.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			ch <- url[:len(url)-2]
		}
	}()
	return ch
}
