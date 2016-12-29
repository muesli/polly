package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// SessionResource is the resource responsible for /sessions
type SessionResource struct {
	smolder.Resource
}

// SessionResponse is the common response to 'session' requests
type SessionResponse struct {
	smolder.Response

	IDToken string `json:"id_token"`
	UserID  int64  `json:"user_id"`
}

// SessionPostStruct holds all values of an incoming POST request
type SessionPostStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// Init a new response
func (r *SessionResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context
}

// Register this resource with the container to setup all the routes
func (r *SessionResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "SessionResource"
	r.TypeName = "session"
	r.Endpoint = "sessions"
	r.Doc = "Manage sessions"

	r.Config = config
	r.Context = context

	log.WithField("Resource", r.Name).Info("Registering Resource")

	ws := new(restful.WebService)
	ws.Path("/" + r.Config.PathPrefix + r.Endpoint).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	route := ws.POST("/create").To(r.Post)

	route.Param(restful.QueryParameter("username", "username").
		DataType("string").
		Required(true).
		AllowMultiple(false))
	route.Param(restful.QueryParameter("password", "password").
		DataType("string").
		Required(true).
		AllowMultiple(false))
	route.Param(restful.QueryParameter("token", "token").
		DataType("string").
		Required(true).
		AllowMultiple(false))

	ws.Route(route)

	container.Add(ws)
}

// PostDoc returns the description of this API endpoint
func (r *SessionResource) PostDoc() string {
	return "create a new user invitation"
}

// PostParams returns the parameters supported by this API endpoint
func (r *SessionResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *SessionResource) Post(request *restful.Request, response *restful.Response) {
	context := r.Context.NewAPIContext()
	resp := SessionResponse{}
	resp.Init(context)

	sps := SessionPostStruct{}
	err := request.ReadEntity(&sps)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"SessionResource PUT"))
		return
	}

	user := DbUser{}
	if len(sps.Token) > 0 {
		auth, aerr := context.(*PollyContext).GetUserByAccessToken(sps.Token)
		if aerr != nil {
			r.NotFound(request, response)
			return
		}
		user = auth.(DbUser)

		if len(sps.Password) > 0 {
			user.UpdatePassword(context.(*PollyContext), sps.Password)
		}
	} else {
		user, err = context.(*PollyContext).GetUserByNameAndPassword(sps.Username, sps.Password)
		if err != nil {
			r.NotFound(request, response)
			return
		}
	}

	uuid, err := UUID()
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't create authtoken",
			"SessionResource PUT"))
		return
	}

	user.AuthToken = uuid
	err = user.Update(context.(*PollyContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't update user session",
			"SessionResource POST"))
		return
	}

	resp.IDToken = user.AuthToken
	resp.UserID = user.ID
	response.WriteHeaderAndEntity(http.StatusOK, resp)
}
