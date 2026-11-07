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
	NPMScripts      map[string]bool
	HasRequirements bool
	Requirements    map[string]bool
	HasVenv         bool
	HasCompose      bool
	ComposeServices []string
	MakeTargets     []string
	GitBranch       string
	HasGitRemote    bool
	HasEnvExample   bool
	HasEnvFile      bool
}
