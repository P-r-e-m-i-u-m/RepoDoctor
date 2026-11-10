package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var pep508Name = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*`)

// inspectPythonDependencies reads declarations only. It never invokes a
// package manager or attempts to resolve import names to distribution names.
func inspectPythonDependencies(root string, p *Project) {
	p.PythonInstallCommand = make(map[string]string)
	p.PythonSources = make(map[string]string)

	for _, filename := range []string{"requirements.txt", "requirements-dev.txt"} {
		data, err := os.ReadFile(filepath.Join(root, filename))
		if err != nil {
			continue
		}
		p.HasRequirements = true
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
				continue
			}
			name := pep508Name.FindString(line)
			if name != "" {
				recordPythonDependency(p, name, filename, "python -m pip install -r "+filename)
			}
		}
	}

	if data, err := os.ReadFile(filepath.Join(root, "pyproject.toml")); err == nil {
		var document map[string]any
		if err := toml.Unmarshal(data, &document); err != nil {
			p.PythonMetadataError = append(p.PythonMetadataError, "pyproject.toml")
		} else {
			inspectPyProject(document, p)
		}
	}
	if data, err := os.ReadFile(filepath.Join(root, "Pipfile")); err == nil {
		var document map[string]any
		if err := toml.Unmarshal(data, &document); err != nil {
			p.PythonMetadataError = append(p.PythonMetadataError, "Pipfile")
		} else {
			for _, section := range []string{"packages", "dev-packages"} {
				if packages, ok := document[section].(map[string]any); ok {
					for name := range packages {
						recordPythonDependency(p, name, "Pipfile", "pipenv install")
					}
				}
			}
		}
	}
}

func inspectPyProject(document map[string]any, p *Project) {
	if project, ok := document["project"].(map[string]any); ok {
		addPEP621Dependencies(project["dependencies"], "pyproject.toml", "python -m pip install .", p)
		if optional, ok := project["optional-dependencies"].(map[string]any); ok {
			for _, dependencies := range optional {
				addPEP621Dependencies(dependencies, "pyproject.toml", "python -m pip install .", p)
			}
		}
	}

	tool, _ := document["tool"].(map[string]any)
	poetry, _ := tool["poetry"].(map[string]any)
	if dependencies, ok := poetry["dependencies"].(map[string]any); ok {
		for name := range dependencies {
			if !strings.EqualFold(name, "python") {
				recordPythonDependency(p, name, "pyproject.toml", "poetry install")
			}
		}
	}
	groups, _ := poetry["group"].(map[string]any)
	for _, groupValue := range groups {
		group, _ := groupValue.(map[string]any)
		dependencies, _ := group["dependencies"].(map[string]any)
		for name := range dependencies {
			recordPythonDependency(p, name, "pyproject.toml", "poetry install")
		}
	}
}

func addPEP621Dependencies(value any, source, command string, p *Project) {
	dependencies, ok := value.([]any)
	if !ok {
		return
	}
	for _, item := range dependencies {
		spec, ok := item.(string)
		if !ok {
			continue
		}
		name := pep508Name.FindString(spec)
		if name != "" {
			recordPythonDependency(p, name, source, command)
		}
	}
}

func recordPythonDependency(p *Project, name, source, command string) {
	normalized := normalizePythonName(name)
	if normalized == "" {
		return
	}
	p.Requirements[normalized] = true
	p.PythonSources[normalized] = source
	p.PythonInstallCommand[normalized] = command
}

func normalizePythonName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.NewReplacer("-", "_", ".", "_").Replace(name)
	return strings.Trim(name, "_")
}

func pythonDependencyExplanation(importName, source string) string {
	return fmt.Sprintf("%q is declared in %s, but it is unavailable in the active Python environment. Import names can differ from distribution names, so RepoDoctor only matches exact normalized names.", importName, source)
}
