package dashboard

import (
	"context"
	"fmt"
	"io"

	tea "charm.land/bubbletea/v2"
)

type RunOptions struct {
	Model ModelOptions
}

func Run(ctx context.Context, input io.Reader, output io.Writer, source Source, options RunOptions) error {
	if source == nil {
		return fmt.Errorf("dashboard source is unavailable")
	}
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	snapshot, err := source.Initial(runCtx)
	if err != nil {
		return fmt.Errorf("collect initial dashboard snapshot: %w", err)
	}
	model := NewModel(runCtx, source, snapshot, options.Model)
	program := tea.NewProgram(
		model,
		tea.WithContext(runCtx),
		tea.WithInput(input),
		tea.WithOutput(output),
	)
	if _, err := program.Run(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return fmt.Errorf("run terminal dashboard: %w", err)
	}
	return nil
}
