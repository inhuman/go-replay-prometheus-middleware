package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"os"
)

func main() {

	m := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "goreplay_mirroring",
	}, []string{
		"payload",
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

		switch payloadType {
		case '1': // Request
			m.With(prometheus.Labels{"payload": "request"}).Add(1)

		case '2': // Original response
			m.With(prometheus.Labels{"payload": "original_response"}).Add(1)

		case '3': // Replayed response
			m.With(prometheus.Labels{"payload": "replayed_response"}).Add(1)
		}
	}
}

func listener() error {

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	p := ginprom.New(
		ginprom.Engine(r),
	)
	r.Use(p.Instrument())

	return r.Run(":9876")
}
