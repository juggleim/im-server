# Promotion Metrics

Use `scripts/collect-promotion-metrics.sh` to create comparable snapshots for the
`juggleim_oss_launch` campaign. The collector reads public GitHub and DEV data and, when the
authenticated GitHub account has repository access, includes the rolling 14-day traffic window.

## Capture a snapshot

```bash
./scripts/collect-promotion-metrics.sh \
  --output ".promotion-metrics/$(date -u +%Y-%m-%dT%H%M%SZ).json"
```

The `.promotion-metrics/` directory is intentionally ignored by Git. Repository traffic is useful
for campaign decisions, but it does not need to be published in source control.

Compare any two retained snapshots:

```bash
./scripts/compare-promotion-metrics.sh \
  .promotion-metrics/launch.json \
  .promotion-metrics/24-hours.json
```

The comparison separates human pull requests from bot pull requests, so Dependabot updates are not
mistaken for community contribution growth. GitHub visitor and clone comparisons describe two
rolling 14-day windows rather than lifetime totals.

For public data only:

```bash
./scripts/collect-promotion-metrics.sh --public-only
```

The public DEV API reports reactions and comments. Article owners can also include page views
without putting a credential in a file or snapshot:

```bash
DEV_API_KEY="..." ./scripts/collect-promotion-metrics.sh --public-only
```

## Reddit and LinkedIn

Reddit and LinkedIn do not expose all post analytics through reliable public endpoints. Collect
their logged-in metrics as a JSON object and merge it into the same snapshot:

```json
{
  "reddit": [
    {
      "community": "r/golang",
      "url": "https://www.reddit.com/r/golang/comments/...",
      "score": 0,
      "comments": 0,
      "views": null,
      "status": "published"
    }
  ],
  "linkedin": {
    "url": "https://www.linkedin.com/posts/...",
    "impressions": 0,
    "reactions": 0,
    "comments": 0,
    "reposts": 0
  }
}
```

```bash
./scripts/collect-promotion-metrics.sh \
  --community-file .promotion-metrics/community.json \
  --output .promotion-metrics/combined.json
```

Do not put session cookies, API keys, email addresses, or private commenter data in the community
file.

## Checkpoints

Capture snapshots at these campaign-relative checkpoints:

- Launch baseline
- 24 hours
- 7 days
- 30 days
- Weekly thereafter while promotion remains active

Compare at least Stars, Forks, unique visitors, unique clones, DEV reactions/comments/views,
GitHub Discussion activity, organization metadata coverage, community engagement, new Issues, and
new pull requests. GitHub traffic is a rolling window, so retain each local snapshot rather than
expecting the API to provide long-term history.

Channel decisions should be evidence-based:

- Continue a channel when it produces qualified repository visits, Stars, useful comments, Issues,
  or pull requests.
- Change the angle when impressions are present but repository visits are absent.
- Stop repeating a promotional post when it is removed, attracts low-quality traffic, or conflicts
  with community rules.
- Convert recurring technical questions into documentation, examples, or tracked Issues.
