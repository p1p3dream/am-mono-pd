package gconf

type AWSVpcConfiguration struct {
	SecurityGroups []string `json:"security_groups,omitempty" yaml:"security_groups,omitempty"`
	Subnets        []string `json:"subnets,omitempty" yaml:"subnets,omitempty"`
	AssignPublicIp bool     `json:"assign_public_ip,omitempty" yaml:"assign_public_ip,omitempty"`
}

type AWSTaskNetworkConfiguration struct {
	AWSVpcConfiguration *AWSVpcConfiguration `json:"vpc_configuration,omitempty" yaml:"vpc_configuration,omitempty"`
}

type AWSTaskContainerOverride struct {
	Command     []string          `json:"command,omitempty" yaml:"command,omitempty"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
}
