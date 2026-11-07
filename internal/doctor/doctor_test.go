package doctor

import (
 "os"
 "path/filepath"
 "testing"
)

func writeFile(t *testing.T,root,name,content string) {
 t.Helper(); path:=filepath.Join(root,name)
 if err:=os.MkdirAll(filepath.Dir(path),0755); err!=nil { t.Fatal(err) }
 if err:=os.WriteFile(path,[]byte(content),0644); err!=nil { t.Fatal(err) }
}
func TestMissingNodeModulesUsesDetectedPackageManager(t *testing.T) {
 root:=t.TempDir(); writeFile(t,root,"package.json", "{\"scripts\":{\"dev\":\"vite\"}}"); writeFile(t,root,"pnpm-lock.yaml","lockfileVersion: '9.0'")
 got:=Analyze(root,"pnpm dev","Error: ENOENT: no such file or directory, node_modules")
 if len(got)!=1 || got[0].Command!="pnpm install" { t.Fatalf("unexpected diagnosis: %#v",got) }
}
func TestDeclaredPythonModuleSuggestsRequirementsInstall(t *testing.T) {
 root:=t.TempDir(); writeFile(t,root,"requirements.txt","fastapi==0.115.0\n")
 got:=Analyze(root,"python app.py","ModuleNotFoundError: No module named 'fastapi'")
 if len(got)!=1 || got[0].Command!="python -m pip install -r requirements.txt" { t.Fatalf("unexpected diagnosis: %#v",got) }
}
func TestUnknownPythonModuleIsManualOnly(t *testing.T) {
 got:=Analyze(t.TempDir(),"python app.py","ModuleNotFoundError: No module named 'mystery_pkg'")
 if len(got)!=1 || !got[0].ManualOnly || got[0].Command!="" { t.Fatalf("unexpected diagnosis: %#v",got) }
}
func TestDestructiveCommandsAreDetected(t *testing.T) {
 for _,command:=range []string{"rm -rf .","git reset --hard","git push origin main --force","DROP TABLE users"} { if !IsDestructive(command) { t.Errorf("expected %q to be blocked",command) } }
}
func TestSafeCommandIsNotDestructive(t *testing.T) { if IsDestructive("npm install") { t.Fatal("npm install must not be classified as destructive") } }
