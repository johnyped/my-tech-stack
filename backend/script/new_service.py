#!/usr/bin/env python3
"""Generate a new service from template-service/, driven by the checklist
in document/SCAFFOLD_CHECKLIST.md (both relative to backend/, where this
script lives). Called by script/new-service.sh and by `make new-service`;
the Kiro steering flow (.kiro/steering/scaffold/, at the workspace root)
also ends by invoking this same script, so there is exactly one
implementation of "how a new service gets assembled."

Usage:
    new_service.py --name my-service \\
        --module gitlab.awesome-poc-th.com/libero-engineering/core/my-service \\
        [--port 8080] [--out PATH] [--force] \\
        [--cache redis|valkey|no] [--db postgres|no] [--kafka yes|no] \\
        [--storage s3|no] [--firebase yes|no] [--mockoon yes|no] [--mocktcp yes|no]
"""
import argparse
import re
import shutil
import sys
from pathlib import Path

SENTINEL_START = re.compile(r"^(?P<indent>\s*)(?:#|//)\s*>>>\s*(?P<name>\w+)\s*$")
SENTINEL_END = re.compile(r"^(?P<indent>\s*)(?:#|//)\s*<<<\s*(?P<name>\w+)\s*$")

TEXT_SUFFIXES = {
    ".go", ".yml", ".yaml", ".md", ".mod", ".sh", ".json", ".txt", ".py",
}


def parse_args():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--name", required=True, help="kebab-case service name")
    p.add_argument("--module", required=True, help="full gitlab module path")
    p.add_argument("--port", default="8080")
    p.add_argument("--out", default=None, help="output dir, defaults to a sibling of this repo")
    p.add_argument("--force", action="store_true", help="overwrite --out if it already exists")
    p.add_argument("--cache", default="no", choices=["redis", "valkey", "no"])
    p.add_argument("--db", default="no", choices=["postgres", "no"])
    p.add_argument("--kafka", default="no", choices=["yes", "no"])
    p.add_argument("--storage", default="no", choices=["s3", "no"])
    p.add_argument("--firebase", default="no", choices=["yes", "no"])
    p.add_argument("--mockoon", default="no", choices=["yes", "no"])
    p.add_argument("--mocktcp", default="no", choices=["yes", "no"])
    return p.parse_args()


def toggle_sentinel_blocks(text: str, enabled: dict) -> str:
    """For each >>> NAME ... <<< NAME block, uncomment its lines if
    enabled[NAME] is True, otherwise drop the block entirely. The sentinel
    lines themselves are always removed."""
    out_lines = []
    lines = text.split("\n")
    i = 0
    while i < len(lines):
        m = SENTINEL_START.match(lines[i])
        if not m:
            out_lines.append(lines[i])
            i += 1
            continue
        name = m.group("name")
        i += 1
        block = []
        while i < len(lines) and not SENTINEL_END.match(lines[i]):
            block.append(lines[i])
            i += 1
        i += 1  # skip the <<< NAME line
        if enabled.get(name):
            for line in block:
                # strip exactly one leading "# " / "// " comment leader (or
                # just the bare marker at end of indent), whichever this
                # file uses
                stripped = re.sub(r"^(\s*)(?:#|//)\s?", r"\1", line)
                out_lines.append(stripped)
        # else: drop the block entirely (both yes/no and leave-commented
        # cases collapse to "don't keep dead config" per the earlier
        # decision — see document/CICD_GUIDE.md style notes)
    return "\n".join(out_lines)


def process_file(path: Path, replacements: list, enabled: dict):
    if path.suffix not in TEXT_SUFFIXES:
        return
    try:
        text = path.read_text(encoding="utf-8")
    except UnicodeDecodeError:
        return
    text = toggle_sentinel_blocks(text, enabled)
    for old, new in replacements:
        text = text.replace(old, new)
    path.write_text(text, encoding="utf-8")


def main():
    args = parse_args()
    # This script lives at <workspace>/backend/script/new_service.py.
    # backend_root is <workspace>/backend/ (contains template-service/);
    # the default output dir is a sibling of <workspace> itself — i.e.
    # <workspace>/../<name> — matching where the other service repos
    # (backend-for-frontend, push-notification, ...) actually live.
    backend_root = Path(__file__).resolve().parent.parent
    template_dir = backend_root / "template-service"
    out_dir = Path(args.out).resolve() if args.out else (backend_root.parent.parent / args.name)

    if out_dir.exists():
        if not args.force:
            print(f"error: {out_dir} already exists (use --force to overwrite)", file=sys.stderr)
            sys.exit(1)
        shutil.rmtree(out_dir)

    shutil.copytree(template_dir, out_dir)

    enabled = {
        "CACHE": args.cache in ("redis", "valkey"),
        "DB": args.db == "postgres",
        "KAFKA": args.kafka == "yes",
        "STORAGE": args.storage == "s3",
        "FIREBASE": args.firebase == "yes",
        "MOCKOON": args.mockoon == "yes",
    }

    replacements = [
        (f"gitlab.awesome-poc-th.com/libero-engineering/TODO-group/TODO-service-name", args.module),
        ("TODO-service-name", args.name),
    ]
    if args.port != "8080":
        replacements.append(("8080", args.port))

    for path in out_dir.rglob("*"):
        if path.is_file():
            process_file(path, replacements, enabled)

    if args.storage == "no":
        shutil.rmtree(out_dir / "internal" / "storage", ignore_errors=True)
    if args.mocktcp != "yes":
        shutil.rmtree(out_dir / "mocktcp", ignore_errors=True)
    if args.mockoon != "yes":
        shutil.rmtree(out_dir / "test" / "integration", ignore_errors=True)

    print_next_steps(out_dir, args)


def print_next_steps(out_dir: Path, args):
    lines = [
        "",
        f"Generated {out_dir}",
        "",
        "Next steps:",
        "  cd " + str(out_dir),
    ]
    go_gets = []
    if args.cache in ("redis", "valkey"):
        go_gets.append("github.com/redis/go-redis/v9@v9.17.1")
    if args.db == "postgres":
        go_gets.append("github.com/jackc/pgx/v5@v5.11.0")
    if args.kafka == "yes":
        go_gets.append("github.com/segmentio/kafka-go@v0.4.50")
    if args.storage == "s3":
        go_gets.append("github.com/aws/aws-sdk-go-v2/service/s3@v1.96.0")
    if args.firebase == "yes":
        go_gets.append("firebase.google.com/go/v4@v4.23.0")
    if go_gets:
        lines.append("  go get " + " ".join(go_gets))
    lines.append("  make vendor")
    lines.append("  make setup && make config-init && make run")
    lines.append("")
    lines.append("  curl localhost:" + args.port + "/health/liveness   # expect 200")
    lines.append("")
    lines.append("Once you have real routes, load skill::bruno-service-docs and")
    lines.append("skill::bruno-integration-test to replace the bruno/ skeleton and the")
    lines.append("placeholder scripts in script/integration/ with the real thing.")
    print("\n".join(lines))


if __name__ == "__main__":
    main()
