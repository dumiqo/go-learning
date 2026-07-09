package main

import (
	"fmt"
	"net/http"
	"sync"
)

type CheckResult struct {
	IsValid bool
	Url     string
}

func main() {
	urls := []string{
		"https://golang.org/",
		"https://golang.org/pkg/",
		"https://golang.org/cmd/",
		"https://golang.org/pkg/fmt/",
		"https://golang.org/pkg/os/",
		"https://github.com/",
		"https://github.com/explore",
		"https://github.com/topics",
		"https://github.com/topicsfdsf",
		"https://stackoverflow.com/",
		"https://stackovesffrflow.com/",
		"https://stackoverflow.com/questions",
		"https://en.wikipedia.org/wiki/Main_Page",
		"https://en.wikipedia.org/wiki/Go_(programming_language)",
		"https://www.google.com/",
		"https://www.google.com/search?q=golang",
		"https://www.youtube.com/",
		"https://www.youtube.com/results?search_query=golang",
		"https://www.reddit.com/",
		"https://www.reddit.com/r/golang/",
		"https://news.ycombinator.com/",
		"https://news.ycombinator.com/newest",
		"https://news.ycombinator.com/random_not_exist_url_adfbuicx",
		"https://www.facebook.com/",
		"https://twitter.com/",
		"https://www.instagram.com/",
		"https://www.linkedin.com/",
		"https://www.amazon.com/",
		"https://www.ebay.com/",
		"https://www.apple.com/",
		"https://www.microdsfffffffsoft.con/",
		"https://www.microsoft.com/",
		"https://www.ubuntu.com/",
		"https://www.debian.org/",
		"https://www.archlinux.org/",
		"https://www.kernel.org/",
		"https://www.docker.com/",
		"https://kubernetes.io/",
		"https://www.python.org/",
		"https://www.ruby-lang.org/",
		"https://www.php.net/",
		"https://www.perl.org/",
		"https://golang.org/doc/",
		"https://golang.org/blog/",
	}

	jobs := make(chan string, 3)
	result := make(chan CheckResult, 3)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go check(i, jobs, result)
	}

	go func() {
		for {
			r, ok := <-result
			if !ok {
				break
			}
			wg.Done()
			fmt.Printf("Result. url: %s, isValid: %t\n", r.Url, r.IsValid)
		}
	}()

	for _, v := range urls {
		jobs <- v
	}
	close(jobs)
	wg.Wait()
	close(result)
}

func check(id int, jobs <-chan string, result chan<- CheckResult) {
	for url := range jobs {
		fmt.Printf("check start. worker: %d, url: %s\n", id, url)
		req, err := http.Get(url)
		if err != nil {
			result <- CheckResult{false, url}
			fmt.Printf("check ended. worker: %d, url: %s\n", id, url)
		} else {
			valid := req.StatusCode == 200

			result <- CheckResult{valid, url}
			fmt.Printf("check ended. worker: %d, url: %s\n", id, url)
		}
	}
}
