package cardinfo

import (
	"context"
	"fmt"
	"io"

	"gitlab.cmpayments.local/creditcard/platform"
)

func NewService(c *Collection, s source, l platform.Logger, filetypes map[string]string, addVisaTestPans bool) Service {
	return Service{
		collection:      c,
		source:          s,
		log:             l,
		filetypes:       filetypes,
		addVisaTestPans: addVisaTestPans,
	}
}

type Service struct {
	collection      *Collection
	source          source
	log             platform.Logger
	filetypes       map[string]string
	addVisaTestPans bool
}

func (s Service) LoadBinRanges(ctx context.Context) error {
	s.log.Info(ctx, "start loading BIN range tables")

	for scheme, fileType := range s.filetypes {
		rc, err := s.source.LastFile(ctx, scheme, fileType)
		if err != nil {
			return fmt.Errorf("LoadBinRanges(): %w", err)
		}

		// If there is no new file, don't do anything
		if rc == nil {
			return fmt.Errorf("%s BIN range table not found", scheme)
		}

		defer rc.Close()

		s.log.Debug(ctx, fmt.Sprintf("found %s file, start parsing", s.filetypes))

		binRange := make([]Range, 0)

		switch scheme {
		case "mastercard":
			binRange, err = ParseTr54(ctx, s.log, rc)
			if err != nil {
				return fmt.Errorf("LoadBinRanges(): %w", err)
			}
		case "visa":
			binRange, err = ParseArdef(ctx, s.log, rc, s.addVisaTestPans)
			if err != nil {
				return fmt.Errorf("ParseArdef(): %w", err)
			}
		}

		s.log.Info(ctx, fmt.Sprintf("Found %d card range rules", len(binRange)))
		s.collection.Set(fileType, binRange)
	}

	return nil
}

func (s Service) LoadTest() error {
	testCol := []Range{
		{
			Low:               normalizePAN("4005529999000123"),
			High:              normalizePAN("4005529999000123"),
			Scheme:            "visa",
			ProductID:         "Visa Traditional Credit",
			ProductName:       "Test set",
			ProgramID:         "DEV",
			IssuerID:          "702879",
			IssuerName:        "Test",
			IssuerCountryCode: "840", //USA
		},
		{
			Low:               normalizePAN("2223001760002700"),
			High:              normalizePAN("2223001760002709"),
			Scheme:            "visa",
			ProductID:         "TST",
			ProductName:       "Test set",
			ProgramID:         "DEV",
			IssuerID:          "0987654321",
			IssuerName:        "Test",
			IssuerCountryCode: "528", //NLD
		}, {
			Low:               normalizePAN("5204740000001000"),
			High:              normalizePAN("5204740000001009"),
			Scheme:            "mastercard",
			ProductID:         "MCS",
			ProductName:       "STANDARD",
			ProgramID:         "MCC",
			IssuerID:          "00000999675",
			IssuerName:        "MTF INTERNAL MEMBER ID - LATVIA",
			IssuerCountryCode: "428",
		}, {
			Low:               normalizePAN("4111111111111110"),
			High:              normalizePAN("4111111145551143"),
			Scheme:            "visa",
			ProductID:         "TST",
			ProductName:       "Test set",
			ProgramID:         "DEV",
			IssuerID:          "0987654321",
			IssuerName:        "Test",
			IssuerCountryCode: "528", //NLD
		},
		{
			Low:               normalizePAN("4761349999000039"),
			High:              normalizePAN("4761349999000039"),
			Scheme:            "visa",
			ProductID:         "P",
			ProductName:       "VISA Gold Debit",
			ProgramID:         "CER", // CERTIFICATION
			IssuerID:          "702866",
			IssuerName:        "Testcard 3",
			IssuerCountryCode: "702", //SGD
		},
		{
			Low:               normalizePAN("4229989999000012"),
			High:              normalizePAN("4229989999000012"),
			Scheme:            "visa",
			ProductID:         "S1",
			ProductName:       "Visa Purchasing with Fleet Charge",
			ProgramID:         "CER", // CERTIFICATION
			IssuerID:          "702865",
			IssuerName:        "Testcard 1",
			IssuerCountryCode: "458", //MYS
		},
		{
			Low:               normalizePAN("4012001037141112"),
			High:              normalizePAN("4012001037141112"),
			Scheme:            "visa",
			ProductID:         "P",
			ProductName:       "PAN - ECI 5",
			ProgramID:         "CER", // CERTIFICATION
			IssuerID:          "476134",
			IssuerName:        "Testcard E1",
			IssuerCountryCode: "840", //USA
		},
		{
			Low:               normalizePAN("4012001037167778"),
			High:              normalizePAN("4012001037167778"),
			Scheme:            "visa",
			ProductID:         "P",
			ProductName:       "PAN - ECI 6",
			ProgramID:         "CER", // CERTIFICATION
			IssuerID:          "476134",
			IssuerName:        "Testcard E2",
			IssuerCountryCode: "840", //USA
		},
		{
			Low:               normalizePAN("4005520000000129"),
			High:              normalizePAN("4005520000000129"),
			Scheme:            "visa",
			ProductID:         "P",
			ProductName:       "PAN - ECI 7",
			ProgramID:         "CER", // CERTIFICATION
			IssuerID:          "702613",
			IssuerName:        "Testcard E3",
			IssuerCountryCode: "840", //USA
		},
		{
			Low:               normalizePAN("4176669999000104"),
			High:              normalizePAN("4176669999000104"),
			Scheme:            "visa",
			ProductID:         "P",
			ProductName:       "PAN - ECI 13",
			ProgramID:         "CER", // CERTIFICATION
			IssuerID:          "417666",
			IssuerName:        "Testcard E13",
			IssuerCountryCode: "826", //GBR
		},
	}

	s.collection.Set("test", testCol)

	return nil
}

type source interface {
	LastFile(ctx context.Context, scheme string, fileType string) (io.ReadCloser, error)
}
