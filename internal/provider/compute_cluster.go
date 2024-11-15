package provider

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &computeClusterResource{}
	_ resource.ResourceWithConfigure = &computeClusterResource{}
)

// ComputeClusterResource implements the resource.Resource interface.
type computeClusterResource struct {
	db *sql.DB
}
type ComputeClusterResourceModel struct {
	compute_profile_name   types.String `tfsdk:"compute_profile_name"`
	compute_group_nameName types.String `tfsdk:"compute_group_name"`
	query_strategy         types.String `tfsdk:"query_strategy"`
	compute_map            types.String `tfsdk:"compute_map"`
	compute_attribute      types.String `tfsdk:"compute_attribute"`
	timeout                types.Int64  `tfsdk:"timeout"`
}

func ManageComputeClusterResource() resource.Resource {
	return &computeClusterResource{}
}

// Metadata returns the resource type name.
func (r *computeClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_computecluster"
}

// Schema defines the schema for the resource.
func (r *computeClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"compute_profile_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Compute Profile to manage.",
			},
			"compute_group_name": schema.StringAttribute{
				Optional:    true,
				Description: "The name of compute group to which compute profile belongs.",
			},
			"query_strategy": schema.StringAttribute{
				Optional:    true,
				Description: "Query Strategy to use.",
			},
			"compute_map": schema.StringAttribute{
				Optional:    true,
				Description: "The ComputeMapName of the compute map.",
			},
			"compute_attribute": schema.StringAttribute{
				Optional:    true,
				Description: "Optional attributes of compute profile.",
			},
			"timeout": schema.StringAttribute{
				Optional:    true,
				Description: "Time elapsed before the task times out and fails.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *computeClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	db, ok := req.ProviderData.(*sql.DB)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sql.DB, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.db = db
}

// Create creates the resource and sets the initial Terraform state.
func (r *computeClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Info(ctx, "Creating New Vatnage Compute Cluster")

	var plan ComputeClusterResourceModel

	tflog.Info(ctx, "Before Mapping")
	diags := req.Plan.Get(ctx, &plan)
	tflog.Info(ctx, "After Mapping")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "After Error", map[string]interface{}{"ProfileName": plan.compute_profile_name.ValueString(),
		"region": plan.compute_group_nameName.ValueString(), "state": plan.query_strategy.ValueString(), "ip": plan.compute_map.ValueString(),
		"dnsname": plan.compute_attribute.ValueString(),
		"owner":   plan.timeout.ValueInt64()})

	var create_cp_query = "CREATE COMPUTE PROFILE " + plan.compute_profile_name.ValueString()
	if len(strings.TrimSpace(plan.compute_group_nameName.ValueString())) == 0 {
		create_cp_query = create_cp_query + " IN " + plan.compute_group_nameName.ValueString()
	}
	if len(strings.TrimSpace(plan.compute_map.ValueString())) == 0 {
		create_cp_query = create_cp_query + ", INSTANCE = " + plan.compute_map.ValueString()
	}
	if len(strings.TrimSpace(plan.query_strategy.ValueString())) == 0 {
		create_cp_query = create_cp_query + ", INSTANCE TYPE = " + plan.query_strategy.ValueString()
	}
	if len(strings.TrimSpace(plan.compute_attribute.ValueString())) == 0 {
		create_cp_query = create_cp_query + " USING " + plan.compute_attribute.ValueString()
	}
	_, err := r.db.Exec(create_cp_query)

	tflog.Info(ctx, "Compute Cluster with profile name %s Created.", map[string]interface{}{"ProfileName": plan.compute_profile_name.ValueString()})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Compute Cluster", err.Error())
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *computeClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Read Vatnage Compute Cluster")
	var state ComputeClusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *computeClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Update Vatnage Compute Cluster")
	var plan ComputeClusterResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *computeClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "Delete Vatnage Compute Cluster")
	var state ComputeClusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
