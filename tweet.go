package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func main() {

	f := createFile("/tmp/test11.txt")
	defer closeFile(f)
	writeFile(f)
}

func createFile(p string) *os.File {
	fmt.Println("creating")
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}

func writeFile(f *os.File) {
	fmt.Println("writing")
	config := oauth1.NewConfig("nHg431t3NqMCNev71bU1GRmYm", "Jcxid8DIbmfDFTpmVN7ng6HA8R0TQCyO1nCLwAPH7747pIb1Ot")
	token := oauth1.NewToken("929970417787322369-369Fm9R8A3AyM5an2JYiR8Omt74GjiP", "cmAK1SVhTaIkm0gjJ5P8qRCZHbBgKWxL3Ksp3dbREckpL")
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// sample returning

	params := &twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
	}
	stream, _ := client.Streams.Sample(params)

	// rceiving messages
	for message := range stream.Messages {
		fmt.Println(message)
	}

	// demux
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	for message := range stream.Messages {
		demux.Handle(message)
	}
	// stopping the stream

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	stream.Stop()

	//msg streams
	// msg := make(chan string)
	// go func() { msg <- "stream.Messages" }()
	//	fmt.Fprintln(f, log)
	//msgs := <-msg
	//fmt.Fprintln(f, ch)
	// Get Unicode code points.
	//n := 0
	//rune := make([]rune, len(input))
	//for _, r := range input {
	//	rune[n] = r
	//	n++
	//}
	//rune = rune[0:n]
	// Reverse
	//for i := 0; i < n/2; i++ {
	//	rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	//}
	// Convert back to UTF-8.
	//output := string(rune)
	//fmt.Fprintln(f, output)

}
func closeFile(f *os.File) {
	fmt.Println("closing")
	f.Close()
}
