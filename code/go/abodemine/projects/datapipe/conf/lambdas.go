package conf

import "abodemine/lib/gconf"

type Lambdas struct {
	TaskLauncher *TaskLauncher `json:"task_launcher,omitempty" yaml:"task_launcher,omitempty"`
}

type TaskLauncher struct {
	SqsQueueUrl string             `json:"sqs_queue_url,omitempty" yaml:"sqs_queue_url,omitempty"`
	Tasks       *TaskLauncherTasks `json:"tasks,omitempty" yaml:"tasks,omitempty"`
}

type TaskLauncherTasks struct {
	Fetcher  *TaskLauncherFetcherTask  `json:"fetcher,omitempty" yaml:"fetcher,omitempty"`
	Loader   *TaskLauncherLoaderTask   `json:"loader,omitempty" yaml:"loader,omitempty"`
	Osloader *TaskLauncherOsloaderTask `json:"osloader,omitempty" yaml:"osloader,omitempty"`
	Synther  *TaskLauncherSyntherTask  `json:"synther,omitempty" yaml:"synther,omitempty"`
}

type TaskLauncherFetcherTask struct {
	ClusterArn           string                                       `json:"cluster_arn,omitempty" yaml:"cluster_arn,omitempty"`
	ContainerOverrides   map[string][]*gconf.AWSTaskContainerOverride `json:"container_overrides,omitempty" yaml:"container_overrides,omitempty"`
	NetworkConfiguration *gconf.AWSTaskNetworkConfiguration           `json:"network_configuration,omitempty" yaml:"network_configuration,omitempty"`
	TaskDefinitionArn    string                                       `json:"task_definition_arn,omitempty" yaml:"task_definition_arn,omitempty"`
}

type TaskLauncherLoaderTask struct {
	ClusterArn           string                                       `json:"cluster_arn,omitempty" yaml:"cluster_arn,omitempty"`
	ContainerOverrides   map[string][]*gconf.AWSTaskContainerOverride `json:"container_overrides,omitempty" yaml:"container_overrides,omitempty"`
	NetworkConfiguration *gconf.AWSTaskNetworkConfiguration           `json:"network_configuration,omitempty" yaml:"network_configuration,omitempty"`
	TaskDefinitionArn    string                                       `json:"task_definition_arn,omitempty" yaml:"task_definition_arn,omitempty"`
}

type TaskLauncherOsloaderTask struct {
	ClusterArn           string                                       `json:"cluster_arn,omitempty" yaml:"cluster_arn,omitempty"`
	ContainerOverrides   map[string][]*gconf.AWSTaskContainerOverride `json:"container_overrides,omitempty" yaml:"container_overrides,omitempty"`
	NetworkConfiguration *gconf.AWSTaskNetworkConfiguration           `json:"network_configuration,omitempty" yaml:"network_configuration,omitempty"`
	TaskDefinitionArn    string                                       `json:"task_definition_arn,omitempty" yaml:"task_definition_arn,omitempty"`
}

type TaskLauncherSyntherTask struct {
	ClusterArn           string                                       `json:"cluster_arn,omitempty" yaml:"cluster_arn,omitempty"`
	ContainerOverrides   map[string][]*gconf.AWSTaskContainerOverride `json:"container_overrides,omitempty" yaml:"container_overrides,omitempty"`
	NetworkConfiguration *gconf.AWSTaskNetworkConfiguration           `json:"network_configuration,omitempty" yaml:"network_configuration,omitempty"`
	TaskDefinitionArn    string                                       `json:"task_definition_arn,omitempty" yaml:"task_definition_arn,omitempty"`
}
