package main

import (
	"context"
	"github.com/nkorange/stock-memo/internal"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os"
)

var cfg *config

type config struct {
	addr string
}

func init() {
	cfg = &config{}
	initCmd(rootCmd)
}

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:  "stock profit analyzer",
	Args: cobra.NoArgs,
	Long: `Stock Profit Analyzer`,
	RunE: runE,
}

// Execute runs the root command
func Execute() {
	// Outputs cmd.Print to stdout.
	rootCmd.SetOut(os.Stdout)
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErr(err.Error() + "\n")
		os.Exit(1)
	}
}

func runE(cmd *cobra.Command, args []string) error {
	server, err := internal.NewServer(cfg.addr)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg, cctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return server.Run(cctx)
	})
	return wg.Wait()
}

func initCmd(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cfg.addr, "addr", "0.0.0.0:10086", "http address")
}
