package helpers

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/debu99/cicd-runner/commands/helpers/archive"
	"github.com/debu99/cicd-runner/common"
	"github.com/debu99/cicd-runner/log"
	"github.com/debu99/cicd-runner/network"
)

//nolint:lll
type ArtifactsDownloaderCommand struct {
	common.JobCredentials
	retryHelper
	network common.Network

	DirectDownload bool `long:"direct-download" env:"FF_USE_DIRECT_DOWNLOAD" description:"Support direct download for data stored externally to GitLab"`
}

func (c *ArtifactsDownloaderCommand) directDownloadFlag(retry int) *bool {
	// We want to send `?direct_download=true`
	// Use direct download only on a first attempt
	if c.DirectDownload && retry == 0 {
		return &c.DirectDownload
	}

	// We don't want to send `?direct_download=false`
	return nil
}

func (c *ArtifactsDownloaderCommand) download(file string, retry int) error {
	switch c.network.DownloadArtifacts(c.JobCredentials, file, c.directDownloadFlag(retry)) {
	case common.DownloadSucceeded:
		return nil
	case common.DownloadNotFound:
		return os.ErrNotExist
	case common.DownloadForbidden:
		return os.ErrPermission
	case common.DownloadFailed:
		return retryableErr{err: os.ErrInvalid}
	default:
		return os.ErrInvalid
	}
}

func (c *ArtifactsDownloaderCommand) Execute(cliContext *cli.Context) {
	log.SetRunnerFormatter()

	wd, err := os.Getwd()
	if err != nil {
		logrus.Fatalln("Unable to get working directory")
	}

	if c.URL == "" || c.Token == "" {
		logrus.Fatalln("Missing runner credentials")
	}
	if c.ID <= 0 {
		logrus.Fatalln("Missing build ID")
	}

	// Create temporary file
	file, err := ioutil.TempFile("", "artifacts")
	if err != nil {
		logrus.Fatalln(err)
	}
	_ = file.Close()
	defer func() { _ = os.Remove(file.Name()) }()

	// Download artifacts file
	err = c.doRetry(func(retry int) error {
		return c.download(file.Name(), retry)
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	f, size, err := openZip(file.Name())
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	extractor, err := archive.NewExtractor(archive.Zip, f, size, wd)
	if err != nil {
		logrus.Fatalln(err)
	}

	// Extract artifacts file
	err = extractor.Extract(context.Background())
	if err != nil {
		logrus.Fatalln(err)
	}
}

func openZip(filename string) (*os.File, int64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, 0, err
	}

	return f, fi.Size(), nil
}

func init() {
	common.RegisterCommand2(
		"artifacts-downloader",
		"download and extract build artifacts (internal)",
		&ArtifactsDownloaderCommand{
			network: network.NewGitLabClient(),
			retryHelper: retryHelper{
				Retry:     2,
				RetryTime: time.Second,
			},
		},
	)
}
