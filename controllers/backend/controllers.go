package backend

import (
	"context"

	"github.com/rancher/rio/controllers/backend/data"
	"github.com/rancher/rio/controllers/backend/envoy"
	"github.com/rancher/rio/controllers/backend/node"
	"github.com/rancher/rio/controllers/backend/service"
	"github.com/rancher/rio/controllers/backend/stack"
	"github.com/rancher/rio/controllers/backend/stackdeploy"
	"github.com/rancher/rio/types"
)

func Register(ctx context.Context, rContext *types.Context) error {
	if err := data.AddData(rContext); err != nil {
		return err
	}

	stack.Register(ctx, rContext)
	stackdeploy.Register(ctx, rContext)
	service.Register(ctx, rContext)
	node.Register(ctx, rContext)
	envoy.Register(ctx, rContext)
	return nil
}
