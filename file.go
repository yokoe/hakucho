package hakucho

import "google.golang.org/api/drive/v3"

func (c *Client) DeleteFile(fileID string) error {
	_, err := c.driveService.Files.Update(fileID, &drive.File{Trashed: true}).Do()
	return err
}
