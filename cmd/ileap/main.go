package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/way-platform/ileap-go"
	"github.com/way-platform/ileap-go/cmd/ileap/internal/auth"
)

func main() {
	if err := fang.Execute(
		context.Background(),
		newRootCommand(),
		fang.WithColorSchemeFunc(func(c lipgloss.LightDarkFunc) fang.ColorScheme {
			base := c(lipgloss.Black, lipgloss.White)
			baseInverted := c(lipgloss.White, lipgloss.Black)
			return fang.ColorScheme{
				Base:         base,
				Title:        base,
				Description:  base,
				Comment:      base,
				Flag:         base,
				FlagDefault:  base,
				Command:      base,
				QuotedString: base,
				Argument:     base,
				Help:         base,
				Dash:         base,
				ErrorHeader:  [2]color.Color{baseInverted, base},
				ErrorDetails: base,
			}
		}),
	); err != nil {
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ileap",
		Short: "iLEAP CLI",
	}
	cmd.AddGroup(&cobra.Group{
		ID:    "pcf",
		Title: "Product Carbon Footprints",
	})
	cmd.AddCommand(newGetFootprintCommand())
	cmd.AddCommand(newListFootprintsCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "tad",
		Title: "Transport Activity Data",
	})
	cmd.AddCommand(newListTADsCommand())
	cmd.AddGroup(&cobra.Group{
		ID:    "auth",
		Title: "Authentication",
	})
	authCmd := auth.NewCommand()
	authCmd.GroupID = "auth"
	cmd.AddCommand(authCmd)
	cmd.AddGroup(&cobra.Group{
		ID:    "utils",
		Title: "Utils",
	})
	cmd.SetHelpCommandGroupID("utils")
	cmd.SetCompletionCommandGroupID("utils")
	return cmd
}

func newGetFootprintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "footprint",
		Short: "Get a product carbon footprint",
		GroupID: "pcf",
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
		Use:   "footprints",
		Short: "List product carbon footprints",
		GroupID: "pcf",
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

func newListTADsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tad",
		Short: "List transport activity data (TAD)",
		GroupID: "tad",
	}
	limit := cmd.Flags().Int("limit", 100, "max TADs queried")
	cmd.RunE = func(cmd *cobra.Command, _ []string) error {
		client, err := auth.NewClient()
		if err != nil {
			return err
		}
		response, err := client.ListTADs(cmd.Context(), &ileap.ListTADsRequest{
			Limit: *limit,
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
	fmt.Println(string(data))
	return nil
}
