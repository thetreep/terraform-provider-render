package types

import "github.com/jackall3n/render-go"

type Context struct {
	Client *render.ClientWithResponses
	Owner  *render.Owner
}
