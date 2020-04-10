package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"net/http"
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

		r := bytes.NewReader(buf)

		_, err := ReadHTTPFromFile(r)
		if err != nil {
			fmt.Printf("error read from stdin: %s\n", err)
		}

		//// First byte indicate payload type, possible values:
		////  1 - Request
		////  2 - Response
		////  3 - ReplayedResponse
		//payloadType := buf[0]
		//
		//switch payloadType {
		//case '1': // Request
		//	m.With(prometheus.Labels{"payload": "request"}).Add(1)
		//	os.Stdout.Write(encode(buf))
		//
		//case '2': // Original response
		//	m.With(prometheus.Labels{"payload": "original_response"}).Add(1)
		//
		//case '3': // Replayed response
		//	m.With(prometheus.Labels{"payload": "replayed_response"}).Add(1)
		//}
	}
}

type Connection struct {
	Request  *http.Request
	Response *http.Response
}

func ReadHTTPFromFile(r io.Reader) ([]Connection, error) {
	buf := bufio.NewReader(r)
	stream := make([]Connection, 0)

	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return stream, err
		}

		fmt.Printf("req: %+v\n", req)

		//resp, err := http.ReadResponse(buf, req)
		//if err != nil {
		//	return stream, err
		//}
		//
		////save response body
		//b := new(bytes.Buffer)
		//io.Copy(b, resp.Body)
		//resp.Body.Close()
		//resp.Body = ioutil.NopCloser(b)
		//
		//stream = append(stream, Connection{Request: req, Response: resp})
	}
	return stream, nil

}

func encode(buf []byte) []byte {
	dst := make([]byte, len(buf)*2+1)
	hex.Encode(dst, buf)
	dst[len(dst)-1] = '\n'

	return dst
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
