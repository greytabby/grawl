package fetcher

import (
	"context"
	"unsafe"

	"github.com/chromedp/chromedp"
)

type HeadlessChrome struct{}

func (hc *HeadlessChrome) Fetch(URL string) (body []byte, err error) {
	ctx, cansel := chromedp.NewContext(context.Background())
	defer cansel()

	var content string
	err = chromedp.Run(ctx,
		chromedp.Navigate(URL),
		chromedp.OuterHTML(`html`, &content, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}
	return *(*[]byte)(unsafe.Pointer(&content)), nil
}
