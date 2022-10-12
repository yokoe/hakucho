package hakucho

import "google.golang.org/api/drive/v3"

// CreateFolder creates new folder
func (c *Client) CreateFolder(folderName string) (*drive.File, error) {
	return c.driveService.Files.Create(&drive.File{Name: folderName, MimeType: "application/vnd.google-apps.folder"}).Do()
}

// CreateSubFolder creates new folder and set parent folder id of it
func (c *Client) CreateSubFolder(parentFolderID string, folderName string) (*drive.File, error) {
	folder, err := c.CreateFolder(folderName)
	if err != nil {
		return nil, err
	}

	if _, err := c.driveService.Files.Update(folder.Id, &drive.File{}).AddParents(parentFolderID).Do(); err != nil {
		return nil, err
	}
	return folder, nil
}
