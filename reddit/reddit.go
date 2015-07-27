package main // import "github.com/omie/shruti-providers/reddit"

import (
	"log"
	"os"
	"time"

	"github.com/SlyMarbo/rss"
	shruti "github.com/omie/shruti-go"
)

const (
	PROVIDER_NAME = "reddit"
)

var (
	sClient shruti.Client
	feedURL string
	visited map[string]bool
)

func Register() (err error) {

	p := shruti.Provider{Name: PROVIDER_NAME,
		DisplayName: "reddit",
		Description: "the front page of the internet",
		WebURL:      "https://reddit.com",
		IconURL:     "http://i.imgur.com/F0T020X.png",
		Voice:       "Raveena",
	}

	err = sClient.RegisterProvider(p)

	return

}

func doWork() (err error) {

	feed, err := rss.Fetch(feedURL)
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
		time.Sleep(30 * time.Minute)
		log.Println("woke up")

		err = feed.Update()
		if err != nil {
			return err
		}

	}
}

func main() {

	log.Println("Starting reddit scraper")

	shrutiServer := os.Getenv("SHRUTI_SERVER")
	if shrutiServer == "" {
		log.Println("SHRUTI_SERVER not set")
		return
	}

	feedURL = os.Getenv("REDDIT_FEED_URL")
	if feedURL == "" {
		log.Println("REDDIT_FEED_URL not set")
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
