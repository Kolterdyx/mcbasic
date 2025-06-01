package frontend

import (
	"embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
	"text/template"
)

func loadStdlib(f *Frontend, stdlib embed.FS) {
	stdlibFiles, err := findStdlibFiles(stdlib, "stdlib")
	if err != nil {
		log.Fatalf("Failed to read stdlib files: %v", err)
	}
	for _, filepath := range stdlibFiles {
		processStdlibTemplate(f, stdlib, filepath)
	}
}

func findStdlibFiles(stdlib embed.FS, dirName string) ([]string, error) {
	files, err := stdlib.ReadDir(dirName)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirName, err)
	}

	var stdlibFiles []string
	for _, file := range files {
		if !file.IsDir() {
			if strings.HasSuffix(file.Name(), ".mcb") {
				stdlibFiles = append(stdlibFiles, path.Join(dirName, file.Name()))
			}
		} else {
			subFiles, err := findStdlibFiles(stdlib, path.Join(dirName, file.Name()))
			if err != nil {
				return nil, err
			}
			stdlibFiles = append(stdlibFiles, subFiles...)
		}
	}
	return stdlibFiles, nil
}

func processStdlibTemplate(f *Frontend, stdlib embed.FS, filepath string) {
	file, err := stdlib.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read stdlib file %s: %v", filepath, err)
	}
	tmpl, err := template.New("stdlib").Parse(string(file))
	if err != nil {
		log.Fatalf("Failed to parse stdlib file %s: %v", filepath, err)
	}
	var content strings.Builder
	err = tmpl.Execute(&content, f.config.Project)
	if err != nil {
		log.Fatalf("Failed to execute template for stdlib file %s: %v", filepath, err)
	}
	fileModulePath, found := strings.CutPrefix(filepath, "stdlib/")
	if !found {
		log.Fatalf("Failed to find module path for stdlib file %s", filepath)
	}
	fileModulePath, _ = strings.CutSuffix(path.Join("@mcb", fileModulePath), ".mcb")
	f.stdlib[fileModulePath] = content.String()
	log.Debugf("Added stdlib %s", fileModulePath)
}
