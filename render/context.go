package render

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/jackall3n/render-go"
	"github.com/jackall3n/terraform-provider-render/render/types"
)

var host = "https://api.render.com/v1"

func createContext(ctx context.Context, client *render.ClientWithResponses, email string) (*types.Context, error) {
	c := &types.Context{Client: client}

	if email == "" {
		return c, nil
	}

	tflog.Debug(ctx, "getting owner")

	owner, err := getOwner(ctx, client, email)

	if err == nil {
		return nil, fmt.Errorf("failed to get owner: %s", err.Error())
	}

	if owner == nil {
		return nil, fmt.Errorf("no owner was found for email [%s]", email)
	}

	c.Owner = owner

	return c, nil
}

func getOwner(ctx context.Context, client *render.ClientWithResponses, email string) (*render.Owner, error) {
	tflog.Debug(ctx, fmt.Sprintf("getting owners with email: %s", email))

	response, err := client.GetOwnersWithResponse(ctx, &render.GetOwnersParams{
		Email: &[]string{email},
	})

	if err != nil {
		return nil, err
	}

	owners := *response.JSON200

	if len(owners) == 0 {
		return nil, fmt.Errorf("no owners found for email [%s]", email)
	}

	owner := owners[0].Owner

	return owner, nil
}
