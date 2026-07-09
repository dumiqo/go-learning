package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CheckResult struct {
	IsValid bool
	Url     string
}

type ValidationClient interface {
	Check(url string) (bool, error)
}

type httpValidationClient struct {
	client *http.Client
}

func (c httpValidationClient) Check(url string) (bool, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200, nil
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

	jobs := make(chan string, 10)
	result := make(chan CheckResult)
	var wg sync.WaitGroup
	var wg1 sync.WaitGroup
	client := httpValidationClient{client: &http.Client{
		Timeout: 1 * time.Second,
	}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < 3; i++ {

		wg.Add(1)
		go check(client, i, jobs, result, &wg, ctx)
	}
	wg1.Add(1)
	go func() {
		for {
			r, ok := <-result
			if !ok {
				wg1.Done()
				break
			}
			fmt.Printf("Result. url: %s, isValid: %t\n", r.Url, r.IsValid)
		}
	}()

	for _, v := range urls {
		jobs <- v
	}
	close(jobs)

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		close(result)
		wg1.Wait()
	case <-time.After(1 * time.Second):
		cancel()
		fmt.Printf("Timeout\n")
	}
}

func check(client ValidationClient, id int, jobs <-chan string, result chan<- CheckResult, wg *sync.WaitGroup, ctx context.Context) {
	for {
		select {
		case url, ok := <-jobs:
			if !ok {
				wg.Done()
				break
			}
			fmt.Printf("check start. worker: %d, url: %s\n", id, url)
			valid, err := client.Check(url)
			if err != nil {
				result <- CheckResult{false, url}
				fmt.Printf("check ended. worker: %d, url: %s, %s\n", id, url, err)
			} else {
				result <- CheckResult{valid, url}
				fmt.Printf("check ended. worker: %d, url: %s\n", id, url)
			}
		case <-ctx.Done():
			wg.Done()
		}
	}
}
