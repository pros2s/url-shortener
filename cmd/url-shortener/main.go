package main

import (
	"fmt"

	"url-shortener/internal"
)

func main() {
	cfg := internal.MustLoad()

	fmt.Println(cfg)
}
