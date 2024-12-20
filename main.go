package main

import "log/slog"

func main() {
	logger := slog.Default()
	logger.Info("hello world", "name", "joe")
}
