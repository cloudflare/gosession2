package main

import (
	"flag"
	"gophq.io/gophq"
	"gophq.io/tls"
	"log"
	"strconv"
	"time"
)

var (
	addrFlag   = flag.String("addr", "127.0.0.1:9092", "broker address")
	topicFlag  = flag.String("topic", "topic", "topic to use")
	modeFlag   = flag.String("mode", "produce", "produce|consume")
	offsetFlag = flag.Int64("offset", 0, "first offset to fetch")

	caPath   = flag.String("tls.ca", "", "CA file")
	certPath = flag.String("tls.cert", "", "certificate file")
	keyPath  = flag.String("tls.key", "", "key file")
)

func main() {
	flag.Parse()

	// helper to deal with TLS files
	tlsConf := tls.NewTLSConfig(*caPath, *certPath, *keyPath)

	switch *modeFlag {
	case "produce":
		produce(*addrFlag, tlsConf, *topicFlag)
	case "consume":
		consume(*addrFlag, tlsConf, *topicFlag, *offsetFlag)
	}
}

func produce(addr string, tlsConf *tls.TLSConfig, topic string) {
	producer, err := gophq.NewProducer("tcp", addr, tlsConf)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	for i := 1; true; i++ {
		key := []byte("key" + strconv.Itoa(i))
		value := []byte("value" + strconv.Itoa(i))

		err = producer.SendMessage(topic, key, value)
		if err != nil {
			panic("SendMessage: " + err.Error())
		}
		log.Printf("sleeping after %d messages...", i)
		time.Sleep(1 * time.Second)
	}
}

func consume(addr string, tlsConf *tls.TLSConfig, topic string, offset int64) {
	// Use gophq.Consumer to fetch from the queue
	panic("something for you to implement")
}
