package main

import (
    "fmt"
    "sync"
    "time"
    "math/rand"
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
    v   map[string]int
    mux sync.Mutex
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(key string, ms int, id int) {
    fmt.Println(id,": Sleeping for", ms, "ms at",time.Now()," ...")
    time.Sleep(time.Duration(ms) * time.Millisecond)
    fmt.Println(id,": Woke up at", time.Now()," ...")
    fmt.Println(id, ":", c.v[key],": Locking before incr ...")
    c.mux.Lock()
    fmt.Println(id,": Got lock at", time.Now()," ...")
    hangPeriod := 50
    fmt.Println(id, ": Hanging onto lock for", hangPeriod, "ms")
    time.Sleep(time.Duration(hangPeriod) * time.Millisecond)
    // Lock so only one goroutine at a time can access the map c.v.
    c.v[key]++
    fmt.Println(id, ":", c.v[key],": Unlocking after incr at", time.Now())
    c.mux.Unlock()
    fmt.Println(id,": Released lock at", time.Now()," ...")
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value(key string) int {
    fmt.Println("Value: Attempting to get lock ...")
    c.mux.Lock()
    // Lock so only one goroutine at a time can access the map c.v.
    defer c.mux.Unlock()
    return c.v[key]
}

func HelpScreen() {
    fmt.Println("This program illustrates the following behaviours:")
    fmt.Println(" * Lock requests are queued and granted FIFO")
    fmt.Println(" * go routines block on lock requests")
    fmt.Println(" * go routines are fast, reduce the hangPeriod to zero")
    fmt.Println("   to see that go routines waking up at the same time")
    fmt.Println("   still don't exhibit concurrent lock requests. Then")
    fmt.Println("   try 1 ms.")
    fmt.Println("Program output .....")
}

func main() {
    HelpScreen()
    c := SafeCounter{v: make(map[string]int)}
    for i := 0; i < 5; i++ {
        go c.Inc("somekey", rand.Intn(100), i)
    }

    time.Sleep(time.Duration(100)*time.Millisecond)
    fmt.Println("Result:", c.Value("somekey"))
}

