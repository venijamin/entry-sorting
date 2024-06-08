package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	file, err := os.OpenFile("data-500000.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for i := 0; i < 500000; i++ {
		file.Write([]byte(ranentry()))
	}
}

func ranentry() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s,%d\n", randate(), rand.Int63n(1023)))
	return sb.String()
}

func randate() string {
	min := time.Date(1900, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2100, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}
