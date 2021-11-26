package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// CloudWatchEventDetails...jsonで返ってくるeventで必要な情報をストアする
type CloudWatchEventDetails struct {
	Containers []struct {
		Name     string `json:"name"`
		ExitCode int    `json:"exitCode"`
	}
	ClusterArn        string `json:"clusterArn"`
	TaskDefinitionArn string `json:"taskDefinitionArn"`
}

// ToCloudWatchEventDetailsStruct...jsonで返ってくるeventを、必要な情報だけstructに入れる
func ToCloudWatchEventDetailsStruct(event events.CloudWatchEvent) (*CloudWatchEventDetails, error) {
	e := &eventDetails
	err := json.Unmarshal(event.Detail, e)
	if err != nil {
		return e, fmt.Errorf("could not unmarshal cloudwatch events: %v\n", err)
	}
	return e, nil
}

// FormatToSlackMessage...slack用にformatする
func (c *CloudWatchEventDetails) FormatToSlackMessage(containerStatus map[string]int) string {
	// account idからnameを出す
	awsAccountName := AwsAccountIDs[strings.Split(c.ClusterArn, ":")[4]]

	// clusterの名前だけ取得する
	clusterName := strings.Split(c.ClusterArn, "/")[1]

	// task名を取得する
	taskName := strings.Split(c.TaskDefinitionArn, "/")[1]
	taskName = strings.Split(taskName, ":")[0]

	var msg string

	for k, v := range containerStatus {
		msg += fmt.Sprintf("container: %s, exit code: %d\n", k, v)
	}
	fmt.Println(msg)

	return fmt.Sprintf("```\n"+
		"aws account: %s\n"+
		"cluster:     %s\n"+
		"task:        %s\n\n"+
		"%s"+
		"```\n", awsAccountName, clusterName, taskName, msg)
}
