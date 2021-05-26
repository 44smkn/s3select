package aws

import (
	"fmt"

	"github.com/44smkn/s3selecgo/pkg/version"

	"github.com/aws/aws-sdk-go/aws/request"
)

const appName = "s3selecgo"

// injectUserAgent will inject app specific user-agent into awsSDK
func injectUserAgent(handlers *request.Handlers) {
	handlers.Build.PushFrontNamed(request.NamedHandler{
		Name: fmt.Sprintf("%s/user-agent", appName),
		Fn:   request.MakeAddToUserAgentHandler(appName, version.GitVersion),
	})
}
