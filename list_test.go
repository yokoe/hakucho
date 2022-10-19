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
		{"Name contains keyword", args{[]option.ListOption{option.NameContains("keyword")}}, "name contains 'keyword'"},
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

	files, err := c.ListFiles([]string{"id", "createdTime", "name"}, 20, option.OnlyFolders(), option.NameContains("hakucho-test-"))
	if err != nil {
		t.Fatalf("Failed to get list of files: %s", err)
	}

	for _, f := range files {
		t.Logf("File %s", f.Name)
	}

	// t.Fail() // for debugging

}
