package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "MiAppWails",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func (a *App) FetchApiMessage() (string, error) {
	// Crear un cliente HTTP con timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Realizar la solicitud GET
	resp, err := client.Get("https://NotAnAPI.net")
	if err != nil {
		return "", fmt.Errorf("error al conectar: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("código de estado inesperado: %d", resp.StatusCode)
	}

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al leer respuesta: %v", err)
	}

	return string(body), nil
}

// CheckCombosFolder verifica si la carpeta combos existe y tiene archivos
func (a *App) CheckCombosFolder() (string, error) {
	folderName := "combos"

	// Verificar si la carpeta existe
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		return "", fmt.Errorf("la carpeta 'combos' no existe. Por favor créala y añade archivos .txt")
	}

	// Verificar si está vacía
	dir, err := os.Open(folderName)
	if err != nil {
		return "", fmt.Errorf("error al verificar la carpeta: %v", err)
	}
	defer dir.Close()

	files, err := dir.Readdirnames(1) // Solo necesitamos saber si hay al menos 1 archivo
	if err != nil {
		return "", fmt.Errorf("error al leer la carpeta: %v", err)
	}

	if len(files) == 0 {
		return "", errors.New("la carpeta 'combos' está vacía. Añade archivos .txt para buscar")
	}

	return "Carpeta combos verificada correctamente", nil
}
