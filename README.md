# 🔧 LazySys

A beautiful Terminal User Interface (TUI) for managing systemd services, built with Go and BubbleTea.

![Demo](assets/lazysys-vid.gif)




## ✨ Features

- **Split View Interface**: Two windows showing all services and running services
- **Real-time Service Management**: Start, stop, restart, enable, and disable services
- **Interactive Navigation**: Use Shift+H/Shift+L to navigate between windows
- **Search Functionality**: Search services by name or description
- **Loading Animations**: Beautiful spinner while loading services
- **Focus Indicators**: Clear visual feedback for active windows
- **Help System**: Built-in help with keybindings
- **About Window**: Coffee donation information with ASCII art
- **Service Counts**: Display total and running service counts
- **Status Icons**: Visual indicators for service status (🟢 running, 🔴 failed, ⚪ inactive)

## 🚀 Installation

### Prerequisites

- Go 1.21 or later
- Linux system with systemd
- Sudo privileges (required for service management)

### Build from Source

```bash
# Clone or download the source
cd lazysys
# Build the application
make build
# Run (requires sudo)
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

## ☕ Support

If you find this tool helpful, consider buying me a coffee! ☕

- **Bitcoin**: `bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh`
- **Ethereum**: `0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6`
- **Ko-fi**: https://ko-fi.com/lazysys

## 📝 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## ⚠️ Disclaimer

This tool requires sudo privileges to manage systemd services. Use with caution and ensure you understand the implications of starting, stopping, or modifying system services. 
