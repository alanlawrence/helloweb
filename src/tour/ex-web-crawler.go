package main

import (
    "fmt"
    "sync"
)

type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(wgroup *sync.WaitGroup, url string, depth int, fetcher Fetcher) {
    fmt.Println("Crawl called for:", url, "at depth", depth)
    if depth <= 0 {
        defer WaitGroupDone(wgroup)
        return
    }
    body, urls, err := fetcher.Fetch(url)
    if err != nil {
        fmt.Println(err)
        defer WaitGroupDone(wgroup)
        return
    }
    fmt.Printf("found: %s %q\n", url, body)
    for _, u := range urls {
        WaitGroupAdd(wgroup)  // Go routine about to be added to group.
        go Crawl(wgroup, u, depth-1, fetcher)
    }
    // Go routine about to return and leave the group.
    defer WaitGroupDone(wgroup)
    return
}

// This counter is for diagnostic purposes only, it is not
// required for synchronisation of goroutines and main().
var wgCounter = 0

// The following wrapper functions are for diagnostic purposes
// only. The WaitGroup calls would normally be used directly at
// the call sites of these functions.
func WaitGroupAdd(wgroup *sync.WaitGroup) {
    wgCounter++
    fmt.Println("Calling wgroup.Add(1):", wgCounter-1,"->",wgCounter)
    wgroup.Add(1)
}

func WaitGroupDone(wgroup *sync.WaitGroup) {
    wgCounter--
    fmt.Println("Calling wgroup.Done():", wgCounter+1,"->",wgCounter)
    wgroup.Done()
}

// Data structure to hold map of fetched urls.
type FetchedUrls struct {
    fUrls   map[string]int
    mux     sync.Mutex
}

func (f *FetchedUrls) Inc(url string) {
    f.mux.Lock()
    f.fUrls[url]++
    f.mux.Unlock()
}

func (f *FetchedUrls) Value(url string) int {
    f.mux.Lock()
    defer f.mux.Unlock()
    return f.fUrls[url]
}

func (f *FetchedUrls) Fetched(url string) bool {
    f.mux.Lock()
    _, fetched := f.fUrls[url]
    f.mux.Unlock()
    return fetched

}

func (f *FetchedUrls) PrettyPrint() {
    f.mux.Lock()
    fmt.Println(len(f.fUrls), "distinct urls were fetched as follows:")
    for url, count := range f.fUrls {
        fmt.Printf("   Url: %-30s attempted fetch %2d times.\n", url, count)
    }
}

// Global variable for the fetched urls map.
var fetchedUrls FetchedUrls

// Run the program with:
// $ go run ex-web-crawler.go
func main() {
    // Instantiate the global fetched urls map.
    fetchedUrls = FetchedUrls{fUrls: make(map[string]int)}

    // Set up a synchronised wait group to block until
    // all parallel goroutines have returned.
    // Create the WaitGroup as a pointer as it will be passed
    // through successive go routines and we don't want copies.
    wgroup := new(sync.WaitGroup)
    WaitGroupAdd(wgroup)  // Necessary to match the first defer when the 
                          // following Crawl top level call returns.
    Crawl(wgroup, "https://golang.org/", 4, fetcher)
    fmt.Println("Calling wgroup.Wait()")
    // Wait until all parallelised goroutines have completed.
    // Without the waitgroup the fetched URL map would be printed
    // before the recursive Crawl goroutines complete leaving the map
    // only partially populated.
    wgroup.Wait()
    // Then print the resulting fetched URL map.
    fetchedUrls.PrettyPrint()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
    body string
    urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
    fmt.Println("Fetch called for:", url)
    if res, ok := f[url]; ok {
        var newUrls []string
        for _, url := range res.urls {
            if !fetchedUrls.Fetched(url) {
                newUrls= append(newUrls, url)
            }
            fmt.Println("Incrementing url:", url)
            fetchedUrls.Inc(url)
        }
        fmt.Println("Returning newUrls =", newUrls)
        return res.body, newUrls, nil
    }
    return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
    "https://golang.org/": &fakeResult{
        "The Go Programming Language",
        []string{
            "https://golang.org/pkg/",
            "https://golang.org/cmd/",
        },
    },
    "https://golang.org/pkg/": &fakeResult{
        "Packages",
        []string{
            "https://golang.org/",
            "https://golang.org/cmd/",
            "https://golang.org/pkg/fmt/",
            "https://golang.org/pkg/os/",
        },
    },
    "https://golang.org/pkg/fmt/": &fakeResult{
        "Package fmt",
        []string{
            "https://golang.org/",
            "https://golang.org/pkg/",
        },
    },
    "https://golang.org/pkg/os/": &fakeResult{
        "Package os",
        []string{
            "https://golang.org/",
            "https://golang.org/pkg/",
        },
    },
}

