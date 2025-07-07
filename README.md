# LazySys

A beautiful Terminal User Interface (TUI) for managing systemd services, built with Go and BubbleTea.

![Demo](assets/lazysys-vid.gif)
## ✨ Features

- **TUI for systemd**: Manage services from your terminal  
- **Split View**: See all and active services side-by-side  
- **Service Control**: Start, stop, restart, enable, disable  
- **Fast Navigation**: Keyboard-driven workflow  
- **Search**: Filter by name or description

## 🚀 Installation

### For Archlinux users

```bash
yay -S lazysys-git
```
### Build from Source

```bash
# Clone the repo
cd lazysys
# Build the application
make build
# Run
make run
```

## 🎮 Usage

### Keybindings

| Key | Action |
|-----|--------|
| `Shift+H` / `Shift+L` | Navigate between windows |
| `j` / `k` | Navigate up/down in lists |
| `Number` | Select service for action |
| `s` | Search services |
| `?` | Toggle help |
| `P` | Show about |
| `q` / `Ctrl+C` | Quit |

### Service Actions

**All Services Window:**
- `1` - Start service
- `2` - Restart service
- `3` - Stop service
- `4` - Disable service
- `5` - Enable service
## 🤝 Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## ⚠️ Disclaimer

This tool requires sudo privileges to manage systemd services. Use with caution and ensure you understand the implications of starting, stopping, or modifying system services. 
