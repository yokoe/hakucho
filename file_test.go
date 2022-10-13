package hakucho

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"google.golang.org/api/drive/v3"
)

func createTempTextFile(t *testing.T, text string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "hakucho-test")
	if err != nil {
		t.Fatalf("Failed to prepare temp file: %s", err)
		return nil
	}
	_, err = file.WriteString(text)
	if err != nil {
		t.Fatalf("Failed to write string into test file: %s", err)
		return nil
	}
	return file
}

func deleteEntry(t *testing.T, c *Client, id string) {
	t.Helper()
	if err := c.DeleteFile(id); err != nil {
		t.Errorf("Failed to delete file %s: %s", id, err)
	}
}

func createFolder(t *testing.T, c *Client, prefix string) *drive.File {
	folderName := fmt.Sprintf("%s-%s", prefix, time.Now().Format("20060102_150406"))
	folder, err := c.CreateFolder(folderName)
	if err != nil {
		t.Fatalf("Failed to create folder %s: %s", folderName, err)
	}
	return folder
}

func uploadTestFile(t *testing.T, c *Client, filenamePrefix string) (string, *drive.File) {
	filename := fmt.Sprintf("%s-%s.txt", filenamePrefix, time.Now().Format("20060102_150406"))

	testFile := createTempTextFile(t, "hello world")

	uploadedFile, err := c.UploadFile(testFile.Name(), filename)
	if err != nil {
		t.Fatalf("Failed to upload file %s: %s", filename, err)
		return "", nil
	}
	return filename, uploadedFile
}

func TestUploadFileAndDelete(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	testFile := createTempTextFile(t, "hello world")

	uploadName := fmt.Sprintf("hakucho-test-%s.txt", time.Now().Format("20060102_150406"))
	uploadedFile, err := c.UploadFile(testFile.Name(), uploadName)
	if err != nil {
		t.Fatalf("Failed to upload file %s: %s", uploadName, err)
	}

	if err = c.DeleteFile(uploadedFile.Id); err != nil {
		t.Fatalf("Failed to delete file %s: %s", uploadName, err)
	}

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

	folder := createFolder(t, c, "hakucho-test-add-parent")

	_, uploadedFile := uploadTestFile(t, c, "hakucho-test-child-file")
	if err = c.AddParentFolder(uploadedFile.Id, folder.Id); err != nil {
		t.Errorf("Failed to add parent: %s", err)
	}

	deleteEntry(t, c, uploadedFile.Id)
	deleteEntry(t, c, folder.Id)
}

func TestUploadFileToFolder(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	folder := createFolder(t, c, "hakucho-test-upload-to-folder")
	file, err := c.UploadFileToFolder(createTempTextFile(t, "hello").Name(), "hakucho-test-child.txt", folder.Id)

	if err != nil {
		t.Fatalf("Failed to upload file to folder: %s", err)
	}

	deleteEntry(t, c, file.Id)
	deleteEntry(t, c, folder.Id)

}
