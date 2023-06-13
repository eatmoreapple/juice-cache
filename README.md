# juice-cache


## Example

```go

package main

import (
    "fmt"
    "time"

    "github.com/eatmoreapple/juice"
    "github.com/eatmoreapple/juice_cache"
)


func main() {
    var client = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    cfg, err := juice.NewXMLConfiguration("config.xml")
	if err != nil {
		panic(err)
	}
	
	engine, err := juice.DefaultEngine(cfg)
	if err != nil {
		panic(err)
	}
	
	engine.SetCacheFactory(func() cache.Cache { return  juice_cache.NewRedisCache(client) })
}