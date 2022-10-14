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