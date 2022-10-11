package hakucho

import "google.golang.org/api/drive/v3"

// CreateFolder creates new folder
func (srv *Client) CreateFolder(folderName string) (*drive.File, error) {
	return srv.driveService.Files.Create(&drive.File{Name: folderName, MimeType: "application/vnd.google-apps.folder"}).Do()
}
