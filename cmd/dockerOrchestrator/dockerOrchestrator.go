package main

import (
	"crypto/rand"
	"fmt"
	"gopkg.in/yaml.v3"
	"hetlesaether.com/dockerOrchestrator/internal/daemon"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	opts := &slog.HandlerOptions{
		AddSource: true,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Starting application", "APP_ENV", os.Getenv("APP_ENV"))

	if err := setup(); err != nil {
		slog.Error("Setup failed", "error", err)
		os.Exit(1)
	}

	if os.Getenv("APP_ENV") == "dev" {
		if err := setupDev(); err != nil {
			slog.Error("Dev setup failed", "error", err)
			slog.Warn("Using uninitialized dev enviroment")
		}
	}

	daemon.Start()
}

func setup() error {
	envVars := []string{
		"APP_PATH",
		"BLUEPRINT_PATH",
		"DATA_PATH",
		"MEDIA_PATH",
		"PUBLIC_PATH",
	}

	for _, v := range envVars {
		rawPath := os.Getenv(v)
		if rawPath == "" {
			return fmt.Errorf("required environment variable %s is not set", v)
		}

		expandedPath := os.ExpandEnv(rawPath)

		if err := os.MkdirAll(expandedPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s at %q: %w", v, expandedPath, err)
		}
	}

	return nil
}

// Until EOF, dev only initilization values for testing, AI generated for efficency
func pseudoUUID() string {
	b := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func setupDev() error {
	appPath := os.Getenv("APP_PATH")
	if appPath == "" {
		return fmt.Errorf("setupDev failed: APP_PATH environment variable is not set")
	}

	// 1. Create a directory within APP_PATH to hold mock hardware storage files
	mockDeviceDir := filepath.Join(appPath, "dev_mock_disks")
	if err := os.MkdirAll(mockDeviceDir, 0755); err != nil {
		return fmt.Errorf("failed creating mock storage folder: %w", err)
	}

	diskNames := []string{"ssd01", "hdd01", "hdd02"}
	diskUUIDs := make(map[string]string)

	// 2. Provision mock disk image files and generate dynamic development UUIDs
	for _, name := range diskNames {
		diskFile := filepath.Join(mockDeviceDir, fmt.Sprintf("%s.img", name))

		// Write a tiny 1MB empty file if it doesn't exist to simulate block layout
		if _, err := os.Stat(diskFile); os.IsNotExist(err) {
			if err := os.WriteFile(diskFile, make([]byte, 1024*1024), 0644); err != nil {
				return fmt.Errorf("failed creating mock disk file %s: %w", name, err)
			}
		}

		diskUUIDs[name] = pseudoUUID()
		slog.Info("[DEV] Mock device generated: %s -> %s\n", name, diskUUIDs[name])
	}

	// 3. Construct the configuration map matching your required YAML output schema
	configData := map[string]interface{}{
		"disks": map[string]interface{}{
			"ssd01": map[string]interface{}{"uuid": diskUUIDs["ssd01"], "ro": false, "ao": false},
			"hdd01": map[string]interface{}{"uuid": diskUUIDs["hdd01"]},
		},
		"mounts": []map[string]string{
			{"source": fmt.Sprintf("$(%s)/media", "ssd01"), "destination": "/srv/media"},
			{"source": fmt.Sprintf("$(%s)/data", "ssd01"), "destination": "/srv/dockerOrchistrator/data"},
			{"source": fmt.Sprintf("$(%s)/blueprints", "ssd01"), "destination": "/srv/dockerOrchistrator/blueprints"},
			{"source": fmt.Sprintf("$(%s)/movies", "hdd01"), "destination": "/srv/media/movies"},
			{"source": fmt.Sprintf("$(%s)/tvshows", "hdd01"), "destination": "/srv/media/tvshows"},
		},
		"networks": map[string]interface{}{
			"database-nw": map[string]interface{}{
				"driver":      "bridge",
				"internal":    true,
				"external":    false,
				"attachable":  true,
				"ipv4":        true,
				"ipv6":        false,
				"description": "",
			},
			"proxy-nw": map[string]interface{}{
				"external":    true,
				"description": "caddy network w/ TLS & reverse proxy",
			},
		},
	}

	// 4. Marshal data out to YAML syntax bytes
	yamlBytes, err := yaml.Marshal(&configData)
	if err != nil {
		return fmt.Errorf("failed to marshal hardware configuration: %w", err)
	}

	// 5. Save infra.yml directly inside your verified $APP_PATH
	targetYamlPath := filepath.Join(appPath, "infra.yml")
	if err := os.WriteFile(targetYamlPath, yamlBytes, 0644); err != nil {
		return fmt.Errorf("failed to write infra.yml: %w", err)
	}

	slog.Info("[DEV] ✓ Compiled dynamic infra.yml at %s\n", "path", targetYamlPath)

	// ### Test blueprints
	// 1. Get the base blueprint path from environment variables
	blueprintPath := os.Getenv("BLUEPRINT_PATH")
	if blueprintPath == "" {
		return fmt.Errorf("Blueprint Path not setup")
	}

	// 2. Construct the absolute destination path for the compose file
	targetDir := filepath.Join(blueprintPath, "test_project")
	composeFilePath := filepath.Join(targetDir, "compose.yml")

	// 3. Ensure the target directory structure exists
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to make dir test_project: %w", err)
	}

	// 4. Define the exact multi-line YAML content string
	// Using raw string literals (``) preserves indentation and spacing perfectly
	composeContent := `services:
  jellyfin:
    image: lscr.io/linuxserver/jellyfin:latest
    container_name: jellyfin
    restart: unless-stopped
    labels:
      caddy: jellyfin.hetlesaether.com
      caddy.reverse_proxy: "{{upstreams 8096}}"
      homepage.group: "Web Apps"
      homepage.name: "Jellyfin"
      homepage.icon: "jellyfin"
      homepage.href: "https://jellyfin.hetlesaether.com"
      homepage.description: "Open-source media streaming hub"
    volumes:
      - ${DATA_PATH}/jellyfin/config:/config
      - ${DATA_PATH}/media/HDD:/media
`

	// 5. Write the file to disk (Creates or truncates, permissions 0644)
	err = os.WriteFile(composeFilePath, []byte(composeContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write test compose.yml: %w", err)
	}

	fmt.Printf("Successfully generated: %s\n", composeFilePath)
	return nil
}
