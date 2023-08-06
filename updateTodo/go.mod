module updateTodo

go 1.20

require utils v1.0.0

require todoDatabase v1.0.0

require github.com/aws/aws-lambda-go v1.41.0

require (
	github.com/aws/aws-sdk-go v1.44.317 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace utils v1.0.0 => ../Utils

replace todoDatabase v1.0.0 => ../Database
