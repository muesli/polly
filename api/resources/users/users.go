package users

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// UserResource is the resource responsible for /users
type UserResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &UserResource{}
	_ smolder.GetSupported   = &UserResource{}
	_ smolder.PostSupported  = &UserResource{}
)

// Register this resource with the container to setup all the routes
func (r *UserResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "UserResource"
	r.TypeName = "user"
	r.Endpoint = "users"
	r.Doc = "Manage users"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}
