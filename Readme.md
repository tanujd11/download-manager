# download-manager

download-manager is a concurrent downloader which uses go concurrency to features to concurrently download huge files to make it fast.

## Usage

To install the download-manager:

```go install github.com/tanujd11/download-manager@latest```

To download files from download-manager use the following command to download and see the progress bar:

```download-manager download --fileUrl $FILEURL --numConcParts $NUMCONCPARTS --workers $WORKERS --output $OUTPUT ```

## Tools

download-manager uses:
- [cobra](https://github.com/spf13/cobra):   for CLI
- [progress-bar](https://github.com/schollz/progressbar): for interactive progress bar

## Tests

To run tests first populate the testdata using:

```dd if=/dev/random of=`pwd`/pkg/downloader/testdata/1gig.bin bs=1G count=1```