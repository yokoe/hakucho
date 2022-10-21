# hakucho
Google Drive helper methods library for Golang

## How to use

```
c, err := NewClient("/path/to/credentials.json", "/path/to/token.json")
if err != nil {
  log.Fatalf("client init error: %s", err)
}

files, err := c.ListFiles([]string{"id", "name", "createdTime"}, 20)
if err != nil {
  log.Fatalf("list error: %s", err)
}

for _, f := range files {
  log.Printf("%s\n", f.Name)
}

```

### Search files

```
files, err := c.ListFiles([]string{"id", "createdTime", "name"}, 5, option.OnlyFolders(), option.FullTextContains("Guitar"))
if err != nil {
  log.Printf("Failed to get list of files: %s", err)
}

for i, f := range files {
  log.Printf("File %02d: %s(%s)\n", i, f.Name, f.CreatedTime)
}
```

#### Available fields
https://developers.google.com/drive/api/v3/reference/files

#### Available options
* option.OnlyFolders()
* option.ParentFolder("parent folder id")
* option.ParentIn("some folder 1", "some folder 2")
* option.NameContains("memo")
* option.FullTextContains("memo")
* option.MimeType("text/csv")
* option.OrderBy("createdTime")

### Upload File

```
// File will be created in root folder
uploadedFile, err := c.UploadFile("/path/to/file", "sample.txt")

// File will be created under selected folder
uploadedFile, err := c.UploadFileToFolder("/path/to/file", "sample.txt", "parent folder id")
```

### Grant permission to user

```
// User can read the file
err := c.GrantUserReaderPemission(uploadedFile.Id, "1@example.com", false)

// User can read & write the file
err := c.GrantUserWriterPemission(uploadedFile.Id, "1@example.com", false)
```

### Delete File

```
err = c.DeleteFile(uploadedFile.Id)
```

### Create Folder

```
// Folder will be created in root folder
folder, err := c.CreateFolder("folder 1")

// Folder will be created under selected folder
folder, err := c.CreateSubFolder(parentFolder.Id, "folder 1")
```
