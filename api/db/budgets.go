package db

import "fmt"

func (context *PollyContext) remainingMonths(month uint) uint {
	endmonth := context.Config.App.Proposals.StartMonth + context.Config.App.Proposals.TotalRuntimeMonths

	i := endmonth - month
	if i < 0 {
		return 0
	}

	return i
}

func (context *PollyContext) sumAcceptedLargeGrants() uint {
	return 0
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
	return context.remainingMonths(month) * context.Config.App.Proposals.MaxGrantValue * (context.Config.App.Proposals.MaxLargeGrantsPerMonth / 2)
}

func (context *PollyContext) remainingSmallGrantValue(month uint) uint {
	return context.Config.App.Proposals.TotalGrantValue - context.remainingLargeGrantValue(month) - context.sumAcceptedLargeGrants() - context.sumAcceptedSmallGrants(0)
}

// RemainingSmallGrantThisMonth returns the total available budget for small grants this month
func (context *PollyContext) RemainingSmallGrantThisMonth(month uint) uint {
	return context.remainingSmallGrantValue(month) / context.remainingMonths(month)
}

// SmallGrantMaxValue returns the max available value for a micro budget
func (context *PollyContext) SmallGrantMaxValue(month uint) uint {
	/* fmt.Println("remaining months:", context.remainingMonths(month))
	fmt.Println("rem-SMALL:", context.remainingSmallGrantValue(month))
	fmt.Println("rem-LARGE:", context.remainingLargeGrantValue(month))
	fmt.Println("sum-SMALL:", context.sumAcceptedSmallGrants(0))
	fmt.Println("sum-SMALL-MONTH:", context.sumAcceptedSmallGrants(month)) */

	i := int(context.remainingSmallGrantValue(month)/context.remainingMonths(month)) -
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
