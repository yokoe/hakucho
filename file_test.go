package hakucho

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"google.golang.org/api/drive/v3"
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

func uploadTestFile(t *testing.T, c *Client, filenamePrefix string) (string, *drive.File) {
	filename := fmt.Sprintf("%s-%s.txt", filenamePrefix, time.Now().Format("20060102_150406"))

	testFile, err := createTempTextFile("hello world")
	if err != nil {
		t.Fatalf("Test file create error: %s", err)
		return "", nil
	}

	uploadedFile, err := c.UploadFile(testFile.Name(), filename)
	if err != nil {
		t.Fatalf("Failed to upload file %s: %s", filename, err)
		return "", nil
	}
	return filename, uploadedFile
}

func TestGrantPermissions(t *testing.T) {
	cfg := loadTestConfig(t)
	if len(cfg.GrantUserEmail) == 0 {
		t.Skip("no grant user email set")
	}

	c, err := newTestClient(t)
	if err != nil {
		return
	}

	t.Run("reader", func(t *testing.T) {
		_, uploadedFile := uploadTestFile(t, c, "hakucho-test-permission-reader")

		if err = c.GrantUserReaderPemission(uploadedFile.Id, cfg.GrantUserEmail, false); err != nil {
			t.Errorf("Failed to grant reader permission to %s: %s", cfg.GrantUserEmail, err)
		}
	})

	t.Run("writer", func(t *testing.T) {
		_, uploadedFile := uploadTestFile(t, c, "hakucho-test-permission-writer")

		if err = c.GrantUserWriterPemission(uploadedFile.Id, cfg.GrantUserEmail, false); err != nil {
			t.Errorf("Failed to grant writer permission to %s: %s", cfg.GrantUserEmail, err)
		}
	})

}

func TestAddParent(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	folderName := fmt.Sprintf("hakucho-test-%s", time.Now().Format("20060102_150406"))
	folder, err := c.CreateFolder(folderName)
	if err != nil {
		t.Errorf("Failed to create folder %s: %s", folderName, err)
	}

	_, uploadedFile := uploadTestFile(t, c, "hakucho-test-child-file")
	if err = c.AddParentFolder(uploadedFile.Id, folder.Id); err != nil {
		t.Errorf("Failed to add parent: %s", err)
	}

}
