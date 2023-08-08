package stats

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kmjayadeep/restic-monitoring/internal/config"
)

type LastSnapshot struct {
	Time     time.Time `json:"time"`
	ShortID  string    `json:"short_id"`
	HostName string    `json:"hostname"`
}

type Stats struct {
	ObjectsCount int64
	Size         int64
	LastSnapshot LastSnapshot
}

func FetchStats(r config.ResticRepository) (*Stats, error) {
	cmd := exec.Command("restic", "snapshots", "latest", "--json")

	cmd.Env = append(cmd.Environ(), "AWS_ACCESS_KEY_ID="+r.AccessKey, "AWS_SECRET_ACCESS_KEY="+r.SecretKey, "RESTIC_PASSWORD="+r.ResticPassword, "RESTIC_REPOSITORY="+r.Endpoint)

	out, err := cmd.Output()

	if err != nil {
		return nil, nil
	}

	snaps := []LastSnapshot{}

	if err := json.Unmarshal(out, &snaps); err != nil {
		return nil, nil
	}

	last := LastSnapshot{}

	if len(snaps) > 0 {
		last = snaps[0]
	}

	e := strings.TrimPrefix(r.Endpoint, "s3:")

	u, err := url.Parse(e)
	if err != nil {
		return nil, err
	}

	fmt.Printf("proto: %q, endpoint: %q, bucket: %q", u.Scheme, u.Host, u.Path)

	bucket := u.Path

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(r.AccessKey, r.SecretKey, ""),
		Endpoint:    aws.String(u.Host),
		Region:      aws.String("us-west-000"),
	}))

	svc := s3.New(sess)

	in := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	size := int64(0)
	count := int64(0)

	svc.ListObjectsV2Pages(in,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, v := range page.Contents {
				size = size + *v.Size
			}

			count = count + *page.KeyCount
			return !lastPage
		})

	return &Stats{
		ObjectsCount: count,
		Size:         size,
		LastSnapshot: last,
	}, nil
}
