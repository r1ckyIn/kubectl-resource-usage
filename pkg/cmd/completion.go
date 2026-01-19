package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const completionLongDesc = `Generate shell completion scripts for kubectl-resource-usage.

To load completions:

Bash:
  $ source <(kubectl resource-usage completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ kubectl resource-usage completion bash > /etc/bash_completion.d/kubectl-resource_usage
  # macOS:
  $ kubectl resource-usage completion bash > /usr/local/etc/bash_completion.d/kubectl-resource_usage

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ kubectl resource-usage completion zsh > "${fpath[1]}/_kubectl-resource_usage"
  # You may need to start a new shell for this setup to take effect.

Fish:
  $ kubectl resource-usage completion fish | source
  # To load completions for each session, execute once:
  $ kubectl resource-usage completion fish > ~/.config/fish/completions/kubectl-resource_usage.fish

PowerShell:
  PS> kubectl resource-usage completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> kubectl resource-usage completion powershell > kubectl-resource_usage.ps1
  # and source this file from your PowerShell profile.
`

// NewCmdCompletion creates the completion command
func NewCmdCompletion() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate shell completion scripts",
		Long:                  completionLongDesc,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}

	return cmd
}
