package sequences

func maxValRollover(max int) RolloverCondition {
	return func(b Block) bool {
		return b.NextValue > max
	}
}

func diffRolloverRollover(v func() string) RolloverCondition {
	return func(b Block) bool {
		return b.RolloverValue != v()
	}
}

func anyRollover(checks ...RolloverCondition) RolloverCondition {
	return func(b Block) bool {
		for _, check := range checks {
			if check(b) {
				return true
			}
		}
		return false
	}
}
