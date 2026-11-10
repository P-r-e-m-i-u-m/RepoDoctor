package doctor

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func writeFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
func TestMissingNodeModulesUsesDetectedPackageManager(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "package.json", "{\"scripts\":{\"dev\":\"vite\"}}")
	writeFile(t, root, "pnpm-lock.yaml", "lockfileVersion: '9.0'")
	got := Analyze(root, "pnpm dev", "Error: ENOENT: no such file or directory, node_modules")
	if len(got) != 1 || got[0].Command != "pnpm install" {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}
func TestDeclaredPythonModuleSuggestsRequirementsInstall(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "requirements.txt", "fastapi==0.115.0\n")
	got := Analyze(root, "python app.py", "ModuleNotFoundError: No module named 'fastapi'")
	if len(got) != 1 || got[0].Command != "python -m pip install -r requirements.txt" {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}
func TestUnknownPythonModuleIsManualOnly(t *testing.T) {
	got := Analyze(t.TempDir(), "python app.py", "ModuleNotFoundError: No module named 'mystery_pkg'")
	if len(got) != 1 || !got[0].ManualOnly || got[0].Command != "" {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}
func TestDestructiveCommandsAreDetected(t *testing.T) {
	for _, command := range []string{"rm -rf .", "git reset --hard", "git push origin main --force", "DROP TABLE users"} {
		if !IsDestructive(command) {
			t.Errorf("expected %q to be blocked", command)
		}
	}
}
func TestSafeCommandIsNotDestructive(t *testing.T) {
	if IsDestructive("npm install") {
		t.Fatal("npm install must not be classified as destructive")
	}
}

func TestUnknownScriptDiagnosedWhenPackageHasNoScripts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "package.json", "{\"name\":\"sample\",\"version\":\"1.0.0\"}")
	got := Analyze(root, "npm run dev", "npm ERR! Missing script: \"dev\"")
	if len(got) != 1 || got[0].Title != "Unknown package script" || !strings.Contains(got[0].Explanation, "no scripts configured") {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}

func TestUnknownScriptListsAvailableScripts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "package.json", "{\"scripts\":{\"start\":\"node app.js\"}}")
	got := Analyze(root, "npm run dev", "npm ERR! Missing script: \"dev\"")
	if len(got) != 1 || !strings.Contains(got[0].Explanation, "start") {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}

func TestAdditionalDestructiveCommandsAreDetected(t *testing.T) {
	for _, command := range []string{
		"git clean -fd",
		"git clean -xdf",
		"git clean --force",
		"docker compose down --volumes",
		"docker-compose down -v",
		"Remove-Item -Recurse -Force .\\build",
		"Remove-Item -Force -Recurse .\\build",
	} {
		if !IsDestructive(command) {
			t.Errorf("expected %q to be blocked", command)
		}
	}
}

func TestAdditionalSafeCommandsAreNotDestructive(t *testing.T) {
	for _, command := range []string{
		"git clean -n",
		"docker compose down",
		"docker compose ps",
		"Remove-Item .\\one-file.txt",
	} {
		if IsDestructive(command) {
			t.Errorf("expected %q to remain allowed", command)
		}
	}
}

func TestGitPushForceShorthandIsDetected(t *testing.T) {
	for _, command := range []string{
		"git push -f origin main",
		"git push origin main -f",
		"git push -uf origin main",
		"git push origin main -fu",
	} {
		if !IsDestructive(command) {
			t.Errorf("expected %q to be blocked", command)
		}
	}
}

func TestOrdinaryGitPushIsNotDestructive(t *testing.T) {
	for _, command := range []string{
		"git push",
		"git push origin main",
		"git push --set-upstream origin feature",
	} {
		if IsDestructive(command) {
			t.Errorf("expected %q to remain allowed", command)
		}
	}
}

func TestParseComposeServicesWithArbitraryIndentationAndComments(t *testing.T) {
	contents := `name: demo
services:
    # App services
    api:
        image: example/api
        environment:
            MODE: test
    worker: # background worker
        image: example/worker
volumes:
    data:
        driver: local
networks:
    internal:
        driver: bridge
`

	want := []string{"api", "worker"}
	got, err := parseComposeServices(contents)
	if err != nil {
		t.Fatalf("parseComposeServices() returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseComposeServices() = %#v, want %#v", got, want)
	}
}

func TestParseComposeServicesDoesNotTreatOtherTopLevelKeysAsServices(t *testing.T) {
	contents := `services:
  web:
    image: nginx
volumes:
  cache:
networks:
  backend:
`

	want := []string{"web"}
	got, err := parseComposeServices(contents)
	if err != nil {
		t.Fatalf("parseComposeServices() returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseComposeServices() = %#v, want %#v", got, want)
	}
}

func TestParseComposeServicesSupportsQuotedServiceNames(t *testing.T) {
	contents := `services:
  "web":
    image: nginx
  'worker':
    image: example/worker
`

	want := []string{"web", "worker"}
	got, err := parseComposeServices(contents)
	if err != nil {
		t.Fatalf("parseComposeServices() returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseComposeServices() = %#v, want %#v", got, want)
	}
}

func TestParseComposeServicesRejectsMalformedYAML(t *testing.T) {
	contents := "services:\n  web:\n    image: nginx\n  worker\n    image: example/worker\n"
	if _, err := parseComposeServices(contents); err == nil {
		t.Fatal("expected malformed YAML to return an error")
	}
}

func TestMalformedComposeFileProducesUsefulDiagnosis(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "compose.yaml", "services:\n  web:\n    image: nginx\n  worker\n    image: example/worker\n")
	got := Analyze(root, "docker compose up missing", "no such service: missing")
	if len(got) != 1 || got[0].Title != "Docker Compose file could not be parsed" || !strings.Contains(got[0].Explanation, "Fix the YAML syntax") {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}

func TestPythonProjectMetadataDependencyFormats(t *testing.T) {
	tests := []struct {
		name, file, contents, importName, wantCommand string
	}{
		{"PEP 621", "pyproject.toml", "[project]\nname = \"demo\"\ndependencies = [\"fastapi>=0.1\"]\n", "fastapi", "python -m pip install ."},
		{"Poetry", "pyproject.toml", "[tool.poetry]\nname = \"demo\"\n[tool.poetry.dependencies]\npython = \"^3.12\"\nhttpx = \"^0.27\"\n", "httpx", "poetry install"},
		{"Poetry group", "pyproject.toml", "[tool.poetry.group.dev.dependencies]\npytest = \"^8.0\"\n", "pytest", "poetry install"},
		{"Pipfile", "Pipfile", "[packages]\nrequests = \"*\"\n", "requests", "pipenv install"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, tt.file, tt.contents)
			got := Analyze(root, "python app.py", "ModuleNotFoundError: No module named '"+tt.importName+"'")
			if len(got) != 1 || got[0].Command != tt.wantCommand || !strings.Contains(got[0].Explanation, tt.file) {
				t.Fatalf("unexpected diagnosis: %#v", got)
			}
		})
	}
}

func TestMalformedPythonMetadataDoesNotGuessInstallCommand(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "pyproject.toml", "[project\ndependencies = [\"requests\"]")
	got := Analyze(root, "python app.py", "ModuleNotFoundError: No module named 'requests'")
	if len(got) != 1 || !got[0].ManualOnly || got[0].Command != "" || !strings.Contains(got[0].Explanation, "pyproject.toml") {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}

func TestRequirementsDevUsesMatchingInstallFile(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "requirements-dev.txt", "pytest>=8\n")
	got := Analyze(root, "python -m pytest", "ModuleNotFoundError: No module named 'pytest'")
	if len(got) != 1 || got[0].Command != "python -m pip install -r requirements-dev.txt" {
		t.Fatalf("unexpected diagnosis: %#v", got)
	}
}
