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
	_ smolder.GetSupported = &VoteResource{}
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
