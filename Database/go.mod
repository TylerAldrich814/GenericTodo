module todoDatabase

go 1.20

require github.com/aws/aws-sdk-go v1.44.317

require utils v1.0.0

require (
	github.com/aws/aws-lambda-go v1.41.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace utils v1.0.0 => ../Utils
