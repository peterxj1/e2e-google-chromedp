package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	var username, password string
	flag.StringVar(&username, "username", "test", "login username")
	flag.StringVar(&password, "password", "password", "login password")
	flag.Parse()

	fmt.Println(username)
	fmt.Println(password)

	actx, acancel := chromedp.NewRemoteAllocator(context.Background(), "http://localhost:9222")
	defer acancel()

	ctx, cancel := chromedp.NewContext(
		actx,
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// find and click "Example" link
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		// retrieve the text of the textarea
		chromedp.Value(`#example-After textarea`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\n%s", example)
}
