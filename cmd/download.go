/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/tanujd11/download-manager/pkg/downloader"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "This command will call download on your file",
	Long: ` This command will take multiple arguments like
	fileUrl, NumConcParts... and will download the file 
	 `,
	Run: func(cmd *cobra.Command, args []string) {

		numConcParts, _ := strconv.Atoi(cmd.Flag("numConcParts").Value.String())
		workers, _ := strconv.Atoi(cmd.Flag("workers").Value.String())
		fileUrl := cmd.Flag("fileUrl").Value.String()

		opts := downloader.DownloadOptions{
			DownloadPath: cmd.Flag("output").Value.String(),
			NumConcParts: numConcParts,
			Workers:      workers,
		}

		downloadClient := downloader.Downloader{
			DownloadOptions: opts,
			Progress:        0,
			Mutex:           &sync.Mutex{},
		}

		bar := progressbar.Default(100)
		go func() {
			for {
				bar.Set(downloadClient.GetDownloadProgress(fileUrl))
				time.Sleep(3 * time.Millisecond)
			}
		}()
		downloadPath, err := downloadClient.Download(fileUrl)
		if err != nil {
			fmt.Println("error downloading the file: " + err.Error())
		} else {
			fmt.Println("file downloaded at: " + downloadPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.PersistentFlags().String("fileUrl", "", "http path of the file to download")
	downloadCmd.PersistentFlags().String("output", "", "the path on your machine where the file will be downloaded")
	downloadCmd.PersistentFlags().Int("numConcParts", 10, "number of concurrent paths to be donwloaded")
	downloadCmd.PersistentFlags().Int("workers", 5, "size of the worker pool")

	if err := downloadCmd.MarkPersistentFlagRequired("fileUrl"); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := downloadCmd.MarkPersistentFlagRequired("output"); err != nil {
		fmt.Println(err.Error())
		return
	}
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
