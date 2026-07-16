# JuggleIM English Promotion Publishing Guide

This directory contains the source article, channel-specific launch copy, and final image assets for the English JuggleIM promotion campaign.

Nothing in this directory publishes content automatically. A maintainer must review each post, sign in to the target platform, and submit it manually or explicitly authorize an automated publication.

## Source of truth

- Article: [`introducing-juggleim.md`](./introducing-juggleim.md)
- Campaign: `juggleim_oss_launch`
- Primary repository: `https://github.com/juggleim/im-server`

The article was published on DEV Community on 2026-07-16:

https://dev.to/yuwnloyblog/juggleim-an-open-source-self-hosted-messaging-backend-built-in-go-43nh

## Assets

| Asset                                                                               |     Size | Use                                                                     |
| ----------------------------------------------------------------------------------- | -------: | ----------------------------------------------------------------------- |
| [`juggleim-cover-1600x840.jpg`](./assets/juggleim-cover-1600x840.jpg)               | 1600×840 | General article and Hashnode cover                                      |
| [`juggleim-cover-social-1200x630.jpg`](./assets/juggleim-cover-social-1200x630.jpg) | 1200×630 | LinkedIn, X, and generic Open Graph sharing                             |
| [`juggleim-cover-devto.jpg`](./assets/juggleim-cover-devto.jpg)                     | 1000×420 | DEV Community cover                                                     |
| [`system-overview.png`](./assets/system-overview.png)                               | 1584×747 | Raster architecture image for platforms that do not render SVG reliably |

The covers are deterministic crops of the approved JuggleIM social preview. They preserve the approved logo, copy, palette, and network visual.

## Canonical URL strategy

Use one canonical source to avoid splitting search authority:

1. If an official `juggle.im` blog is available, publish there first and set that URL as the canonical URL everywhere else.
2. The initial campaign uses DEV Community as the canonical source because no official blog URL was selected before publication.
3. When syndicating to Hashnode or another platform, set its canonical URL to the primary article.
4. Do not publish identical copies on several platforms without a canonical URL.

Before publication, add `canonical_url` to the article front matter only when the official source URL exists.

## UTM convention

Use the following campaign name everywhere:

```text
utm_campaign=juggleim_oss_launch
```

| Channel           | Source              | Medium      |
| ----------------- | ------------------- | ----------- |
| DEV Community     | `devto`             | `article`   |
| Hacker News       | `hackernews`        | `community` |
| Reddit selfhosted | `reddit_selfhosted` | `community` |
| Reddit Go         | `reddit_golang`     | `community` |
| LinkedIn          | `linkedin`          | `social`    |
| X                 | `x`                 | `social`    |

Example:

```text
https://github.com/juggleim/im-server?utm_source=devto&utm_medium=article&utm_campaign=juggleim_oss_launch
```

## Pre-publication checklist

- [ ] Run the Docker Quick Start from a clean environment.
- [ ] Confirm ports `9001`, `9002`, `9003`, and `8090` are reachable.
- [ ] Change or clearly label all development-only credentials.
- [ ] Recheck GitHub Stars, Forks, latest release, and links on publication day.
- [ ] Preview the article on the target platform for table, code block, and image rendering.
- [ ] Confirm the cover is not cropped around the logo or product name.
- [ ] Confirm the architecture PNG loads without GitHub authentication.
- [ ] Add the final canonical URL when applicable.
- [x] Set `published: true` after maintainer approval.
- [ ] Prepare at least one maintainer to answer comments for the first two hours.

Update repository metrics with:

```bash
gh repo view juggleim/im-server --json stargazerCount,forkCount,latestRelease
```

## Recommended sequence

### Day 0

1. Publish the canonical article on the official blog or DEV Community.
2. Verify every image and link in the public post.
3. Publish the LinkedIn and X announcements.

### Day 1

1. Submit the project to Hacker News using the prepared Show HN copy.
2. Add the technical context as the first comment.
3. Stay available to answer questions and acknowledge limitations directly.

### Day 2–3

1. Post the self-hosting version to `r/selfhosted` after checking its current rules.
2. Post the Go architecture version to `r/golang` after checking its current rules.
3. Do not cross-post identical bodies; use the channel-specific drafts.

### After the initial launch

1. Syndicate the article to Hashnode with the canonical URL.
2. Publish a Chinese adaptation separately rather than machine-translating the English post inline.
3. Convert recurring questions into documentation or FAQ updates.

## Measurement

Record metrics immediately before launch, then after 24 hours, 7 days, and 30 days.

| Metric                       | Baseline | 24 hours | 7 days | 30 days |
| ---------------------------- | -------: | -------: | -----: | ------: |
| GitHub Stars                 |    3,584 |          |        |         |
| GitHub Forks                 |      363 |          |        |         |
| Unique repository visitors   |          |          |        |         |
| Unique clones                |          |          |        |         |
| Article views                |          |          |        |         |
| Discussions / Issues created |          |          |        |         |
| New contributors / PRs       |          |          |        |         |

GitHub traffic data is available to repository maintainers under Insights → Traffic or through the traffic API.
Use [`METRICS.md`](./METRICS.md) and `scripts/collect-promotion-metrics.sh` to retain comparable
local snapshots, including the rolling traffic window and optional Reddit/LinkedIn analytics.

## Response principles

- Answer technical questions with links to code or documentation.
- State the single-node community-edition boundary clearly.
- Do not repeat performance claims without published benchmark conditions.
- Do not claim end-to-end encryption unless a documented client and key-management design supports it.
- Thank users for concrete criticism and convert reproducible problems into Issues.
- Avoid arguing about competitor comparisons; focus on JuggleIM's actual design and use cases.
