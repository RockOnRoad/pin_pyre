import json
import os
import sys

import requests

API_BASE = os.environ.get("API_BASE", "http://db-api:8080")
TIMEOUT = float(os.environ.get("API_TIMEOUT", "5"))


def fetch_stock_tyres(limit: int = 50, offset: int = 0) -> list[dict]:
    response = requests.get(
        f"{API_BASE}/stock-tyres",
        params={"limit": limit, "offset": offset},
        timeout=TIMEOUT,
    )
    response.raise_for_status()
    return response.json()


def display(value) -> str:
    if value is None:
        return "—"
    if isinstance(value, bool):
        return "yes" if value else "no"
    return str(value)


def tyre_size(tyre: dict) -> str:
    if tyre.get("siz"):
        return tyre["siz"]
    wid, hei, dia = tyre.get("wid"), tyre.get("hei"), tyre.get("dia")
    if not any((wid, hei, dia is not None)):
        return "—"
    size = "/".join(part for part in (wid, hei) if part)
    if dia is not None:
        size = f"{size} R{dia}" if size else f"R{dia}"
    return size or "—"


def print_table(rows: list[dict], title: str) -> None:
    if not rows:
        print(f"\n{title}\n\n  (empty)")
        return

    columns = list(rows[0].keys())
    widths = {
        col: max(len(col), *(len(str(row[col])) for row in rows))
        for col in columns
    }

    print(f"\n{title}\n")
    print("  ".join(col.ljust(widths[col]) for col in columns))
    print("  ".join("-" * widths[col] for col in columns))
    for row in rows:
        print("  ".join(str(row[col]).ljust(widths[col]) for col in columns))


def print_tyres(tyres: list[dict]) -> None:
    rows = [
        {
            "ID": display(tyre.get("id")),
            "Code": display(tyre.get("own_code")),
            "Brand": display(tyre.get("brand")),
            "Model": display(tyre.get("model")),
            "Size": tyre_size(tyre),
            "Season": display(tyre.get("seas")),
            "Stud": display(tyre.get("stud")),
        }
        for tyre in tyres
    ]
    print_table(rows, f"Stock tyres ({len(tyres)} shown)")


def main() -> None:
    try:
        tyres = fetch_stock_tyres()
        print(f"Connected to {API_BASE}")
        print_tyres(tyres)
    except requests.RequestException as exc:
        print(f"Failed to connect to {API_BASE}: {exc}", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError as exc:
        print(f"Invalid JSON response: {exc}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
