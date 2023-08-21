package entity

type Source string

const (
	Ecommerce   Source = "ecommerce"
	Moto        Source = "moto"
	CardPresent Source = "cardPresent"
)

var (
	sourceMap = map[string]Source{
		`ecommerce`:   Ecommerce,
		`moto`:        Moto,
		`cardPresent`: CardPresent,
	}
)

func IsValidSource(source string) bool {
	_, ok := sourceMap[source]
	return ok
}

func MapSource(source string) (Source, bool) {
	v, ok := sourceMap[source]
	return v, ok
}
