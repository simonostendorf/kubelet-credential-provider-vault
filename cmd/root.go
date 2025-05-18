package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/communicationInterface"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/config"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/credentialFetcher"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/logger"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/provider"
	"github.com/simonostendorf/kubelet-credential-provider-vault/internal/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "0.0.1"

var (
	// flag variables
	configFile string

	// wait group for graceful shutdown
	wg sync.WaitGroup

	// logger
	log logger.Logger

	// root command
	rootCmd = &cobra.Command{
		Use:     "kubelet-credential-provider-vault",
		Version: version,
		Short:   "A Kubernetes Kubelet Image Credential Provider for HashiCorp Vault",
		Long:    "A Kubernetes Kubelet Image Credential Provider for HashiCorp Vault",
		Run:     executeRootCmd,
	}
)

func executeRootCmd(cmd *cobra.Command, args []string) {
	// create shutdown context handler
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer cancel()

	// setup wait group for graceful shutdown
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		handleShutdown(ctx, shutdownReasonSignal)
	}()

	// setup initial logger (to log errors before config is loaded and logger is initialized)
	initialLogger, err := logger.NewFileLogger(true, logger.DefaultLogFile, "error")
	log = initialLogger
	if err != nil {
		log.Log(ctx, slog.LevelError, "Failed to initialize initial logger", "error", err)
		handleShutdown(ctx, shutdownReasonError)
		return
	}

	// load config
	cfg, err := config.New(ctx, log, configFile)
	if err != nil {
		log.Log(ctx, slog.LevelError, "Failed to load configuration", "error", err)
		handleShutdown(ctx, shutdownReasonError)
		return
	}

	// setup logger
	newLogger, err := logger.NewFileLogger(cfg.Log.Enabled, cfg.Log.File, cfg.Log.Level)
	log = newLogger
	if err != nil {
		log.Log(ctx, slog.LevelError, "Failed to initialize logger", "error", err)
		handleShutdown(ctx, shutdownReasonError)
		return
	}

	// log startup information after logger is initialized
	log.Log(ctx, slog.LevelDebug, "Loaded configuration", "config", *cfg)
	log.Log(ctx, slog.LevelDebug, "Initialized logger", "file", cfg.Log.File, "level", cfg.Log.Level)
	log.Log(ctx, slog.LevelInfo, "Starting kubelet-credential-provider-vault", "version", version)

	// setup communication interface (stdio)
	communicationInterface := communicationInterface.NewStdIOCommunicationInterface()
	log.Log(ctx, slog.LevelDebug, "Initialized communication interface", "interface", "StdIO")

	// setup credential fetcher (vault)
	credentialFetcher := credentialFetcher.NewVaultCredentialFetcher(vault.NewHashicorpClientBuilder(), &cfg.Vault)
	log.Log(ctx, slog.LevelDebug, "Initialized credential fetcher", "fetcher", "Vault")

	// provide credentials to kubelet
	provider := provider.NewKubeletCredentialProvider(communicationInterface, credentialFetcher)
	err = provider.Run(ctx, log)
	if err != nil {
		log.Log(ctx, slog.LevelError, "Failed to run provider", "error", err)
		handleShutdown(ctx, shutdownReasonError)
		return
	}

	// do not wait for shutdown with wg.Wait() because credential provider should exit after responding to the request
	// instead, directly perform shutdown logic because wait group will not be reached
	handleShutdown(ctx, shutdownReasonFinished)
}

type shutdownReason string

const (
	shutdownReasonFinished shutdownReason = "finished"
	shutdownReasonSignal   shutdownReason = "signal"
	shutdownReasonError    shutdownReason = "error"
)

func handleShutdown(ctx context.Context, shutdownReason shutdownReason) {
	if log != nil {
		switch shutdownReason {
		case shutdownReasonFinished:
			log.Log(ctx, slog.LevelInfo, "Application finished successfully, shutting down...")
		case shutdownReasonSignal:
			log.Log(ctx, slog.LevelInfo, "Received shutdown signal, shutting down...")
		case shutdownReasonError:
			log.Log(ctx, slog.LevelInfo, "An error occurred, shutting down...")
		}
	}

	// close logger
	if log != nil {
		err := log.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close logger: %v", err))
		}
		log = nil
	}
}

func Execute() {
	// cmd as command with Run instead of RunE, so no error is expected
	// nolint:errcheck
	rootCmd.Execute() //gosec:disable G104
}

func init() {
	rootCmd.Flags().StringVar(&configFile, "config", "", "configuration file to use. If not set, the application will look for "+config.DefaultConfigFile)
	// no viper bind for config file because it must be handled before viper

	rootCmd.Flags().String("log-file", logger.DefaultLogFile, "file the logger will write to")
	// nolint:errcheck
	viper.BindPFlag("log.file", rootCmd.Flags().Lookup("log-file")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("log.file", "LOG_FILE") //gosec:disable G104

	rootCmd.Flags().String("log-level", "info", "log level to use. Possible values: debug, info, warn, error")
	// nolint:errcheck
	viper.BindPFlag("log.level", rootCmd.Flags().Lookup("log-level")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("log.level", "LOG_LEVEL") //gosec:disable G104

	rootCmd.Flags().Bool("log-enabled", true, "enable or disable logging")
	// nolint:errcheck
	viper.BindPFlag("log.enabled", rootCmd.Flags().Lookup("log-enabled")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("log.enabled", "LOG_ENABLED") //gosec:disable G104

	rootCmd.Flags().String("vault-addr", "", "address of the Vault server")
	// nolint:errcheck
	viper.BindPFlag("vault.address", rootCmd.Flags().Lookup("vault-addr")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("vault.address", "VAULT_ADDR") //gosec:disable G104

	rootCmd.Flags().Bool("vault-insecure-skip-verify", false, "skip TLS verification of the Vault server")
	// nolint:errcheck
	viper.BindPFlag("vault.insecureSkipVerify", rootCmd.Flags().Lookup("vault-insecure-skip-verify")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("vault.insecureSkipVerify", "VAULT_INSECURE_SKIP_VERIFY") //gosec:disable G104

	rootCmd.Flags().String("vault-auth-method", "kubernetes", "name of the auth method to use. Possible values: kubernetes")
	// nolint:errcheck
	viper.BindPFlag("vault.auth.method", rootCmd.Flags().Lookup("vault-auth-method")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("vault.auth.method", "VAULT_AUTH_METHOD") //gosec:disable G104

	rootCmd.Flags().String("vault-auth-mount", "", "name of the auth mount to use")
	// nolint:errcheck
	viper.BindPFlag("vault.auth.mount", rootCmd.Flags().Lookup("vault-auth-mount")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("vault.auth.mount", "VAULT_AUTH_MOUNT") //gosec:disable G104

	rootCmd.Flags().String("vault-auth-role", "", "name of the auth role to use")
	// nolint:errcheck
	viper.BindPFlag("vault.auth.role", rootCmd.Flags().Lookup("vault-auth-role")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("vault.auth.role", "VAULT_AUTH_ROLE") //gosec:disable G104

	rootCmd.Flags().String("vault-secret-mount", "", "name of the secret mount to use")
	// nolint:errcheck
	viper.BindPFlag("vault.secret.mount", rootCmd.Flags().Lookup("vault-secret-mount")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("vault.secret.mount", "VAULT_SECRET_MOUNT") //gosec:disable G104

	rootCmd.Flags().String("vault-secret-path", "", "path of the secret to use")
	// nolint:errcheck
	viper.BindPFlag("vault.secret.path", rootCmd.Flags().Lookup("vault-secret-path")) //gosec:disable G104
	// nolint:errcheck
	viper.BindEnv("vault.secret.path", "VAULT_SECRET_PATH") //gosec:disable G104
}
