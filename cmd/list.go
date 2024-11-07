package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/Shackelford-Arden/hctx/config"
	"github.com/Shackelford-Arden/hctx/types"

	"github.com/urfave/cli/v2"
)

func basicList(currStack string, stacks []config.Stack) {
	fmt.Println("Stacks:")
	for _, stack := range stacks {
		var indicator string
		if currStack != "" && (stack.Name == currStack || stack.Alias == currStack) {
			indicator = "*"
		}
		fmt.Printf("  %s %s\n", stack.Name, indicator)
	}
}

func detailedList(currStack string, stacks []config.Stack) {

	// Colors
	baseFg := lipgloss.Color("#FFFFFF")
	nomadColor := lipgloss.Color("#00CA8E")
	consulColor := lipgloss.Color("#DC477D")
	vaultColor := lipgloss.Color("#FFCF25")
	activeColor := lipgloss.Color("#1cdd24")

	// Base styles
	lg := lipgloss.NewRenderer(os.Stdout)
	baseStyle := lg.NewStyle().Padding(0, 1)
	borderStyle := lipgloss.NewStyle().Foreground(baseFg)
	headerStyle := baseStyle.Foreground(baseFg).Bold(true).Align(lipgloss.Center)
	cellStyle := baseStyle.Foreground(baseFg).Width(25)
	activeStyle := baseStyle.Foreground(activeColor)

	headers := []string{"Stack", "Alias", "Nomad", "Consul", "Vault"}
	stackData := make([][]string, len(stacks))
	var activeIndex int

	for i, stack := range stacks {

		nomadStr := ""
		consulStr := ""
		vaultStr := ""
		stackName := stack.Name

		if stack.Name == currStack {
			// Adding plus 1, because in the data, headers will always be 0
			// so all the data we insert will always be 1 higher.
			activeIndex = i + 1
			stackName = fmt.Sprintf("%s  ✓", stackName)
		}

		if stack.Nomad != nil && stack.Nomad.Address != "" {
			nomadStr = stack.Nomad.Address
		}

		if stack.Nomad != nil && stack.Nomad.Namespace != "" {
			nomadStr = fmt.Sprintf("%s (%s)", nomadStr, stack.Nomad.Namespace)
		}

		if stack.Consul != nil && stack.Consul.Address != "" {
			consulStr = stack.Consul.Address
		}

		if stack.Consul != nil && stack.Consul.Namespace != "" {
			consulStr = fmt.Sprintf("%s (%s)", consulStr, stack.Consul.Namespace)
		}

		if stack.Vault != nil && stack.Vault.Address != "" {
			vaultStr = stack.Vault.Address
		}

		if stack.Vault != nil && stack.Vault.Namespace != "" {
			vaultStr = fmt.Sprintf("%s (%s)", vaultStr, stack.Vault.Namespace)
		}

		stackData[i] = []string{
			stackName,
			stack.Alias,
			nomadStr,
			consulStr,
			vaultStr,
		}
	}

	configTable := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(borderStyle).
		Headers(headers...).
		Rows(stackData...).
		StyleFunc(func(row, col int) lipgloss.Style {
			// Column Reference
			// Nomad: 2
			// Consul: 3
			// Vault: 4

			// Nomad Column Header
			if row == 0 && col == 2 {
				return cellStyle.Foreground(nomadColor).Bold(true).Align(lipgloss.Center)
			}

			// Consul Column Header
			if row == 0 && col == 3 {
				return cellStyle.Foreground(consulColor).Bold(true).Align(lipgloss.Center)
			}

			// Nomad Column Header
			if row == 0 && col == 4 {
				return cellStyle.Foreground(vaultColor).Bold(true).Align(lipgloss.Center)
			}

			// Handle all other headers
			if row == 0 {
				return headerStyle
			}

			if row == activeIndex && col == 0 {
				// Add any formatting for the current stack here.
				return activeStyle
			}

			// All other unspecified row/column data
			return baseStyle.Foreground(baseFg)
		})

	fmt.Println(configTable)
}

func List(ctx *cli.Context) error {

	if len(AppConfig.Stacks) == 0 {
		fmt.Fprintf(ctx.App.Writer, "No stacks!\n")
		return nil
	}

	currStack := os.Getenv(types.StackNameEnv)
	detailed := ctx.Bool("detailed")

	if detailed {
		detailedList(currStack, AppConfig.Stacks)
		return nil
	}

	basicList(currStack, AppConfig.Stacks)

	return nil
}
