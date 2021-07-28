package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type YarnApplicationList struct {
	Apps struct {
		App []YarnApplication `json:"app,omitempty"`
	} `json:"apps"`
}

type YarnApplication struct {
	ID                         string  `json:"id"`
	User                       string  `json:"user"`
	Name                       string  `json:"name"`
	Queue                      string  `json:"queue"`
	State                      string  `json:"state"`
	FinalStatus                string  `json:"finalStatus"`
	Progress                   float64 `json:"progress"`
	TrackingUI                 string  `json:"trackingUI"`
	TrackingURL                string  `json:"trackingUrl"`
	Diagnostics                string  `json:"diagnostics"`
	ClusterID                  int64   `json:"clusterId"`
	ApplicationType            string  `json:"applicationType"`
	ApplicationTags            string  `json:"applicationTags"`
	Priority                   int     `json:"priority"`
	StartedTime                int64   `json:"startedTime"`
	FinishedTime               int64   `json:"finishedTime"`
	ElapsedTime                int     `json:"elapsedTime"`
	AmContainerLogs            string  `json:"amContainerLogs"`
	AmHostHTTPAddress          string  `json:"amHostHttpAddress"`
	AmRPCAddress               string  `json:"amRPCAddress"`
	AllocatedMB                int     `json:"allocatedMB"`
	AllocatedVCores            int     `json:"allocatedVCores"`
	RunningContainers          int     `json:"runningContainers"`
	MemorySeconds              int     `json:"memorySeconds"`
	VcoreSeconds               int     `json:"vcoreSeconds"`
	QueueUsagePercentage       float64 `json:"queueUsagePercentage"`
	ClusterUsagePercentage     float64 `json:"clusterUsagePercentage"`
	PreemptedResourceMB        int     `json:"preemptedResourceMB"`
	PreemptedResourceVCores    int     `json:"preemptedResourceVCores"`
	NumNonAMContainerPreempted int     `json:"numNonAMContainerPreempted"`
	NumAMContainerPreempted    int     `json:"numAMContainerPreempted"`
	PreemptedMemorySeconds     int     `json:"preemptedMemorySeconds"`
	PreemptedVcoreSeconds      int     `json:"preemptedVcoreSeconds"`
	LogAggregationStatus       string  `json:"logAggregationStatus"`
	UnmanagedApplication       bool    `json:"unmanagedApplication"`
	AmNodeLabelExpression      string  `json:"amNodeLabelExpression"`
}

type YarnClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewYarnClient(host string) *YarnClient {
	return &YarnClient{
		BaseURL: fmt.Sprintf("http://%s:8088/ws/v1", host),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// GetNotebooks returns a list of Zepplin notebooks
// state can be one of [NEW, NEW_SAVING, SUBMITTED, ACCEPTED, RUNNING, FINISHED, FAILED, KILLED]
func (c *YarnClient) GetJobsByState(ctx context.Context, state string) (*YarnApplicationList, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cluster/apps?state=%s", c.BaseURL, state), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := YarnApplicationList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *YarnClient) sendRequest(req *http.Request, v interface{}) error {
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
