package igdl

import (
	"fmt"
	"log"
	"strconv"
	"time"
	//"github.com/siongui/instago"
)

type TrayInfo struct {
	Id       int64
	Username string
	Layer    int
	//Tray     instago.IGReelTray
}

func IsTrayInfoInQueue(queue []TrayInfo, ti TrayInfo) bool {
	for _, t := range queue {
		if t.Id == ti.Id {
			return true
		}
	}
	return false
}

func GetTrayInfoFromQueue(queue []TrayInfo, id int64) (ti TrayInfo, ok bool) {
	for _, t := range queue {
		if t.Id == id {
			ti = t
			ok = true
			return
		}
	}
	return
}

func (m *IGDownloadManager) DownloadTrayInfos(tis []TrayInfo, c chan TrayInfo, interval int, getTime map[string]time.Time, verbose bool) {
	downloadIds := []string{}
	for _, ti := range tis {
		id := strconv.FormatInt(ti.Id, 10)
		username := ti.Username
		if verbose {
			PrintUsernameIdMsg(username, id, " to be batch downloaded")
		}
		downloadIds = append(downloadIds, id)
	}

	// wait at least *interval* seconds until next private API access
	d := time.Now().Sub(getTime["last"])
	for d < time.Duration(interval)*time.Second {
		time.Sleep(time.Duration(interval)*time.Second - d)
		d = time.Now().Sub(getTime["last"])
	}
	// get stories of multiple users at one API access
	trays, err := m.GetMultipleReelsMedia(downloadIds)
	getTime["last"] = time.Now()
	if err != nil {
		log.Println(err)
		// sent back to channel to re-download
		for _, ti := range tis {
			c <- ti
		}
		return
	}

	// we have story trays of multiple users now. start to download
	for _, tray := range trays {
		ti, ok := GetTrayInfoFromQueue(tis, tray.Id)
		if !ok {
			log.Println("cannot find", tray.Id, "in queue (impossible to happen)")
			continue
		}

		username := tray.User.GetUsername()
		id := tray.User.GetUserId()
		if err != nil {
			PrintUsernameIdMsg(username, id, "downloading ...")
			return
		}
		for _, item := range tray.Items {
			_, err = getStoryItem(item, username)
			if err != nil {
				PrintUsernameIdMsg(username, id, err)
				return
			}

			if ti.Layer-1 < 1 {
				continue
			}
			for _, rm := range item.ReelMentions {
				rmti := TrayInfo{}
				rmti.Layer = ti.Layer - 1
				rmti.Id = rm.User.Pk
				rmti.Username = rm.GetUsername()
				c <- rmti
				if verbose {
					PrintUsernameIdMsg(rmti.Username, rmti.Id, "sent to channel (reel mention)")
				}
			}
		}
	}
}

func (m *IGDownloadManager) TrayDownloader(c chan TrayInfo, interval int, getTime map[string]time.Time, verbose bool) {
	queue := []TrayInfo{}
	for {
		select {
		case ti := <-c:
			// append to queue if not exist
			id := ti.Id
			username := ti.Username
			if verbose {
				UsernameIdColorPrint(username, id)
				log.Println("legnth of channel:", len(c))
			}
			if IsTrayInfoInQueue(queue, ti) {
				if verbose {
					PrintUsernameIdMsg(username, id, "exist. ignore.")
				}
			} else {
				queue = append(queue, ti)
				if verbose {
					PrintUsernameIdMsg(username, id, "appended")
				}
			}
		default:
			if len(queue) > 0 {
				tis := []TrayInfo{}
				if len(queue) > 20 {
					tis = queue[0:20]
					queue = queue[20:]
				} else {
					tis = queue
					queue = []TrayInfo{}
				}

				if len(tis) > 0 {
					m.DownloadTrayInfos(tis, c, interval, getTime, verbose)
				}
			}

			restInterval := 1
			if verbose {
				log.Println("current queue length: ", len(queue))
				PrintMsgSleep(restInterval, "TrayDownloader: ")
			} else {
				SleepSecond(restInterval)
			}
		}
	}
}

func (m *IGDownloadManager) AccessReelsTrayOnce2(c chan TrayInfo, ignoreMuted, verbose bool) (err error) {
	rt, err := m.GetReelsTray()
	if err != nil {
		log.Println(err)
		return
	}

	go PrintLiveBroadcasts(rt.Broadcasts)

	for index, tray := range rt.Trays {
		fmt.Print(index, ":")

		username := tray.GetUsername()
		id := tray.Id
		//items := tray.GetItems()

		if ignoreMuted && tray.Muted {
			if verbose {
				PrintUsernameIdMsg(username, id, " is muted && ignoreMuted set. no download")
			}
			continue
		}

		if IsLatestReelMediaDownloaded(username, tray.LatestReelMedia) {
			if verbose {
				PrintUsernameIdMsg(username, id, " all downloaded")
			}
			continue
		}

		if tray.HasBestiesMedia {
			PrintUsernameIdMsg(username, id, "has close friend (besties) story item(s)")
		}

		if verbose {
			UsernameIdColorPrint(username, id)
			fmt.Println(" has undownloaded items")
		}

		ti := TrayInfo{}
		ti.Layer = 2 // 2: also download its reel mentions in story item
		ti.Username = username
		ti.Id = id
		c <- ti
		/*
			items := tray.GetItems()
			if len(items) > 0 {
				for _, item := range items {
					_, err = GetStoryItem(item, username)
					if err != nil {
						PrintUsernameIdMsg(username, id, err)
					}
				}
			} else {
				emptyids = append(emptyids, strconv.FormatInt(id, 10))
			}

			if len(emptyids) > 20 {
				break
			}
		*/
	}

	return
}

func (m *IGDownloadManager) DownloadStoryForever2(interval1, interval2 int, ignoreMuted, verbose bool) {
	// usually there are at most 150 trays in reels_tray.
	// double the buffer to 300. 160 or 200 may be ok as well.
	c := make(chan TrayInfo, 300)

	getTime := make(map[string]time.Time)
	getTime["last"] = time.Now()
	go m.TrayDownloader(c, interval2, getTime, verbose)

	for {
		err := m.AccessReelsTrayOnce2(c, ignoreMuted, verbose)
		if err != nil {
			log.Println(err)
		}
		PrintMsgSleep(interval1, "DownloadStoryForever2: ")
	}
}
