package main

import (
    "log"
    "time"
    "sync"
    "runtime"
)

var (
    dsn       = "127.0.0.1:9876"
    cacheItem = &CacheItem{Key: "some key", Value: "some value"}
    wg sync.WaitGroup
)

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func newClient() (*Client, error) {
    return NewClient(dsn, time.Millisecond*500)
}

func TestColdGet(c *Client) {
    item, _ := c.Get(cacheItem.Key)
    if item != nil {
        log.Printf("Cache key should not exist: %s\n", cacheItem.Key)
    }
}

func TestPut(c *Client) {
    _, err := c.Put(cacheItem)
    if err != nil {
        log.Printf("[ERR] %v", err)
    }
}

func TestWarmGet(c *Client) {
    item, _ := c.Get(cacheItem.Key)
    if item == nil {
        log.Printf("Cache key should exist: %s\n", cacheItem.Key)
    }
    if item.Value != cacheItem.Value {
        log.Printf("Cache value expected %s got %s\n", cacheItem.Value, item.Value)
    }
}

func TestDelete(c *Client) {
    _, err := c.Delete(cacheItem.Key)
    if err != nil {
        log.Printf("[ERR] %v", err)
    }
    item, _ := c.Get(cacheItem.Key)
    if item != nil {
        log.Printf("Cache key should not exist: %s\n", cacheItem.Key)
    }
}

func TestClear(c *Client) {
    _, err := c.Clear()
    if err != nil {
        log.Printf("[ERR] %v", err)
    }
}

func TestStats(c *Client) {
    stats, err := c.Stats()
    if err != nil {
        log.Printf("[ERR] %v", err)
    }
    if stats.Get != 1 {
        log.Printf("Get: expected 1, got %d\n", stats.Get)
    }
    if stats.Put != 1 {
        log.Printf("Put: expected 1, got %d\n", stats.Put)
    }
    if stats.Delete != 1 {
        log.Printf("Delete: expected 1, got %d\n", stats.Delete)
    }
    if stats.Clear != 1 {
        log.Printf("Clear: expected 1, got %d\n", stats.Clear)
    }
}

func main() {

    wg.Add(2)
    go runServer(&wg)
    go func(w *sync.WaitGroup) {
        defer wg.Done()

        log.Print("Client cycle start ...")

        for {
            c, err := newClient()
            if err != nil {
                log.Print("ERR] %v", err)
                time.Sleep(time.Second)
                continue
            }

            TestPut(c)
            TestColdGet(c)
            TestWarmGet(c)
            TestDelete(c)
            TestClear(c)
            TestStats(c)

            log.Println("a cycle completed\n")
            time.Sleep(time.Second * 2)
        }
    }(&wg)
    wg.Wait()
}