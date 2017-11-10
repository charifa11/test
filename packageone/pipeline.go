package packageone
 
import"log"

func TextSocketCollector(name string, out chan streamer.Message) {
    listener, _ := net.Listen("tcp", ":9999")
    conn, _ := listener.Accept()

    for {
        line, _ := bufio.NewReader(conn).ReadString('\n')
        line = strings.TrimSuffix(line, "\n")

        out_message := streamer.NewMessage()
        out_message.Put("line", line)

        log.Printf("Received raw message from socket: %s\n", out_message)

        out <- out_message
    }
}

func WordExtractor(name string, input streamer.Message, out chan streamer.Message) {
    line, _ := input.Get("line").(string)

    words := strings.Split(line, " ")

    for _, word := range words {
        out_message := streamer.NewMessage()
        out_message.Put("word", word)
        log.Printf("Extracted word: %s\n", word)
        out <- out_message
    }
}

func HashTagFilter(name string, input streamer.Message, out chan streamer.Message) {
    word, _ := input.Get("word").(string)

    if (strings.HasPrefix(word, "#")) {
        out_message := streamer.NewMessage()
        out_message.Put("hashtag", word)
        log.Printf("Filtered hashtag %s\n", word)
        out <- out_message
    }
}

func RunPipeline() {
    // build pipeline elements
    collector := streamer.NewCollector("collector",
	    TextSocketCollector)

    extractor := streamer.NewProcessor("extractor",
        WordExtractor, streamer.NewIndexedChannelDemux(2, streamer.RandomIndex))

    filter := streamer.NewProcessor("filter",
        HashTagFilter, streamer.NewIndexedChannelDemux(2, streamer.RandomIndex))

    publisher := streamer.NewProcessor("publisher",
        HashTagPublisher, streamer.NewIndexedChannelDemux(5, streamer.NewGroupDemux("hashtag").GroupIndex))

    // execute pipeline
    sequence := collector.Execute()
    extracted := extractor.Execute(sequence)
    filtered := filter.Execute(extracted)
    <-publisher.Execute(filtered)
}
