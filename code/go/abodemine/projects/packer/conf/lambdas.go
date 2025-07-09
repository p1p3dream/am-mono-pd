package conf

type Lambdas struct {
	SecureDownload *SecureDownload `json:"secure_download,omitempty" yaml:"secure_download,omitempty"`
}

type SecureDownload struct {
	DynamodbTables map[string]*DynamodbTable `json:"dynamodb_tables,omitempty" yaml:"dynamodb_tables,omitempty"`
	S3Buckets      map[string]*S3Bucket      `json:"s3_buckets,omitempty" yaml:"s3_buckets,omitempty"`
}

type DynamodbTable struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type S3Bucket struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
