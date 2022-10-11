package hakucho

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// https://developers.google.com/drive/api/quickstart/go

// Client provides apis for Google Drive
type Client struct {
	driveService *drive.Service
}

// NewClient creates a new client.
func NewClient(credPath string, tokenPath string) (*Client, error) {
	srv, err := newDriveService(credPath, tokenPath)
	if err != nil {
		return nil, err
	}
	return &Client{
		driveService: srv,
	}, nil
}

func newDriveService(credPath string, tokenPath string) (*drive.Service, error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client secret file: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}
	client, err := getClient(config, tokenPath)
	if err != nil {
		return nil, err
	}

	return drive.NewService(ctx, option.WithHTTPClient(client))
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokenPath string) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		return nil, err
	}
	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
