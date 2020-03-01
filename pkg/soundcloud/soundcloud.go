package main

import (
	"fmt"
	types "github.com/treethought/amnisiac/pkg/types"
	"github.com/yanatan16/golang-soundcloud/soundcloud"
	"net/url"
	"os"
	"strings"
)

func Client() (client soundcloud.Api) {

	scID := os.Getenv("SC_CLIENT_ID")
	scSecret := os.Getenv("SC_CLIENT_SECRET")

	client = soundcloud.Api{
		ClientId:     scID,
		ClientSecret: scSecret,
	}

	return client

}

func GetUser(username string) (user soundcloud.User, err error) {
	client := Client()

	users, err := client.Users(url.Values{"q": []string{"treethought"}})
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		// fmt.Printf("User %s, avatar: %s", u.Username, u.AvatarUrl)
		if strings.EqualFold(username, u.Username) {
			return *u, nil

		}
	}
	return
	// fmt.Println(user.Username)

	// fmt.Printf("User: %s, tracks: %s", user.Username, user.Id)
}

func UserTracks(username string) (tracks []*soundcloud.Track, err error) {
	client := Client()
	user, err := GetUser(username)

	track_list, err := client.User(user.Id).Tracks(nil)
	if err != nil {
		return nil, err
	}
	for _, t := range track_list {
		fmt.Println(t.Title)
		if t.Streamable {
			tracks = append(tracks, t)

		}
	}

	// fmt.Println(user)
	return
}

func FetchItemsFromSoundcloud() (sc_items []*types.Item, err error) {
	tracks, _ := UserTracks("treethought")
	for _, t := range tracks {
		item, _ := types.MakeItemFromSoundcloudTrack(t)
		fmt.Println(item)
		sc_items = append(sc_items, &item)
	}
	return sc_items, nil

}

func main() {
	GetUser("treethought")
	UserTracks("treethought")
}
