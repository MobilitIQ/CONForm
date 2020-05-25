package main

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
	"net/http"
	"os"
	"time"
)

var PORT = os.Getenv("PORT")

func checkEnv() {
	if PORT == "" {
		PORT = "8812"
	}
}

func main() {

	checkEnv()

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LFU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		fmt.Println(err)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10 * time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
	}

	box := packr.New("Assets", "./assets")
	fileServer := http.FileServer(box)
	cached := cacheClient.Middleware(fileServer)

	http.Handle("/", cached)

	if err := http.ListenAndServe(":" + PORT, nil); err != nil {
		panic(err)
	}
}
