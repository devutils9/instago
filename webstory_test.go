package instago

import (
	"fmt"
	"os"
	"testing"
)

func ExampleGetWebFeedReelsTray(t *testing.T) {
	mgr, err := NewInstagramApiManager("auth.json")
	if err != nil {
		t.Error(err)
		return
	}

	url, err := mgr.GetGetWebFeedReelsTrayUrl()
	if err != nil {
		t.Error(err)
		return
	}
	//t.Log(url)

	rms, err := mgr.GetWebFeedReelsTray(url)
	if err != nil {
		t.Error(err)
		return
	}

	//t.Log(rms)
	for _, rm := range rms {
		t.Log(rm.User.Username, rm.User.Id)
	}
}

func ExampleGetUserStoryByWebGraphql(t *testing.T) {
	mgr, err := NewInstagramApiManager("auth.json")
	if err != nil {
		t.Error(err)
		return
	}
	rm, err := mgr.GetUserStoryByWebGraphql(os.Getenv("id"), os.Getenv("storyhash"))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rm)
	for _, item := range rm.Items {
		t.Log(item.GetUsername(), item.GetUserId(), item.GetMediaUrl(), FormatTimestamp(item.GetTimestamp()))
	}
}

func ExampleGetIdFromWebStoryUrl() {
	m := NewApiManager(nil, nil)
	id, err := m.GetIdFromWebStoryUrl("https://www.instagram.com/stories/highlights/17862445040107085/")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(id)
	// Output:
	// 25025320
}

func ExampleGetInfoFromWebStoryUrl() {
	m := NewApiManager(nil, nil)
	info, err := m.GetInfoFromWebStoryUrl("https://www.instagram.com/stories/highlights/17862445040107085/")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(info.User.Username)
	// Output:
	// instagram
}
