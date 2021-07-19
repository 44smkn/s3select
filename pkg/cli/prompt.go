package cli

import (
	"reflect"

	"github.com/AlecAivazis/survey/v2"
)

var SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(p, response, opts...)
}

func Select(msg string, options []string, defaultVal string) string {
	prompt := &survey.Select{
		Message: msg,
		Options: options,
	}
	var val string
	_ = SurveyAskOne(prompt, &val)
	return val
}

func StubSurveyAskOne(val string) func() {
	original := SurveyAskOne
	SurveyAskOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
		target := reflect.ValueOf(response)
		elem := target.Elem()
		elem.Set(reflect.ValueOf(val))
		return nil
	}
	return func() {
		SurveyAskOne = original
	}
}
