# kubectl-abriment Plugin

A kubectl plugin for authenticating with the Abriment backend service and automatically configuring your kubeconfig file.

## Overview

The `kubectl-abriment` plugin simplifies the process of connecting to Kubernetes clusters managed by Abriment. It handles authentication, retrieves your personalized kubeconfig, and seamlessly integrates it with your existing kubectl configuration.

## Installation

### Prerequisites

- kubectl installed and configured
- Go 1.19+ (if building from source)

### Install from Release (Recommended)

1. Download the latest release for your platform from the [releases page](https://github.com/abrimentcloud/kubectl-abriment-plugin/releases)
2. Extract the binary and place it in your PATH
3. Rename the binary to `kubectl-abriment` (or `kubectl-abriment.exe` on Windows)
4. Make it executable (Unix/Linux/macOS): `chmod +x kubectl-abriment`

### Build from Source

```bash
git clone https://github.com/abriemntcloud/kubectl-abriment-plugin.git
cd kubectl-abriment-plugin
go build -o kubectl-abriment
```

Move the binary to a directory in your PATH, such as `/usr/local/bin` (Unix/Linux/macOS) or a directory in your Windows PATH.

### Verify Installation

```bash
kubectl abriment help
```

## Usage

### Quick Start

The simplest way to get started is with interactive mode:

```bash
kubectl abriment
```

This will guide you through the authentication process step by step.

### Command Line Usage

#### Username and Password Authentication

```bash
kubectl abriment login -u your-username -p your-password
```

#### Token Authentication

```bash
kubectl abriment login -t your-authentication-token
```

#### Dry Run (Preview Only)

```bash
kubectl login -u your-username -p your-password --dry-run client
```

#### Logout

```bash
kubectl abriment logout
```

### Available Commands

| Command | Description |
|---------|-------------|
| `kubectl abriment` | Interactive mode with guided prompts |
| `kubectl abriment login` | Main login command with flags |
| `kubectl abriment logout` | Main logout command with flags |
| `kubectl abriment help` | Display detailed help information |

### Login Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--username` | `-u` | Your username for authentication |
| `--password` | `-p` | Your password for authentication |
| `--token` | `-t` | Your authentication token |
| `--dry-run` | | Options: `client` (prints config without saving) |

## Configuration

### Environment Variables

You can customize the plugin behavior using these environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `LOGIN_ENDPOINT` | `https://backend.abriment.com/dashboard/api/login/` | Backend login endpoint |
| `CONFIG_ENDPOINT` | `https://backend.abriment.com/api/v1/paas/kubeconfig/` | Backend config endpoint |
| `KUBECONFIG` | `~/.kube/config` | Custom path for kubeconfig file |

### Example Environment Setup

```bash
# Custom endpoints
export LOGIN_ENDPOINT="https://custom.backend.com/api/login/"
export CONFIG_ENDPOINT="https://custom.backend.com/api/v1/kubeconfig/"

# Custom kubeconfig location
export KUBECONFIG="/path/to/my/kubeconfig"
```

## How It Works

1. **Authentication**: The plugin authenticates with the Abriment backend using your credentials
2. **Token Retrieval**: Upon successful authentication, it receives an authentication token
3. **Config Retrieval**: Uses the token to fetch your personalized kubeconfig from the backend
4. **Config Merging**: Intelligently merges the new configuration with your existing kubeconfig
5. **Resource Addition**: Adds the following resources to your kubeconfig:
   - **Cluster**: `abriment-cluster`
   - **Context**: `abriment-context`
   - **User**: `abriment-user`

### Config Merging Behavior

- If no kubeconfig exists, creates a new one
- If kubeconfig exists, preserves all existing configurations
- Only adds/updates Abriment-specific resources
- Never removes or modifies existing clusters, contexts, or users

## Examples

### Basic Login

```bash
# Using username and password
kubectl abriment login -u john.doe -p mypassword

# Using token
kubectl abriment login -t eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Interactive Mode

```bash
kubectl abriment
```

### Preview Configuration (Dry Run)

```bash
kubectl abriment login -u john.doe -p mypassword --dry-run client
```

This will display the kubeconfig content without saving it to disk.

## Switching to Abriment Context

After successful login, switch to the Abriment context:

```bash
kubectl config use-context abriment-context
```

Verify the connection:

```bash
kubectl cluster-info
kubectl get nodes
```

## Troubleshooting

### Common Issues

#### Authentication Failed
- Verify your credentials are correct
- Check if the login endpoint is accessible
- Ensure your account is not locked or disabled

#### Config Retrieval Failed
- Verify your token is valid and not expired
- Check if the config endpoint is accessible
- Ensure you have permissions to access kubeconfig

#### Permission Denied When Saving
- Check write permissions on the kubeconfig directory
- Verify the KUBECONFIG path is correct
- Try running with appropriate permissions

#### Plugin Not Found
- Ensure the binary is named `kubectl-abriment` (or `kubectl-abriment.exe` on Windows)
- Verify the binary is in your PATH
- Check that the binary has execute permissions

### Debug Mode

For detailed error information, you can examine the HTTP responses by modifying the source code to include debug logging.

### Getting Help

```bash
kubectl abriment help
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
