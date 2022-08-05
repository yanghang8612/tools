package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func sol() {
	repData := doGet("https://api.github.com/repos/tronprotocol/solidity/releases")
	var releases []Releases
	if err := json.Unmarshal(repData, &releases); err == nil {
		var releaseStr string
		for _, release := range releases {
			rCommittedTime, _ := time.Parse("2006-01-02T15:04:05Z", release.CreatedAt)
			rPublishedTime, _ := time.Parse("2006-01-02T15:04:05Z", release.PublishedAt)
			releaseStr = fmt.Sprintf("%s,%s,%s,%s",
				release.TagName,
				release.CreatedAt,
				release.PublishedAt,
				convertSecToReadable(rPublishedTime.Unix()-rCommittedTime.Unix()))

			for _, asset := range release.Assets {
				aCreatedTime, _ := time.Parse("2006-01-02T15:04:05Z", asset.CreatedAt)
				aUpdatedTime, _ := time.Parse("2006-01-02T15:04:05Z", asset.UpdatedAt)
				fmt.Printf("%s,%s,%s,%s,%s,%s,%s\n",
					releaseStr,
					asset.Name,
					asset.Uploader.Login,
					asset.CreatedAt,
					convertSecToReadable(aCreatedTime.Unix()-rPublishedTime.Unix()),
					asset.UpdatedAt,
					convertSecToReadable(aUpdatedTime.Unix()-rPublishedTime.Unix()))
			}
		}
	}
}

func convertSecToReadable(sec int64) string {
	timeStr := time.Unix(sec, 0).Format("15事04分05秒")
	day := sec / (24 * 60 * 60)
	if day > 0 {
		timeStr = strconv.Itoa(int(day)) + "天" + timeStr
	}
	if day > 30 {
		timeStr = timeStr + " (WARNING)"
	}
	return timeStr
}

type Releases struct {
	Id          uint64
	TagName     string `json:"tag_name"`
	CreatedAt   string `json:"created_at"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name          string
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
		DownloadCount uint64 `json:"download_count"`
		Uploader      struct {
			Login string
		}
	}
}
