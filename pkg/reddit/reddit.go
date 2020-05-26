package reddit

import (
	"fmt"

	"github.com/gocolly/colly"
	"github.com/jzelinskie/geddit"

	types "github.com/treethought/amnisiac/pkg/types"
	"strings"
)

func FetchSubmissions(subreddit string) (submissions []*geddit.Submission, err error) {
	clean_subname := subreddit[3:]
	session := geddit.NewSession("testagent")
	if err != nil {
		return submissions, err

	}

	subOpts := geddit.ListingOptions{
		Limit: 10,
	}

	submissions, err = session.SubredditSubmissions(clean_subname, geddit.DefaultPopularity, subOpts)
	if err != nil {
		return nil, err
	}
	if len(submissions) == 0 {
		// fmt.Println("No submissions found!")

	}

	// for _, s := range submissions {
	return

}

// Checks if a URL is from a music-related domain
func IsMediaURL(url string) bool {
	media_domains := []string{"youtube", "vimeo"}
	for _, d := range media_domains {
		if strings.Contains(url, d) {
			return true

		}

	}
	return false

}

func ParseMediaURLs(submissions []*geddit.Submission) (mediaURLs []string) {
	// fmt.Println("PARSING MEDIA URLS")
	for _, s := range submissions {
		println(s.URL)
		if IsMediaURL(s.URL) {
			// fmt.Printf("Media URL found: %s", s.URL)
			mediaURLs = append(mediaURLs, s.URL)
		}

	}
	return

}

// Scrapes music relate subreddits from promiment wikis
// Returns all subreddits amnsiac will use for aggregating
func SubRedditsFromWiki(subreddit string, wikiname string) (subreddits []string, err error) {
	c := colly.NewCollector()

	wikiUrl := fmt.Sprintf("https://reddit.com/r/%s/wiki/%s", subreddit, wikiname)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if strings.HasPrefix(url, "/r/") {
			subreddits = append(subreddits, url)
		}
	})

	c.Visit(wikiUrl)

	return

}

func FetchItemsFromReddit(query string) (reddit_items []*types.Item, err error) {
	subreddits, _ := SubRedditsFromWiki("Music", "musicsubreddits")
	var submissions []*geddit.Submission
	for _, s := range subreddits {
		if strings.Contains(s, query) {
			// log.Println("fetching", s)
			posts, _ := FetchSubmissions(s)
			for _, p := range posts {
				if IsMediaURL(p.URL) {
					submissions = append(submissions, p)
				}
			}
		}
	}
	for _, s := range submissions {
		item, err := types.MakeItemFromRedditPost(s)
		if err != nil {
			panic(err)
		}

		reddit_items = append(reddit_items, &item)
	}
	return reddit_items, nil

}
