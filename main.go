package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func createAllocatorContext(ctx context.Context, addr string) (context.Context, context.CancelFunc) {
	if addr == "" {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("enable-automation", false),
			chromedp.Flag("disable-extensions", false),
		)
		return chromedp.NewExecAllocator(ctx, opts...)
	}

	return chromedp.NewRemoteAllocator(ctx, addr)
}

func main() {

	var chromeUrl, username, password string
	flag.StringVar(&username, "username", "test@test", "login username")
	flag.StringVar(&password, "password", "password", "login password")
	flag.StringVar(&chromeUrl, "chromeUrl", "ws://localhost:9222", "chrome url")
	flag.Parse()

	actx, acancel := createAllocatorContext(context.Background(), chromeUrl)
	defer acancel()

	ctx, cancel := chromedp.NewContext(
		actx,
		chromedp.WithErrorf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Minute)
	defer cancel()

	got := ""
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://test.apply.se/soeb`),
		chromedp.SendKeys("//input[@type='email']", username, chromedp.BySearch),
		chromedp.SendKeys("//input[@type='password']", password, chromedp.BySearch),
		chromedp.Click("//button[@type='submit']", chromedp.BySearch),
		chromedp.Text(`//p[contains(@class, 'text-error')]`, &got),
	)
	if err != nil {
		log.Fatal(err)
	}

	want := "Felaktigt användarnamn eller lösenord"
	if got != want {
		log.Fatalf("expected %s got %s\n", want, got)
	}

}
