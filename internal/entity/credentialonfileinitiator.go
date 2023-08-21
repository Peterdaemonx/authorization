package entity

type InitiatedBy string

const (
	Cardholder          InitiatedBy = "cardholder"
	MITRecurring        InitiatedBy = "mitRecurring"
	MITIndustryPractice InitiatedBy = "mitIndustryPractice"
)

var initiatedBys = map[string]InitiatedBy{
	"cardholder":          Cardholder,
	"mitRecurring":        MITRecurring,
	"mitIndustryPractice": MITIndustryPractice,
}

func IsValidInitiatedBy(i string) bool {
	_, ok := initiatedBys[i]
	return ok
}

func MapInitiatedByFromStr(i string) InitiatedBy {
	return initiatedBys[i]
}

type SubCategory string

const (
	CredentialOnFile            SubCategory = "credentialOnFile"
	StandingOrder               SubCategory = "standingOrder"
	Subscription                SubCategory = "subscription"
	Installment                 SubCategory = "installment"
	UnscheduledCredentialOnFile SubCategory = "unscheduledCredentialOnFile"
	PartialShipment             SubCategory = "partialShipment"
	DelayedCharge               SubCategory = "delayedCharge"
	NoShow                      SubCategory = "noShow"
	Resubmission                SubCategory = "resubmission"
)

var subCategories = map[string]SubCategory{
	"credentialOnFile":            CredentialOnFile,
	"standingOrder":               StandingOrder,
	"subscription":                Subscription,
	"installment":                 Installment,
	"unscheduledCredentialOnFile": UnscheduledCredentialOnFile,
	"partialShipment":             PartialShipment,
	"delayedCharge":               DelayedCharge,
	"noShow":                      NoShow,
	"resubmission":                Resubmission,
}

func IsValidSubCategory(i string) bool {
	_, ok := subCategories[i]
	return ok
}

func MapSubCategoryFromStr(i string) SubCategory {
	return subCategories[i]
}

type CitMitIndicator struct {
	InitiatedBy InitiatedBy
	SubCategory SubCategory
}

type schemeConfig struct {
	initiatedBy   string
	subCategories []string
}

var config = map[string][]schemeConfig{
	"mastercard": {
		schemeConfig{
			initiatedBy:   "cardholder",
			subCategories: []string{"credentialOnFile", "standingOrder", "subscription", "installment"},
		},
		schemeConfig{
			initiatedBy:   "mitRecurring",
			subCategories: []string{"unscheduledCredentialOnFile", "standingOrder", "subscription", "installment"},
		},
		schemeConfig{
			initiatedBy:   "mitIndustryPractice",
			subCategories: []string{"partialShipment", "delayedCharge", "noShow", "resubmission"},
		},
	},
	"visa": {
		schemeConfig{
			initiatedBy:   "mitIndustryPractice",
			subCategories: []string{"unscheduledCredentialOnFile", "resubmission", "delayedCharge", "noShow"},
		},
		schemeConfig{
			initiatedBy:   "cardholder",
			subCategories: []string{"credentialOnFile", "standingOrder", "subscription", "installment"},
		},
	},
}

func IsValidCitMit(scheme, initiatedBy, subCategory string) bool {
	sc, ok := config[scheme]
	if ok {
		for _, v := range sc {
			if v.initiatedBy == initiatedBy {
				for _, sc := range v.subCategories {
					if sc == subCategory {
						return true
					}
				}
			}
		}
	}

	return false
}
