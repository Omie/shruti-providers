package main // import "github.com/omie/shruti-providers/hn"

import (
	"log"
	"os"
	"strconv"

	shruti "github.com/omie/shruti-go"
	"github.com/peterhellberg/hn"
)

const (
	PROVIDER_NAME = "hackernews"
)

var (
	sClient shruti.Client
	visited map[int]bool
)

func Register() (err error) {

	p := shruti.Provider{Name: PROVIDER_NAME,
		DisplayName: "Hacker News",
		Description: "updates around the globe, mostly tech related",
		WebURL:      "https://news.ycombinator.com/",
		IconURL:     "https://news.ycombinator.com/favicon.ico",
		Voice:       "Emma",
	}

	err = sClient.RegisterProvider(p)

	return

}

func doWork() (err error) {
	hn := hn.NewClient(nil)

	ids, err := hn.TopStories()
	if err != nil {
		return err
	}

	for _, id := range ids[:10] {
		item, err := hn.Item(id)
		if _, ok := visited[id]; ok {
			continue
		}
		if err == nil {
			n := shruti.Notification{
				Title:        item.Title,
				Url:          item.URL,
				Key:          PROVIDER_NAME + strconv.Itoa(item.ID),
				Priority:     shruti.PRIO_MED,
				Action:       shruti.ACT_PUSH,
				ProviderName: PROVIDER_NAME,
			}
			err = sClient.PushNotification(n)
			msg := "submitted"
			if err != nil {
				msg = err.Error()
			}
			log.Println(msg)
			visited[id] = true
		}
	}
	return nil
}

func main() {

	log.Println("Starting Hacker News Scraper")

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

	visited = make(map[int]bool)

	err = doWork()
	if err != nil {
		log.Println(err)
	}
}
