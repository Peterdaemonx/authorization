package processing

import (
	"context"
)

type CardInfoFetcher struct {
	cis CardInfoService
}

func NewCardInfoFetcher(cis CardInfoService) *CardInfoFetcher {
	fetcher := new(CardInfoFetcher)
	fetcher.cis = cis
	return fetcher
}

type CardInfoService interface {
	FetchCardInfo(ctx context.Context, pan string) (CardInfo, error)
}

func (pt CardInfoFetcher) FetchCardInfo(ctx context.Context, pan string) (CardInfo, error) {
	return CardInfo{
		Info: Info{
			Low:    "000000000000",
			High:   "999999999999",
			Scheme: "mastercard",
			Programs: []Program{
				{
					ID:       "1",
					Priority: 1,
					Issuer: Issuer{
						Name:        "Test",
						Country:     "NLD",
						Website:     "www.test.nl",
						Phonenumber: "0987654321",
					},
				},
			},
		},
	}, nil
	//info, err := pt.cis.FetchCardInfo(ctx, pan)
	//if err != nil {
	//	return CardInfo{}, err
	//}
	//return info, nil
}
