package hakucho

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/yokoe/hakucho/option"
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

func createTempCSVFile(t *testing.T) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "hakucho-test")
	if err != nil {
		t.Fatalf("Failed to prepare temp file: %s", err)
		return nil
	}
	_, err = file.WriteString("id,name")
	if err != nil {
		t.Fatalf("Failed to write string into test file: %s", err)
		return nil
	}
	_, err = file.WriteString("1,Peter")
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

		permissions, err := c.ListPermissions(uploadedFile.Id)
		if err != nil {
			t.Errorf("Failed to get permissions list: %s", err)
		}
		for _, p := range permissions {
			if len(p.EmailAddress) == 0 {
				t.Errorf("Permission object had no email address")
			}
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
	txtFile, err := c.UploadFileToFolder(createTempTextFile(t, "hello").Name(), "hakucho-test-child.txt", folder.Id)

	if err != nil {
		t.Fatalf("Failed to upload file to folder: %s", err)
	}

	csvFile, err := c.UploadFileToFolder(createTempCSVFile(t).Name(), "hakucho-test-child.csv", folder.Id, option.ContentType("text/csv"))

	if err != nil {
		t.Fatalf("Failed to upload file to folder: %s", err)
	}

	deleteEntry(t, c, txtFile.Id)
	deleteEntry(t, c, csvFile.Id)
	deleteEntry(t, c, folder.Id)

}

func TestGetFile(t *testing.T) {
	c, err := newTestClient(t)
	if err != nil {
		return
	}

	files, err := c.ListFiles([]string{"id"}, 1)
	if err != nil {
		t.Fatalf("Failed to get file list: %s", err)
	}
	if len(files) == 0 {
		t.Skip("no files")
	}

	f, err := c.GetFile(files[0].Id, []string{"id", "createdTime"})
	if err != nil {
		t.Fatalf("Failed to get file: %s", err)
	}
	t.Logf("Root ID: %s %s", f.Id, f.CreatedTime)
}
