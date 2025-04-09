package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"

    "github.com/Anandhu3301/Websrapergo/internalValues"

)

func main() {
	var links []string
	var webaddress string
	var wg sync.WaitGroup
	pattern := regexp.MustCompile(`^(https?:\/\/)?([\w-]+\.)+[\w-]+(\/[\w-.\/]*)*$`)
	fmt.Printf("Enter the url: \n")
	fmt.Scan(&webaddress)
	u, err := url.Parse(webaddress)
	if err != nil && u.Host != "" && u.Scheme != "" {
		log.Fatal("Dont try to fool with me you sucker👿")
	}

	res, err := http.Get(webaddress)
	if err != nil {
		log.Fatal("Error")
	}
	defer res.Body.Close()
	if res.StatusCode != internalValues.StatuscodeOk {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		links = append(links, href)
	})
	var httpsresult []string = validLinkChecker(&links, pattern)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Link", "Status"})

	for i := internalValues.LoopStarter; i < len(httpsresult); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer res.Body.Close()
			res, _ := http.Get(httpsresult[i])
			table.Append([]string{httpsresult[i], res.Status})
		}(i)
	}
	wg.Wait()
	table.Render()
}

func validLinkChecker(urlLinks *[]string, pattern *regexp.Regexp) []string {
	var wg sync.WaitGroup
	validLinks := make(chan string)
	var validLinkCollection []string

	for _, link := range *urlLinks {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			if pattern.Match([]byte(link)) {
				validLinks <- link
			}
		}(link)
	}
	go func() {
		wg.Wait()
		close(validLinks)
	}()

	for link := range validLinks {
		validLinkCollection = append(validLinkCollection, link)
	}
	return validLinkCollection
}
