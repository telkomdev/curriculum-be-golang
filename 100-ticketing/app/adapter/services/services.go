package services

import (
	"100-ticketing/app/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Services struct {
	Client         *http.Client
	EndpointUser   string
	EndpointRoute  string
	EndpointTicket string
}

var services *Services

func init() {
	services = &Services{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		EndpointUser:   "http://localhost:8000/api/v1/user",
		EndpointRoute:  "http://localhost:8000/api/v1/route",
		EndpointTicket: "http://localhost:8000/api/v1/ticket",
	}

	// load from env if available
	if len(os.Getenv("SERVICE_ENDPOINT_USER")) != 0 {
		services.EndpointUser = os.Getenv("SERVICE_ENDPOINT_USER")
	}

	if len(os.Getenv("SERVICE_ENDPOINT_ROUTE")) != 0 {
		services.EndpointRoute = os.Getenv("SERVICE_ENDPOINT_ROUTE")
	}

	if len(os.Getenv("SERVICE_ENDPOINT_TICKET")) != 0 {
		services.EndpointRoute = os.Getenv("SERVICE_ENDPOINT_TICKET")
	}
}

// Load services object
func Load() *Services {
	return services
}

func (c *Services) Get(token, url string, query map[string]string, result interface{}) (err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	return c.Do(token, query, req, result)
}

func (c *Services) Post(token, url string, query map[string]string, data interface{}, result interface{}) (err error) {
	dat, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewReader(dat))
	if err != nil {
		return err
	}
	return c.Do(token, query, req, result)
}

func (c *Services) Do(token string, query map[string]string, req *http.Request, result interface{}) (err error) {
	if len(token) != 0 {
		req.Header.Set("Authorization", token)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("empty response object")
	}

	if resp.Body == nil {
		return fmt.Errorf("empty body")
	}

	if resp.StatusCode != 200 {

		all, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		var response models.ResponseItem
		err = json.Unmarshal(all, &response)
		if err != nil {
			return err
		}

		return fmt.Errorf(response.Message)
	}

	all, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	err = json.Unmarshal(all, result)
	if err != nil {
		return err
	}

	return nil
}
