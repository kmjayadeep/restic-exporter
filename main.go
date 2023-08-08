package main

import (
	"fmt"
	"time"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
	"github.com/kmjayadeep/restic-monitoring/internal/config"
	"github.com/kmjayadeep/restic-monitoring/internal/stats"
)

type LastSnapshot struct {
	Time     time.Time `json:"time"`
	ShortID  string    `json:"short_id"`
	HostName string    `json:"hostname"`
}

func main() {
	c, err := config.ParseConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	for _, repo := range c.Repos {
		s, err := stats.FetchStats(repo)
		if err != nil {
			fmt.Println("unable to fetch stats for ", repo.Name)
			continue
		}
		fmt.Printf("%+v", s)
	}

}
