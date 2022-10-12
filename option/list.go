package option

import (
	"fmt"
	"strings"
)

type ListOption interface {
	QueryString() string
	OrderString() string
}

type parentFolder struct {
	folderID string
}

func ParentFolder(folderID string) ListOption {
	return parentFolder{folderID: folderID}
}

func (o parentFolder) QueryString() string {
	return fmt.Sprintf("'%s' in parents", strings.ReplaceAll(o.folderID, "'", ""))
}
func (o parentFolder) OrderString() string { return "" }

type nameContains struct {
	keyword string
}

func NameContains(keyword string) ListOption {
	return nameContains{keyword: keyword}
}

func (o nameContains) QueryString() string {
	return fmt.Sprintf("name contains '%s'", strings.ReplaceAll(o.keyword, "'", ""))
}
func (o nameContains) OrderString() string { return "" }

type orderBy struct {
	order string
}

func OrderBy(order string) ListOption {
	return orderBy{order: order}
}

func (o orderBy) QueryString() string {
	return ""
}

func (o orderBy) OrderString() string { return o.order }

type mimeTypeFilter struct {
	mimeType string
}

func OnlyFolders() ListOption {
	return mimeTypeFilter{mimeType: "application/vnd.google-apps.folder"}
}

func (o mimeTypeFilter) QueryString() string {
	return fmt.Sprintf("mimeType = '%s'", strings.ReplaceAll(o.mimeType, "'", ""))
}

func (o mimeTypeFilter) OrderString() string {
	return ""
}
