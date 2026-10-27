package checker

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/nxvrlxv/hosts_monitor/internal/config"
)

type Checker struct {
	Client *http.Client
}

type Result struct {
	IsUp         bool
	Host         string
	Status       int
	ResponseTime time.Duration
	Error        error
}

func NewChecker() *Checker {
	return &Checker{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Checker) Check(ctx context.Context, target config.Target) *Result {
	start := time.Now()
	resp, err := c.Client.Get(target.URL)
	Url, _ := url.Parse(target.URL)
	if err != nil {
		res := &Result{
			IsUp:         false,
			Host:         Url.Host,
			ResponseTime: time.Since(start),
			Error:        err,
		}
		return res
	}
	defer resp.Body.Close()

	respTime := time.Since(start)
	res := Result{}

	res.ResponseTime = respTime
	res.Status = resp.StatusCode
	res.Host = Url.Host
	if resp.StatusCode < 400 {
		res.IsUp = true
	} else {
		res.IsUp = false
	}

	return &res

}

type Semaphor struct {
	C chan struct{}
}

func (s *Semaphor) Acquire() {
	s.C <- struct{}{}
}

func (s *Semaphor) Release() {
	<-s.C
}

func CheckerPool(maxCheckers int, data []config.Target, ctx context.Context, c Checker) []*Result {
	sem := Semaphor{
		C: make(chan struct{}, maxCheckers),
	}
	var wg sync.WaitGroup
	resultCh := make(chan *Result, len(data))

	for _, target := range data {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem.Acquire()
			defer sem.Release()

			res := c.Check(ctx, target)
			//fmt.Println(res)
			resultCh <- res
		}()
	}

	wg.Wait()
	close(resultCh)
	resArr := make([]*Result, 0)
	for res := range resultCh {
		resArr = append(resArr, res)
	}
	return resArr
}
