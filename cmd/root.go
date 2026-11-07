package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/P-r-e-m-i-u-m/RepoDoctor/internal/doctor"
	"github.com/P-r-e-m-i-u-m/RepoDoctor/internal/history"
	"github.com/spf13/cobra"
)

var commandText, errorText, repoPath string
var autoApprove bool

var rootCmd = &cobra.Command{
	Use:   "repodoctor",
	Short: "Diagnose failed developer commands and suggest safe fixes",
	Long:  "RepoDoctor inspects your repository and explains likely fixes. Nothing is executed without your confirmation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(commandText) == "" || strings.TrimSpace(errorText) == "" {
			return fmt.Errorf("both --command and --error are required")
		}
		abs, err := filepath.Abs(repoPath)
		if err != nil {
			return err
		}
		results := doctor.Analyze(abs, commandText, errorText)
		if len(results) == 0 {
			fmt.Println("No supported diagnosis matched this error. Nothing was executed.")
			return nil
		}
		for i, item := range results {
			fmt.Printf("\nDiagnosis %d: %s\n  %s\n", i+1, item.Title, item.Explanation)
			if item.Command == "" {
				fmt.Printf("  Safety: %s\n", item.Safety)
				continue
			}
			fmt.Printf("  Suggested: %s\n  Safety: %s\n", item.Command, item.Safety)
			if item.ManualOnly || doctor.IsDestructive(item.Command) {
				fmt.Println("  Manual review only. RepoDoctor will not execute this command.")
				continue
			}
			suggested := item.Command
			approved := autoApprove
			if !approved {
				var edited bool
				suggested, approved, edited = confirm(suggested)
				if edited {
					fmt.Printf("  Updated command: %s\\n", suggested)
				}
			}
			if !approved {
				fmt.Println("  Skipped.")
				continue
			}
			if doctor.IsDestructive(suggested) {
				fmt.Println("  That command matches a destructive-command safety rule. RepoDoctor will not run it.")
				continue
			}
			if err := runSuggested(abs, suggested); err != nil {
				fmt.Printf("  Command failed: %v\n", err)
				continue
			}
			fmt.Println("  Command completed.")
			recordFix(item.Title, suggested)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&commandText, "command", "c", "", "Command that failed")
	rootCmd.Flags().StringVarP(&errorText, "error", "e", "", "Error output")
	rootCmd.Flags().StringVarP(&repoPath, "dir", "d", ".", "Repository directory to inspect")
	rootCmd.Flags().BoolVar(&autoApprove, "yes", false, "Skip confirmation for safe suggested commands")
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
func confirm(command string) (string, bool, bool) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Run this command? [y]es / [n]o / [e]dit: ")
	line, _ := reader.ReadString('\\n')
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "y", "yes":
		return command, true, false
	case "e", "edit":
		fmt.Print("New command: ")
		edited, err := reader.ReadString('\\n')
		edited = strings.TrimSpace(edited)
		if err != nil && edited == "" {
			return command, false, false
		}
		if edited == "" {
			fmt.Println("  Empty command. Skipped.")
			return command, false, false
		}
		return edited, true, true
	default:
		return command, false, false
	}
}

func runSuggested(dir, command string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.CommandContext(ctx, "cmd.exe", "/C", command)
	} else {
		c = exec.CommandContext(ctx, "sh", "-c", command)
	}
	c.Dir = dir
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
func recordFix(title, command string) {
	path, err := history.DefaultPath()
	if err != nil {
		return
	}
	db, err := history.Open(path)
	if err != nil {
		return
	}
	defer db.Close()
	_ = db.Record(title, command, time.Now())
}
