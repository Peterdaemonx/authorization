package cardinfo

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gitlab.cmpayments.local/creditcard/platform"
	"gitlab.cmpayments.local/creditcard/platform/countrycode"
)

type visaBinTableRecord struct {
	lowPrimaryAccountRange  uint64
	highPrimaryAccountRange uint64
	customerID              string
	issuingCountry          string
	brandProductCode        string
}

func parseVisaBinTableResourceFileLine(l string) visaBinTableRecord {
	btr := visaBinTableRecord{}

	btr.lowPrimaryAccountRange, _ = strconv.ParseUint(strings.TrimSpace(l[12:21]), 10, 64)
	btr.highPrimaryAccountRange, _ = strconv.ParseUint(strings.TrimSpace(l[0:9]), 10, 64)
	btr.brandProductCode = strings.TrimSpace(l[58:60])
	btr.customerID = strings.TrimSpace(l[25:30])
	btr.issuingCountry = strings.TrimSpace(l[43:45])

	return btr
}

func ParseArdef(ctx context.Context, logger platform.Logger, r io.Reader, addVisaTestPans bool) ([]Range, error) {
	scanner := bufio.NewScanner(r)

	var btrs []visaBinTableRecord
	for scanner.Scan() {
		line := scanner.Text()
		btrs = append(btrs, parseVisaBinTableResourceFileLine(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ParseArdef(): %w", err)
	}

	ranges := map[uint64]Range{}
	var rng Range
	var rngs []Range

	for _, record := range btrs {
		if _, ok := ranges[record.lowPrimaryAccountRange]; ok {
			logger.Error(ctx, fmt.Sprintf("found duplicate visa bin: %d", record.lowPrimaryAccountRange))
			continue
		}
		country, err := countrycode.VisaFromAlpha2(record.issuingCountry)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("countrycode.FromAlpha2(): %s", err.Error()))
		}
		rng = Range{
			Low:               strconv.FormatUint(record.lowPrimaryAccountRange, 10),
			High:              strconv.FormatUint(record.highPrimaryAccountRange, 10),
			Scheme:            "visa",
			ProductID:         record.brandProductCode,
			IssuerID:          record.customerID,
			IssuerCountryCode: country.Visa.Numeric(),
		}

		ranges[record.lowPrimaryAccountRange] = rng

		rngs = append(rngs, rng)
	}

	if addVisaTestPans {
		rngs = addCertificationRanges(rngs)
	}

	return rngs, nil
}

func addCertificationRanges(rngs []Range) []Range {
	rngsWithTestCards := rngs

	rngsWithTestCards = append(rngsWithTestCards, Range{
		Low:               "4229989999000012",
		High:              "4229989999000012",
		Scheme:            "visa",
		ProductID:         "P",
		IssuerID:          "702866",
		IssuerCountryCode: "702",
	})
	rngsWithTestCards = append(rngsWithTestCards, Range{
		Low:               "4012001037141112",
		High:              "4012001037141112",
		Scheme:            "visa",
		IssuerID:          "476134",
		IssuerCountryCode: "840",
	})
	rngsWithTestCards = append(rngsWithTestCards, Range{
		Low:               "4012001037167778",
		High:              "4012001037167778",
		Scheme:            "visa",
		IssuerID:          "476134",
		IssuerCountryCode: "840",
	})
	rngsWithTestCards = append(rngsWithTestCards, Range{
		Low:               "4005529999000123",
		High:              "4005529999000123",
		Scheme:            "visa",
		ProductID:         "A",
		IssuerID:          "702879",
		IssuerCountryCode: "840",
	})
	rngsWithTestCards = append(rngsWithTestCards, Range{
		Low:               "4005520000000129",
		High:              "4005520000000129",
		Scheme:            "visa",
		IssuerID:          "702613",
		IssuerCountryCode: "840",
	})
	rngsWithTestCards = append(rngsWithTestCards, Range{
		Low:               "4176669999000104",
		High:              "4176669999000104",
		Scheme:            "visa",
		ProductID:         "F",
		IssuerID:          "702872",
		IssuerCountryCode: "826",
	})
	rngsWithTestCards = append(rngsWithTestCards, Range{
		Low:               "4761349999000039",
		High:              "4761349999000039",
		Scheme:            "visa",
		ProductID:         "P",
		IssuerID:          "702866",
		IssuerCountryCode: "702",
	})

	return rngsWithTestCards
}
