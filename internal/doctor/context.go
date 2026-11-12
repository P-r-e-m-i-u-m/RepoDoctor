package doctor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

var makeTarget = regexp.MustCompile("^([A-Za-z0-9_.-]+):")

func inspectProjectFiles(root string) Project {
	p := Project{NPMScripts: map[string]bool{}, Requirements: map[string]bool{}}
	exists := func(name string) bool { _, err := os.Stat(filepath.Join(root, name)); return err == nil }
	lockfileManagers := make([]string, 0, 3)
	if exists("pnpm-lock.yaml") {
		lockfileManagers = append(lockfileManagers, "pnpm")
	}
	if exists("yarn.lock") {
		lockfileManagers = append(lockfileManagers, "yarn")
	}
	if exists("package-lock.json") {
		lockfileManagers = append(lockfileManagers, "npm")
	}
	p.HasNodeModules = exists("node_modules")
	if data, err := os.ReadFile(filepath.Join(root, "package.json")); err == nil {
		p.HasPackageJSON = true
		var manifest struct {
			Scripts        map[string]json.RawMessage `json:"scripts"`
			PackageManager string                     `json:"packageManager"`
		}
		if json.Unmarshal(data, &manifest) == nil {
			for name := range manifest.Scripts {
				p.NPMScripts[name] = true
			}
			if manifest.PackageManager != "" {
				manager := strings.SplitN(manifest.PackageManager, "@", 2)[0]
				if manager == "npm" || manager == "pnpm" || manager == "yarn" {
					p.PackageManager = manager
				} else {
					p.PackageManagerConflict = "package.json declares an unsupported packageManager value: " + manifest.PackageManager
				}
			}
		}
	}
	if p.PackageManagerConflict == "" {
		if p.PackageManager != "" {
			for _, manager := range lockfileManagers {
				if manager != p.PackageManager {
					p.PackageManagerConflict = "package.json declares " + p.PackageManager + " but the repository also contains a " + manager + " lockfile"
					break
				}
			}
		} else if len(lockfileManagers) == 1 {
			p.PackageManager = lockfileManagers[0]
		} else if len(lockfileManagers) > 1 {
			p.PackageManagerConflict = "multiple package-manager lockfiles are present (" + strings.Join(lockfileManagers, ", ") + ") and package.json does not declare packageManager"
		}
	}
	inspectPythonDependencies(root, &p)
	p.HasVenv = exists(".venv") || exists("venv")
	for _, name := range []string{"compose.yaml", "compose.yml", "docker-compose.yml", "docker-compose.yaml"} {
		if data, err := os.ReadFile(filepath.Join(root, name)); err == nil {
			p.HasCompose = true
			services, parseErr := parseComposeServices(string(data))
			p.ComposeServices = services
			if parseErr != nil {
				p.ComposeError = parseErr.Error()
			}
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
	sort.Strings(services)
	return services, nil
}
