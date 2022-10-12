package hakucho

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateFolder(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	folderName := fmt.Sprintf("hakucho-test-%s", time.Now().Format("20060102_150406"))
	_, err = c.CreateFolder(folderName)
	if err != nil {
		t.Errorf("Failed to create folder %s: %s", folderName, err)
	}
}

func TestCreateSubFolder(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	dateStr := time.Now().Format("20060102_150406")
	folderName := fmt.Sprintf("hakucho-test-parent-%s", dateStr)
	parent, err := c.CreateFolder(folderName)
	if err != nil {
		t.Errorf("Failed to create folder %s: %s", folderName, err)
	}

	subFolderName := fmt.Sprintf("hakucho-test-child-%s", dateStr)
	_, err = c.CreateSubFolder(parent.Id, subFolderName)
	if err != nil {
		t.Errorf("Failed to create sub folder %s: %s", subFolderName, err)
	}

}
