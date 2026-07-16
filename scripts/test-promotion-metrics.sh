#!/usr/bin/env bash

set -euo pipefail

script_dir="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

create_snapshot() {
  captured_at="$1"
  org_stars="$2"
  repo_stars="$3"
  forks="$4"
  watchers="$5"
  visitors="$6"
  cloners="$7"
  reactions="$8"
  comments="$9"
  output="${10}"

  jq -n \
    --arg captured_at "$captured_at" \
    --argjson org_stars "$org_stars" \
    --argjson repo_stars "$repo_stars" \
    --argjson forks "$forks" \
    --argjson watchers "$watchers" \
    --argjson visitors "$visitors" \
    --argjson cloners "$cloners" \
    --argjson reactions "$reactions" \
    --argjson comments "$comments" '
    {
      schema_version: 1,
      captured_at: $captured_at,
      github: {
        organization: {total_stars: $org_stars},
        primary_repository: {
          stars: $repo_stars,
          forks: $forks,
          watchers: $watchers,
          star_growth: {last_24h: 2, last_7d: 4, last_30d: 14, daily_30d: []},
          campaign_activity: {
            issues: {created: 0, items: []},
            pull_requests: {created: 1, human_created: 1, bot_created: 0, items: []},
            discussions: {
              total: 5,
              comments_total: $comments,
              created_since_campaign: [],
              updated_since_campaign: []
            }
          },
          traffic: {
            views_14d: {uniques: $visitors},
            clones_14d: {uniques: $cloners},
            referrers_14d: []
          }
        }
      },
      dev: {reactions: $reactions, comments: $comments, page_views: null},
      community: {}
    }
  ' >"$output"
}

create_snapshot "2026-07-16T06:05:11Z" 3802 3584 363 240 79 160 0 0 "$tmp_dir/before.json"
create_snapshot "2026-07-17T06:05:11Z" 3805 3586 364 241 84 162 1 2 "$tmp_dir/after.json"

"$script_dir/compare-promotion-metrics.sh" "$tmp_dir/before.json" "$tmp_dir/after.json" \
  | jq -e '
      .github.organization_stars.delta == 3 and
      .github.primary_repository.stars.delta == 2 and
      .github.primary_repository.forks.delta == 1 and
      .github.primary_repository.watchers.delta == 1 and
      .github.primary_repository.traffic_14d.unique_visitors.delta == 5 and
      .github.primary_repository.traffic_14d.unique_cloners.delta == 2 and
      .dev.reactions.delta == 1 and
      .dev.comments.delta == 2 and
      .dev.page_views.delta == null
    ' >/dev/null

bash -n "$script_dir/collect-promotion-metrics.sh" "$script_dir/compare-promotion-metrics.sh"
echo "promotion metrics tests passed"
