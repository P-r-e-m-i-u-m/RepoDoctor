package doctor

import (
 "fmt"
 "path/filepath"
 "strings"
 "regexp"
)

var moduleMissing=regexp.MustCompile("No module named ['\"]([^'\"]+)['\"]")
var serviceMissing=regexp.MustCompile("(?i)no such service: ([A-Za-z0-9_.-]+)")
var unknownScript=regexp.MustCompile("(?i)missing script: [\"']?([A-Za-z0-9_.-]+)")

func Analyze(root,command,output string) []Diagnosis {
 p:=inspectProjectFiles(root); cmd,out:=strings.ToLower(command),strings.ToLower(output)
 var found []Diagnosis
 if (strings.Contains(out,"enoent") && strings.Contains(out,"node_modules")) || (strings.Contains(out,"cannot find module") && strings.Contains(cmd,"npm")) {
  pm:=p.PackageManager; if pm=="" { pm="npm" }
  found=append(found,Diagnosis{"Node.js dependencies are missing","The repository appears to be missing installed Node.js dependencies.",pm+" install","Safe after confirmation",false})
 }
 if m:=unknownScript.FindStringSubmatch(output); len(m)==2 {
  available:=make([]string,0,len(p.NPMScripts)); for name:=range p.NPMScripts { available=append(available,name) }
  if len(available)>0 { found=append(found,Diagnosis{"Unknown package script",fmt.Sprintf("Script %q is not declared in package.json. Available scripts: %s.",m[1],strings.Join(available,", ")),"","Information only",false}) }
 }
 if m:=moduleMissing.FindStringSubmatch(output); len(m)==2 {
  name:=strings.ToLower(strings.ReplaceAll(m[1],"-","_"))
  if p.Requirements[name] { found=append(found,Diagnosis{"Python dependency is missing",fmt.Sprintf("%q is declared in requirements but is unavailable in the active environment.",m[1]),"python -m pip install -r requirements.txt","Safe after confirmation",false})
  } else { found=append(found,Diagnosis{"Python module is missing",fmt.Sprintf("Python cannot import %q, and it was not found in the inspected requirements files.",m[1]),"","Manual review required",true}) }
 }
 if strings.Contains(out,"no such file or directory") && (strings.Contains(out,".venv") || strings.Contains(out,"activate")) && !p.HasVenv {
  found=append(found,Diagnosis{"Python virtual environment is missing","No .venv or venv directory was found.","python -m venv .venv","Safe after confirmation",false})
 }
 if strings.Contains(out,"no such file or directory") && strings.Contains(out,".env") && !p.HasEnvFile && p.HasEnvExample {
  cp:="cp .env.example .env"; if isWindows() { cp="Copy-Item .env.example .env" }
  found=append(found,Diagnosis{"Environment file is missing","A sample environment file exists. Copy it, then fill in the required values.",cp,"Creates a local config file. Review values before use.",false})
 }
 if m:=serviceMissing.FindStringSubmatch(output); len(m)==2 && p.HasCompose {
  found=append(found,Diagnosis{"Docker Compose service was not found",fmt.Sprintf("Unknown service %q. Services detected: %s.",m[1],strings.Join(p.ComposeServices,", ")),"","Information only",false})
 }
 if strings.Contains(out,"address already in use") || strings.Contains(out,"eaddrinuse") {
  found=append(found,Diagnosis{"Port is already in use","Another process is listening on the requested port. Identify it and stop it manually, or configure the app to use another port.","","Manual review required",true})
 }
 if strings.Contains(out,"no upstream branch") {
  if p.GitBranch!="" && p.HasGitRemote { found=append(found,Diagnosis{"Git upstream is not configured","The current branch has no upstream tracking branch.","git push -u origin "+p.GitBranch,"Network operation. Explicit confirmation required.",false})
  } else { found=append(found,Diagnosis{"Git upstream is not configured","Verify the intended remote and branch, then set an upstream manually.","","Manual review required",true}) }
 }
 if strings.Contains(out,"no rule to make target") && len(p.MakeTargets)>0 {
  found=append(found,Diagnosis{"Unknown Make target","Available targets: "+strings.Join(p.MakeTargets,", "),"","Information only",false})
 }
 if strings.Contains(out,"unknown command") && strings.Contains(cmd,"git") {
  found=append(found,Diagnosis{"Possible Git command typo","Check the subcommand spelling and consult git help.","","Information only",false})
 }
 if strings.Contains(out,"command not found") && strings.Contains(out,"docker") && !p.HasCompose {
  found=append(found,Diagnosis{"Docker Compose project not detected","No Compose configuration file was found in "+filepath.Base(root)+".","","Information only",false})
 }
 return found
}
