package cardinfo

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gitlab.cmpayments.local/creditcard/platform"
)

func ParseTr54(ctx context.Context, logger platform.Logger, r io.Reader) ([]Range, error) {
	scanner := bufio.NewScanner(r)

	var btrs []binTableRecord
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "H"):
			continue
		case strings.HasPrefix(line, "T"):
			break
		default:
			btrs = append(btrs, parseMastercardBINTableResourceFileLine(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ParseTr54(): %w", err)
	}

	ranges := map[uint64]Range{}
	var rng Range
	var rngs []Range

	for _, record := range btrs {
		if _, ok := ranges[record.lowPrimaryAccountRange]; ok {
			logger.Error(ctx, fmt.Sprintf("found duplicate mastercard bin: %d", record.lowPrimaryAccountRange))
			continue
		}

		rng = Range{
			Low:               strconv.FormatUint(record.lowPrimaryAccountRange, 10),
			High:              strconv.FormatUint(record.highPrimaryAccountRange, 10),
			Scheme:            "mastercard",
			ProgramID:         record.acceptanceBrand,
			ProductID:         record.brandProductCode,
			ProductName:       record.brandProductDescription,
			IssuerID:          record.customerID,
			IssuerName:        record.customerName,
			IssuerCountryCode: strconv.Itoa(record.issuingCountryCode),
		}

		ranges[record.lowPrimaryAccountRange] = rng

		rngs = append(rngs, rng)
	}

	return rngs, nil
}

type binTableRecord struct {
	recordTypeIdentifier      string
	lowPrimaryAccountRange    uint64
	highPrimaryAccountRange   uint64
	acceptanceBrand           string
	customerID                string
	customerName              string
	issuingCountryCode        int
	localUse                  string
	authorizationOnly         string
	brandProductCode          string
	brandProductDescription   string
	nonReloadableIndicator    int
	anonymousPrepaidIndicator string
}

func NewBINTableRecord() binTableRecord {
	return binTableRecord{}
}

func parseMastercardBINTableResourceFileLine(l string) binTableRecord {

	btr := NewBINTableRecord()

	btr.recordTypeIdentifier = strings.TrimSpace(l[0:1])
	btr.lowPrimaryAccountRange, _ = strconv.ParseUint(strings.TrimSpace(l[1:20]), 10, 64)
	btr.highPrimaryAccountRange, _ = strconv.ParseUint(strings.TrimSpace(l[20:39]), 10, 64)
	btr.acceptanceBrand = strings.TrimSpace(l[39:42])
	btr.customerID = strings.TrimSpace(l[42:53])
	btr.customerName = strings.TrimSpace(l[53:123])
	btr.issuingCountryCode, _ = strconv.Atoi(strings.TrimSpace(l[123:126]))
	btr.localUse = strings.TrimSpace(l[126:127])
	btr.authorizationOnly = strings.TrimSpace(l[127:128])
	btr.brandProductCode = strings.TrimSpace(l[128:131])
	btr.brandProductDescription = strings.TrimSpace(l[131:331])
	btr.nonReloadableIndicator, _ = strconv.Atoi(strings.TrimSpace(l[331:333]))
	btr.anonymousPrepaidIndicator = strings.TrimSpace(l[333:334])

	return btr
}
