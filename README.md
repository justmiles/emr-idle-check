# emr-idle-check

Send EMR activity metrics to CloudWatch.

## Usage

```text
Usage:
  emr-idle-check [command]

Available Commands:
  help        Help about any command
  push        push the idle-check metrics to AWS CloudWatch
  view        view a snapshot of the current idle-check metrics

Flags:
  -h, --help      help for emr-idle-check
      --version   version for emr-idle-check

Use "emr-idle-check [command] --help" for more information about a command.
```

The below metrics are pushed to CloudWatch. All metrics are pushed to the `EMRActivity` namespace with a dimensions for the `JobFlowId` (EMR cluster ID).

| Metric                | Description                                                         |
| --------------------- | ------------------------------------------------------------------- |
| ActiveSSMSessions     | The number of active AWS SSM sessions                               |
| ActiveSSHSessions     | The number of active SSH sessions                                   |
| ActiveZepplinSessions | The number of active Zepplin sessions                               |
| ActiveYarnJobs        | The number of running or recently completed Yarn jobs               |
| ActiveEMRCluster      | Whether or not the cluster is active based on EMR's `isIdle` status |
