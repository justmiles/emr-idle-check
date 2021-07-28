package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ZepplinNotebook struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type ZepplinResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ZepplinNotebookList struct {
	ZepplinResponse
	Body []ZepplinNotebook `json:"body"`
}

type ZepplinNotebookJob struct {
	Progress string `json:"progress"`
	ID       string `json:"id"`
	Status   string `json:"status"`
}

type ZepplinNotebookListJobs struct {
	ZepplinResponse
	Body []ZepplinNotebookJob `json:"body"`
}

type ZepplinClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewZepplinClient(host string) *ZepplinClient {
	return &ZepplinClient{
		BaseURL: fmt.Sprintf("http://%s:8890/api", host),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetNotebooks returns a list of Zepplin notebooks
func (c *ZepplinClient) GetNotebooks(ctx context.Context) (*ZepplinNotebookList, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/notebook", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := ZepplinNotebookList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ListNotebooks returns a list of Zepplin notebooks
func (c *ZepplinClient) GetNotebookJobs(ctx context.Context, notebookID string) (*ZepplinNotebookListJobs, error) {
	endpoint := fmt.Sprintf("%s/notebook/job/%s", c.BaseURL, notebookID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := ZepplinNotebookListJobs{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *ZepplinClient) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return fmt.Errorf("Request failed: %s\n,%v", req.URL, err)
	}

	return nil
}

func getZepplinNotebooks() {

	url := "http://ip-10-140-1-107.prd.i-edo.net:8890/api/notebook"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
