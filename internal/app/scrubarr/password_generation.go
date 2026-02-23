package scrubarr

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"syscall"

	"github.com/almanac1631/scrubarr/internal/app/auth"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func generatePasswordHash(cmd *cobra.Command, args []string) {
	slog.Info("Starting scrubarr password generation")
	slog.Info("This utility will generate a random salt and a password hash from a given password.")
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		slog.Error("Error reading password.", "error", err)
		os.Exit(1)
	}
	fmt.Print("Re-enter Password: ")
	bytePasswordCheck, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		slog.Error("Error reading password.", "error", err)
		os.Exit(1)
	}
	if !bytes.Equal(bytePasswordCheck, bytePassword) {
		slog.Error("Passwords do not match.")
		os.Exit(1)
	}
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		slog.Error("Error generating salt", "error", err)
		os.Exit(1)
	}
	saltString := hex.EncodeToString(salt)
	passwordHash := auth.GenerateHash(bytePassword, salt)
	passwordHashString := hex.EncodeToString(passwordHash)
	slog.Info("Successfully generated password.", "salt", saltString, "passwordHash", passwordHashString)
	fmt.Println("Salt:", saltString)
	fmt.Println("Password hash:", passwordHashString)
}
