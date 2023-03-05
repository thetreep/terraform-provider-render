package render

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/types"
)

var host = "https://api.render.com/v1"

func createContext(ctx context.Context, client *render.ClientWithResponses, email string) *types.Context {
	c := &types.Context{Client: client}

	if email == "" {
		return c
	}

	tflog.Debug(ctx, "getting owner")

	owner := getOwner(ctx, client, email)

	if owner == nil {
		return nil
	}

	c.Owner = owner

	return c
}

func getOwner(ctx context.Context, client *render.ClientWithResponses, email string) *render.Owner {
	tflog.Debug(ctx, fmt.Sprintf("getting owners with email: %s", email))

	response, err := client.GetOwnersWithResponse(ctx, &render.GetOwnersParams{
		Email: &[]string{email},
	})

	if err != nil {
		return nil
	}

	owner := (*response.JSON200)[0].Owner

	return owner
}
