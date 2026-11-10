package doctor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var makeTarget = regexp.MustCompile("^([A-Za-z0-9_.-]+):")

func inspectProjectFiles(root string) Project {
	p := Project{NPMScripts: map[string]bool{}, Requirements: map[string]bool{}}
	exists := func(name string) bool { _, err := os.Stat(filepath.Join(root, name)); return err == nil }
	if exists("pnpm-lock.yaml") {
		p.PackageManager = "pnpm"
	} else if exists("yarn.lock") {
		p.PackageManager = "yarn"
	} else if exists("package-lock.json") {
		p.PackageManager = "npm"
	}
	p.HasNodeModules = exists("node_modules")
	if data, err := os.ReadFile(filepath.Join(root, "package.json")); err == nil {
		p.HasPackageJSON = true
		var manifest struct {
			Scripts map[string]json.RawMessage `json:"scripts"`
		}
		if json.Unmarshal(data, &manifest) == nil {
			for name := range manifest.Scripts {
				p.NPMScripts[name] = true
			}
		}
	}
	for _, name := range []string{"requirements.txt", "requirements-dev.txt"} {
		if data, err := os.ReadFile(filepath.Join(root, name)); err == nil {
			p.HasRequirements = true
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.FieldsFunc(line, func(r rune) bool { return r == '=' || r == '<' || r == '>' || r == '[' || r == ';' })
				if len(parts) > 0 {
					p.Requirements[strings.ToLower(strings.ReplaceAll(parts[0], "-", "_"))] = true
				}
			}
		}
	}
	p.HasVenv = exists(".venv") || exists("venv")
	for _, name := range []string{"compose.yaml", "compose.yml", "docker-compose.yml", "docker-compose.yaml"} {
		if data, err := os.ReadFile(filepath.Join(root, name)); err == nil {
			p.HasCompose = true
			p.ComposeServices, p.ComposeError = parseComposeServices(string(data))
			break
		}
	}
	if data, err := os.ReadFile(filepath.Join(root, "Makefile")); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if m := makeTarget.FindStringSubmatch(line); len(m) == 2 && !strings.HasPrefix(m[1], ".") {
				p.MakeTargets = append(p.MakeTargets, m[1])
			}
		}
	}
	p.HasEnvExample = exists(".env.example") || exists(".env.sample")
	p.HasEnvFile = exists(".env")
	p.GitBranch, p.HasGitRemote = gitInfo(root)
	return p
}

// parseComposeServices decodes Compose YAML and reads only the top-level
// services mapping. YAML syntax errors are returned to the caller.
func parseComposeServices(contents string) ([]string, error) {
	var document struct {
		Services map[string]yaml.Node `yaml:"services"`
	}
	if err := yaml.Unmarshal([]byte(contents), &document); err != nil {
		return nil, err
	}
	services := make([]string, 0, len(document.Services))
	for name := range document.Services {
		services = append(services, name)
	}
	return services, nil
}
