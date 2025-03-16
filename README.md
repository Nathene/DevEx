# DevEx - Developer Experience Toolkit

## Overview

DevEx is a comprehensive desktop application designed to enhance the developer experience by providing a unified interface for monitoring and managing various development tools and resources. Built with Wails and Svelte, it offers a modern, responsive UI with native integration for macOS.

## Features

### System Monitoring
- **Real-time metrics**: Monitor CPU, RAM, and disk usage with live updates
- **Historical data**: View performance trends over time with interactive charts
- **Process management**: View running processes, their resource usage, and manage them

### Docker Integration
- **Container management**: View and manage Docker containers
- **Image management**: List and manage Docker images
- **Resource monitoring**: Track Docker's resource usage

### Git Repository Management
- **Repository tracking**: Add and manage local Git repositories
- **Status monitoring**: View repository status, branch information, and changes
- **Quick actions**: Open repositories in VS Code with a single click
- **Visual feedback**: See at-a-glance status with color-coded indicators

### Development Servers
- **Server management**: Start, stop, and monitor development servers
- **Quick access**: Open server URLs directly from the interface

### Database Connections
- **Connection management**: Store and manage database connections
- **Connection testing**: Test connections before use
- **Quick connect**: Connect to databases with a single click

### API Testing
- **Request builder**: Create and send API requests
- **Response viewer**: View formatted API responses
- **Request history**: Save and reuse API requests

## Planned Features

### Enhanced Git Integration
- Commit history visualization
- Branch management
- Pull request integration

### Expanded Docker Capabilities
- Container logs viewer
- Docker Compose integration
- Container health monitoring

### Development Workflow Improvements
- Project templates and scaffolding
- Integrated terminal
- Task automation

### Debugging Tools
- Log aggregation and analysis
- Error tracking
- Performance profiling

### CI/CD Integration
- Pipeline status monitoring
- Build and deployment triggers
- Test result visualization

## Getting Started

### Prerequisites
- Go 1.18 or later
- Node.js 14 or later
- Docker (for Docker features)
- Git (for repository management)
- VS Code (for repository opening feature)

### Installation
1. Clone the repository
2. Run `wails dev` for development mode
3. Run `wails build` to create a production build

## Development

To run in live development mode:
```bash
wails dev
```

This will run a Vite development server with hot reload for frontend changes. A dev server also runs on http://localhost:34115 for browser-based development.

## Building

To build a redistributable, production package:
```bash
wails build
```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues to improve the application.

## License

[MIT License](LICENSE)
