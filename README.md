**[English](README.md) | [中文](README_CN.md)**

---

LiSteward — Development Environment Switcher
=========================================

Overview
--------
LiSteward is a lightweight development-environment configuration switch tool. It provides a simple GUI to visualize and switch `hosts` entries and manage Node versions (NVM) profiles. It streamlines development workflows by centralizing configuration management for multiple environments. Future versions will integrate AI-assisted features for smarter suggestions and automation.

Features and Capabilities
----------

### Hosts Configuration Management
LiSteward provides a comprehensive solution for managing your system's hosts file:

- **Visual hosts editor** - Edit hosts entries with a clean, user-friendly interface instead of manually editing text files
- **Configuration profiles** - Create and manage multiple hosts configurations for different environments, switch between them with a single click
- **Backup and restore** - Automatic backups protect your configurations; restore to any previous version instantly
- **Import and export** - Share configurations across machines or version control your setups
- **Enable/disable entries** - Toggle individual hosts entries on/off without deleting them
- **Search and filter** - Quickly find specific entries in large hosts lists

#### Use Cases
- Front-end development: Manage local domain mappings and test environment configurations
- Multi-environment management: Maintain separate hosts configurations for development, staging, and production
- Network management: Quickly switch between different network setups and firewall rules

### Node Version Management (NVM)
Effortlessly manage multiple Node.js versions on your system:

- **Node.js version detection** - Automatically detects installed versions and displays the currently active one
- **Quick version switching** - Switch between installed Node.js versions with a single click; changes take effect immediately
- **Interactive version installation** - Quick access to see available Node.js versions and start installation
- **Environment management** - Supports both traditional nvm (Unix/Linux/macOS) and nvm-windows environments
- **Cross-platform support** - Works seamlessly on Windows, macOS, and Linux systems

#### Use Cases
- Multi-project development: Different projects may require different Node.js versions
- Compatibility testing: Quickly test your code against multiple Node.js versions
- Version upgrade management: Smoothly transition from older to newer Node.js versions

Getting Started
----------

### Requirements
- **Go**: 1.23 or higher
- **Node.js**: 16 or higher (for frontend development)
- **Wails CLI**: Required for building the native application (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Installation

Download the latest release for your platform from the [Releases](../../releases) page, or build from source (see below).

### First Run

1. Launch the application
2. The main interface displays two modules:
   - **Hosts Manager**: Left sidebar shows configuration profiles; main area allows editing hosts entries
   - **Node Manager**: Switch between installed Node.js versions and access version management tools
3. Explore both modules to understand available features

Build & Run (development)
----------
Start the frontend dev server and run the Go backend for quick iteration:

**Terminal A** (frontend dev server):
```bash
cd frontend
npm install
npm run dev
```

**Terminal B** (backend):
```bash
cd ..
go run .
```

Then open http://localhost:34115 in your browser (should open automatically).

Build (release)
----------
Produce production frontend assets and build native app (requires `wails` CLI):

**Prerequisites**:
- All requirements from "Getting Started" section above
- `wails` CLI installed globally

**Build steps**:
```bash
cd frontend
npm install
npm run build

# from repository root
wails build
```

The compiled application will be in `build/bin/` directory for your platform.

**Build backend binary only** (without native packaging):
```bash
go build -o LiSteward
```

License
-------
This project is licensed under the GNU General Public License v3.0 (GPL-3.0). See `LICENSE` for details.

Support
-------
If LiSteward helps you, you can support continued development with a one-time donation:

![Donate](#file:wechat-qr.png)

Contributing
------------
Welcome contributions via Issues and Pull Requests. Please keep changes focused and include tests or verification steps where appropriate.
