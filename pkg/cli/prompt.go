package cli

import "github.com/AlecAivazis/survey/v2"

var SurvayAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(p, response, opts...)
}
