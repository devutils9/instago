package igdl

import (
	"testing"

	"github.com/siongui/instago"
)

func TestGetPostFilePath(t *testing.T) {
	path := getPostFilePath("instagram", "25025320", "Bh7kySfDYq8", "https://instagram.fkhh1-2.fna.fbcdn.net/vp/893534d61bdc5ea6911593d3ee0a1922/5B6363AB/t51.2885-19/s320x320/14719833_310540259320655_1605122788543168512_a.jpg?abc=1", 1520056661)
	if path != "Instagram/instagram/posts/instagram-25025320-post-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.jpg" {
		t.Error(path)
		return
	}
}

func TestGetStoryFilePath(t *testing.T) {
	path := getStoryFilePath("instagram", "25025320", "Bh7kySfDYq8", "123.mp4", 1520056661)
	if path != "Instagram/instagram/stories/instagram-25025320-story-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.mp4" {
		t.Error(path)
		return
	}
}

func TestGetStoryFilePath2(t *testing.T) {
	path := getStoryFilePath2("instagram", "25025320", "Bh7kySfDYq8", "123.mp4", 1520056661, nil)
	if path != "Instagram/instagram/stories/instagram-25025320-story-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.mp4" {
		t.Error(path)
		return
	}

	user1 := instago.IGUser{Pk: 12345, Username: "testuser"}
	user2 := instago.IGUser{Pk: 123456, Username: "testuser111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"}
	user3 := instago.IGUser{Pk: 25025320, Username: "instagram"}
	user4 := instago.IGUser{Pk: 12345, Username: "testuser"}

	rms := []instago.ItemReelMention{{User: user1}}
	path = getStoryFilePath2("instagram", "25025320", "Bh7kySfDYq8", "123.mp4", 1520056661, rms)
	if path != "Instagram/instagram/stories/instagram-25025320-testuser-story-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.mp4" {
		t.Error(path)
		return
	}

	// test username more than filename length 256
	rms2 := []instago.ItemReelMention{{User: user2}}
	path = getStoryFilePath2("instagram", "25025320", "Bh7kySfDYq8", "123.mp4", 1520056661, rms2)
	if path != "Instagram/instagram/stories/instagram-25025320-story-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.mp4" {
		t.Error(path)
		return
	}

	// test username more than filename length 256
	rms3 := []instago.ItemReelMention{{User: user2}, {User: user1}}
	path = getStoryFilePath2("instagram", "25025320", "Bh7kySfDYq8", "123.mp4", 1520056661, rms3)
	if path != "Instagram/instagram/stories/instagram-25025320-testuser-story-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.mp4" {
		t.Error(path)
		return
	}

	// test duplicate username
	rms4 := []instago.ItemReelMention{{User: user3}, {User: user2}, {User: user1}}
	path = getStoryFilePath2("instagram", "25025320", "Bh7kySfDYq8", "123.mp4", 1520056661, rms4)
	if path != "Instagram/instagram/stories/instagram-25025320-testuser-story-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.mp4" {
		t.Error(path)
		return
	}

	// test duplicate username
	rms5 := []instago.ItemReelMention{{User: user3}, {User: user2}, {User: user1}, {User: user4}}
	path = getStoryFilePath2("instagram", "25025320", "Bh7kySfDYq8", "123.mp4", 1520056661, rms5)
	if path != "Instagram/instagram/stories/instagram-25025320-testuser-story-2018-03-03T13:57:41+08:00-Bh7kySfDYq8-1520056661.mp4" {
		t.Error(path)
		return
	}
}

func TestGetPostLiveFilePath(t *testing.T) {
	path := getPostLiveFilePath("instagram", "25025320", "123.mp4", "video", 1520056661)
	if path != "Instagram/instagram/postlives/instagram-25025320-postlive-video-2018-03-03T13:57:41+08:00-1520056661.mp4" {
		t.Error(path)
		return
	}
}

func TestGetPostLiveMergedFilePath(t *testing.T) {
	vpath := "Instagram/instagram/postlives/instagram-25025320-postlive-video-2018-03-03T13:57:41+08:00-1520056661.mp4"
	apath := "Instagram/instagram/postlives/instagram-25025320-postlive-audio-2018-03-03T13:57:41+08:00-1520056661.mp4"
	path := getPostLiveMergedFilePath(vpath, apath)
	if path != "Instagram/instagram/postlives/instagram-25025320-postlive-2018-03-03T13:57:41+08:00-1520056661.mp4" {
		t.Error(path)
		return
	}
}

func TestAppendIndexToFilename(t *testing.T) {
	nf := appendIndexToFilename("instagram-25025320-post-Bh7kySfDYq8-2018-03-03T13:57:41+08:00-1520056661.mp4", 0)
	if nf != "instagram-25025320-post-Bh7kySfDYq8-2018-03-03T13:57:41+08:00-1520056661-0.mp4" {
		t.Error(nf)
	}
	nf = appendIndexToFilename("instagram-25025320-post-Bh7kySfDYq8-2018-03-03T13:57:41+08:00-1520056661.mp4", 1)
	if nf != "instagram-25025320-post-Bh7kySfDYq8-2018-03-03T13:57:41+08:00-1520056661-1.mp4" {
		t.Error(nf)
	}
}

func TestGetUserProfilPicFilePath(t *testing.T) {
	path := getUserProfilPicFilePath("instagram", "25025320", "https://instagram.fkhh1-2.fna.fbcdn.net/vp/893534d61bdc5ea6911593d3ee0a1922/5B6363AB/t51.2885-19/s320x320/14719833_310540259320655_1605122788543168512_a.jpg", 1520056661)
	if path != "Instagram/instagram/instagram-25025320-profile_pic-1520056661.jpg" {
		t.Error(path)
		return
	}
}

func TestGetIdUsernamePath(t *testing.T) {
	path := getIdUsernamePath("25025320", "instagram")
	if path != "Data/ID-USERNAME/25025320-instagram" {
		t.Error(path)
		return
	}
}
