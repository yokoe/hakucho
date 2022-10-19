package hakucho

import (
	"strings"

	"github.com/yokoe/hakucho/option"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const listMaxPageSize = 20

func queryFromListOptions(options []option.ListOption) string {
	queries := []string{}
	for _, o := range options {
		q := o.QueryString()
		if len(q) > 0 {
			queries = append(queries, q)
		}
	}
	return strings.Join(queries, " and ")
}

func orderFromListOptions(options []option.ListOption) string {
	for _, o := range options {
		if len(o.OrderString()) > 0 {
			return o.OrderString()
		}
	}
	return ""
}

func (c *Client) ListFiles(fileFields []string, limit int64, options ...option.ListOption) ([]*drive.File, error) {
	pageSize := limit
	if pageSize >= listMaxPageSize {
		pageSize = listMaxPageSize
	}

	call := c.driveService.Files.List().Fields(escapedFileFields(fileFields), googleapi.Field("nextPageToken")).PageSize(pageSize)

	query := queryFromListOptions(options)
	if len(query) > 0 {
		call = call.Q(query)
	}

	order := orderFromListOptions(options)
	if len(order) > 0 {
		call = call.OrderBy(order)
	}

	files := []*drive.File{}
	nextPageToken := ""
	for {
		l, err := call.PageToken(nextPageToken).Do()
		if err != nil {
			return nil, err
		}
		files = append(files, l.Files...)
		nextPageToken = l.NextPageToken

		if len(files) > int(limit) || len(nextPageToken) == 0 {
			break
		}
	}

	if len(files) > int(limit) {
		return files[0:limit], nil
	}
	return files, nil
}
