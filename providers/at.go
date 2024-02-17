package providers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	ATProvider interface{}

	atProvider struct {
		atBaseURL,
		username,
		key string
		client *http.Client
	}
)

func NewATProvider() ATProvider {
	return newATProviderWithCredentials(os.Getenv("AT_BASE_URL"), os.Getenv("AT_USERNAME"), os.Getenv("AT_KEY"))
}

func newATProviderWithCredentials(atBaseURL, username, key string) ATProvider {
	return &atProvider{
		atBaseURL: atBaseURL,
		username:  username,
		key:       key,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *atProvider) Send() error {

	atRequest := map[string]string{
		"username": p.username,
		"message":  "",
		"to":       "",
		"from":     "",
		"enqueue":  "1",
	}

	form := url.Values{}

	for key, value := range atRequest {
		form.Add(key, value)
	}

	request, err := http.NewRequest(http.MethodPost, p.atBaseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	header := make(http.Header)

	header.Set("Content-Length", strconv.Itoa(len(form)))
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Add("apikey", p.key)
	header.Add("Accept", "application/json")

	request.Header = header
	request.Close = true

	response, err := p.client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("at response %v", string(body))

	return nil
}
