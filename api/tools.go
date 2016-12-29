package main

import (
	"net/url"
	"strconv"
)

func buildURL(items ...interface{}) string {
	proposal := ""
	user := ""

	for _, item := range items {
		switch v := item.(type) {
		case DbProposal:
			proposal = url.QueryEscape(strconv.FormatInt(v.ID, 10))
		case DbUser:
			user = url.QueryEscape(v.Username)
		}
	}

	switch items[0].(type) {
	case DbProposal:
		return config.Web.BaseURL + "proposal/" + proposal
	case DbUser:
		return config.Web.BaseURL + "user/" + user
	}

	return ""
}
