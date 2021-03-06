package config

import (
	"context"

	"github.com/rancher/norman/api/builtin"
	"github.com/rancher/norman/pkg/subscribe"
	"github.com/rancher/norman/store/crd"
	"github.com/rancher/norman/store/proxy"
	normantypes "github.com/rancher/norman/types"
	"github.com/rancher/rio/api/named"
	"github.com/rancher/rio/api/space"
	"github.com/rancher/rio/api/stack"
	"github.com/rancher/rio/types"
	"github.com/rancher/rio/types/apis/rio.cattle.io/v1beta1/schema"
	spaceSchema "github.com/rancher/rio/types/apis/space.cattle.io/v1beta1/schema"
	"github.com/rancher/rio/types/client/rio/v1beta1"
	spaceClient "github.com/rancher/rio/types/client/space/v1beta1"
)

func SetupTypes(ctx context.Context, context *types.Context) error {
	factory := crd.NewFactoryFromClientGetter(context.ClientGetter)
	factory.BatchCreateCRDs(ctx, normantypes.DefaultStorageContext, context.Schemas,
		&schema.Version,
		client.ServiceType,
		client.StackType)
	factory.BatchCreateCRDs(ctx, normantypes.DefaultStorageContext, context.Schemas,
		&spaceSchema.Version,
		spaceClient.ListenConfigType)
	factory.BatchWait()

	setupSpaces(ctx, factory.ClientGetter, context)
	setupNodes(ctx, factory.ClientGetter, context)
	setupPods(ctx, factory.ClientGetter, context)
	setupServices(ctx, context)
	setupStacks(ctx, context)

	subscribe.Register(&builtin.Version, context.Schemas)
	subscribe.Register(&schema.Version, context.Schemas)
	subscribe.Register(&spaceSchema.Version, context.Schemas)

	return nil
}

func setupServices(ctx context.Context, rContext *types.Context) {
	ef := &stack.ExportFormatter{}
	s := rContext.Schemas.Schema(&schema.Version, client.ServiceType)
	s.Formatter = ef.FormatService
	s.Store = named.New(s.Store)
}

func setupStacks(ctx context.Context, rContext *types.Context) {
	ef := &stack.ExportFormatter{}
	s := rContext.Schemas.Schema(&schema.Version, client.StackType)
	s.Formatter = ef.Format
	s.Store = named.New(s.Store)
}

func setupNodes(ctx context.Context, clientGetter proxy.ClientGetter, rContext *types.Context) {
	s := rContext.Schemas.Schema(&spaceSchema.Version, spaceClient.NodeType)
	s.Store = proxy.NewProxyStore(ctx,
		clientGetter,
		normantypes.DefaultStorageContext,
		[]string{"/api"},
		"",
		"v1",
		"Node",
		"nodes")
}

func setupPods(ctx context.Context, clientGetter proxy.ClientGetter, rContext *types.Context) {
	s := rContext.Schemas.Schema(&spaceSchema.Version, spaceClient.PodType)
	s.Store = proxy.NewProxyStore(ctx,
		clientGetter,
		normantypes.DefaultStorageContext,
		[]string{"/api"},
		"",
		"v1",
		"Pod",
		"pods")
}

func setupSpaces(ctx context.Context, clientGetter proxy.ClientGetter, rContext *types.Context) {
	s := rContext.Schemas.Schema(&spaceSchema.Version, spaceClient.SpaceType)
	s.Store = proxy.NewProxyStore(ctx,
		clientGetter,
		normantypes.DefaultStorageContext,
		[]string{"/api"},
		"",
		"v1",
		"Namespace",
		"namespaces")
	s.Store = space.New(s.Store)
}
