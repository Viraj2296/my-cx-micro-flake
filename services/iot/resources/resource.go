package resources

import (
	"embed"
	"io/fs"
)

//go:embed component
var resource embed.FS

func GetComponents() ([]byte, error) {
	// Define the path to the component.json file within the embedded filesystem.
	filePath := "component/component.json"

	// Check if the file exists in the embedded FS
	if _, err := fs.Stat(resource, filePath); err != nil {
		return nil, err
	}

	// Read the file content
	content, err := resource.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}
