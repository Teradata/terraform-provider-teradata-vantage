// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure TeradataVantageProvider satisfies various provider interfaces.
var _ provider.Provider = &TeradataVantageProvider{}
var _ provider.ProviderWithFunctions = &TeradataVantageProvider{}

// TeradataVantageProvider defines the provider implementation.
type TeradataVantageProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	db      *sql.DB
}

// TeradataVantageProviderModel describes the provider data model.
type TeradataVantageProviderModel struct {
	DBUser     types.String `tfsdk:"db_user"`
	DBPassword types.String `tfsdk:"db_password"`
	DBName     types.String `tfsdk:"db_name"`
	DBHost     types.String `tfsdk:"db_host"`
}

func (p *TeradataVantageProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "teradata-vantage"
	resp.Version = p.version
}

func (p *TeradataVantageProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"db_host": schema.StringAttribute{
				Required: true,
			},
			"db_user": schema.StringAttribute{
				Required: true,
			},
			"db_password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"db_name": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *TeradataVantageProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TeradataVantageProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		data.DBHost.ValueString(),
		data.DBUser.ValueString(),
		data.DBPassword.ValueString(),
		data.DBName.ValueString(),
	)

	db, err := sql.Open("teradata", connStr)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create database connection",
			err.Error(),
		)
		return
	}

	p.db = db
	resp.DataSourceData = p.db
	resp.ResourceData = p.db
}

func (p *TeradataVantageProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		ManageComputeClusterResource,
	}
}

func (p *TeradataVantageProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *TeradataVantageProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TeradataVantageProvider{
			version: version,
		}
	}
}
