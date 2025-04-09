package internalValues

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
)

func internalSolution1() {
	var links []string
	var webaddress string
	var wg sync.WaitGroup
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
	if res.StatusCode != StatuscodeOk {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		links = append(links, href)
	})
	urlChecker(&links)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Link", "Status"})

	for i := LoopStarter; i < len(links); i++ {
		wg.Add(1)
		go func(i int) {
			defer res.Body.Close()
			defer wg.Done()
			res, _ := http.Get(links[i])
			table.Append([]string{links[i], res.Status})
		}(i)
	}
	wg.Wait()
	table.Render()

}

func urlChecker(urlLinks *[]string) {
	pattern := regexp.MustCompile(`^(https?:\/\/)?([\w-]+\.)+[\w-]+(\/[\w-.\/]*)*$`)
	for i := 0; i < len(*urlLinks); i++ {
		if !pattern.Match([]byte((*urlLinks)[i])) {
			*urlLinks = append((*urlLinks)[:i], (*urlLinks)[i+1:]...)
			i--
		}
	}
}
