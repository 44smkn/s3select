package aws

import (
	"fmt"

	"github.com/44smkn/s3select/pkg/build"

	awsrequest "github.com/aws/aws-sdk-go/aws/request"
)

const appName = "s3select"

// injectUserAgent will inject app specific user-agent into awsSDK
func injectUserAgent(handlers *awsrequest.Handlers) {
	handlers.Build.PushFrontNamed(awsrequest.NamedHandler{
		Name: fmt.Sprintf("%s/user-agent", appName),
		Fn:   awsrequest.MakeAddToUserAgentHandler(appName, build.Version),
	})
}
