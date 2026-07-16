#!/usr/bin/env bash

set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "Usage: scripts/compare-promotion-metrics.sh BEFORE.json AFTER.json" >&2
  exit 2
fi

before="$1"
after="$2"

for file in "$before" "$after"; do
  [ -f "$file" ] || { echo "Snapshot not found: $file" >&2; exit 1; }
  jq -e '.schema_version == 1 and (.captured_at | type == "string")' "$file" >/dev/null || {
    echo "Invalid promotion snapshot: $file" >&2
    exit 1
  }
done

jq -n --slurpfile before "$before" --slurpfile after "$after" '
  def delta($old; $new):
    if $old == null or $new == null then null else $new - $old end;

  ($before[0]) as $b |
  ($after[0]) as $a |
  {
    from: $b.captured_at,
    to: $a.captured_at,
    github: {
      organization_stars: {
        before: $b.github.organization.total_stars,
        after: $a.github.organization.total_stars,
        delta: delta($b.github.organization.total_stars; $a.github.organization.total_stars)
      },
      primary_repository: {
        stars: {
          before: $b.github.primary_repository.stars,
          after: $a.github.primary_repository.stars,
          delta: delta($b.github.primary_repository.stars; $a.github.primary_repository.stars)
        },
        forks: {
          before: $b.github.primary_repository.forks,
          after: $a.github.primary_repository.forks,
          delta: delta($b.github.primary_repository.forks; $a.github.primary_repository.forks)
        },
        watchers: {
          before: $b.github.primary_repository.watchers,
          after: $a.github.primary_repository.watchers,
          delta: delta($b.github.primary_repository.watchers; $a.github.primary_repository.watchers)
        },
        current_star_growth: $a.github.primary_repository.star_growth,
        campaign_activity: $a.github.primary_repository.campaign_activity,
        traffic_14d: {
          unique_visitors: {
            before: $b.github.primary_repository.traffic.views_14d.uniques,
            after: $a.github.primary_repository.traffic.views_14d.uniques,
            delta: delta($b.github.primary_repository.traffic.views_14d.uniques; $a.github.primary_repository.traffic.views_14d.uniques)
          },
          unique_cloners: {
            before: $b.github.primary_repository.traffic.clones_14d.uniques,
            after: $a.github.primary_repository.traffic.clones_14d.uniques,
            delta: delta($b.github.primary_repository.traffic.clones_14d.uniques; $a.github.primary_repository.traffic.clones_14d.uniques)
          },
          current_referrers: $a.github.primary_repository.traffic.referrers_14d
        }
      }
    },
    dev: {
      reactions: {
        before: $b.dev.reactions,
        after: $a.dev.reactions,
        delta: delta($b.dev.reactions; $a.dev.reactions)
      },
      comments: {
        before: $b.dev.comments,
        after: $a.dev.comments,
        delta: delta($b.dev.comments; $a.dev.comments)
      },
      page_views: {
        before: $b.dev.page_views,
        after: $a.dev.page_views,
        delta: delta($b.dev.page_views; $a.dev.page_views)
      }
    },
    community: $a.community
  }
'
