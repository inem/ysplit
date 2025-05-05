# ysplit

> ⛏️ *Split YAML-with-blank‑lines into an array under a single parent key.*

`ysplit` is a tiny Go CLI that takes a YAML file where **blank lines separate logical blocks** and converts it into one YAML (or JSON) document shaped like

```yaml
parent:
  - block‑1
  - block‑2
  - ...
```

Originally built to process deployment manifests like:

```yaml
test-docs:
  domains:
    - docs.test.com
    - internal-docs.test.com

  registry: https://gitlab.com/inem/docs/container_registry
  repo: https://github.com/inem/test-docs
  ...
```

which should become:

```yaml
test-docs:
  - domains:
      - docs.test.com
      - internal-docs.test.com
  - registry: https://gitlab.com/inem/docs/container_registry
    repo: https://github.com/inem/test-docs
    ...
```

---

## ✨ Features

* **Keeps a single top‑level key** – every blank‑separated chunk slides under it as a list item.
* **Two output formats** – YAML (default) or JSON (`-json`).
* **Streams** – read from `stdin`, write to `stdout`; easy to chain in pipes.
* **Zero deps** – only needs `gopkg.in/yaml.v3`.
* **Tiny** – single file, <200 LoC.

---

## 🚀 Install

```bash
# replace <user> with your GitHub user or use go install once repo is public
go install github.com/inem/ysplit@latest
```

Or clone and build locally:

```bash
git clone https://github.com/inem/ysplit.git
cd ysplit && go build -o ysplit ./cmd/ysplit
```

Requires Go 1.22+.

---

## 🔧 Usage

### Basic

```bash
cat cfg.yml | ysplit                 # YAML ➜ YAML
cat cfg.yml | ysplit -json > cfg.json  # YAML ➜ JSON
```

### Options

| Flag    | Default | Description                           |
| ------- | ------- | ------------------------------------- |
| `-json` | `false` | Emit a JSON document instead of YAML. |
| `-h`    |         | Show built‑in help.                   |

### Example workflow

```bash
# 1. Split into patches
ysplit -json < cfg.yml > blocks.json

# 2. Feed each block to your deploy script
jq -c '."test-docs"[]' blocks.json | while read blk; do
  deploy_one "$blk"
done
```

---

## 🧠 How it works (TL;DR)

1. Reads the first non‑empty, zero‑indent line → that’s the **parent key**.
2. Scans the rest line‑by‑line.
3. Every blank line **closes** the current block.
4. For each block, removes the parent’s indent (assumes 2‑space indent) so the chunk is valid YAML on its own.
5. Parses all blocks with `yaml.Decoder` → builds a `[]any` slice.
6. Wraps the slice under the parent key and encodes as YAML or JSON.

*Algorithm lives in `cmd/ysplit/main.go`; tweak indent/constants as you like.*

---

## 🚧 Limitations & TODO

* Assumes the parent block uses **2‑space** indent. Make it dynamic if you mix tabs/spaces.
* Comments and exact blank‑line counts inside a block are not preserved.
* Only handles a **single** top‑level key. Support for multi‑root docs is on the wishlist.

Feel free to open issues / PRs.

---

## ⚖️ License

MIT – do what you want but don’t blame me.
