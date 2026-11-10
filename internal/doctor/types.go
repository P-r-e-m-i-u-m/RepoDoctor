package doctor

type Diagnosis struct {
	Title       string
	Explanation string
	Command     string
	Safety      string
	ManualOnly  bool
}

type Project struct {
	PackageManager  string
	HasNodeModules  bool
	HasPackageJSON  bool
	NPMScripts      map[string]bool
	HasRequirements bool
	Requirements         map[string]bool
	PythonInstallCommand map[string]string
	PythonSources        map[string]string
	PythonMetadataError  []string
	HasVenv         bool
	HasCompose      bool
	ComposeServices []string
	ComposeError    string
	MakeTargets     []string
	GitBranch       string
	HasGitRemote    bool
	HasEnvExample   bool
	HasEnvFile      bool
}
