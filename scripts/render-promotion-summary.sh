#!/usr/bin/env bash

set -euo pipefail

if [ "$#" -lt 1 ] || [ "$#" -gt 2 ]; then
  echo "Usage: scripts/render-promotion-summary.sh SNAPSHOT.json [COMPARISON.json]" >&2
  exit 2
fi

snapshot="$1"
comparison="${2:-}"

jq -e '
  .schema_version == 1 and
  (.captured_at | type == "string") and
  (.github.primary_repository.stars | type == "number")
' "$snapshot" >/dev/null || {
  echo "Invalid promotion snapshot: $snapshot" >&2
  exit 1
}

if [ -n "$comparison" ]; then
  jq -e '
    (.from | type == "string") and
    (.to | type == "string") and
    (.github.primary_repository.stars.delta | type == "number")
  ' "$comparison" >/dev/null || {
    echo "Invalid promotion comparison: $comparison" >&2
    exit 1
  }
fi

jq -r '
  def display($value): if $value == null then "n/a" else ($value | tostring) end;

  "# JuggleIM promotion metrics\n",
  "Captured at **\(.captured_at)** for `\(.campaign.name // "promotion")`.\n",
  "| Metric | Current |",
  "| --- | ---: |",
  "| im-server Stars | \(.github.primary_repository.stars) |",
  "| Organization Stars | \(.github.organization.total_stars) |",
  "| Forks | \(.github.primary_repository.forks) |",
  "| Watchers | \(.github.primary_repository.watchers) |",
  "| Stars, last 24 hours | \(display(.github.primary_repository.star_growth.last_24h)) |",
  "| Stars, last 7 days | \(display(.github.primary_repository.star_growth.last_7d)) |",
  "| Stars, last 30 days | \(display(.github.primary_repository.star_growth.last_30d)) |",
  "| GitHub Discussions | \(.github.primary_repository.campaign_activity.discussions.total) |",
  "| Community discussion comments | \(display(.github.primary_repository.campaign_activity.discussions.community_comments_total // .github.primary_repository.campaign_activity.discussions.comments_total)) |",
  "| Maintainer discussion comments | \(display(.github.primary_repository.campaign_activity.discussions.maintainer_comments_total)) |",
  "| DEV reactions | \(display(.dev.reactions)) |",
  "| DEV comments | \(display(.dev.comments)) |",
  "| DEV page views | \(display(.dev.page_views)) |",
  "\n## Campaign activity\n",
  "- Issues created: \(.github.primary_repository.campaign_activity.issues.created)",
  "- Human pull requests: \(.github.primary_repository.campaign_activity.pull_requests.human_created)",
  "- Bot pull requests: \(.github.primary_repository.campaign_activity.pull_requests.bot_created)",
  "- Repository: \(.github.primary_repository.url)",
  if .dev.url then "- DEV article: \(.dev.url)" else empty end
' "$snapshot"

if [ -n "$comparison" ]; then
  jq -r '
    def signed($value):
      if $value == null then "n/a"
      elif $value > 0 then "+\($value)"
      else ($value | tostring)
      end;

    "\n## Change since previous automated snapshot\n",
    "| Metric | Change |",
    "| --- | ---: |",
    "| im-server Stars | \(signed(.github.primary_repository.stars.delta)) |",
    "| Organization Stars | \(signed(.github.organization_stars.delta)) |",
    "| Forks | \(signed(.github.primary_repository.forks.delta)) |",
    "| Watchers | \(signed(.github.primary_repository.watchers.delta)) |",
    "| DEV reactions | \(signed(.dev.reactions.delta)) |",
    "| DEV comments | \(signed(.dev.comments.delta)) |"
  ' "$comparison"
fi
