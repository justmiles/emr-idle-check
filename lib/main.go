package lib

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/keybase/go-ps"
)

var (
	instanceID, jobFlowId, hostname, availablityZone, region string
	sess                                                     *session.Session
	ec2metadataService                                       *ec2metadata.EC2Metadata
	cloudwatchService                                        *cloudwatch.CloudWatch
)

type IdleMetrics struct {
	ActiveSSMSessions     int
	ActiveSSHSessions     int
	ActiveZepplinSessions int
	ActiveYarnJobs        int
	ActiveEMRCluster      int
}

func (s IdleMetrics) String() string {
	return fmt.Sprintf(
		"Active SSM Sessions: %d\nActive SSH Sessions: %d\nActive Zepplin Sessions: %d\nActive Yarn Jobs: %d\nActive EMR Cluster: %d",
		s.ActiveSSMSessions,
		s.ActiveSSHSessions,
		s.ActiveZepplinSessions,
		s.ActiveYarnJobs,
		s.ActiveEMRCluster,
	)
}

// GetIdleCheckMetrics returns a struct of idleCheck metrics
func PushIdleCheckMetricsToCloudWatch() error {
	idleMetrics := GetIdleCheckMetrics()

	emrInfo, err := GetEMRInfo()
	if err != nil {
		return err
	}

	dimensions := []*cloudwatch.Dimension{
		{
			Name:  aws.String("JobFlowId"),
			Value: &emrInfo.JobFlowID,
		},
	}

	_, err = cloudwatchService.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("EMRActivity"),
		MetricData: []*cloudwatch.MetricDatum{
			{
				Dimensions: dimensions,
				MetricName: aws.String("ActiveSSMSessions"),
				Value:      aws.Float64(float64(idleMetrics.ActiveSSMSessions)),
			},
			{
				Dimensions: dimensions,
				MetricName: aws.String("ActiveSSHSessions"),
				Value:      aws.Float64(float64(idleMetrics.ActiveSSHSessions)),
			},
			{
				Dimensions: dimensions,
				MetricName: aws.String("ActiveZepplinSessions"),
				Value:      aws.Float64(float64(idleMetrics.ActiveZepplinSessions)),
			},
			{
				Dimensions: dimensions,
				MetricName: aws.String("ActiveYarnJobs"),
				Value:      aws.Float64(float64(idleMetrics.ActiveYarnJobs)),
			},
			{
				Dimensions: dimensions,
				MetricName: aws.String("ActiveEMRCluster"),
				Value:      aws.Float64(float64(idleMetrics.ActiveEMRCluster)),
			},
		},
	})

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

// GetIdleCheckMetrics returns a struct of idleCheck metrics
func GetIdleCheckMetrics() IdleMetrics {

	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2metadataService = ec2metadata.New(sess)
	region, _ = ec2metadataService.Region()

	cloudwatchService = cloudwatch.New(sess, aws.NewConfig().WithRegion(region))

	instanceID, _ = ec2metadataService.GetMetadata("instance-id")
	hostname, _ = ec2metadataService.GetMetadata("hostname")
	availablityZone, _ = ec2metadataService.GetMetadata("placement/availability-zone")

	var idleMetrics IdleMetrics
	var err error

	idleMetrics.ActiveSSMSessions, err = CheckActiveSSMSessions()
	if err != nil {
		fmt.Println(err)
	}

	idleMetrics.ActiveZepplinSessions, err = CheckActiveZepplinSessions()
	if err != nil {
		fmt.Println(err)
	}

	idleMetrics.ActiveYarnJobs, err = CheckActiveYarnJobs()
	if err != nil {
		fmt.Println(err)
	}

	idleMetrics.ActiveSSHSessions, err = CheckActiveSSHSessions()
	if err != nil {
		fmt.Println(err)
	}

	idleMetrics.ActiveEMRCluster, err = CheckActiveEMRCluster()
	if err != nil {
		fmt.Println(err)
	}

	return idleMetrics

}

// CheckActiveEMRCluster returns 1 if the EMR cluster is in use, as determined by the
//  AWS/ElasticMapReduce/IsIdle metric stored in CloudWatch
func CheckActiveEMRCluster() (int, error) {
	now := time.Now()
	then := now.Add(time.Duration(-5) * time.Minute)
	emrInfo, err := GetEMRInfo()
	if err != nil {
		return 1, err
	}

	result, err := cloudwatchService.GetMetricData(&cloudwatch.GetMetricDataInput{
		EndTime:   &now,
		StartTime: &then,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{
			{
				Id: aws.String("q1"),
				MetricStat: &cloudwatch.MetricStat{
					Period: aws.Int64(300),
					Stat:   aws.String("Minimum"),
					Metric: &cloudwatch.Metric{
						MetricName: aws.String("IsIdle"),
						Namespace:  aws.String("AWS/ElasticMapReduce"),
						Dimensions: []*cloudwatch.Dimension{
							{
								Name:  aws.String("JobFlowId"),
								Value: &emrInfo.JobFlowID,
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		return 1, err
	}

	for _, metric := range result.MetricDataResults {
		for _, value := range metric.Values {
			if *value == 0 {
				return 1, nil
			}
		}
	}

	return 0, nil
}

// CheckActiveYarnJobs returns the number of running Yarn jobs
func CheckActiveYarnJobs() (int, error) {
	var count int
	yarnClient := NewYarnClient(hostname)

	ctx := context.Background()

	// Get the count of Running Jobs
	jobs, err := yarnClient.GetJobsByState(ctx, "RUNNING")
	if err != nil {
		return 0, err
	}
	count = len(jobs.Apps.App)

	// Get all the Finished jobs
	jobs, err = yarnClient.GetJobsByState(ctx, "FINISHED")
	if err != nil {
		return 0, err
	}

	// If a job finished within the last 5 minutes, consider it a running job
	now := time.Now()
	for _, job := range jobs.Apps.App {
		durationSinceFinished := now.Sub(time.Unix(job.FinishedTime/1000, 0))
		if durationSinceFinished < time.Duration(5*time.Minute) {
			fmt.Println(durationSinceFinished)
			count++
		}
	}

	return count, nil
}

// CheckActiveZepplinSessions returns the number of act,ive Zepplin Sessions
func CheckActiveZepplinSessions() (int, error) {
	var count int

	zepplinClient := NewZepplinClient(hostname)

	ctx := context.Background()
	notebooks, err := zepplinClient.GetNotebooks(ctx)
	if err != nil {
		return count, err
	}

	for _, notebook := range notebooks.Body {

		notebookJob, err := zepplinClient.GetNotebookJobs(ctx, notebook.ID)
		if err != nil {
			fmt.Println(err)
		}

		for _, job := range notebookJob.Body {
			if job.Status == "RUNNING" {
				count++
			}
		}
	}

	return count, nil
}

// CheckActiveSSHSessions returns the number of active SSH sessions
func CheckActiveSSHSessions() (int, error) {
	return countMatchingProcesses("ssh")
}

// CheckActiveSSMSessions returns the number of active SSM sessions
func CheckActiveSSMSessions() (int, error) {
	return countMatchingProcesses("ssm-session-worker")
}

// countProcessRunning returns the number of processes running that matches the provided input
func countMatchingProcesses(s string) (int, error) {
	var processCount int
	processList, err := ps.Processes()
	if err != nil {
		return processCount, fmt.Errorf("unable to list running processes")
	}

	for x := range processList {
		var process ps.Process
		process = processList[x]
		if process.Executable() == s {
			processCount++
		}
	}

	return processCount, nil
}
