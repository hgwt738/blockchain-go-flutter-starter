package main

import (
	"image"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-flutter-desktop/go-flutter"
	"github.com/ground-x/blockchain-go-flutter-starter/go/db"
	"github.com/ground-x/blockchain-go-flutter-starter/go/models/settings"
	"github.com/ground-x/blockchain-go-flutter-starter/go/tlog"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// vmArguments may be set by hover at compile-time
var vmArguments string

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	initDatabase()
	tlog.ReplaceLogger(logger)

	tlog.Info("Starting blockchain starter...")
	// DO NOT EDIT, add options in options.go
	mainOptions := []flutter.Option{
		flutter.OptionVMArguments(strings.Split(vmArguments, ";")),
		flutter.WindowIcon(iconProvider),
	}
	err := flutter.Run(append(options, mainOptions...)...)
	if err != nil {
		tlog.Error(err.Error())
		os.Exit(1)
	}
}

func iconProvider() ([]image.Image, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve executable path")
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to eval symlinks for executable path")
	}
	imgFile, err := os.Open(filepath.Join(filepath.Dir(execPath), "assets", "icon.png"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open assets/icon.png")
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode image")
	}
	return []image.Image{img}, nil
}

func initDatabase() {
	d, err := db.InitDatabase()
	if err != nil {
		panic(err)
	}
	if err := d.AutoMigrate(&settings.Settings{}); err != nil {
		panic(err)
	}
}
