package handlers

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"runtime/debug"

	"github.com/titpetric/etl/model"
)

func Version(ctx context.Context, command *model.Command, _ io.Reader) error {
	// Print the Go runtime version.
	fmt.Printf("Go Runtime Version: %s\n", runtime.Version())

	// Read the build information.
	if info, ok := debug.ReadBuildInfo(); ok {
		fmt.Println("Build Info:")
		// Print main module info.
		fmt.Printf("  Main Module: %s %s\n", info.Main.Path, info.Main.Version)

		// Print any settings (these include VCS details if available).
		for _, setting := range info.Settings {
			fmt.Printf("  %s: %s\n", setting.Key, setting.Value)
		}
	} else {
		fmt.Println("No build info available.")
	}

	return nil
}
