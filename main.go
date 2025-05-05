package main

import (
	"log/slog"
)

func main() {
	logger := slog.Default()

	if err := Render(); err != nil {
		// TODO: display error properly formatted
		logger.Error(err.Error())
	}
}
