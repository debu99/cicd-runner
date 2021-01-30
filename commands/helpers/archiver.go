package helpers

import (
	"os"

	"github.com/debu99/cicd-runner/commands/helpers/archive"
	"github.com/debu99/cicd-runner/commands/helpers/archive/fastzip"
	"github.com/debu99/cicd-runner/helpers/featureflags"

	// auto-register default archivers/extractors
	_ "github.com/debu99/cicd-runner/commands/helpers/archive/gziplegacy"
	_ "github.com/debu99/cicd-runner/commands/helpers/archive/raw"
	_ "github.com/debu99/cicd-runner/commands/helpers/archive/ziplegacy"

	"github.com/sirupsen/logrus"
)

func init() {
	// enable fastzip archiver/extractor
	if on, _ := featureflags.IsOn(os.Getenv(featureflags.UseFastzip)); on {
		archive.Register(archive.Zip, fastzip.NewArchiver, fastzip.NewExtractor)
	}
}

// getCompressionLevel converts the compression level name to compression level type
// https://docs.gitlab.com/ee/ci/runners/README.html#artifact-and-cache-settings
func getCompressionLevel(name string) archive.CompressionLevel {
	switch name {
	case "fastest":
		return archive.FastestCompression
	case "fast":
		return archive.FastCompression
	case "slow":
		return archive.SlowCompression
	case "slowest":
		return archive.SlowestCompression
	case "default", "":
		return archive.DefaultCompression
	}

	logrus.Warningf("compression level %q is invalid, falling back to default", name)

	return archive.DefaultCompression
}
