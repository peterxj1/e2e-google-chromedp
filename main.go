package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
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

	url := "https://test.apply.se"
	var statusCode int64
	responseHeaders := map[string]interface{}{}

	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch responseReceivedEvent := event.(type) {
		case *network.EventResponseReceived:
			response := responseReceivedEvent.Response
			if response.URL == url {
				statusCode = response.Status
				responseHeaders = response.Headers
				fmt.Printf("%+v\n", responseHeaders)
			}
		}
	})

	fmt.Println(statusCode)

	stuff := ""
	// navigate to a page, wait for an element, click
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://test.apply.se/soeb`),
		chromedp.SendKeys("//input[@type='email']", username, chromedp.BySearch),
		chromedp.SendKeys("//input[@type='password']", password, chromedp.BySearch),
		chromedp.Click("//*button[@type='submit']", chromedp.BySearch),
		chromedp.Value(`//p[contains(@class, 'text-error')]`, &stuff),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("stuff:\n%s", stuff)
}
