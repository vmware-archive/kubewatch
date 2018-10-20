package icq

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

// HTTP Client interface
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// API
type API struct {
	token   string
	baseUrl string
	client  Doer
	mu      *sync.Mutex
}

// NewAPI constructor of API object
func NewAPI(token string) *API {
	return &API{
		token:   token,
		baseUrl: "https://botapi.icq.net",
		mu:      new(sync.Mutex),
		client:  http.DefaultClient,
	}
}

// SendMessage with `message` text to `to` participant
func (a *API) SendMessage(to string, message string) (*MessageResponse, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	v := url.Values{
		"aimsid":  []string{a.token},
		"r":       []string{strconv.FormatInt(time.Now().Unix(), 10)},
		"t":       []string{to},
		"message": []string{message},
	}
	req, err := http.NewRequest(http.MethodGet, a.baseUrl+"/im/sendIM?"+v.Encode(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	r := &Response{}
	if err := json.Unmarshal(b, r); err != nil {
		return nil, err
	}
	if r.Response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to send message: %s", r.Response.StatusText)
	}
	return r.Response.Data, nil
}

// UploadFile to ICQ servers and returns URL to file
func (a *API) UploadFile(fileName string, r io.Reader) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	v := url.Values{"aimsid": []string{a.token}, "filename": []string{fileName}}
	req, err := http.NewRequest(http.MethodPost, a.baseUrl+"/im/sendFile?"+v.Encode(), r)
	if err != nil {
		return "", err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	rb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	file := struct {
		Data struct {
			StaticUrl string `json:"static_url"`
		} `json:"data"`
	}{}

	if err := json.Unmarshal(rb, &file); err != nil {
		return "", err
	}
	return file.Data.StaticUrl, nil
}

// GetWebhookHandler returns http.HandleFunc that parses webhooks
func (a *API) GetWebhookHandler(cu chan<- Update, e chan<- error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			e <- fmt.Errorf("incorrect method: %s", r.Method)
			return
		}
		wr := &WebhookRequest{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			e <- err
			return
		}
		if err := json.Unmarshal(b, wr); err != nil {
			e <- err
			return
		}
		for _, u := range wr.Updates {
			cu <- u
		}
	}
}
