package main

import (
	"log"
	"os"
	"time"

	"github.com/SlyMarbo/rss"
	shruti "github.com/omie/shruti-go"
)

const (
	PROVIDER_NAME = "bbc"
	BBC_FEED_URL  = "http://feeds.bbci.co.uk/news/world/rss.xml?edition=uk"
)

var (
	sClient shruti.Client
	visited map[string]bool
)

func Register() (err error) {

	p := shruti.Provider{Name: PROVIDER_NAME,
		DisplayName: "BBC",
		Description: "world news headlines from BBC",
		WebURL:      "http://bbc.com",
		IconURL:     "http://static.bbci.co.uk/frameworks/barlesque/2.83.10/desktop/3.5/img/blq-blocks_grey_alpha.png",
		Voice:       "Brian",
	}

	err = sClient.RegisterProvider(p)

	return

}

func doWork() (err error) {

	feed, err := rss.Fetch(BBC_FEED_URL)
	if err != nil {
		return err
	}

	for {
		for _, item := range feed.Items[:10] {
			if _, ok := visited[item.ID]; ok {
				continue
			}
			n := shruti.Notification{
				Title:        item.Title,
				Url:          item.Link,
				Key:          PROVIDER_NAME + item.ID,
				Priority:     shruti.PRIO_MED,
				Action:       shruti.ACT_PUSH,
				ProviderName: PROVIDER_NAME,
			}
			err = sClient.PushNotification(n)
			msg := "submitted"
			if err != nil {
				msg = err.Error()
			}
			visited[item.ID] = true
			log.Println(msg)
		}
		log.Println("Sleeping")
		time.Sleep(15 * time.Minute)
		log.Println("woke up")

		err = feed.Update()
		if err != nil {
			return err
		}

	}
}

func main() {

	log.Println("Starting BBC News Scraper")

	shrutiServer := os.Getenv("SHRUTI_SERVER")
	if shrutiServer == "" {
		log.Println("SHRUTI_SERVER not set")
		return
	}

	sClient = shruti.Client{"http://", shrutiServer}

	err := Register()
	if err != nil {
		log.Println(err)
		return
	}

	visited = make(map[string]bool)

	err = doWork()
	if err != nil {
		log.Println(err)
	}
}
