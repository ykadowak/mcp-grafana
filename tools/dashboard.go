package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"github.com/grafana/grafana-openapi-client-go/models"
	mcpgrafana "github.com/grafana/mcp-grafana"
)

type GetDashboardByUIDParams struct {
	UID string `json:"uid" jsonschema:"required,description=The UID of the dashboard"`
}

func getDashboardByUID(ctx context.Context, args GetDashboardByUIDParams) (*models.DashboardFullWithMeta, error) {
	c := mcpgrafana.GrafanaClientFromContext(ctx)

	dashboard, err := c.Dashboards.GetDashboardByUID(args.UID)
	if err != nil {
		return nil, fmt.Errorf("get dashboard by uid %s: %w", args.UID, err)
	}
	return dashboard.Payload, nil
}

var GetDashboardByUID = mcpgrafana.MustTool(
	"get_dashboard_by_uid",
	"Get dashboard by uid",
	getDashboardByUID,
)

// TODO: Implement delete/restore

type PostDashboardParams struct {
	Dashboard models.JSON `json:"dashboard" jsonschema:"required,description=The JSON object of the Grafana dashboard POST request."`
	FolderUID string      `json:"folderUid" jsonschema:"required,description=The UID of the folder"`
	Overwrite bool        `json:"overwrite" jsonschema:"description=Whether to overwrite the existing dashboard. If true, the existing dashboard will be replaced with the new one. If false, a new dashboard will be created."`
	IsFolder  bool        `json:"isFolder" jsonschema:"description=Whether the dashboard is a folder. If true, the dashboard will be created as a folder. If false, the dashboard will be created as a regular dashboard."`
}

func postDashboard(ctx context.Context, args PostDashboardParams) (*models.PostDashboardOKBody, error) {
	c := mcpgrafana.GrafanaClientFromContext(ctx)

	response, err := c.Dashboards.PostDashboard(&models.SaveDashboardCommand{
		Dashboard: args.Dashboard,
		FolderUID: args.FolderUID,
		Overwrite: args.Overwrite,
	})
	if err != nil {
		return nil, fmt.Errorf("post dashboard failed: JSON: %v: %w", args.Dashboard, err)
	}
	return response.Payload, nil
}

var postDashboardDesc = `
Tool to post a Grafana Dashboard.
Always try to update the existing panel if the dashboard is being modified in the context.
Expand the panels as much as possible in the window and usually it's a good practice to have 2 panels in 1 row.
When updating the dashboard, always avoid removing the existing panels unless it's requested.
To do so, get the current dashboard JSON first then append the new panels to it. Then, post the updated dashboard JSON.

Use "rate" for counters
Use "histogram_quantile" with 0.5/0.8/0.9/0.99 for buckets
`

var PostDashboard = mcpgrafana.MustTool(
	"post_dashboard",
	postDashboardDesc,
	postDashboard,
)

func AddDashboardTools(mcp *server.MCPServer) {
	GetDashboardByUID.Register(mcp)
	PostDashboard.Register(mcp)
}
