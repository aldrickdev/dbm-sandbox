package helpers

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
)

type fileType struct {
	name     string
	isDir    bool
	childern []fileType
}

// CheckDirectory will check to see if there is a directory with the same name
// as argument passed in. Returns an error if a directory is found.
func CheckDirectory(name string) (err error) {
	_, err = os.Stat(name)
	if err == nil {
		return fmt.Errorf("Found file/directory with the same name as %s", name)
	}
	return nil
}

// CreateDirectory creates a directory that will be used for the project.
func CreateDirectory(name string) (err error) {
	err = os.Mkdir(name, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create the desired directory '%s', error: %q", name, err)
	}
	return err
}

// GetFSTree will create a slice of fileType that represents the file
// structure in the embedded filesystem.
func GetFSTree(eFileSystem embed.FS, source string) ([]fileType, error) {
	// Example source: embed/docker/postgres
	//
	// Given the file structure:
	// embed/docker/postgres
	// |- docker-compose.yaml.example
	// |- postgres
	// |  |- postgresql.conf
	// |	|- init.sql
	// |
	// |- conf.d
	// |	|- postgres.d
	// |	|  |- conf.yaml
	//
	// Output:
	// []fileType{
	// 	{
	// 		name:     "docker-compose.yaml.example",
	// 		isDir:    false,
	// 		childern: []fileType{},
	// 	},
	// 	{
	// 		name:  "postgres",
	// 		isDir: true,
	// 		childern: []fileType{
	// 			{
	// 				name:     "postgresql.conf",
	// 				isDir:    false,
	// 				childern: []fileType{},
	// 			},
	// 			{
	// 				name:     "init.sql",
	// 				isDir:    false,
	// 				childern: []fileType{},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		name:  "conf.d",
	// 		isDir: false,
	// 		childern: []fileType{
	// 			{
	// 				name:  "postgres.d",
	// 				isDir: true,
	// 				childern: []fileType{
	// 					{
	// 						name:     "conf.yaml",
	// 						isDir:    false,
	// 						childern: []fileType{},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	//

	baseDirectory, err := fs.ReadDir(eFileSystem, source)
	if err != nil {
		log.Fatalf("Failed to read directory: %q", err)
	}

	var allFiles []fileType

	for _, file := range baseDirectory {
		f := fileType{
			name: file.Name(),
		}

		if file.IsDir() {
			f.isDir = true
			f.childern, err = GetFSTree(eFileSystem, source+"/"+file.Name())
			if err != nil {
				return nil, err
			}
		}
		allFiles = append(allFiles, f)
	}

	return allFiles, nil
}

// CreateProjectTree creates the directory structure in the destination directory 
// that matches the embed FS 
func CreateProjectTree(eFileSystem embed.FS, dbms, destination string, tree []fileType) error {
	for _, element := range tree {
		if element.isDir {
			foundDirName := dbms + "/" + element.name
			newDirName := destination + "/" + element.name

			if err := CreateDirectory(newDirName); err != nil {
				return err
			}

			if element.childern != nil {
				if err := CreateProjectTree(eFileSystem, foundDirName, newDirName, element.childern); err != nil {
					return err
				}
			}
		} else {
			foundFile := dbms + "/" + element.name

			fileContent, err := fs.ReadFile(eFileSystem, foundFile)
			if err != nil {
				return err
			}

			if err := os.WriteFile(destination+"/"+element.name, fileContent, 0644); err != nil {
				return nil
			}
		}
	}

	return nil
}

// CopyDirectoryFS will create a copy of the embedded file system, on the
// users machine for the selected DBMS.
func CopyDirectoryFS(eFileSystem embed.FS, dbms, destination string) error {
	var sourceTree []fileType

	tree, err := GetFSTree(eFileSystem, dbms)
	if err != nil {
		return err
	}

	sourceTree = tree
	if err := CreateProjectTree(eFileSystem, dbms, destination, sourceTree); err != nil {
		return err
	}

	return nil
}

