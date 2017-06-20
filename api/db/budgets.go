package db

import (
	"fmt"
	"math"
)

func (context *PollyContext) remainingMonths(month uint) uint {
	endmonth := context.Config.App.Proposals.StartMonth + context.Config.App.Proposals.TotalRuntimeMonths

	i := endmonth - month
	if i < 0 {
		return 0
	}

	return i
}

func (context *PollyContext) remainingLargeGrantPeriods(month uint) uint {
	return uint(math.Ceil(float64(context.remainingMonths(month)) / 2.0))
}

func (context *PollyContext) sumAcceptedLargeGrants() uint {
	proposals, _ := context.LoadAllProposals()

	i := uint(0)
	for _, p := range proposals {
		if !p.Ended(context) {
			continue
		}
		if p.Value < uint64(context.Config.App.Proposals.SmallGrantValueThreshold) {
			continue
		}
		if !p.Accepted(context) {
			continue
		}

		i += uint(p.Value)
	}

	return i
}

func (context *PollyContext) sumAcceptedSmallGrants(month uint) uint {
	proposals, _ := context.LoadAllProposals()

	i := uint(0)
	for _, p := range proposals {
		if month > 0 && uint(p.Ends(context).Month()) != month {
			continue
		}
		if !p.Ended(context) {
			continue
		}
		if p.Value >= uint64(context.Config.App.Proposals.SmallGrantValueThreshold) {
			continue
		}
		fmt.Println("FOUND:", p.Title, p.Starts, p.Value, p.Votes)
		if !p.Accepted(context) {
			continue
		}
		fmt.Println("ACCEPTING:", p.Title, p.Starts, p.Value, p.Votes)

		i += uint(p.Value)
	}

	return i
}

func (context *PollyContext) remainingLargeGrantValue(month uint) uint {
	return context.remainingLargeGrantPeriods(month) * context.Config.App.Proposals.MaxGrantValue * context.Config.App.Proposals.MaxLargeGrantsPerMonth
}

func (context *PollyContext) remainingSmallGrantValue(month uint) uint {
	return context.Config.App.Proposals.TotalGrantValue - context.remainingLargeGrantValue(month) - context.sumAcceptedLargeGrants() - context.sumAcceptedSmallGrants(0)
}

// RemainingSmallGrantThisMonth returns the total available budget for small grants this month
func (context *PollyContext) RemainingSmallGrantThisMonth(month uint) uint {
	remmonths := context.remainingMonths(month)
	if remmonths < 1 {
		remmonths = 1
	}
	return context.remainingSmallGrantValue(month) / remmonths
}

// SmallGrantMaxValue returns the max available value for a micro budget
func (context *PollyContext) SmallGrantMaxValue(month uint) uint {
	/* fmt.Println("remaining months:", context.remainingMonths(month))
	fmt.Println("rem-SMALL:", context.remainingSmallGrantValue(month))
	fmt.Println("rem-LARGE:", context.remainingLargeGrantValue(month))
	fmt.Println("sum-SMALL:", context.sumAcceptedSmallGrants(0))
	fmt.Println("sum-SMALL-MONTH:", context.sumAcceptedSmallGrants(month)) */

	remmonths := context.remainingMonths(month)
	if remmonths < 1 {
		remmonths = 1
	}

	i := int(context.remainingSmallGrantValue(month)/remmonths) -
		int(context.sumAcceptedSmallGrants(month))

	// fmt.Println("smallgrantmax:", i)

	if i < 0 {
		return 0
	}
	if i >= int(context.Config.App.Proposals.SmallGrantValueThreshold) {
		return context.Config.App.Proposals.SmallGrantValueThreshold - 1
	}

	return uint(i)
}

// GrantMaxValue returns the max allowed grant value
func (context *PollyContext) GrantMaxValue() uint {
	return context.Config.App.Proposals.MaxGrantValue
}
