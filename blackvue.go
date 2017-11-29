package blackvue

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// New initializes a fresh blackvue client
func New(ip string) *Client {
	return &Client{
		ip: ip,
	}
}

type Client struct {
	ip string
}

type Videos struct {
	Front, Rear, Unknown []Video
}

// Video represents the base path to a video or gps file
// ie. {Video}.thm
type Video string

func (v Video) MP4() string {
	return fmt.Sprintf("%s.mp4", v)
}

func (v Video) THM() string {
	return fmt.Sprintf("%s.thm", v)
}

func (c *Client) list() (Videos, error) {
	var vids Videos
	path := fmt.Sprintf("http://%s/blackvue_vod.cgi", c.ip)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return vids, err
	}

	httpCli := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	resp, err := httpCli.Do(req)
	if err != nil {
		return vids, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return vids, err
	}

	lines := strings.Split(string(bs), "\r\n")
	for _, l := range lines {
		base := strings.TrimPrefix(strings.TrimSuffix(l, ".mp4,s:1000000"), "n:/Record/")

		if strings.HasSuffix(base, "R") {
			vids.Rear = append(vids.Rear, Video(base))
		} else if strings.HasSuffix(base, "F") {
			vids.Front = append(vids.Front, Video(base))
		} else if base == "v:1.00" {
			// ignore this
		} else {
			vids.Unknown = append(vids.Unknown, Video(base))
		}
	}

	return vids, nil
}

// List enumerates all the video files on Dashcam
func (c *Client) List() (Videos, error) {
	return c.list()
}

// Sync pulls all the video files not found in path
func (c *Client) Sync(path string) error {
	rearDir := filepath.Join(path, "rear")
	if err := os.MkdirAll(rearDir, 744); err != nil {
		log.Fatal(err)
	}

	frontDir := filepath.Join(path, "front")
	if err := os.MkdirAll(frontDir, 744); err != nil {
		log.Fatal(err)
	}

	vids, err := c.list()
	if err != nil {
		return err
	}

	return c.sync(path, vids)
}

func (c *Client) sync(path string, vids Videos) error {

	rearDir := filepath.Join(path, "rear")
	frontDir := filepath.Join(path, "front")

	for _, rear := range vids.Rear {
		path := filepath.Join(rearDir, rear.MP4())
		if _, err := os.Stat(path); err != nil {
			if err := c.fetchVideo(rearDir, rear); err != nil {
				log.Printf("failed to fetch video %s: %s\n", rear, err)
			}
		}
	}

	for _, front := range vids.Front {
		path := filepath.Join(frontDir, front.MP4())
		if _, err := os.Stat(path); err == nil {
			if err := c.fetchVideo(frontDir, front); err != nil {
				log.Printf("failed to fetch video %s: %s\n", front, err)
			}
		}
	}
	return nil
}

func (c *Client) fetchVideo(path string, vid Video) error {
	outMP4, err := os.Create(filepath.Join(path, vid.MP4()))
	if err != nil {
		return err
	}
	outTHM, err := os.Create(filepath.Join(path, vid.THM()))
	if err != nil {
		return err
	}

	log.Printf("fetching %s\n", vid)

	// mp4
	uri := fmt.Sprintf("http://%s/Record/%s", c.ip, vid.MP4())
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(outMP4, resp.Body); err != nil {
		return err
	}

	// thm
	path = fmt.Sprintf("http://%s/Record/%s", c.ip, vid.THM())
	resp, err = http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(outTHM, resp.Body); err != nil {
		return err
	}

	return nil
}
