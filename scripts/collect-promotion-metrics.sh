#!/usr/bin/env bash

set -euo pipefail

ORG="${JUGGLEIM_ORG:-juggleim}"
REPO="${JUGGLEIM_REPO:-im-server}"
CAMPAIGN="${JUGGLEIM_CAMPAIGN:-juggleim_oss_launch}"
CAMPAIGN_PUBLISHED_AT="${JUGGLEIM_CAMPAIGN_PUBLISHED_AT:-2026-07-16T06:05:11Z}"
DEV_ARTICLE_SLUG="${JUGGLEIM_DEV_ARTICLE_SLUG:-juggleim-an-open-source-self-hosted-messaging-backend-built-in-go-43nh}"
DEV_USERNAME="${JUGGLEIM_DEV_USERNAME:-yuwnloyblog}"

OUTPUT=""
COMMUNITY_FILE=""
PUBLIC_ONLY=false

usage() {
  cat <<'EOF'
Collect a comparable JuggleIM promotion snapshot from GitHub and DEV.

Usage:
  scripts/collect-promotion-metrics.sh [options]

Options:
  --output PATH          Write JSON to PATH instead of stdout.
  --community-file PATH  Merge a manually collected Reddit/LinkedIn JSON object.
  --public-only          Skip maintainer-only GitHub traffic endpoints.
  -h, --help             Show this help.

Environment:
  DEV_API_KEY                    Optional. Adds owner-only DEV page views.
  JUGGLEIM_ORG                   GitHub organization (default: juggleim).
  JUGGLEIM_REPO                  Primary repository (default: im-server).
  JUGGLEIM_CAMPAIGN              Campaign name.
  JUGGLEIM_CAMPAIGN_PUBLISHED_AT Campaign publication time in ISO 8601 UTC.
  JUGGLEIM_DEV_USERNAME          DEV username.
  JUGGLEIM_DEV_ARTICLE_SLUG      DEV article slug.

The script never writes credentials to the snapshot.
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --output)
      [ "$#" -ge 2 ] || { echo "--output requires a path" >&2; exit 2; }
      OUTPUT="$2"
      shift 2
      ;;
    --community-file)
      [ "$#" -ge 2 ] || { echo "--community-file requires a path" >&2; exit 2; }
      COMMUNITY_FILE="$2"
      shift 2
      ;;
    --public-only)
      PUBLIC_ONLY=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

for command_name in gh jq curl date mktemp; do
  command -v "$command_name" >/dev/null 2>&1 || {
    echo "Required command not found: $command_name" >&2
    exit 1
  }
done

if ! gh auth status >/dev/null 2>&1; then
  echo "GitHub CLI is not authenticated. Run: gh auth login" >&2
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

captured_at="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
now_epoch="$(date -u +%s)"
cutoff_24h=$((now_epoch - 86400))
cutoff_7d=$((now_epoch - 604800))
cutoff_30d=$((now_epoch - 2592000))

gh api "orgs/${ORG}/repos?type=public&per_page=100" --paginate \
  | jq -s 'add' >"$tmp_dir/repos.json"

gh api "repos/${ORG}/${REPO}" >"$tmp_dir/primary.json"

star_total="$(jq -r '.stargazers_count' "$tmp_dir/primary.json")"
last_star_page=$(((star_total + 99) / 100))
star_page="$last_star_page"

while [ "$star_page" -ge 1 ]; do
  star_file="$tmp_dir/stars-${star_page}.json"
  gh api -H 'Accept: application/vnd.github.star+json' \
    "repos/${ORG}/${REPO}/stargazers?per_page=100&page=${star_page}" >"$star_file"

  oldest_is_before_cutoff="$(jq --argjson cutoff "$cutoff_30d" \
    'if length == 0 then true else (.[0].starred_at | fromdateiso8601) <= $cutoff end' \
    "$star_file")"
  if [ "$oldest_is_before_cutoff" = "true" ] || [ "$star_page" -eq 1 ]; then
    break
  fi
  star_page=$((star_page - 1))
done

jq -s 'add' "$tmp_dir"/stars-*.json >"$tmp_dir/recent-stars.json"
jq --argjson cutoff_24h "$cutoff_24h" \
  --argjson cutoff_7d "$cutoff_7d" \
  --argjson cutoff_30d "$cutoff_30d" '
  {
    last_24h: ([.[] | select((.starred_at | fromdateiso8601) >= $cutoff_24h)] | length),
    last_7d: ([.[] | select((.starred_at | fromdateiso8601) >= $cutoff_7d)] | length),
    last_30d: ([.[] | select((.starred_at | fromdateiso8601) >= $cutoff_30d)] | length),
    daily_30d: ([.[] | select((.starred_at | fromdateiso8601) >= $cutoff_30d) | .starred_at[0:10]]
      | group_by(.) | map({date: .[0], stars: length}))
  }' "$tmp_dir/recent-stars.json" >"$tmp_dir/star-growth.json"

gh api "repos/${ORG}/${REPO}/issues?state=all&since=${CAMPAIGN_PUBLISHED_AT}&per_page=100" --paginate \
  | jq -s --arg published "$CAMPAIGN_PUBLISHED_AT" '
      add | [.[] | select(has("pull_request") | not) | select(.created_at >= $published)]
      | {created: length, items: map({number, title, state, comments, created_at})}
    ' >"$tmp_dir/issues.json"

gh api "repos/${ORG}/${REPO}/pulls?state=all&sort=created&direction=desc&per_page=100" --paginate \
  | jq -s --arg published "$CAMPAIGN_PUBLISHED_AT" '
      add | [.[] | select(.created_at >= $published)]
      | {
          created: length,
          human_created: ([.[] | select(.user.type != "Bot")] | length),
          bot_created: ([.[] | select(.user.type == "Bot")] | length),
          items: map({number, title, state, author: .user.login, author_type: .user.type, created_at})
        }
    ' >"$tmp_dir/pulls.json"

gh api graphql \
  -f owner="$ORG" \
  -f name="$REPO" \
  -f query='query($owner:String!,$name:String!){
    repository(owner:$owner,name:$name){
      discussions(first:100,orderBy:{field:UPDATED_AT,direction:DESC}){
        totalCount
        nodes{
          number
          title
          url
          createdAt
          updatedAt
          comments{totalCount}
        }
      }
    }
  }' \
  | jq --arg published "$CAMPAIGN_PUBLISHED_AT" '
      .data.repository.discussions as $discussions |
      {
        total: $discussions.totalCount,
        comments_total: ($discussions.nodes | map(.comments.totalCount) | add // 0),
        created_since_campaign: ($discussions.nodes
          | map(select(.createdAt >= $published))
          | map({number, title, url, created_at: .createdAt})),
        updated_since_campaign: ($discussions.nodes
          | map(select(.updatedAt >= $published))
          | map({number, title, url, updated_at: .updatedAt, comments: .comments.totalCount}))
      }
    ' >"$tmp_dir/discussions.json"

dev_public_url="https://dev.to/api/articles/${DEV_USERNAME}/${DEV_ARTICLE_SLUG}"
if curl --fail --silent --show-error --max-time 30 "$dev_public_url" >"$tmp_dir/dev-public.json"; then
  dev_available=true
else
  dev_available=false
  printf '%s\n' '{}' >"$tmp_dir/dev-public.json"
  echo "Warning: DEV public metrics were unavailable" >&2
fi

printf '%s\n' 'null' >"$tmp_dir/dev-owner.json"
if [ -n "${DEV_API_KEY:-}" ]; then
  if printf 'header = "api-key: %s"\n' "$DEV_API_KEY" \
    | curl --config - --fail --silent --show-error --max-time 30 \
      'https://dev.to/api/articles/me/published?per_page=1000' \
    | jq --arg slug "$DEV_ARTICLE_SLUG" '[.[] | select(.slug == $slug)][0] // null' \
      >"$tmp_dir/dev-owner.json"; then
    :
  else
    printf '%s\n' 'null' >"$tmp_dir/dev-owner.json"
    echo "Warning: authenticated DEV metrics were unavailable" >&2
  fi
fi

for metric in views clones referrers paths; do
  printf '%s\n' 'null' >"$tmp_dir/traffic-${metric}.json"
done

traffic_available=false
if [ "$PUBLIC_ONLY" = false ]; then
  traffic_available=true
  for metric in views clones popular/referrers popular/paths; do
    metric_file="$(printf '%s' "$metric" | tr '/' '-')"
    output_name="$(printf '%s' "$metric" | sed 's#popular/##')"
    if gh api "repos/${ORG}/${REPO}/traffic/${metric}" >"$tmp_dir/traffic-${output_name}.json" 2>"$tmp_dir/traffic-${metric_file}.err"; then
      :
    else
      traffic_available=false
      printf '%s\n' 'null' >"$tmp_dir/traffic-${output_name}.json"
      echo "Warning: GitHub traffic metric unavailable: $metric" >&2
    fi
  done
fi

printf '%s\n' '{}' >"$tmp_dir/community.json"
if [ -n "$COMMUNITY_FILE" ]; then
  jq -e 'type == "object"' "$COMMUNITY_FILE" >"$tmp_dir/community.json" || {
    echo "Community file must contain a JSON object: $COMMUNITY_FILE" >&2
    exit 1
  }
fi

jq -n \
  --arg schema_version "1" \
  --arg captured_at "$captured_at" \
  --arg campaign "$CAMPAIGN" \
  --arg campaign_published_at "$CAMPAIGN_PUBLISHED_AT" \
  --arg org "$ORG" \
  --arg repo "$REPO" \
  --argjson dev_available "$dev_available" \
  --argjson traffic_available "$traffic_available" \
  --slurpfile repos "$tmp_dir/repos.json" \
  --slurpfile primary "$tmp_dir/primary.json" \
  --slurpfile growth "$tmp_dir/star-growth.json" \
  --slurpfile issues "$tmp_dir/issues.json" \
  --slurpfile pulls "$tmp_dir/pulls.json" \
  --slurpfile discussions "$tmp_dir/discussions.json" \
  --slurpfile dev_public "$tmp_dir/dev-public.json" \
  --slurpfile dev_owner "$tmp_dir/dev-owner.json" \
  --slurpfile views "$tmp_dir/traffic-views.json" \
  --slurpfile clones "$tmp_dir/traffic-clones.json" \
  --slurpfile referrers "$tmp_dir/traffic-referrers.json" \
  --slurpfile paths "$tmp_dir/traffic-paths.json" \
  --slurpfile community "$tmp_dir/community.json" '
  {
    schema_version: ($schema_version | tonumber),
    captured_at: $captured_at,
    campaign: {name: $campaign, published_at: $campaign_published_at},
    github: {
      organization: {
        login: $org,
        public_repositories: ($repos[0] | length),
        total_stars: ($repos[0] | map(.stargazers_count) | add),
        metadata_coverage: {
          descriptions: ($repos[0] | map(select((.description // "") != "")) | length),
          homepages: ($repos[0] | map(select((.homepage // "") != "")) | length),
          topics: ($repos[0] | map(select((.topics // []) | length > 0)) | length),
          detected_licenses: ($repos[0] | map(select(.license != null)) | length),
          complete: ($repos[0] | map(select(
            ((.description // "") != "") and
            ((.homepage // "") != "") and
            (((.topics // []) | length) > 0) and
            (.license != null)
          )) | length)
        },
        top_repositories: ($repos[0] | sort_by(-.stargazers_count) | .[0:10]
          | map({name, stars: .stargazers_count, forks: .forks_count, url: .html_url}))
      },
      primary_repository: {
        name: ($org + "/" + $repo),
        url: $primary[0].html_url,
        stars: $primary[0].stargazers_count,
        forks: $primary[0].forks_count,
        open_issues_and_pull_requests: $primary[0].open_issues_count,
        watchers: $primary[0].subscribers_count,
        star_growth: $growth[0],
        campaign_activity: {
          issues: $issues[0],
          pull_requests: $pulls[0],
          discussions: $discussions[0]
        },
        traffic: {
          available: $traffic_available,
          views_14d: $views[0],
          clones_14d: $clones[0],
          referrers_14d: $referrers[0],
          popular_paths_14d: $paths[0]
        }
      }
    },
    dev: {
      available: $dev_available,
      url: ($dev_public[0].url // null),
      published_at: ($dev_public[0].published_at // null),
      reactions: ($dev_public[0].public_reactions_count // null),
      comments: ($dev_public[0].comments_count // null),
      reading_time_minutes: ($dev_public[0].reading_time_minutes // null),
      page_views: ($dev_owner[0].page_views_count // null)
    },
    community: $community[0]
  }' >"$tmp_dir/snapshot.json"

jq -e '
  .schema_version == 1 and
  (.captured_at | fromdateiso8601 | type == "number") and
  (.github.organization.total_stars | type == "number") and
  (.github.primary_repository.stars | type == "number") and
  (.github.primary_repository.star_growth.last_30d | type == "number")
' "$tmp_dir/snapshot.json" >/dev/null

if [ -n "$OUTPUT" ]; then
  mkdir -p "$(dirname "$OUTPUT")"
  cp "$tmp_dir/snapshot.json" "$OUTPUT"
  echo "Wrote promotion metrics: $OUTPUT" >&2
else
  cat "$tmp_dir/snapshot.json"
fi
