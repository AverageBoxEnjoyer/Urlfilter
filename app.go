package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type App struct {
	ctx     context.Context
	Version string
}

func NewApp() *App {
	return &App{
		Version: "1.0.0",
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// Elimina el 'return nil, err' ya que startup no retorna valores
	if err := ensureCombosFolder(); err != nil {
		fmt.Printf("Advertencia: %v\n", err) // Solo loguea el error
	}
}

func ensureCombosFolder() error {
	folderName := "combos"

	// Verificar si la carpeta existe
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		if err := os.Mkdir(folderName, 0755); err != nil {
			return fmt.Errorf("no se pudo crear la carpeta 'combos': %v", err)
		}
		return nil // No devolver error aquí, solo cuando se verifique el contenido
	}

	// Verificar si está vacía solo si es necesario
	return checkFolderContents(folderName)
}

func checkFolderContents(folderName string) error {
	dir, err := os.Open(folderName)
	if err != nil {
		return fmt.Errorf("error al abrir la carpeta: %v", err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(1)
	if err != nil {
		return fmt.Errorf("error al leer la carpeta: %v", err)
	}

	if len(files) == 0 {
		return errors.New("la carpeta 'combos' está vacía - coloca tus archivos .txt allí")
	}
	return nil
}

func (a *App) SearchCredentials(keyword string) (map[string]interface{}, error) {
	start := time.Now()
	result := make(map[string]interface{})

	if err := checkFolderContents("combos"); err != nil {
		return nil, err
	}
	if keyword == "" {
		return nil, fmt.Errorf("el término de búsqueda no puede estar vacío")
	}

	outputFilename := fmt.Sprintf("%s_results.txt", strings.ToLower(keyword))
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return nil, fmt.Errorf("error al crear archivo de resultados: %v", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	files, err := os.ReadDir("combos")
	if err != nil {
		return nil, fmt.Errorf("error al leer la carpeta combos: %v", err)
	}

	totalCount := 0
	var hasStructured, hasSimple bool

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}

		inputPath := filepath.Join("combos", file.Name())
		format, err := detectFileFormat(inputPath)
		if err != nil {
			return nil, err
		}

		var count int
		if format == "structured" {
			hasStructured = true
			count, err = processStructuredFormat(keyword, inputPath, writer)
		} else {
			hasSimple = true
			count, err = processSimpleFormat(keyword, inputPath, writer)
		}

		if err != nil {
			return nil, err
		}

		totalCount += count
	}

	// Determinar formato para el resultado
	switch {
	case hasStructured && hasSimple:
		result["format"] = "mixed"
	case hasStructured:
		result["format"] = "structured"
	case hasSimple:
		result["format"] = "simple"
	default:
		result["format"] = "none"
	}

	result["count"] = totalCount
	result["duration"] = time.Since(start).String()
	result["filename"] = outputFilename

	return result, nil
}

// Resto de funciones permanecen iguales excepto por los parámetros en processStructuredFormat y processSimpleFormat

func detectFileFormat(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	for i := 0; i < 5; i++ {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		if strings.HasPrefix(line, "browser:") {
			return "structured", nil
		}
		if strings.Count(line, ":") >= 2 {
			return "simple", nil
		}
	}

	return "simple", nil
}

func processStructuredFormat(keyword string, inputFilename string, writer *bufio.Writer) (int, error) {
	file, err := os.Open(inputFilename)
	if err != nil {
		return 0, fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	credentialMap := make(map[string]bool)

	var currentURL, currentLogin, currentPassword string
	searchKeyword := strings.ToLower(keyword)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		switch {
		case strings.HasPrefix(line, "url:"):
			if strings.Contains(strings.ToLower(currentURL), searchKeyword) {
				if currentLogin != "" && currentPassword != "" {
					creds := fmt.Sprintf("%s:%s", currentLogin, currentPassword)
					if !credentialMap[creds] {
						writer.WriteString(creds + "\n")
						count++
						credentialMap[creds] = true
					}
				}
			}
			currentURL = strings.TrimSpace(strings.TrimPrefix(line, "url:"))
			currentLogin = ""
			currentPassword = ""

		case strings.HasPrefix(line, "login:"):
			currentLogin = strings.TrimSpace(strings.TrimPrefix(line, "login:"))

		case strings.HasPrefix(line, "password:"):
			currentPassword = strings.TrimSpace(strings.TrimPrefix(line, "password:"))
		}
	}

	if strings.Contains(strings.ToLower(currentURL), searchKeyword) {
		if currentLogin != "" && currentPassword != "" {
			creds := fmt.Sprintf("%s:%s", currentLogin, currentPassword)
			if !credentialMap[creds] {
				writer.WriteString(creds + "\n")
				count++
			}
		}
	}

	return count, scanner.Err()
}

func processSimpleFormat(keyword string, inputFilename string, writer *bufio.Writer) (int, error) {
	file, err := os.Open(inputFilename)
	if err != nil {
		return 0, fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	searchKeyword := strings.ToLower(keyword)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), searchKeyword) {
			creds := extractCredentials(line)
			if creds != "" {
				writer.WriteString(creds + "\n")
				count++
			}
		}
	}

	return count, scanner.Err()
}

func extractCredentials(line string) string {
	parts := strings.Split(line, ":")
	if len(parts) < 3 {
		return ""
	}
	return fmt.Sprintf("%s:%s", parts[len(parts)-2], parts[len(parts)-1])
}

func (a *App) GetVersion() string {
	return a.Version
}
