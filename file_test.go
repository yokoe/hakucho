package hakucho

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func createTempTextFile(text string) (*os.File, error) {
	file, err := ioutil.TempFile(os.TempDir(), "hakucho-test")
	if err != nil {
		return nil, err
	}
	_, err = file.WriteString(text)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func TestUploadFileAndDelete(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	testFile, err := createTempTextFile("hello world")
	if err != nil {
		t.Fatalf("Test file create error: %s", err)
	}

	uploadName := fmt.Sprintf("hakucho-test-%s.txt", time.Now().Format("20060102_150406"))
	uploadedFile, err := c.UploadFile(testFile.Name(), uploadName)
	if err != nil {
		t.Fatalf("Failed to upload file %s: %s", uploadName, err)
	}

	if err = c.DeleteFile(uploadedFile.Id); err != nil {
		t.Fatalf("Failed to delete file %s: %s", uploadName, err)
	}

}
