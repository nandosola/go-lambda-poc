module main

go 1.19

replace (
	service => ../../service
	transport => ../../transport
)

require (
	github.com/aws/aws-lambda-go v1.41.0
	transport v0.0.0-00010101000000-000000000000
)
