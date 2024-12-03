# pareto-linux
Automatically audit your Linux machine for basic security hygiene.




## Installation

### Using Nix

The recommended way to install Pareto Linux is through the Nix package manager:

```bash
nix profile install --accept-flake-config github:paretosecurity/pareto-linux
```

#### Quick Start

To run a one-time security audit without installation:

```bash
nix run github:paretosecurity/pareto-linux -- check
```

or if running from local repo with JSON reporter

```bash
nix run . -- check --json
```

This will analyze your system and provide a security report highlighting potential improvements and vulnerabilities.

### Using Debian/Ubuntu/Pop_OS!

See [https://apt.paretosecurity.com](https://apt.paretosecurity.com) for install steps.


#### Quick Start

To run a one-time security audit:

```bash
paretosecurity check
```

or with JSON reporter

```bash
paretosecurity check --json
```
