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
        organization: {
          total_stars: $org_stars,
          metadata_coverage: {
            descriptions: 30,
            homepages: 24,
            topics: 23,
            detected_licenses: 29,
            complete: 20
          }
        },
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
              maintainer_comments_total: 1,
              bot_comments_total: 0,
              community_comments_total: $comments,
              community_comments_since_campaign: [],
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

jq -n '{reddit: [{community: "r/selfhosted", status: "published"}]}' \
  >"$tmp_dir/community.json"
jq -e 'select(type == "object")' "$tmp_dir/community.json" \
  >"$tmp_dir/community-normalized.json"
jq -e '.reddit[0].community == "r/selfhosted"' \
  "$tmp_dir/community-normalized.json" >/dev/null

"$script_dir/compare-promotion-metrics.sh" "$tmp_dir/before.json" "$tmp_dir/after.json" \
  | tee "$tmp_dir/comparison.json" \
  | jq -e '
      .github.organization_stars.delta == 3 and
      .github.organization_metadata.changes.complete == 0 and
      .github.primary_repository.stars.delta == 2 and
      .github.primary_repository.forks.delta == 1 and
      .github.primary_repository.watchers.delta == 1 and
      .github.primary_repository.traffic_14d.unique_visitors.delta == 5 and
      .github.primary_repository.traffic_14d.unique_cloners.delta == 2 and
      .dev.reactions.delta == 1 and
      .dev.comments.delta == 2 and
      .dev.page_views.delta == null
    ' >/dev/null

"$script_dir/render-promotion-summary.sh" \
  "$tmp_dir/after.json" \
  "$tmp_dir/comparison.json" >"$tmp_dir/summary.md"

grep -F '| im-server Stars | 3586 |' "$tmp_dir/summary.md" >/dev/null
grep -F '| im-server Stars | +2 |' "$tmp_dir/summary.md" >/dev/null
grep -F '| Community discussion comments | 2 |' "$tmp_dir/summary.md" >/dev/null
grep -F '| Maintainer discussion comments | 1 |' "$tmp_dir/summary.md" >/dev/null
grep -F -- '- Human pull requests: 1' "$tmp_dir/summary.md" >/dev/null

jq '.github.primary_repository.star_growth |= {
  available: false,
  last_24h: null,
  last_7d: null,
  last_30d: null,
  daily_30d: []
}' "$tmp_dir/after.json" >"$tmp_dir/no-star-history.json"

"$script_dir/render-promotion-summary.sh" \
  "$tmp_dir/no-star-history.json" >"$tmp_dir/no-star-history.md"
grep -F '| Stars, last 24 hours | n/a |' "$tmp_dir/no-star-history.md" >/dev/null

bash -n \
  "$script_dir/collect-promotion-metrics.sh" \
  "$script_dir/compare-promotion-metrics.sh" \
  "$script_dir/render-promotion-summary.sh"
echo "promotion metrics tests passed"
