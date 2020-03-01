package types
import (
	"github.com/jzelinskie/geddit"
	"github.com/satori/go.uuid"
	"github.com/yanatan16/golang-soundcloud/soundcloud"
)

type Item struct {
	ID              string `json:"id"`
	PlatformTrackID string `json:"platform_track_id"`
	RawTitle        string `json:"raw_title"`
	Title           string `json:"title"`
	Artist          string `json:"artist"`
	URL             string `json:"url"`
	Domain          string `json:"domain"`
	SourcePlatform  string `json:"source_platform"`
	SubReddit       string `json:"sub_reddit"`
	StreamURL       string `json:"stream_url"`
	Duration        string `json:"duration"`
	ArtworkURL      string `json:"artwork_url"`
	streamable      bool   `json:"streamable"`
	Genre           string `json:"genre"`

	// DatePublished string
	// SourcePlatform string
	// DateProcessed string

}

func MakeItemFromRedditPost(submission *geddit.Submission) (item Item, err error) {

	item.ID = string(uuid.NewV4().String())
	item.PlatformTrackID = submission.FullID
	item.RawTitle = submission.Title

    item.SubReddit = submission.Subreddit
    item.URL = submission.URL
    item.Domain = submission.Domain

    item.SourcePlatform = "reddit"




	return item, nil

}


func MakeItemFromSoundcloudTrack (track *soundcloud.Track) (item Item, err error) {

	item.ID = string(uuid.NewV4().String())
    item.PlatformTrackID = string(track.Id)
    item.RawTitle = track.User.Username + " - " + track.Title
    item.Artist = track.User.Username

    item.SourcePlatform = "sc"
    item.URL = track.User.PermalinkUrl
    item.Domain = "soundcloud.com"
    item.StreamURL = track.StreamUrl


    return item, nil


}


