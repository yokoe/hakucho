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
