package hakucho

import (
	"fmt"
	"os"

	"github.com/yokoe/hakucho/option"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

func (c *Client) GetFile(fileID string, fileFields []string) (*drive.File, error) {
	return c.driveService.Files.Get(fileID).Fields(escapedFields(fileFields)).Do()
}

func (c *Client) DeleteFile(fileID string) error {
	_, err := c.driveService.Files.Update(fileID, &drive.File{Trashed: true}).Do()
	return err
}

// UploadFile uploads file to Google Drive
func (c *Client) UploadFile(localFile string, uploadFilename string, options ...option.UploadOption) (*drive.File, error) {
	uploadFile, err := os.Open(localFile)
	if err != nil {
		return nil, err
	}
	f := &drive.File{Name: uploadFilename}

	mediaOptions := []googleapi.MediaOption{}
	for _, option := range options {
		mediaOptions = append(mediaOptions, option.MediaOption())
	}

	res, err := c.driveService.Files.Create(f).Media(uploadFile, mediaOptions...).Do()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) UploadFileToFolder(localFile string, uploadFilename string, folderID string, options ...option.UploadOption) (*drive.File, error) {
	f, err := c.UploadFile(localFile, uploadFilename, options...)
	if err != nil {
		return nil, fmt.Errorf("upload error: %w", err)
	}
	if err = c.AddParentFolder(f.Id, folderID); err != nil {
		return nil, fmt.Errorf("add parents error: %w", err)
	}
	return f, nil
}

func (c *Client) AddParentFolder(fileID string, folderID string) error {
	_, err := c.driveService.Files.Update(fileID, &drive.File{}).AddParents(folderID).Do()
	return err
}

func (c *Client) createPermission(fileID string, email string, role string, sendEmail bool) error {
	_, err := c.driveService.Permissions.Create(fileID, &drive.Permission{EmailAddress: email, Role: role, Type: "user"}).SendNotificationEmail(sendEmail).Do()
	return err
}

func (c *Client) GrantUserReaderPemission(fileID string, email string, sendEmail bool) error {
	return c.createPermission(fileID, email, "reader", sendEmail)
}

func (c *Client) GrantUserWriterPemission(fileID string, email string, sendEmail bool) error {
	return c.createPermission(fileID, email, "writer", sendEmail)
}
