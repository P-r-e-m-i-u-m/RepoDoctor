package doctor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
			p.ComposeServices = parseComposeServices(string(data))
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

// parseComposeServices reads service keys only from the top-level services
// mapping. It supports ordinary block-style Compose YAML with arbitrary
// indentation; it is intentionally not a general-purpose YAML validator.
func parseComposeServices(contents string) []string {
	servicesIndent := -1
	serviceIndent := -1
	inServices := false
	var services []string

	for _, line := range strings.Split(contents, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))
		if strings.HasPrefix(strings.TrimLeft(line, " "), "\t") {
			continue
		}
		entry := strings.TrimSpace(line)
		key, _, hasValue := strings.Cut(entry, ":")
		if !hasValue {
			continue
		}
		key = strings.TrimSpace(key)
		key = strings.Trim(key, "\"'")

		if !inServices {
			if indent == 0 && key == "services" {
				inServices = true
				servicesIndent = indent
			}
			continue
		}

		if indent <= servicesIndent {
			break
		}
		if serviceIndent == -1 {
			serviceIndent = indent
		}
		if indent == serviceIndent && key != "" {
			services = append(services, key)
		}
	}

	return services
}
