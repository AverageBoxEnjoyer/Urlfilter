# URL Filter

## Description
URL Filter is a desktop application designed to filter and process URLs, logs, and credentials. The application supports multiple input formats and is built as a personal proof of concept using the Svelte framework.

## Features
- Intuitive and user-friendly interface
- Support for multiple input formats
- URL, log, and credential processing
- Modern and responsive design

## Prerequisites
- [Go](https://golang.org/dl/) (version 1.23)
- [Node.js](https://nodejs.org/) (version 20 or higher)
- [Wails](https://wails.io/) (desktop application development framework)

## Installation

### 1. Install Wails
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. Clone the Repository
```bash
git clone https://github.com/AverageBoxEnjoyer/Urlfilter
cd Urlfilter
```

### 3. Install Frontend Dependencies
```bash
cd frontend
npm install
```

### 4. Install Go Dependencies
```bash
cd ..
go mod download
```

## Development

### Run in Development Mode
```bash
wails dev
```
This command will start the application in development mode with hot reloading.

### Build Executable
```bash
wails build
```
This command will generate an executable in the `build` folder.

## Project Structure
```
Urlfilter/
├── frontend/          # Svelte source code
├── build/            # Build files
├── main.go           # Main entry point
├── app.go            # Application logic
├── go.mod            # Go dependencies
└── wails.json        # Wails configuration
```

## Contributing
If you find any issues or have suggestions to improve the project, please create an issue in the repository.

## License
This project is open source and available for use and modification. Users are free to modify and edit the code according to their needs.

## Disclaimer
The author is not responsible for any use or misuse of this software. Each user is responsible for their own use of the application.
