# Restic Exporter

A prometheus exporter for monitoring restic repository stats. Currently limited to repos using S3 compatible buckets as backend. The functionality can be easily extended to other providers if there is enough interest. This project was mainly done to monitor my [homelab](https://github.com/kmjayadeep/homelab-k8s) backups and fire alerts if the backups are not up-to-date

## Exposed metrics example

```
# TYPE restic_repo_last_snapshot gauge
restic_repo_last_snapshot{host="jayadeep-nuc",repo="jd-backup-nuc",shortId="2f9462c6"} 1.691492001e+09
# HELP restic_repo_stats_fetch_duration Amount of time taken to fetch restic repo stats
# TYPE restic_repo_stats_fetch_duration histogram
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="0.1"} 0
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="0.2"} 0
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="0.4"} 0
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="1"} 0
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="3"} 0
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="8"} 0
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="20"} 1
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="60"} 1
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="120"} 1
restic_repo_stats_fetch_duration_bucket{repo="jd-backup-nuc",le="+Inf"} 1
restic_repo_stats_fetch_duration_sum{repo="jd-backup-nuc"} 9.926568872
restic_repo_stats_fetch_duration_count{repo="jd-backup-nuc"} 1
# HELP restic_repo_stats_refresh_success_total Number of times restic stats are successfuly refreshed in the cache for the repo
# TYPE restic_repo_stats_refresh_success_total counter
restic_repo_stats_refresh_success_total{repo="jd-backup-nuc"} 1
# HELP restic_s3_object_count Number of objects in s3 bucket
# TYPE restic_s3_object_count gauge
restic_s3_object_count{repo="jd-backup-nuc"} 4956
# HELP restic_s3_size_total Total Size of objects in s3 bucket
# TYPE restic_s3_size_total gauge
restic_s3_size_total{repo="jd-backup-nuc"} 5.497566033e+10
# HELP restic_stats_refresh_total Number of times restic stats are refreshed in the cache
# TYPE restic_stats_refresh_total counter
restic_stats_refresh_total 1
```

## Configuration and Running

Refer to `config.yaml.example` to configure your repositories. Then run using docker as below


```
docker run -p 18090:18090  -v ./config.yaml:/app/config.yaml kmjayadeep/restic-exporter:latest
```

It will expose the metrics on port 18090 which can be scrapped by prometheus. The stats are refreshed every 10 mins by default and can be customized in the config file.

### Using docker-compose

```
docker-compose up
```

### Kubernetes

You can find the Kubernetes kustomize manifests under `k8s` folder

```
kubectl apply -k k8s/
```
