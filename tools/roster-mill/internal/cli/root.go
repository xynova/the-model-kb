package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/xynova/roster-mill/internal/aa"
	"github.com/xynova/roster-mill/internal/apply"
	derrors "github.com/xynova/roster-mill/internal/errors"
	"github.com/xynova/roster-mill/internal/generate"
	"github.com/xynova/roster-mill/internal/limits"
	"github.com/xynova/roster-mill/internal/opconfig"
	"github.com/xynova/roster-mill/internal/polypus"
	"github.com/xynova/roster-mill/internal/repo"
	"github.com/xynova/roster-mill/internal/reviewbody"
)

// Version is set at link time; default for go run.
var Version = "dev"

func repoRootFlag(cmd *cobra.Command) string {
	v, _ := cmd.Flags().GetString("repo-root")
	return v
}

func resolveRoot(cmd *cobra.Command) (string, error) {
	flag := repoRootFlag(cmd)
	if flag != "" {
		return repo.FindRoot(flag)
	}
	return repo.FindRoot("")
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "roster-mill",
		Short:         "Compile roster cards and retrieve host limits / Artificial Analysis scores",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			printAgentGuide(cmd.OutOrStdout())
			return nil
		},
	}
	root.PersistentFlags().String("repo-root", "", "repository root (default: walk up to roster/taxonomy.md)")
	root.AddCommand(
		newVersionCmd(),
		newGenerateCmd(),
		newRetrieveLimitsCmd(),
		newApplyMineCmd(),
		newRetrieveAACmd(),
		newApplyAACmd(),
		newPolypusDiffCmd(),
		newInitCmd(),
		newStoreAAKeyCmd(),
		newReviewBodyCmd(),
	)
	return root
}

func printAgentGuide(w io.Writer) {
	fmt.Fprintf(w, "roster-mill %s\n", Version)
	fmt.Fprintln(w, "Compile roster cards and retrieve host limits / Artificial Analysis scores.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Documentation for agents")
	fmt.Fprintln(w, "  tools/roster-mill/README.md")
	fmt.Fprintln(w, "  .cursor/skills/model-roster/SKILL.md")
	fmt.Fprintln(w, "  .cursor/skills/providers-mint/SKILL.md")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Read-only")
	fmt.Fprintln(w, "  version, help, polypus-diff, review-body")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Execute")
	fmt.Fprintln(w, "  generate          Validate cards; write Hugo pages, catalog.json, INDEX.md")
	fmt.Fprintln(w, "  retrieve-limits   Fetch public host pages into limits-review.json")
	fmt.Fprintln(w, "  apply-mine        Apply limits review + capacity-decisions.yaml")
	fmt.Fprintln(w, "  retrieve-aa       Fetch Artificial Analysis scores")
	fmt.Fprintln(w, "  apply-aa          File or clear aa_intelligence on cards")
	fmt.Fprintln(w, "  init              Write ~/.config/the-model-kb/config.yaml")
	fmt.Fprintln(w, "  store-aa-key      Copy AA_API_KEY from env into the OS keyring")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Full command list: roster-mill help")
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print roster-mill version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), Version)
		},
	}
}

func newGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Validate roster cards and write Hugo pages plus catalog.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			_, err = generate.Run(root, cmd.OutOrStdout())
			return err
		},
	}
}

func newRetrieveLimitsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "retrieve-limits",
		Short: "Fetch public host pages into roster/retrieval/limits-review.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			return limits.Run(root, cmd.OutOrStdout())
		},
	}
}

func newApplyMineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply-mine",
		Short: "Apply limits-review.json onto roster cards",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			return apply.RunMine(root)
		},
	}
}

func newRetrieveAACmd() *cobra.Command {
	return &cobra.Command{
		Use:   "retrieve-aa",
		Short: "Fetch Artificial Analysis Intelligence Index scores",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			env := opconfig.EnvironFromOS()
			opconfig.LoadDotenvFiles(root, env)
			for k, v := range env {
				if os.Getenv(k) == "" {
					_ = os.Setenv(k, v)
				}
			}
			key, err := opconfig.ResolveAAAPIKey(opconfig.ResolveOptions{
				Environ:  env,
				RepoRoot: root,
			})
			if err != nil {
				return err
			}
			return aa.RunRetrieve(root, key, cmd.OutOrStdout())
		},
	}
}

func newApplyAACmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply-aa",
		Short: "Apply aa-intelligence-review.json onto roster cards",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			return apply.RunAA(root)
		},
	}
}

func newPolypusDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "polypus-diff",
		Short: "Diff Polypus model ids against roster/catalog.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			catalog, _ := cmd.Flags().GetString("catalog")
			base, _ := cmd.Flags().GetString("base-url")
			full, _ := cmd.Flags().GetString("url")
			asJSON, _ := cmd.Flags().GetBool("json")
			return polypus.Run(root, catalog, base, full, asJSON, cmd.OutOrStdout())
		},
	}
	cmd.Flags().String("catalog", "", "Path to catalog.json")
	cmd.Flags().String("base-url", "", "Polypus base URL")
	cmd.Flags().String("url", "", "Full models URL (skips the default pair)")
	cmd.Flags().Bool("json", false, "Print the report as JSON")
	return cmd
}

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Write ~/.config/the-model-kb/config.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			force, _ := cmd.Flags().GetBool("force")
			exampleSrc := root + "/site/config.yaml.example"
			live, example, err := opconfig.InitUserConfig("", force, exampleSrc)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", example)
			fmt.Fprintf(cmd.OutOrStdout(), "live %s\n", live)
			return nil
		},
	}
	cmd.Flags().Bool("force", false, "overwrite live config.yaml")
	return cmd
}

func newStoreAAKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "store-aa-key",
		Short: "Copy AA_API_KEY from env into the OS keyring",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			env := opconfig.EnvironFromOS()
			opconfig.LoadDotenvFiles(root, env)
			if err := opconfig.StoreAAKey(env[opconfig.SecretAA], nil); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "stored %s in keyring service=%s\n", opconfig.SecretAA, opconfig.AppDir)
			return nil
		},
	}
}

func newReviewBodyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "review-body",
		Short: "Print a markdown PR body from limits and AA review JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := resolveRoot(cmd)
			if err != nil {
				return err
			}
			body, err := reviewbody.Markdown(
				root+"/roster/retrieval/limits-review.json",
				root+"/roster/retrieval/aa-intelligence-review.json",
			)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), body)
			return nil
		},
	}
}

func mapExit(err error) int {
	if derrors.IsCode(err, derrors.CodeInvalidArgument) {
		return 2
	}
	return 1
}

// Execute runs the Cobra root command.
func Execute() int {
	root := newRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return mapExit(err)
	}
	return 0
}
