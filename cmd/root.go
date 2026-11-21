package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/devsy-org/devsy/logging"
	"github.com/devsy-org/devsy/pkg/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// NewRootCmd returns a new root command.
func NewRootCmd() *cobra.Command {
	gcloudCmd := &cobra.Command{
		Use:           "devsy-provider-gcloud",
		Short:         "gcloud Provider commands",
		SilenceErrors: true,
		SilenceUsage:  true,

		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			config, invalid := logging.ResolveProvider(
				os.Getenv("DEVSY_LOG_LEVEL"),
				os.Getenv("DEVSY_LOG_FORMAT"),
			)
			if invalid.Level {
				return errors.New("invalid DEVSY_LOG_LEVEL")
			}
			if invalid.Format {
				return errors.New("invalid DEVSY_LOG_FORMAT")
			}

			log.Init(
				log.Config{
					Level:        logging.LevelName(config.Level),
					DefaultLevel: log.LevelInfoName,
					Format:       string(config.Format),
				},
			)

			return nil
		},
	}

	return gcloudCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// build the root command
	rootCmd := BuildRoot()

	// execute command
	err := rootCmd.Execute()
	if err != nil {
		var sshExitErr *ssh.ExitError
		if errors.As(err, &sshExitErr) {
			os.Exit(sshExitErr.ExitStatus())
		}
		var execExitErr *exec.ExitError
		if errors.As(err, &execExitErr) {
			if len(execExitErr.Stderr) > 0 {
				log.Error(string(execExitErr.Stderr))
			}
			os.Exit(execExitErr.ExitCode())
		}

		log.Fatal(err)
	}
}

// BuildRoot creates a new root command from the available subcommands.
func BuildRoot() *cobra.Command {
	rootCmd := NewRootCmd()

	rootCmd.AddCommand(NewCreateCmd())
	rootCmd.AddCommand(NewStatusCmd())
	rootCmd.AddCommand(NewDeleteCmd())
	rootCmd.AddCommand(NewStartCmd())
	rootCmd.AddCommand(NewStopCmd())
	rootCmd.AddCommand(NewCommandCmd())
	rootCmd.AddCommand(NewTokenCmd())
	rootCmd.AddCommand(NewInitCmd())
	rootCmd.AddCommand(NewDescribeCmd())
	return rootCmd
}
