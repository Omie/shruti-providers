package main

import (
	"log"
	"os"
	"time"

	//"github.com/SlyMarbo/rss"
	shruti "github.com/omie/shruti-go"
)

const (
	PROVIDER_NAME = "<PROVIDER_NAME>"
)

var (
	sClient shruti.Client
	visited map[string]bool
)

func Register() (err error) {

	p := shruti.Provider{Name: PROVIDER_NAME,
		DisplayName: "<DisplayName>",
		Description: "<Discription>",
		WebURL:      "<WebURL>",
		IconURL:     "<IconURL>",
		Voice:       "<Voice>",
	}

	err = sClient.RegisterProvider(p)

	return

}

func doWork() (err error) {

	for {
		for _, item := range items[:10] {
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

	}
}

func main() {

	log.Println("Starting Scraper")

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
