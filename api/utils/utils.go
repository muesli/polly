package utils

import (
	"net/url"
	"strconv"

	"github.com/muesli/polly/api/db"
)

// BuildURL constructs a url from one or many items
func BuildURL(base string, items ...interface{}) string {
	proposal := ""
	user := ""

	for _, item := range items {
		switch v := item.(type) {
		case db.Proposal:
			proposal = url.QueryEscape(strconv.FormatInt(v.ID, 10))
		case db.User:
			user = url.QueryEscape(v.Username)
		}
	}

	switch items[0].(type) {
	case db.Proposal:
		return base + "proposal/" + proposal
	case db.User:
		return base + "user/" + user
	}

	return ""
}
