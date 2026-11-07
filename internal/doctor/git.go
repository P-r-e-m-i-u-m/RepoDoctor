package doctor

import (
 "os"
 "path/filepath"
 "strings"
)

func gitInfo(root string) (string,bool) {
 head,err:=os.ReadFile(filepath.Join(root,".git","HEAD")); if err!=nil { return "",false }
 value:=strings.TrimSpace(string(head)); branch:=""
 if strings.HasPrefix(value,"ref: refs/heads/") { branch=strings.TrimPrefix(value,"ref: refs/heads/") }
 config,_:=os.ReadFile(filepath.Join(root,".git","config"))
 return branch,strings.Contains(string(config),"[remote ")
}
