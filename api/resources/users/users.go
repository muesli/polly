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

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *UserResource) Reads() interface{} {
	return UserPostStruct{}
}

// Returns returns the model that will be returned
func (r *UserResource) Returns() interface{} {
	return UserResponse{}
}

func (r *UserResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	return nil
}
