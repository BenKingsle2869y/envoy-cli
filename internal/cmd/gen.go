package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	charsetAlpha   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetDigits  = "0123456789"
	charsetSymbols = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

var genCmd = &cobra.Command{
	Use:   "gen KEY",
	Short: "Generate a random secret value and store it",
	Args:  cobra.ExactArgs(1),
	RunE:  runGen,
}

var (
	genLength  int
	genSymbols bool
	genDigits  bool
	genPrint   bool
)

func init() {
	genCmd.Flags().IntVarP(&genLength, "length", "l", 32, "Length of generated value")
	genCmd.Flags().BoolVarP(&genSymbols, "symbols", "s", false, "Include symbols")
	genCmd.Flags().BoolVarP(&genDigits, "digits", "d", true, "Include digits")
	genCmd.Flags().BoolVarP(&genPrint, "print", "p", false, "Print generated value to stdout")
	rootCmd.AddCommand(genCmd)
}

func runGen(cmd *cobra.Command, args []string) error {
	key := args[0]

	passphrase := resolvePassphrase(cmd)
	if passphrase == "" {
		return fmt.Errorf("passphrase is required")
	}

	value := generateSecret(genLength, genDigits, genSymbols)

	ctx := ActiveContext()
	path := StorePathForContext(ctx)
	st, err := storeLoad(path, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	st.Entries[key] = value
	if err := storeSave(st, path, passphrase); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	if genPrint {
		fmt.Fprintln(cmd.OutOrStdout(), value)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Generated and stored key %q\n", key)
	}
	return nil
}

func generateSecret(length int, digits, symbols bool) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := charsetAlpha
	if digits {
		charset += charsetDigits
	}
	if symbols {
		charset += charsetSymbols
	}
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(charset[rng.Intn(len(charset))])
	}
	return sb.String()
}
