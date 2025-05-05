# 🧹 NoxDir

**NoxDir** is a fast, cross-platform CLI tool for scanning your file system.
It detects mounted drives or volumes and displays space usage in a clean,
interactive terminal UI. Ideal for figuring out what’s eating up your disk space
and where to clean up first.

## 🚀 Features

- ✅ Detects drives/mount points on **Windows**, **macOS**, and **Linux**
- 📊 Shows disk usage (used, free, total, percentage usage) per drive
- 🖥️ Interactive terminal UI (keyboard navigation)
- ⚡ Fast and efficient scanning, using native system calls
- 🔒 No telemetry, no BS — just your drives and the data

## 📸 Preview

### Drives list

![The San Juan Mountains are beautiful!](/img/drives.png "drives list")

### Directories list

![The San Juan Mountains are beautiful!](/img/dirs.png "directories list")

## 📦 Installation

### Download Binaries

Grab the latest binary from the [Releases](https://github.com/crumbyte/noxdir/releases) page for your
platform. Run it from wherever you want.

### Or build from source (Go 1.24+)

```bash
git clone https://github.com/crumbyte/noxdir.git
cd noxdir
make build

./noxdir
```

## 🛠 Usage

Just run in the terminal:

```bash
noxdir
```

No flags, no fuss. It starts the interactive UI immediately.

## ⚙️ How It Works

- **Windows:** Uses `GetLogicalDrives` and `GetDiskFreeSpaceExW` through direct
  syscalls for optimal performance.
- **Linux/macOS:** Uses `statfs` and parses `/proc/mounts` or `mount` command
  output to find mounted volumes.

## 🧩 Planned Features

- [ ] Listen for FS event for rendering
- [ ] Dirs/files management (delete, rename, add, etc.)
- [ ] Sort directories by usage, free space, etc. (already done for
  drives)
- [ ] Theming / color customization

## 🙋 FAQ

- **Q:** Can I use this in scripts or headless environments?
- **A:** Not yet — it's designed for interactive use.


- **Q:** Is this safe to run?
- **A:** Yes — it’s strictly read-only and does not touch any files.


- **Q:** Can I delete dirs/files from the application?
- **A:** Not yet. Already in the roadmap.

- **Q:** I don't see the icons and everything looks ugly.
- **A:** It depends solely on your terminal's settings and fonts. Theming your
terminal application is another topic.

## 🧪 Contributing

Pull requests are welcome! If you’d like to add features or report bugs, please
open an issue first to discuss.

## 📝 License

MIT © [crumbyte](https://github.com/crumbyte)

---

> NoxDir is built with 💻 and ❤️ to help you take back control of your
> storage.
