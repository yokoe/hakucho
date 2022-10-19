package hakucho

import (
	"testing"

	"github.com/yokoe/hakucho/option"
)

func Test_queryFromListOptions(t *testing.T) {
	type args struct {
		options []option.ListOption
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"No options", args{}, ""},
		{"Parent folder", args{[]option.ListOption{option.ParentFolder("folder")}}, "'folder' in parents"},
		{"Parent folders", args{[]option.ListOption{option.ParentIn("a", "b", "c")}}, "('a' in parents or 'b' in parents or 'c' in parents)"},
		{"Name contains keyword", args{[]option.ListOption{option.NameContains("keyword")}}, "name contains 'keyword'"},
		{"FullText contains keyword", args{[]option.ListOption{option.FullTextContains("keyword")}}, "fullText contains 'keyword'"},
		{"Parent folder and name contains keyword", args{[]option.ListOption{option.ParentFolder("folder"), option.NameContains("keyword")}}, "'folder' in parents and name contains 'keyword'"},
		{"Parent folder order by createdTime", args{[]option.ListOption{option.ParentFolder("folder"), option.OrderBy("createdTime")}}, "'folder' in parents"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := queryFromListOptions(tt.args.options); got != tt.want {
				t.Errorf("queryFromListOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ListFiles(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	files, err := c.ListFiles([]string{"id", "createdTime", "name"}, 5, option.OnlyFolders(), option.FullTextContains("hakucho-test-"))
	if err != nil {
		t.Fatalf("Failed to get list of files: %s", err)
	}

	for i, f := range files {
		t.Logf("File %02d: %s", i, f.Name)
	}

	// t.Fail() // for debugging

}
