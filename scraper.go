package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var asuraUrl = "https://asuratoon.com/"

func request(url string) (*goquery.Document, error) {
	cl := getClient(time.Second * 10)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", opts.UserAgent)

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad HTML status code: %v", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

type Latest struct {
	Link  string
	Title string
}

func getImages(url string) ([]string, error) {
	doc, err := request(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to make request: %v", err)
	}

	var images []string
	doc.Find("div#readerarea.rdminimal img").Each(func(i int, s *goquery.Selection) {
		im, ok := s.Attr("src")
		if ok {
			images = append(images, im)
		}
	})
	return images, nil
}

func getChapterList(url string) (map[string]string, error) {
	doc, err := request(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to make request: %v", err)
	}
	// fi, err := os.Open("out.html")
	// defer fi.Close()

	// doc, err := goquery.NewDocumentFromReader(fi)
	// if err != nil {
	// 	return nil, fmt.Errorf("Unable to make request: %v", err)
	// }

	m := make(map[string]string)

	// html body.darkmode div#content div.wrapper div.postbody article#post-40871.post-40871.hentry div.bixbox.bxcl.epcheck div#chapterlist.eplister ul.clstyle li div.chbox div.eph-num a span.chapternum
	doc.Find(".clstyle a").Each(func(i int, s *goquery.Selection) {
		link, ok := s.Attr("href")
		title := s.Find("span.chapternum").Text()

		if ok {
			if title != "" {
				m[title] = link
			} else {
				m[link] = link
			}
		}
	})

	if m == nil || len(m) == 0 {
		return nil, fmt.Errorf("Couldnt parse chapter list")
	}

	return m, nil
}

func getLatest() (map[string]string, error) {
	doc, err := request(asuraUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to make request: %v", err)
	}

	m := make(map[string]string)

	doc.Find("div.utao.styletwo div.uta div.luf a.series").Each(func(i int, s *goquery.Selection) {
		// div.utao:nth-child(1) > div:nth-child(1) > div:nth-child(2) > a:nth-child(1)
		// html body.darkmode div#content div.wrapper div.postbody div.bixbox div.listupd div.utao.styletwo div.uta div.luf a.series
		l := Latest{}
		link, exists := s.Attr("href")
		if exists {
			l.Link = link
		}

		title, exists1 := s.Attr("title")
		if exists1 {
			l.Title = title
		}

		if l.Link != "" && l.Title != "" {
			m[l.Title] = l.Link
		}

	})

	return m, nil
}
