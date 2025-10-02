# upd8

![upd8 logo](upd8_logo.png)

> Universal package manager update checker — one CLI to rule them all.  

`upd8` scans your system for supported package managers (npm, pip, cargo, brew, snap, flatpak, …), lists outdated packages, and shows you a **single one-liner** to update each.  

Fast, minimal, cross-platform. Written in Go/Rust.

---

## ✨ Features
- 🔍 Auto-detects installed package managers  
- 📋 Lists outdated binaries in a clean table  
- ⚡ Shows update command for each (copy/paste friendly)  
- 🛠️ Runs as a single static binary (no runtime deps)  
- ⏰ `--watch` daemon mode → daily update summary  

---

## 🚀 Installation

### Go
```bash
go install github.com/makalin/upd8@latest
````

### Rust (cargo)

```bash
cargo install upd8
```

### Prebuilt Binary

Download from [Releases](https://github.com/makalin/upd8/releases) and place in `$PATH`.

---

## 🖥️ Usage

Check all updates:

```bash
upd8
```

Watch mode (daily summary):

```bash
upd8 --watch
```

Sample output:

```
📦 npm      5 outdated  →  npm update -g
📦 pip      3 outdated  →  pip install --upgrade -r requirements.txt
📦 brew     7 outdated  →  brew upgrade
📦 cargo    2 outdated  →  cargo install-update -a
📦 flatpak  4 outdated  →  flatpak update
```


### Flags

- `--packages` — include a short list of outdated packages for each manager.
- `--verbose` — show package managers even when no updates are available.
- `--watch` — keep running and re-scan on an interval (default 24h).
- `--interval=<duration>` — change the watch interval (e.g. `--interval=1h`).
- `--no-color` — disable ANSI colors in the output.

---

## ⚙️ Roadmap

* [ ] Add Windows support (choco, winget, scoop)
* [ ] Config file for custom commands
* [ ] JSON/YAML output for automation
* [ ] Notification hooks (Slack, Discord, Email)

---

## 🛡️ License

MIT © 2025 [Mehmet T. AKALIN](https://github.com/makalin)

---

## 🏷️ Badges

![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?style=for-the-badge\&logo=go\&logoColor=white)
![Rust](https://img.shields.io/badge/Rust-%23000000.svg?style=for-the-badge\&logo=rust\&logoColor=white)
![Cross-Platform](https://img.shields.io/badge/OS-Linux%20%7C%20macOS%20%7C%20Windows-blue?style=for-the-badge)
