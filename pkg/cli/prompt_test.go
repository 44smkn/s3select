package cli_test

import (
	"reflect"
	"testing"

	"github.com/44smkn/s3select/pkg/cli"
	"github.com/AlecAivazis/survey/v2"
)

func StubSurveyAskOne(t *testing.T, val string) func() {
	t.Helper()
	original := cli.SurveyAskOne
	cli.SurveyAskOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
		target := reflect.ValueOf(response)
		elem := target.Elem()
		elem.Set(reflect.ValueOf(val))
		return nil
	}
	return func() {
		cli.SurveyAskOne = original
	}
}
