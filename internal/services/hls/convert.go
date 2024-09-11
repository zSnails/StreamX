package hls

import (
	"os"
	"os/exec"
	"path"

	"github.com/zSnails/streamx/pkg/logging"
)

var log = logging.Get().WithField("service", "hls")

func Convert(outdir, hash, input string) error {
	if err := os.MkdirAll(path.Join(outdir, hash), os.ModePerm); err != nil {
		return err
	}
	cmd := exec.Command(
		"ffmpeg",
		"-i",
		input,
		"-c",
		"copy",
		"-f",
		"segment",
		"-segment_time",
		"10",
		"-segment_list",
		path.Join(outdir, hash, "outputlist.m3u8"),
		"-segment_format",
		"mpegts",
		path.Join(outdir, hash, "output%03d.ts"),
	)

	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
