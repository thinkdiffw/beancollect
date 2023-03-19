package types

import (
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// Transaction means a real transaction in beancount.
type Transaction struct {
	Time      time.Time
	Flag      string
	Narration string
	Payee     string
	Tags      []string
	Accounts  []string
	Amount    float64
	Currency  string

	// Metadata map[string]string
}

var regexCache = make(map[string]*regexp.Regexp, 1)

// IsMatch will check whether this transaction match condition.
func (t Transaction) IsMatch(cond map[string]string) bool {
	// Currently, we only support payee.
	return strings.Contains(t.Payee, cond["payee"]) || regexMatch(t.Payee, cond["payee"])
}

func regexMatch(str, pattern string) bool {
	x, ok := regexCache[pattern]
	if ok {
		return x.MatchString(str)
	}
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Errorf("regex parse failed for %s", err)
		return false
	}
	regexCache[pattern] = regex
	return regex.MatchString(str)
}

// Transactions is the array for transactions
type Transactions []Transaction

// Len implement Sorter.Len
func (t Transactions) Len() int {
	return len(t)
}

// Less implement Sorter.Less
func (t Transactions) Less(i, j int) bool {
	return t[i].Time.Before(t[j].Time)
}

// Swap implement Sorter.Swap
func (t Transactions) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
