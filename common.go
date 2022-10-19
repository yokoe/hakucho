package hakucho

import (
	"fmt"
	"strings"

	"google.golang.org/api/googleapi"
)

func escapedFields(fields []string) googleapi.Field {
	escapedFields := []string{}
	for _, f := range fields {
		escapedFields = append(escapedFields, strings.ReplaceAll(f, ")", ""))
	}

	return googleapi.Field(strings.Join(escapedFields, ", "))
}

func escapedFileFields(fields []string) googleapi.Field {
	escapedFields := []string{}
	for _, f := range fields {
		escapedFields = append(escapedFields, strings.ReplaceAll(f, ")", ""))
	}

	return googleapi.Field(fmt.Sprintf("files(%s)", strings.Join(escapedFields, ", ")))
}
