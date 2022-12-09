package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"explore/reader"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "start-delay",
				Value: 0,
				Usage: "Seconds to delay at the beginning of post data",
			},
			&cli.IntFlag{
				Name:  "middle-delay",
				Value: 0,
				Usage: "Seconds to delay in the middle of post data",
			},
			&cli.IntFlag{
				Name:  "end-delay",
				Value: 0,
				Usage: "Seconds to delay at the end of post data",
			},
			&cli.StringFlag{
				Name:  "url",
				Value: "",
				Usage: "URL to issue a POST request against",
				Action: func(ctx *cli.Context, v string) error {
					if v == "" {
						return fmt.Errorf("Required flag `--url` was not provided")
					}
					return nil
				},
			},
		},
		Action: makeRequest,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func makeRequest(cli *cli.Context) error {
	startDelay := cli.Int("start-delay")
	middleDelay := cli.Int("middle-delay")
	endDelay := cli.Int("end-delay")

	body := reader.NewDelayedReadWriter(startDelay, middleDelay, endDelay)

	done := make(chan error)
	go func() {
		fmt.Printf("Writing request body\n")
		err := body.Write(strings.Repeat("request data!!!\n", 64))
		fmt.Printf("Done writing request body\n")
		done <- err
	}()

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	fmt.Printf("Issuing request\n")
	resp, err := http.Post(cli.String("url"), "application/text", body)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}
	fmt.Printf("Got response\n")
	dump, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return fmt.Errorf("Error dumping response: %s", err)
	}
	fmt.Printf("%s\n", dump)

	return <-done
}
