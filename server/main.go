package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	systemPortString := os.Getenv("PORT")
	port, err := strconv.Atoi(systemPortString)
	if err != nil {
		log.Fatal("invalid required env var PORT")
	}

	startDelayString := os.Getenv("START_DELAY")
	if startDelayString == "" {
		startDelayString = "0"
	}
	startDelay, err := strconv.Atoi(startDelayString)
	if err != nil {
		log.Fatal("invalid START_DELAY")
	}

	middleDelayString := os.Getenv("MIDDLE_DELAY")
	if middleDelayString == "" {
		middleDelayString = "0"
	}
	middleDelay, err := strconv.Atoi(middleDelayString)
	if err != nil {
		log.Fatal("invalid MIDDLE_DELAY")
	}
	endDelayString := os.Getenv("END_DELAY")
	if endDelayString == "" {
		endDelayString = "0"
	}
	endDelay, err := strconv.Atoi(endDelayString)
	if err != nil {
		log.Fatal("invalid END_DELAY")
	}

	mux := http.NewServeMux()
	mux.Handle("/", newDelayHandler(startDelay, middleDelay, endDelay))

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mux)
}

type DelayHandler struct {
	startDelay  time.Duration
	middleDelay time.Duration
	endDelay    time.Duration
	payload     string
}

func newDelayHandler(s, m, e int) *DelayHandler {
	p := strings.Repeat("response data!!\n", 64)
	return &DelayHandler{
		startDelay:  time.Duration(s),
		middleDelay: time.Duration(m),
		endDelay:    time.Duration(e),
		payload:     p,
	}
}

func (dh *DelayHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logMsg("Received %d byte request", req.ContentLength)
	data, err := io.ReadAll(req.Body)
	if err != nil {
		logErr("error reading request: %s", err)
		res.WriteHeader(http.StatusBadRequest)
	}
	logMsg("Read %d bytes from request", len(data))

	logMsg("Writing first chunk of response")
	time.Sleep(dh.startDelay * time.Second)
	_, err = res.Write([]byte(dh.payload)[0 : len(dh.payload)/2])
	if err != nil {
		logErr("error writing furst chunk of response: %s", err)
		res.WriteHeader(http.StatusInternalServerError)
	}
	logMsg("Writing second chunk of response")
	time.Sleep(dh.middleDelay * time.Second)
	_, err = res.Write([]byte(dh.payload)[len(dh.payload)/2 : len(dh.payload)])
	if err != nil {
		logErr("error writing furst chunk of response: %s", err)
		res.WriteHeader(http.StatusInternalServerError)
	}
	time.Sleep(dh.endDelay * time.Second)
	logMsg("Done responding")
}

func logMsg(msg string, a ...any) {
	fmt.Printf("%s\n", fmt.Sprintf(msg, a...))
}
func logErr(msg string, a ...any) {
	fmt.Fprintf(os.Stderr, "%s\n", fmt.Sprintf(msg, a...))
}
