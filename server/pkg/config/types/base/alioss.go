package base

type AliOSS struct {
	Endpoint        *string `mapstructure:"endpoint"`
	AccessKeyId     string  `mapstructure:"accessKeyId"`
	AccessKeySecret string  `mapstructure:"accessKeySecret"`
	BucketName      string  `mapstructure:"bucketName"`
	Region          string  `mapstructure:"region"`
}
