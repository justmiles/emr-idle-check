package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type EMRInfo struct {
	JobFlowID            string `json:"jobFlowId"`
	InstanceCount        int    `json:"instanceCount"`
	MasterInstanceID     string `json:"masterInstanceId"`
	MasterPrivateDNSName string `json:"masterPrivateDnsName"`
	MasterInstanceType   string `json:"masterInstanceType"`
	SlaveInstanceType    string `json:"slaveInstanceType"`
	HadoopVersion        string `json:"hadoopVersion"`
	InstanceGroups       []struct {
		InstanceGroupID        string `json:"instanceGroupId"`
		InstanceRole           string `json:"instanceRole"`
		MarketType             string `json:"marketType"`
		InstanceType           string `json:"instanceType"`
		RequestedInstanceCount int    `json:"requestedInstanceCount"`
	} `json:"instanceGroups"`
}

func GetEMRInfo() (*EMRInfo, error) {

	jsonFile, err := os.Open("/mnt/var/lib/info/job-flow.json")
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var emrInfo EMRInfo
	err = json.Unmarshal(byteValue, &emrInfo)
	if err != nil {
		return nil, err
	}

	return &emrInfo, nil

}
