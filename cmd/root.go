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

var flags struct {
	length       int
	special      string
	digits       bool
	noSpecial    bool
	noDigits     bool
	memo         bool
	digitsSet    bool
	specialSet   bool
	noDigitsSet  bool
	noSpecialSet bool
}

var rootCmd = &cobra.Command{
	Use:     "jtgpwgen",
	Short:   "Генератор паролей с настраиваемыми классами символов",
	Version: Version,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		flags.digitsSet = cmd.Flags().Changed("digits")
		flags.specialSet = cmd.Flags().Changed("special")
		flags.noDigitsSet = cmd.Flags().Changed("no-digits")
		flags.noSpecialSet = cmd.Flags().Changed("no-special")

		if flags.digitsSet && flags.noDigitsSet {
			return passgen.ErrConflictDigits
		}
		if flags.specialSet && flags.noSpecialSet {
			return passgen.ErrConflictSpecial
		}

		opts := passgen.DefaultOptions()
		opts.Length = flags.length
		opts.Memo = flags.memo
		if flags.noDigitsSet {
			opts.UseDigits = false
		}
		if flags.noSpecialSet {
			opts.UseSpecial = false
		}
		if flags.specialSet {
			opts.UseSpecial = true
			opts.ExtraSpecial = flags.special
		}
		if flags.memo && (flags.specialSet || flags.noSpecialSet || flags.digitsSet || flags.noDigitsSet) {
			fmt.Fprintln(os.Stderr, "warning: --memo ignores character class flags")
		}

		pw, err := passgen.Generate(opts)
		if err != nil {
			if errors.Is(err, passgen.ErrLengthOutOfRange) || errors.Is(err, passgen.ErrNoCharClasses) {
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(versionTemplate())

	f := rootCmd.Flags()
	f.IntVarP(&flags.length, "length", "l", passgen.DefaultLength, "password length")
	f.StringVarP(&flags.special, "special", "s", "", "extra special characters to include (added to default '@')")
	f.BoolVarP(&flags.digits, "digits", "d", false, "explicitly enable digits (default: enabled)")
	f.BoolVar(&flags.noSpecial, "no-special", false, "disable special characters entirely")
	f.BoolVar(&flags.noDigits, "no-digits", false, "disable digits")
	f.BoolVarP(&flags.memo, "memo", "m", false, "generate a memorable password")
}

func versionTemplate() string {
	return fmt.Sprintf("jtgpwgen %s (commit %s, built %s by %s)\n", Version, Commit, Date, BuiltBy)
}
