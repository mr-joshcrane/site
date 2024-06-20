package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mr-joshcrane/site"
)

func main() {
	lambda.StartHandlerFunc(site.Handler)
}
