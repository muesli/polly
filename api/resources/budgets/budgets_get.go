package budgets

import (
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/muesli/polly/api/db"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *BudgetResource) GetAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *BudgetResource) GetDoc() string {
	return "retrieve budgets"
}

// GetParams returns the parameters supported by this API endpoint
func (r *BudgetResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("month", "budget for this month number").DataType("int"))

	return params
}

// Get sends out items matching the query parameters
func (r *BudgetResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	/*	authUser := db.User{}
		if auth, err := context.Authentication(request); err == nil {
			authUser = auth.(db.User)
		}*/

	ctx := context.(*db.PollyContext)
	resp := BudgetResponse{}
	resp.Init(context)

	month := int(time.Now().Month())
	m := params["month"]
	if len(m) > 0 {
		var err error
		month, err = strconv.Atoi(m[0])
		if err != nil {
			r.NotFound(request, response)
			return
		}
	}

	resp.Budgets = append(resp.Budgets, prepareBudgetResponse(context, uint(month), ctx.SmallGrantMaxValue(uint(month))))

	resp.Send(response)
}
