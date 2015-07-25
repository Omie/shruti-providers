package shrutigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	// Priorities
	PRIO_LOW  = 10
	PRIO_MED  = 20
	PRIO_HIGH = 30

	// Actions
	ACT_POLL = 10
	ACT_PUSH = 20

	// Heard status
	HRD_UNHEARD = 10
	HRD_HEARD   = 20
)

type Notification struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	Url          string    `json:"url, omitempty"`
	Key          string    `json:"key"`
	Heard        int       `json:"heard"`
	Provider     int       `json:"provider"`
	CreatedOn    time.Time `json:"created_on, omitempty"`
	Priority     int       `json:"priority"`
	Action       int       `json:"action"`
	ProviderName string    `json:"provider_name, omitempty"`
}

func (client *Client) GetSingleNotification(id int) (n *Notification, err error) {

	url := client.Protocol + path.Join(client.Host, "notifications", fmt.Sprintf("%d", id))

	_, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		err = errs[0]
		return
	}

	n = new(Notification)
	err = json.Unmarshal([]byte(body), n)

	return
}

func (client *Client) GetNotificationsSince(since *time.Time) (n []*Notification, err error) {

	url := client.Protocol + path.Join(client.Host, "notifications/since", fmt.Sprintf("%d", since.Unix()))

	_, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		err = errs[0]
		return
	}

	n = make([]*Notification, 0)
	err = json.Unmarshal([]byte(body), &n)

	return
}

func (client *Client) GetUnheardNotifications() (n []*Notification, err error) {

	url := client.Protocol + path.Join(client.Host, "notifications/unheard")

	_, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		err = errs[0]
		return
	}

	n = make([]*Notification, 0)
	err = json.Unmarshal([]byte(body), &n)

	return
}

func (client *Client) PushNotification(n Notification) (err error) {

	url := client.Protocol + path.Join(client.Host, "notifications")

	request := gorequest.New()
	resp, body, errs := request.Post(url).
		Send(n).
		End()

	if errs != nil {
		err = errs[0]
	}

	if resp.StatusCode == 500 {
		err = errors.New(body)
	}

	return
}
