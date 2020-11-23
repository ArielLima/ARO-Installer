package compute

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	mgmtcompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-03-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// UsageClient is a minimal interface for azure UsageClient
type UsageClient interface {
	UsageClientAddons
}

type usageClient struct {
	mgmtcompute.UsageClient
}

var _ UsageClient = &usageClient{}

// NewUsageClient creates a new UsageClient
func NewUsageClient(environment *azure.Environment, tenantID string, authorizer autorest.Authorizer) UsageClient {
	client := mgmtcompute.NewUsageClientWithBaseURI(environment.ResourceManagerEndpoint, tenantID)
	client.Authorizer = authorizer

	return &usageClient{
		UsageClient: client,
	}
}
