package hakucho

import (
	"os"

	"google.golang.org/api/drive/v3"
)

func (c *Client) DeleteFile(fileID string) error {
	_, err := c.driveService.Files.Update(fileID, &drive.File{Trashed: true}).Do()
	return err
}

// UploadFile uploads file to Google Drive
func (c *Client) UploadFile(localFile string, uploadFilename string) (*drive.File, error) {
	uploadFile, err := os.Open(localFile)
	if err != nil {
		return nil, err
	}
	f := &drive.File{Name: uploadFilename}
	res, err := c.driveService.Files.Create(f).Media(uploadFile).Do()
	if err != nil {
		return nil, err
	}
	return res, nil
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
