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

	actx, acancel := chromedp.NewRemoteAllocator(context.Background(), "ws://localhost:9222")
	defer acancel()

	ctx, cancel := chromedp.NewContext(
		actx,
		chromedp.WithErrorf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Minute)
	defer cancel()

	stuff := ""
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://test.apply.se/soeb`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("1")
			return nil
		}),
		chromedp.SendKeys("//input[@type='email']", username, chromedp.BySearch),
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("2")
			return nil
		}),
		chromedp.SendKeys("//input[@type='password']", password, chromedp.BySearch),
		chromedp.Click("//button[@type='submit']", chromedp.BySearch),
		chromedp.Value(`//p[contains(@class, 'text-error')]`, &stuff),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("stuff:\n%s", stuff)
}
