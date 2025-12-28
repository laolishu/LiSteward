LiSteward — Development Environment Switcher
=========================================

Overview
--------
LiSteward is a lightweight development-environment configuration switch tool. It provides a simple GUI to visualize and switch `hosts` entries and manage Node versions (NVM) profiles. Future versions will integrate AI-assisted features for smarter suggestions and automation.

Key features
- Visual hosts editor and quick switching
- NVM/profile visualization and switching
- Local backups and profile import/export
- Cross-platform desktop app (Windows / macOS / Linux)

Quick start
1. Explore the repository for build and packaging scripts: see `frontend/` and `build/`.
2. Run the desktop app builds according to platform-specific instructions in `build/`.

Build & Run (development)
- Start the frontend dev server and run the Go backend for quick iteration:

```bash
# Terminal A (frontend dev server)
cd frontend
npm install
npm run dev

# Terminal B (backend)
cd ..
go run .
```

Build (release)
- Produce production frontend assets and build native app (requires `wails` CLI for packaging):

```bash
cd frontend
npm install
npm run build

# from repository root
wails build
# or build backend binary only
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
