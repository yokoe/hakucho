package option

import "google.golang.org/api/googleapi"

type UploadOption interface {
	MediaOption() googleapi.MediaOption
}

type uploadContentType struct {
	contentType string
}

func ContentType(contentType string) UploadOption {
	return uploadContentType{contentType: contentType}
}

func (o uploadContentType) MediaOption() googleapi.MediaOption {
	return googleapi.ContentType(o.contentType)
}
