package hakucho

import (
	"strings"

	"github.com/yokoe/hakucho/option"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

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

func (c *Client) ListFiles(fields []string, limit int64, options ...option.ListOption) ([]*drive.File, error) {
	convertedFields := []googleapi.Field{}
	for _, f := range fields {
		convertedFields = append(convertedFields, googleapi.Field(f))
	}

	call := c.driveService.Files.List().Fields(convertedFields...)

	query := queryFromListOptions(options)
	if len(query) > 0 {
		call = call.Q(query)
	}

	order := orderFromListOptions(options)
	if len(order) > 0 {
		call = call.OrderBy(order)
	}

	l, err := call.PageSize(limit).Do()
	if err != nil {
		return nil, err
	}
	return l.Files, nil
}
