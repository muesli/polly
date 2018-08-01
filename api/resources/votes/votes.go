package votes

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// VoteResource is the resource responsible for /votes
type VoteResource struct {
	smolder.Resource
}

var (
	_ smolder.GetSupported  = &VoteResource{}
	_ smolder.PostSupported = &VoteResource{}
)

// Register this resource with the container to setup all the routes
func (r *VoteResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "VoteResource"
	r.TypeName = "vote"
	r.Endpoint = "votes"
	r.Doc = "Manage votes"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *VoteResource) Reads() interface{} {
	return VotePostStruct{}
}

// Returns returns the model that will be returned
func (r *VoteResource) Returns() interface{} {
	return VoteResponse{}
}

func (r *VoteResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	return nil
}
