package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/jtprogru/jtgpwgen/internal/passgen"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "today"
	BuiltBy = "go build"
)

func newRootCmd() *cobra.Command {
	var flags struct {
		length    int
		special   string
		digits    bool
		noSpecial bool
		noDigits  bool
		memo      bool
	}

	cmd := &cobra.Command{
		Use:     "jtgpwgen",
		Short:   "Генератор паролей с настраиваемыми классами символов",
		Version: Version,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			digitsSet := cmd.Flags().Changed("digits")
			specialSet := cmd.Flags().Changed("special")
			noDigitsSet := cmd.Flags().Changed("no-digits")
			noSpecialSet := cmd.Flags().Changed("no-special")

			if digitsSet && noDigitsSet {
				return passgen.ErrConflictDigits
			}
			if specialSet && noSpecialSet {
				return passgen.ErrConflictSpecial
			}
			if flags.memo && (specialSet || noSpecialSet || digitsSet || noDigitsSet) {
				return passgen.ErrMemoIncompatibleFlag
			}

			opts := passgen.DefaultOptions()
			opts.Length = flags.length
			opts.Memo = flags.memo
			if noDigitsSet {
				opts.UseDigits = false
			}
			if noSpecialSet {
				opts.UseSpecial = false
			}
			if specialSet {
				opts.UseSpecial = true
				opts.ExtraSpecial = flags.special
			}

			pw, err := passgen.Generate(opts)
			if err != nil {
				if errors.Is(err, passgen.ErrLengthOutOfRange) ||
					errors.Is(err, passgen.ErrNoCharClasses) ||
					errors.Is(err, passgen.ErrMemoEntropyTooLow) {
					return err
				}
				return fmt.Errorf("generate: %w", err)
			}
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), pw); err != nil {
				return fmt.Errorf("write output: %w", err)
			}
			return nil
		},
	}
	cmd.SetVersionTemplate(versionTemplate())

	f := cmd.Flags()
	f.IntVarP(&flags.length, "length", "l", passgen.DefaultLength, "password length")
	f.StringVarP(&flags.special, "special", "s", "", "extra special characters to include (added to default '@')")
	f.BoolVarP(&flags.digits, "digits", "d", false, "explicitly enable digits (default: enabled)")
	f.BoolVar(&flags.noSpecial, "no-special", false, "disable special characters entirely")
	f.BoolVar(&flags.noDigits, "no-digits", false, "disable digits")
	f.BoolVarP(&flags.memo, "memo", "m", false, "generate a memorable password")

	return cmd
}

func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func versionTemplate() string {
	return fmt.Sprintf("jtgpwgen %s (commit %s, built %s by %s)\n", Version, Commit, Date, BuiltBy)
}
