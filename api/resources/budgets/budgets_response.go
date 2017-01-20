package budgets

import "github.com/muesli/smolder"

// BudgetResponse is the common response to 'budget' requests
type BudgetResponse struct {
	smolder.Response

	Budgets []budgetInfoResponse `json:"budgets,omitempty"`
}

type budgetInfoResponse struct {
	ID    uint `json:"id"`
	Value uint `json:"value"`
}

// Init a new response
func (r *BudgetResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Budgets = []budgetInfoResponse{}
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *BudgetResponse) EmptyResponse() interface{} {
	if len(r.Budgets) == 0 {
		var out struct {
			Budgets interface{} `json:"budgets"`
		}
		out.Budgets = []budgetInfoResponse{}
		return out
	}
	return nil
}

func prepareBudgetResponse(context smolder.APIContext, month uint, budget uint) budgetInfoResponse {
	//	ctx := context.(*db.PollyContext)
	resp := budgetInfoResponse{
		ID:    month,
		Value: budget,
	}

	return resp
}
