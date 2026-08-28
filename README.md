# 🍌 Banana Software SSH Terminal

> **`ssh bananasoftware.net`** — A simple, clean, terminal-native version of the Banana Software landing page.

Built with [Charm](https://charm.sh) **Wish**, **Bubble Tea**, and **Lipgloss**.

---

## 🌟 Highlights

- **🍌 Minimalist Code-Editor UI**: Matches the exact web landing page layout with line numbers, syntax highlights, and clean typography.
- **📦 Products Overview**:
  - `TaiHoaDon` (`taihoadon.online`) — Tải hóa đơn điện tử hàng loạt
  - `TaiToKhai` (`taitokhai.online`) — Tải hồ sơ đã nộp dichvucong và thuedientu
  - `KeToanONE` (`ketoan.one`) — Phần mềm kế toán
  - `HoaDonGoc` (`hoadongoc.com`) — Tìm và tải hóa đơn PDF gốc từ nhà cung cấp
  - `QuanLyHoaDon` (`quanlyhoadon.com`) — Quản lý hóa đơn
- **🌐 Multilingual**: One-key toggle between **VI** and **EN** (`l`).
- **⌨ Keyboard Navigation**: Arrow keys / `j`/`k`, `Enter` to open link, `q` to quit.

---

## 🚀 Quick Start

### 1. Connect via SSH

```bash
# Connect to production SSH server
ssh bananasoftware.net

# Connect to local test server
ssh localhost -p 2222
```

### 2. Run Local TUI Mode

```bash
cd ssh
go run ./cmd/banana-ssh -tui
# or
make tui
```

### 3. Run SSH Server

```bash
cd ssh
go run ./cmd/banana-ssh -port 2222
# or
make run
```

---

## ⌨ Keybindings

| Key | Action |
| --- | --- |
| `↑` / `↓` or `k` / `j` | Select product line |
| `Enter` | Open / Show product URL |
| `l` | Toggle language (`VI` ⇄ `EN`) |
| `q` or `Esc` / `Ctrl+C` | Quit |

---

## 🐳 Docker

```bash
cd ssh
docker compose up -d
```
