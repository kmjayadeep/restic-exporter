package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/kmjayadeep/restic-monitoring/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	resticRefreshTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "restic_stats_refresh_total",
		Help: "Number of times restic stats are refreshed in the cache",
	})

	resticRepoRefreshTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "restic_repo_stats_refresh_success_total",
		Help: "Number of times restic stats are successfuly refreshed in the cache for the repo",
	}, []string{"repo"})

	resticRepoRefreshFailTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "restic_repo_stats_refresh_fail_total",
		Help: "Number of times restic stats are unsuccessfuly refreshed in the cache for the repo",
	}, []string{"repo"})

	resticS3ObjectCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_s3_object_count",
		Help: "Number of objects in s3 bucket",
	}, []string{"repo"})

	resticS3Size = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_s3_size_total",
		Help: "Total Size of objects in s3 bucket",
	}, []string{"repo"})

	resticRepoLastSnapshot = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "restic_repo_last_snapshot",
		Help: "Last snapshot in the restic repo",
	}, []string{"repo", "host", "shortId"})

	resticRepoStatsDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "restic_repo_stats_fetch_duration",
		Help:    "Amount of time taken to fetch restic repo stats",
		Buckets: []float64{.1, .2, .4, 1, 3, 8, 20, 60, 120},
	}, []string{"repo"})
)

func RefreshMetrics(c *config.Config) {
	for _, repo := range c.Repos {

		go func(repo config.ResticRepository) {
			t := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			s, err := FetchStats(ctx, repo)
			if err != nil {
				fmt.Println("unable to fetch stats for ", repo.Name, err)
				resticRepoRefreshFailTotal.WithLabelValues(repo.Name).Inc()
				return
			}
			fmt.Println(repo.Name, s)
			resticRepoRefreshTotal.WithLabelValues(repo.Name).Inc()
			resticS3ObjectCount.WithLabelValues(repo.Name).Set(float64(s.ObjectsCount))
			resticS3Size.WithLabelValues(repo.Name).Set(float64(s.Size))
			resticRepoLastSnapshot.WithLabelValues(repo.Name, s.LastSnapshot.HostName, s.LastSnapshot.ShortID).Set(float64(s.LastSnapshot.Time.Unix()))
			resticRepoStatsDuration.WithLabelValues(repo.Name).Observe(float64(time.Since(t).Seconds()))

		}(repo)
	}
	resticRefreshTotal.Inc()
}

func Run(c *config.Config) {
	ticker := time.NewTicker(time.Duration(c.RefreshMinutes) * time.Minute)

	RefreshMetrics(c)
	for range ticker.C {
		RefreshMetrics(c)
	}
}
