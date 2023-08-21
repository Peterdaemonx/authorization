package cardinfo

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/btree"
)

type Scheme string

const (
	Mastercard Scheme = "mastercard"
	Visa       Scheme = "visa" // nolint:deadcode
)

// IssuerIdentificationNumber is a 8 character numberical string
type PermanentAccountNumber string

func normalizePAN(pan string) string {
	if len(pan) < 19 { //nolint:gomnd
		//nolint:gofmt,gofumpt,gomnd,goimports
		return pan + fmt.Sprintf("%0"+fmt.Sprintf("%d", 19-len(pan))+"d", 0)
	} else {
		return pan
	}
}

func NewCollection(blockedBins []string) *Collection {
	return &Collection{
		blockedBins: blockedBins,
		btree:       btree.New(2), //nolint:gomnd
		setMutex:    &sync.Mutex{},
	}
}

type Collection struct {
	blockedBins []string
	btree       *btree.BTree
	setMutex    *sync.Mutex
}

func (c *Collection) Set(source string, ranges []Range) {
	// Ensure only one Set process happens at the same time
	c.setMutex.Lock()
	defer c.setMutex.Unlock()

	// Clone the tree, so the change-over is atomic
	bt := c.btree.Clone()

	// Delete everything for the scheme that is being updated
	bt.Ascend(func(i btree.Item) bool {
		if i.(Range).Source == source {
			bt.Delete(i)
		}
		return true
	})

	// Add all scheme nodes
	for _, r := range ranges {
		// Ensure all the data is as we need it
		r.Low = normalizePAN(r.Low)
		r.High = normalizePAN(r.High)
		r.Source = source
		r.IsBlocked = c.isBlocked(r.Low)
		bt.ReplaceOrInsert(r)
	}

	c.btree = bt
}

func (c Collection) isBlocked(low string) bool {
	for _, bin := range c.blockedBins {
		if strings.HasPrefix(low, bin) {
			return true
		}
	}
	return false
}

func (c Collection) Find(pan string) (r Range, ok bool) {
	c.btree.DescendLessOrEqual(Range{Low: normalizePAN(pan)}, func(i btree.Item) bool {
		fr := i.(Range)
		if fr.High > pan {
			r = fr
			ok = true

			return false
		}
		return true
	})

	return
}

type Range struct {
	Low               string
	High              string
	Source            string
	Scheme            string
	ProductID         string
	ProductName       string
	ProgramID         string
	IssuerID          string
	IssuerName        string
	IssuerCountryCode string
	IsBlocked         bool
}

func (r Range) Less(than btree.Item) bool {
	return than.(Range).Low > r.Low
}
