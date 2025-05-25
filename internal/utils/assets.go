package utils

import (
	"os"
	"path/filepath"
	"strings"
)

var cssFileName string

// InitAssets inicializa los assets al arrancar el servidor
func InitAssets() error {
	// Buscar el archivo CSS más reciente en el directorio
	files, err := os.ReadDir(filepath.Join("static", "css"))
	if err != nil {
		return err
	}

	// Buscar el archivo que empieza por output_
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "output_") {
			cssFileName = file.Name()
			return nil
		}
	}

	return nil
}

// GetCSSFileName devuelve el nombre completo del archivo CSS
func GetCSSFileName() string {
	return cssFileName
} 