package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/buger/goreplay/proto"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"strings"
)

func main() {

	conf, err := Config()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "goreplay_mirroring",
	}, []string{
		"payload",
		"type",
		"ctype",
	})

	prometheus.MustRegister(m)

	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		if err := listener(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	for scanner.Scan() {
		encoded := scanner.Bytes()
		buf := make([]byte, len(encoded)/2)
		hex.Decode(buf, encoded)

		// First byte indicate payload type, possible values:
		//  1 - Request
		//  2 - Response
		//  3 - ReplayedResponse
		payloadType := buf[0]
		headerSize := bytes.IndexByte(buf, '\n') + 1
		payload := buf[headerSize:]

		switch payloadType {
		case '1': // Request

			t := getUrlType(conf, payload)

			m.With(prometheus.Labels{"payload": "request", "type": t.Type, "ctype": t.CType}).Add(1)

			os.Stdout.Write(encode(buf))

		case '2': // Original response
			m.With(prometheus.Labels{"payload": "original_response"}).Add(1)

		case '3': // Replayed response
			m.With(prometheus.Labels{"payload": "replayed_response"}).Add(1)
		}
	}
}

func getUrlType(conf *AppConfig, payload []byte) UrlType {

	for _, c := range conf.UrlTypes {
		if strings.Contains(string(proto.Path(payload)), c.Url) {
			return c
		}
	}
	return UrlType{}
}

func encode(buf []byte) []byte {
	dst := make([]byte, len(buf)*2+1)
	hex.Encode(dst, buf)
	dst[len(dst)-1] = '\n'

	return dst
}
