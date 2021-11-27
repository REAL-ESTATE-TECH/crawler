package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sync"
)

type Fetcher interface {
	// GetUrls returns a slice of URLs found on that page.
	GetUrls(u string) (urls []string, err error)
}

func (crawer *Crawler) GetUrls(u string) (urls []string, err error) {
	var (
		urlParsed *url.URL
		body      string
	)
	if urlParsed, err = url.Parse(u); err != nil {
		return nil, err
	}
	if body, err = getBody(u); err != nil {
		return nil, err
	}
	urls, err = parseUrls(urlParsed, body)
	return urls, err
}

func parseUrls(u *url.URL, body string) (urls []string, err error) {
	regularExp := regexp.MustCompile("<a.*?href=\"(.*?)\"")
	matches := regularExp.FindAllStringSubmatch(body, -1)
	for _, val := range matches {
		var href *url.URL
		if href, err = url.Parse(val[1]); err != nil {
			continue
		}
		uri, _ := url.QueryUnescape(href.String())
		if href.IsAbs() {
			urls = append(urls, uri)
		} else {
			if len(uri) > 0 && uri[0:1] != "/" { // some hrefs are missing a leading front slash
				uri = "/" + uri
			}
			urls = append(urls, u.Scheme+"://"+u.Host+uri)
		}
	}
	return urls, nil
}

func getBody(url string) (string, error) {
	var (
		err  error
		body []byte
		resp *http.Response
	)
	if resp, err = http.Get(url); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	}
	return string(body), nil
}

type Crawler struct {
	crawled map[string]bool
	done    chan struct{}
	mux     sync.Mutex
	wg      sync.WaitGroup
}

func New() *Crawler {
	return &Crawler{
		crawled: make(map[string]bool),
	}
}

func (c *Crawler) Done() <-chan struct{} {
	return c.done
}

func (c *Crawler) visit(url string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.crawled[url]
	if ok {
		return true
	}
	c.crawled[url] = true
	return false
}

// Crawl recursively crawls pages for urls,
// to a maximum of depth. And returns all urls found.
func (c *Crawler) Crawl(u string, depth int, threadLimit int) (result []string) {
	c.done = make(chan struct{})
	sem := make(chan int, threadLimit)
	var recursiveCrawl func(u string, depth int)
	recursiveCrawl = func(u string, depth int) {
		sem <- 1
		v := c.visit(u)
		if v || depth <= 0 {
			return
		}
		urls, err := c.GetUrls(u)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, u := range urls {

			c.wg.Add(1)
			go func(u string) {
				defer c.wg.Done()
				recursiveCrawl(u, depth-1)
			}(u)
			<-sem
		}

	}
	recursiveCrawl(u, depth)
	c.wg.Wait()
	close(c.done)
	return Keys(c.crawled)
}

func Keys(m map[string]bool) []string {
	keys := make([]string, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func main() {
	u := flag.String("url", "https://filipdahlberg.dev/", "url that is to be crawled recursivley for urls")
	depth := flag.Int("depth", 2, "depth of the recursive crawl")
	threadLimit := flag.Int("threadLimit", 8, "threadLimit defines the number of thread allowed to run in parallel")
	crawler := New()
	res := crawler.Crawl(*u, *depth, *threadLimit)
	fmt.Println(res)
}
