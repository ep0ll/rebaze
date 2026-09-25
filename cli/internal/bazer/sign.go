package bazer

import (
	"fmt"

	"github.com/ep0ll/rebaze/internal/sign"
	"github.com/spf13/cobra"
)

var (
	signKeygen     bool
	signVerify     bool
	signPriv       string
	signPub        string
	signOut        string
)

var signCmd = &cobra.Command{
	Use:   "sign [file]",
	Short: "Sign or verify a mutation plan with an ed25519 key",
	Long: `Generate a key pair, sign a plan file, or verify an existing signature.

  rebaze sign --keygen --private ed25519.key --public ed25519.pub
  rebaze sign plan.json --private ed25519.key --out plan.json.sig
  rebaze sign plan.json --verify --public ed25519.pub --out plan.json.sig`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSign,
}

func init() {
	signCmd.Flags().BoolVar(&signKeygen, "keygen", false, "Generate an ed25519 key pair")
	signCmd.Flags().BoolVar(&signVerify, "verify", false, "Verify a signature envelope")
	signCmd.Flags().StringVar(&signPriv, "private", "ed25519.key", "Private key path")
	signCmd.Flags().StringVar(&signPub, "public", "ed25519.pub", "Public key path")
	signCmd.Flags().StringVar(&signOut, "out", "", "Signature envelope path")
}

func runSign(cmd *cobra.Command, args []string) error {
	if signKeygen {
		if err := sign.GenerateKeyPair(signPriv, signPub); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "wrote %s and %s\n", signPriv, signPub)
		return nil
	}
	if len(args) != 1 {
		return fmt.Errorf("file argument is required unless --keygen is set")
	}
	file := args[0]
	if signOut == "" {
		signOut = file + ".sig"
	}
	if signVerify {
		if err := sign.VerifyFile(file, signPub, signOut); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "signature ok")
		return nil
	}
	env, err := sign.SignFile(file, signPriv, signOut)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "signed %s -> %s (%s)\n", file, signOut, env.Algorithm)
	return nil
}

var sign = signCmd
