package risk

import (
	"context"
)

type Assessor struct {
	assessmentRules []Rule
}

func NewAssessor(assessmentRules []Rule) *Assessor {
	return &Assessor{
		assessmentRules: assessmentRules,
	}
}

func (a Assessor) Process(ctx context.Context, arg interface{}) error {
	var assessment Assessment
	for _, assessmentRule := range a.assessmentRules {
		if assessmentRule.Support(arg) {
			err := assessmentRule.Eval(ctx, arg, &assessment)
			if err != nil {
				return err
			}
		}

		if assessment.Completed {
			break
		}
	}

	return nil
}
