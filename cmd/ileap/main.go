package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/cmd/ileap/internal/auth"
)

func main() {
	if err := newRootCommand().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ileap",
		Short: "iLEAP CLI",
	}
	cmd.AddCommand(auth.NewCommand())
	cmd.AddCommand(newGetFootprintCommand())
	cmd.AddCommand(newListFootprintsCommand())
	return cmd
}

func newGetFootprintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-footprint",
		Short: "Get a product carbon footprint",
		Args:  cobra.ExactArgs(1),
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		footprint, err := client.GetFootprint(cmd.Context(), &ileap.GetFootprintRequest{
			ID: args[0],
		})
		if err != nil {
			return err
		}
		printJSON(cmd, footprint)
		return nil
	}
	return cmd
}

func newListFootprintsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-footprints",
		Short: "List product carbon footprints",
	}
	limit := cmd.Flags().Int("limit", 100, "max footprints queried")
	filter := cmd.Flags().String("filter", "", "filter footprints by OData filter")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.ListFootprints(cmd.Context(), &ileap.ListFootprintsRequest{
			Limit:  *limit,
			Filter: *filter,
		})
		if err != nil {
			return err
		}
		printJSON(cmd, response)
		return nil
	}
	return cmd
}

func printJSON(cmd *cobra.Command, msg any) error {
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		return err
	}
	cmd.Println(string(data))
	return nil
}
