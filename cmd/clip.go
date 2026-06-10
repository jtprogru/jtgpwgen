package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
)

// clipboardCmd describes how to write to and clear the system clipboard on a
// given OS.
type clipboardCmd struct {
	name      string   // executable to lookup in PATH
	args      []string // arguments passed when writing via stdin
	clearArgs []string // arguments for an explicit clear; empty means "write empty string"
}

// clipboardCmds returns the candidate commands for the current OS, in order
// of preference. The first one whose binary is found in PATH is used.
func clipboardCmds() []clipboardCmd {
	switch runtime.GOOS {
	case "darwin":
		return []clipboardCmd{{name: "pbcopy"}}
	case "linux":
		return []clipboardCmd{
			{name: "wl-copy", clearArgs: []string{"--clear"}},
			{name: "xclip", args: []string{"-selection", "clipboard"}},
			{name: "xsel", args: []string{"--clipboard", "--input"}, clearArgs: []string{"--clipboard", "--delete"}},
		}
	default:
		return nil
	}
}

// firstAvailable returns the first candidate whose binary is found in PATH
// along with its resolved path.
func firstAvailable() (clipboardCmd, string, error) {
	cands := clipboardCmds()
	if len(cands) == 0 {
		return clipboardCmd{}, "", fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
	var tried []string
	for _, c := range cands {
		path, err := exec.LookPath(c.name)
		if err != nil {
			tried = append(tried, c.name)
			continue
		}
		return c, path, nil
	}
	return clipboardCmd{}, "", fmt.Errorf("no clipboard tool found in PATH (tried: %v)", tried)
}

// writeViaStdin runs the tool at path with args and feeds s to its stdin.
func writeViaStdin(name, path string, args []string, s string) error {
	// #nosec G204 -- path is resolved via LookPath from a fixed allowlist of
	// tool names; args are static, no user-controlled input reaches the command.
	cmd := exec.Command(path, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("%s stdin pipe: %w", name, err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s start: %w", name, err)
	}
	if _, err := stdin.Write([]byte(s)); err != nil {
		_ = stdin.Close()
		_ = cmd.Wait()
		return fmt.Errorf("%s write: %w", name, err)
	}
	if err := stdin.Close(); err != nil {
		_ = cmd.Wait()
		return fmt.Errorf("%s close stdin: %w", name, err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s wait: %w", name, err)
	}
	return nil
}

// copyToClipboard writes s to the system clipboard. It returns the name of
// the underlying tool used (for logging) or an error if no supported tool
// is available.
func copyToClipboard(s string) (string, error) {
	c, path, err := firstAvailable()
	if err != nil {
		return "", err
	}
	if err := writeViaStdin(c.name, path, c.args, s); err != nil {
		return "", err
	}
	return c.name, nil
}

// clearClipboard removes the previously copied value from the clipboard.
// Tools with an explicit clear operation (wl-copy --clear, xsel --delete)
// use it; the rest get an empty string written instead.
func clearClipboard() (string, error) {
	c, path, err := firstAvailable()
	if err != nil {
		return "", err
	}
	if len(c.clearArgs) > 0 {
		// #nosec G204 -- path is resolved via LookPath from a fixed allowlist
		// of tool names; clearArgs are static.
		if err := exec.Command(path, c.clearArgs...).Run(); err != nil {
			return "", fmt.Errorf("%s clear: %w", c.name, err)
		}
		return c.name, nil
	}
	if err := writeViaStdin(c.name, path, c.args, ""); err != nil {
		return "", err
	}
	return c.name, nil
}
