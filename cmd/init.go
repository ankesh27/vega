package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	common "github.com/srijanone/vega/pkg/common"
	downloader "github.com/srijanone/vega/pkg/downloader"
	vega "github.com/srijanone/vega/pkg/core"
)

const (
	starterKitsRepoName = "git@github.com:Azure/draft.git" // TODO: Change this to vega once make public
	starterKitsDirName = "packs"
)

type initCmd struct {
	in     io.Reader
	out    io.Writer
	home   vega.Home
	dryRun bool
}

func newInitCmd(out io.Writer, in io.Reader) *cobra.Command {
	init := &initCmd{
		in:  in,
		out: out,
	}

	const initDesc = `sets up local configuration in $VEGA_HOME (default ~/.vega/) with default starter-kits`

	initCmd := &cobra.Command{
		Use:   "init",
		Short: initDesc,
		Long:  initDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("init called")
			if len(args) != 0 {
				return errors.New("Command does not accept arguments")
			}
			fmt.Println("Vega Home: ", homePath())
			init.home = vega.Home(homePath())
			return init.execute()
		},
	}

	return initCmd
}

func (iCmd *initCmd) execute() error {

	if !iCmd.dryRun {
		if err := iCmd.setupVegaHome(); err != nil {
			return err
		}
	}

	fmt.Fprintln(iCmd.out, "$VEGA_HOME has been initialized at %s", vegaHome)
	return nil
}

func (iCmd *initCmd) setupVegaHome() error {
	// Ensuring that required directory exists or not
	directories := []string{
		iCmd.home.String(),
		iCmd.home.StarterKits(),
		iCmd.home.Logs(),
	}

	for _, path := range directories {
		// TODO: One liner
		err := common.EnsureDir(path)
		fmt.Fprintln(iCmd.out, "Initializing", path)
		if err != nil {
			return err
		}
	}

	// Ensuring default starter kits exists or not
	d := downloader.Downloader{}
	sourceRepo := fmt.Sprintf("%s//%s", starterKitsRepoName, starterKitsDirName)
	fmt.Fprintln(iCmd.out, "Downloading starterkits...")
	d.Download(sourceRepo, iCmd.home.StarterKits())
	return nil
}