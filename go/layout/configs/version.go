package configs

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

var (
	buildDate string

	gitBranch string
	gitTag    string
	gitCommit string
)

func newVersion() versionInfo {
	return versionInfo{
		buildDate: buildDate,

		gitBranch: gitBranch,
		gitTag:    gitTag,
		gitCommit: gitCommit,

		goVersion: runtime.Version(),

		compiler: runtime.Compiler,
		platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

type versionInfo struct {
	buildDate string

	gitBranch string
	gitTag    string
	gitCommit string

	goVersion string

	compiler string
	platform string
}

func (vi versionInfo) String() string {
	var versionStr strings.Builder

	versionStr.WriteString("Build Date: " + vi.buildDate + "\n")

	versionStr.WriteString("Git Branch: " + vi.gitBranch + "\n")
	versionStr.WriteString("Git Tag: " + vi.gitTag + "\n")
	versionStr.WriteString("Git Commit: " + vi.gitCommit + "\n")

	versionStr.WriteString("Go Version: " + vi.goVersion + "\n")

	versionStr.WriteString("Compiler: " + vi.compiler + "\n")
	versionStr.WriteString("OS/Arch: " + vi.platform + "\n")

	return versionStr.String()
}

func IsPrintVersion() bool {
	if isUseVersion() {
		fmt.Println(newVersion())
		return true
	}

	return false
}

func isUseVersion() bool {
	return viper.GetBool("version")
}
