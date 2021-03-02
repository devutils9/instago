package igdl

import (
	"fmt"
	"log"
	"strconv"
)

func (m *IGDownloadManager) AccessReelsTrayOnce2Chan(cPublicUser, cPrivateUser chan TrayInfo, ignoreMuted, verbose bool) (err error) {
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
			PrintUsernameIdMsg(username, id, " has undownloaded items")
			UsernameIdColorPrint(username, id)
		}

		// 2: also download reel mentions in story item
		if tray.User.IsPrivate {
			cPrivateUser <- setupTrayInfo(id, username, 2, tray.User.IsPrivate)
		} else {
			cPublicUser <- setupTrayInfo(id, username, 2, tray.User.IsPrivate)
		}
	}

	return
}

func (m *IGDownloadManager) TwoAccountDownloadStoryForever(interval1 int, interval2, interval3 int64, ignoreMuted, verbose bool) {
	if !m.IsCleanAccountSet() {
		fmt.Println("clean account not set. exit")
		return
	}

	// usually there are at most 150 trays in reels_tray.
	// double the buffer to 300. 160 or 200 may be ok as well.
	cPublicUser := make(chan TrayInfo, 300)
	cPrivateUser := make(chan TrayInfo, 300)

	go m.GetCleanAccountManager().TrayDownloader(cPublicUser, NewTimeLimiter(interval2), true, verbose)
	go m.PrivateTrayDownloaderViaReelMediaAPI(cPublicUser, cPrivateUser, NewTimeLimiter(interval3), verbose)

	for {
		err := m.AccessReelsTrayOnce2Chan(cPublicUser, cPrivateUser, ignoreMuted, verbose)
		if err != nil {
			log.Println(err)
		}
		PrintMsgSleep(interval1, "TwoAccountDownloadStoryForever: ")
	}
}

func (m *IGDownloadManager) PrivateTrayDownloaderViaReelMediaAPI(cPublicUser, cPrivateUser chan TrayInfo, tl *TimeLimiter, verbose bool) {
	//maxReelsMediaIds := 20
	queue := []TrayInfo{}
	for {
		select {
		case ti := <-cPrivateUser:
			// append to queue if not exist
			id := ti.Id
			username := ti.Username
			if IsTrayInfoInQueue(queue, ti) {
				if verbose {
					PrintUsernameIdMsg(username, id, "exist. ignore.", "len(cPrivateUser):", len(cPrivateUser), "len(queue):", len(queue))
				}
			} else {
				queue = append(queue, ti)
				if verbose {
					PrintUsernameIdMsg(username, id, "appended.", "len(cPrivateUser):", len(cPrivateUser), "len(queue):", len(queue))
				}
			}
		default:
			if len(queue) > 0 {
				ti := queue[0]
				queue = queue[1:]

				// wait at least *interval* seconds until next private API access (prevent http 429)
				tl.WaitAtLeastIntervalAfterLastTime()
				tray, err := m.GetUserReelMedia(strconv.FormatInt(ti.Id, 10))
				tl.SetLastTimeToNow()
				if err != nil {
					log.Println(err)
					cPrivateUser <- ti
				} else {
					id := tray.User.GetUserId()
					username := tray.User.GetUsername()
					for _, item := range tray.GetItems() {
						_, err := getStoryItem(item, tray.GetUsername())
						if err != nil {
							PrintUsernameIdMsg(username, id, err)
							continue
						}

						if ti.Layer-1 < 1 {
							continue
						}
						for _, rm := range item.ReelMentions {
							if !rm.User.IsPrivate {
								cPublicUser <- setupTrayInfo(rm.User.Pk, rm.GetUsername(), ti.Layer-1, rm.User.IsPrivate)
								if verbose {
									PrintUsernameIdMsg(rm.GetUsername(), rm.User.Pk, "sent to public channel (reel mention)")
								}
							}
						}
					}
				}
			}

			restInterval := 1
			if verbose {
				log.Println("current private user queue length: ", len(queue))
				PrintMsgSleep(restInterval, "PrivateTrayDownloader: ")
			} else {
				SleepSecond(restInterval)
			}
		}
	}
}
