package aws

import (
	"fmt"

	"github.com/44smkn/s3select/pkg/build"

	"github.com/aws/aws-sdk-go/aws/request"
)

const appName = "s3select"

// injectUserAgent will inject app specific user-agent into awsSDK
func injectUserAgent(handlers *request.Handlers) {
	handlers.Build.PushFrontNamed(request.NamedHandler{
		Name: fmt.Sprintf("%s/user-agent", appName),
		Fn:   request.MakeAddToUserAgentHandler(appName, build.Version),
	})
}
