package shrutigo

import (
	"encoding/json"
	"path"

	"github.com/parnurzeal/gorequest"
)

type Provider struct {
	Id          int    `json:"id, omitempty"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description, omitempty"`
	WebURL      string `json:"web_url, omitempty"`
	IconURL     string `json:"icon_url, omitempty"`
	Active      bool   `json:"active, omitempty"`
	Voice       string `json:"voice"`
}

func (client *Client) GetAllProviders() (providers []*Provider, err error) {

	url := client.Protocol + path.Join(client.Host, "providers")

	_, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		err = errs[0]
		return
	}

	providers = make([]*Provider, 0)
	err = json.Unmarshal([]byte(body), &providers)

	return
}

func (client *Client) GetSingleProvider(providerName string) (p *Provider, err error) {

	url := client.Protocol + path.Join(client.Host, "providers", providerName)

	_, body, errs := gorequest.New().Get(url).End()
	if errs != nil {
		err = errs[0]
		return
	}

	p = new(Provider)
	err = json.Unmarshal([]byte(body), p)

	return
}

func (client *Client) RegisterProvider(p Provider) (err error) {

	url := client.Protocol + path.Join(client.Host, "providers", p.Name)

	request := gorequest.New()
	_, _, errs := request.Post(url).
		Send(p).
		End()

	if errs != nil {
		err = errs[0]
	}

	return
}
