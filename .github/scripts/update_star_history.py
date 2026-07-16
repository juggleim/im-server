#!/usr/bin/env python3
"""Generate a self-contained star-history chart for the README."""

from __future__ import annotations

import json
import math
import os
import sys
import urllib.error
import urllib.request
from collections import Counter
from concurrent.futures import ThreadPoolExecutor
from datetime import date, datetime, timezone
from pathlib import Path


REPOSITORY = os.environ.get("GITHUB_REPOSITORY", "juggleim/im-server")
OUTPUT = Path(__file__).resolve().parents[2] / "docs" / "assets" / "star-history.svg"
WIDTH, HEIGHT = 900, 520
LEFT, RIGHT, TOP, BOTTOM = 76, 30, 62, 62
PLOT_WIDTH = WIDTH - LEFT - RIGHT
PLOT_HEIGHT = HEIGHT - TOP - BOTTOM


def github_get(path: str) -> tuple[object, dict[str, str]]:
    token = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")
    if not token:
        raise RuntimeError("GITHUB_TOKEN or GH_TOKEN is required")

    request = urllib.request.Request(
        f"https://api.github.com{path}",
        headers={
            "Accept": "application/vnd.github.star+json",
            "Authorization": f"Bearer {token}",
            "User-Agent": "juggleim-star-history-generator",
            "X-GitHub-Api-Version": "2022-11-28",
        },
    )
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            headers = {key.lower(): value for key, value in response.headers.items()}
            return json.load(response), headers
    except urllib.error.HTTPError as error:
        message = error.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"GitHub API returned HTTP {error.code}: {message}") from error


def fetch_stars() -> list[date]:
    repository, _ = github_get(f"/repos/{REPOSITORY}")
    if not isinstance(repository, dict):
        raise RuntimeError("Unexpected response from the GitHub repository API")
    page_count = math.ceil(repository["stargazers_count"] / 100)

    def fetch_page(page: int) -> list[dict[str, object]]:
        data, _ = github_get(f"/repos/{REPOSITORY}/stargazers?per_page=100&page={page}")
        if not isinstance(data, list):
            raise RuntimeError("Unexpected response from the GitHub stargazers API")
        return data

    with ThreadPoolExecutor(max_workers=8) as executor:
        pages = executor.map(fetch_page, range(1, page_count + 1))

    stars: list[date] = []
    for data in pages:
        for item in data:
            starred_at = item.get("starred_at")
            if starred_at:
                stars.append(datetime.fromisoformat(starred_at.replace("Z", "+00:00")).date())
    return sorted(stars)


def nice_ceiling(value: int) -> int:
    if value <= 10:
        return 10
    magnitude = 10 ** (len(str(value)) - 1)
    normalized = value / magnitude
    step = 1 if normalized <= 1 else 2 if normalized <= 2 else 5 if normalized <= 5 else 10
    return step * magnitude


def build_svg(stars: list[date]) -> str:
    if not stars:
        raise RuntimeError("No stargazer timestamps were returned")

    counts = Counter(stars)
    start, end = stars[0], max(stars[-1], datetime.now(timezone.utc).date())
    days = max((end - start).days, 1)
    total = 0
    points: list[tuple[date, int]] = []
    current = start
    while current <= end:
        total += counts[current]
        points.append((current, total))
        current = date.fromordinal(current.toordinal() + 1)

    y_max = nice_ceiling(total)

    def x(day: date) -> float:
        return LEFT + ((day - start).days / days) * PLOT_WIDTH

    def y(value: int) -> float:
        return TOP + PLOT_HEIGHT - (value / y_max) * PLOT_HEIGHT

    line_points = " ".join(f"{x(day):.1f},{y(value):.1f}" for day, value in points)
    area_points = (
        f"{LEFT},{TOP + PLOT_HEIGHT} {line_points} "
        f"{LEFT + PLOT_WIDTH},{TOP + PLOT_HEIGHT}"
    )

    grid: list[str] = []
    for index in range(6):
        value = round(y_max * index / 5)
        py = y(value)
        grid.append(
            f'<line class="grid" x1="{LEFT}" y1="{py:.1f}" x2="{LEFT + PLOT_WIDTH}" y2="{py:.1f}"/>'
            f'<text class="tick" x="{LEFT - 12}" y="{py + 4:.1f}" text-anchor="end">{value:,}</text>'
        )

    x_ticks: list[str] = []
    for index in range(6):
        tick_day = date.fromordinal(start.toordinal() + round(days * index / 5))
        px = x(tick_day)
        x_ticks.append(
            f'<line class="grid" x1="{px:.1f}" y1="{TOP}" x2="{px:.1f}" y2="{TOP + PLOT_HEIGHT}"/>'
            f'<text class="tick" x="{px:.1f}" y="{TOP + PLOT_HEIGHT + 28}" text-anchor="middle">{tick_day:%Y-%m}</text>'
        )

    updated = datetime.now(timezone.utc).strftime("%Y-%m-%d UTC")
    return f'''<svg xmlns="http://www.w3.org/2000/svg" width="{WIDTH}" height="{HEIGHT}" viewBox="0 0 {WIDTH} {HEIGHT}" role="img" aria-labelledby="title description">
  <title id="title">Star history for {REPOSITORY}</title>
  <desc id="description">{total:,} GitHub stars from {start:%Y-%m-%d} to {end:%Y-%m-%d}</desc>
  <style>
    .background {{ fill: #ffffff; }}
    .grid {{ stroke: #e5e7eb; stroke-width: 1; }}
    .tick {{ fill: #6b7280; font: 12px -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }}
    .title {{ fill: #111827; font: 600 22px -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }}
    .subtitle {{ fill: #6b7280; font: 13px -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }}
    .area {{ fill: #2563eb; fill-opacity: .10; }}
    .line {{ fill: none; stroke: #2563eb; stroke-width: 3; stroke-linejoin: round; stroke-linecap: round; }}
    .legend {{ fill: #374151; font: 600 13px -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; }}
    @media (prefers-color-scheme: dark) {{
      .background {{ fill: #0d1117; }} .grid {{ stroke: #30363d; }}
      .tick, .subtitle {{ fill: #8b949e; }} .title {{ fill: #f0f6fc; }} .legend {{ fill: #c9d1d9; }}
      .area {{ fill: #58a6ff; fill-opacity: .12; }} .line {{ stroke: #58a6ff; }}
    }}
  </style>
  <rect class="background" width="{WIDTH}" height="{HEIGHT}" rx="8"/>
  <text class="title" x="{LEFT}" y="32">Star History</text>
  <text class="subtitle" x="{WIDTH - RIGHT}" y="31" text-anchor="end">Updated {updated}</text>
  {''.join(grid)}
  {''.join(x_ticks)}
  <polygon class="area" points="{area_points}"/>
  <polyline class="line" points="{line_points}"/>
  <circle cx="{x(end):.1f}" cy="{y(total):.1f}" r="5" fill="#2563eb"/>
  <line x1="{LEFT + 8}" y1="{TOP + 18}" x2="{LEFT + 34}" y2="{TOP + 18}" class="line"/>
  <text class="legend" x="{LEFT + 42}" y="{TOP + 22}">{REPOSITORY} · {total:,} stars</text>
  <text class="subtitle" x="18" y="{TOP + PLOT_HEIGHT / 2}" text-anchor="middle" transform="rotate(-90 18 {TOP + PLOT_HEIGHT / 2})">GitHub Stars</text>
</svg>
'''


def main() -> int:
    stars = fetch_stars()
    OUTPUT.parent.mkdir(parents=True, exist_ok=True)
    OUTPUT.write_text(build_svg(stars), encoding="utf-8")
    print(f"Wrote {OUTPUT} with {len(stars):,} stars")
    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except RuntimeError as error:
        print(f"error: {error}", file=sys.stderr)
        sys.exit(1)
