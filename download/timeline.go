package igdl

import (
	"fmt"
	"github.com/siongui/instago"
	"os"
	"time"
)

func printDownloadInfo(item instago.IGItem, username string, url, filepath string) {
	fmt.Print("username: ")
	cc.Println(username)
	fmt.Print("time: ")
	cc.Println(formatTimestamp(item.GetTimestamp()))
	fmt.Print("post url: ")
	cc.Println(item.GetPostUrl())

	fmt.Print("Download ")
	rc.Print(url)
	fmt.Print(" to ")
	cc.Println(filepath)
}

// getTimelineItems is obsoleted. Use getPostItem instead.
func getTimelineItems(items []instago.IGItem) {
	for _, item := range items {
		if !item.IsRegularMedia() {
			continue
		}

		urls, err := item.GetMediaUrls()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for index, url := range urls {
			filepath := getPostFilePath(
				item.GetUsername(),
				item.GetUserId(),
				item.GetPostCode(),
				url,
				item.GetTimestamp())
			if index > 0 {
				filepath = appendIndexToFilename(filepath, index)
			}

			CreateFilepathDirIfNotExist(filepath)
			// check if file exist
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				// file not exists
				printDownloadInfo(item, item.GetUsername(), url, filepath)
				err = Wget(url, filepath)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

// download timeline until page n
func (m *IGDownloadManager) DownloadTimeline(n int) {
	sleepInterval := 12 // seconds

	for {
		items, err := m.apimgr.GetTimelineUntilPageN(n)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, item := range items {
				if !item.IsRegularMedia() {
					continue
				}
				m.getPostItem(item)
			}
		}

		// sleep for a while
		fmt.Println("=========================")
		rc.Print(time.Now().Format(time.RFC3339))
		fmt.Print(": sleep ")
		cc.Print(sleepInterval)
		fmt.Println(" second")
		fmt.Println("=========================")
		time.Sleep(time.Duration(sleepInterval) * time.Second)
	}
}
