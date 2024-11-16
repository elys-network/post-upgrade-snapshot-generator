package version

import (
	"fmt"

	appversion "github.com/elys-network/post-upgrade-snapshot-generator/version"
	"github.com/spf13/cobra"
)

// version command shows the version of the application
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(appversion.Version)
		},
	}
}
