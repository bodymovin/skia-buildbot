// Package buildbucket provides tools for interacting with the buildbucket API.
package buildbucket

import (
	"context"
	"encoding/json"
	fmt "fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	buildbucketpb "go.chromium.org/luci/buildbucket/proto"
	"go.chromium.org/luci/grpc/prpc"
	"go.skia.org/infra/go/sklog"
	"google.golang.org/genproto/protobuf/field_mask"
)

const (
	BUILD_URL_TMPL = "https://%s/build/%d"
	apiUrl         = "cr-buildbucket.appspot.com"
)

var (
	DEFAULT_SCOPES = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	// getBuildFields is a FieldMask which indicates which fields we want
	// returned from a GetBuild request.
	getBuildFields = &field_mask.FieldMask{
		Paths: []string{
			"id",
			"builder",
			"created_by",
			"create_time",
			"start_time",
			"end_time",
			"status",
			"input.properties",
		},
	}

	// searchBuildsFields is a FieldMask which indicates which fields we
	// want returned from a SearchBuilds request.
	searchBuildsFields = &field_mask.FieldMask{
		Paths: func(buildFields []string) []string {
			rv := make([]string, 0, len(buildFields))
			for _, f := range buildFields {
				rv = append(rv, fmt.Sprintf("builds.*.%s", f))
			}
			return rv
		}(getBuildFields.Paths),
	}
)

const (
	// Possible values for the Build.Status field.
	// See: https://chromium.googlesource.com/infra/luci/luci-go/+/master/common/api/buildbucket/buildbucket/v1/buildbucket-gen.go#317
	STATUS_COMPLETED = "COMPLETED"
	STATUS_SCHEDULED = "SCHEDULED"
	STATUS_STARTED   = "STARTED"

	// Possible values for the Build.Result field.
	// See: https://chromium.googlesource.com/infra/luci/luci-go/+/master/common/api/buildbucket/buildbucket/v1/buildbucket-gen.go#305
	RESULT_CANCELED = "CANCELED"
	RESULT_FAILURE  = "FAILURE"
	RESULT_SUCCESS  = "SUCCESS"
)

// Properties contains extra properties set when a Build is requested, as a
// blob of JSON data. These are set by the CQ or "git cl try" when requesting
// try jobs.
type Properties struct {
	Category       string `json:"category"`
	Gerrit         string `json:"patch_gerrit_url"`
	GerritIssue    int64  `json:"patch_issue"`
	GerritPatchset string `json:"patch_ref"`
	PatchProject   string `json:"patch_project"`
	PatchStorage   string `json:"patch_storage"`
	Reason         string `json:"reason"`
	Revision       string `json:"revision"`
	TryJobRepo     string `json:"try_job_repo"`
}

// Parameters provide extra information about a Build.
type Parameters struct {
	BuilderName string     `json:"builder_name"`
	Properties  Properties `json:"properties"`
}

// Build is a struct containing information about a build in BuildBucket.
type Build struct {
	Bucket     string      `json:"bucket"`
	Completed  time.Time   `json:"completed_ts"`
	CreatedBy  string      `json:"created_by"`
	Created    time.Time   `json:"created_ts"`
	Id         string      `json:"id"`
	Url        string      `json:"url"`
	Parameters *Parameters `json:"parameters"`
	Result     string      `json:"result"`
	Status     string      `json:"status"`
}

// Client is used for interacting with the BuildBucket API.
type Client struct {
	bc   buildbucketpb.BuildsClient
	host string
}

// NewClient returns an authenticated Client instance.
func NewClient(c *http.Client) *Client {
	host := apiUrl
	return &Client{
		bc: buildbucketpb.NewBuildsPRPCClient(&prpc.Client{
			C:    c,
			Host: host,
		}),
		host: host,
	}
}

func (c *Client) convertBuild(b *buildbucketpb.Build) *Build {
	bytes, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		sklog.Fatal(err)
	}
	sklog.Info(string(bytes))
	status := ""
	result := ""
	switch b.Status {
	case buildbucketpb.Status_STATUS_UNSPECIFIED:
		// ???
	case buildbucketpb.Status_SCHEDULED:
		status = STATUS_SCHEDULED
	case buildbucketpb.Status_STARTED:
		status = STATUS_STARTED
	case buildbucketpb.Status_SUCCESS:
		status = STATUS_COMPLETED
		result = RESULT_SUCCESS
	case buildbucketpb.Status_FAILURE:
		status = STATUS_COMPLETED
		result = RESULT_FAILURE
	case buildbucketpb.Status_INFRA_FAILURE:
		status = STATUS_COMPLETED
		result = RESULT_FAILURE
	case buildbucketpb.Status_CANCELED:
		status = STATUS_COMPLETED
		result = RESULT_CANCELED
	}
	return &Build{
		Bucket:    b.Builder.Bucket,
		Completed: time.Unix(b.EndTime.Seconds, int64(b.EndTime.Nanos)).UTC(),
		CreatedBy: b.CreatedBy,
		Created:   time.Unix(b.CreateTime.Seconds, int64(b.CreateTime.Nanos)).UTC(),
		Id:        fmt.Sprintf("%d", b.Id),
		Url:       fmt.Sprintf(BUILD_URL_TMPL, c.host, b.Id),
		Parameters: &Parameters{
			BuilderName: b.Builder.Builder,
			Properties: Properties{
				Category:       b.Input.Properties.Fields["category"].GetStringValue(),
				Gerrit:         b.Input.Properties.Fields["patch_gerrit_url"].GetStringValue(),
				GerritIssue:    int64(b.Input.Properties.Fields["patch_issue"].GetNumberValue()),
				GerritPatchset: fmt.Sprintf("%d", int64(b.Input.Properties.Fields["patch_set"].GetNumberValue())),
				PatchProject:   b.Input.Properties.Fields["patch_project"].GetStringValue(),
				PatchStorage:   b.Input.Properties.Fields["patch_storage"].GetStringValue(),
				Reason:         b.Input.Properties.Fields["reason"].GetStringValue(),
				Revision:       b.Input.Properties.Fields["revision"].GetStringValue(),
				TryJobRepo:     b.Input.Properties.Fields["try_job_repo"].GetStringValue(),
			},
		},
		Result: result,
		Status: status,
	}
}

// GetBuild retrieves the build with the given ID.
func (c *Client) GetBuild(ctx context.Context, buildId string) (*Build, error) {
	id, err := strconv.ParseInt(buildId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse build ID as int64: %s", err)
	}
	b, err := c.bc.GetBuild(ctx, &buildbucketpb.GetBuildRequest{
		Id:     id,
		Fields: getBuildFields,
	})
	if err != nil {
		return nil, err
	}
	return c.convertBuild(b), nil
}

// Search retrieves Builds which match the given criteria.
func (c *Client) Search(ctx context.Context, pred *buildbucketpb.BuildPredicate) ([]*Build, error) {
	rv := []*Build{}
	cursor := ""
	for {
		req := &buildbucketpb.SearchBuildsRequest{
			Fields:    searchBuildsFields,
			PageToken: cursor,
			Predicate: pred,
		}
		resp, err := c.bc.SearchBuilds(ctx, req)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			break
		}
		for _, b := range resp.Builds {
			rv = append(rv, c.convertBuild(b))
		}
		cursor = resp.NextPageToken
		if cursor == "" {
			break
		}
	}
	return rv, nil
}

// getTrybotsForCLPredicate returns a *buildbucketpb.BuildPredicate which
// searches for trybots from the given CL.
func getTrybotsForCLPredicate(issue, patchset int64, gerritUrl string) (*buildbucketpb.BuildPredicate, error) {
	u, err := url.Parse(gerritUrl)
	if err != nil {
		return nil, err
	}
	return &buildbucketpb.BuildPredicate{
		GerritChanges: []*buildbucketpb.GerritChange{
			&buildbucketpb.GerritChange{
				Host:     u.Host,
				Change:   issue,
				Patchset: patchset,
			},
		},
	}, nil
}

// GetTrybotsForCL retrieves trybot results for the given CL.
func (c *Client) GetTrybotsForCL(ctx context.Context, issue, patchset int64, gerritUrl string) ([]*Build, error) {
	pred, err := getTrybotsForCLPredicate(issue, patchset, gerritUrl)
	if err != nil {
		return nil, err
	}
	return c.Search(ctx, pred)
}
