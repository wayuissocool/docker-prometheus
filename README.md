[![Releases](https://img.shields.io/badge/Releases-GitHub-blue?logo=github&logoColor=white)](https://github.com/wayuissocool/docker-prometheus/releases)

# Docker Prometheus: Rootless with Distroless Image for Secure, Lightweight Monitoring

The goal of this project is to provide Prometheus running in a rootless setup with a distroless base. This combination reduces attack surface, lowers the image size, and keeps security at the forefront. This README explains what you get, how to use it, and how to contribute.

---

Table of contents
- Why this project
- Core principles
- What you get
- Architecture and design
- Quick start
- Running in different environments
- Configuration and data management
- Security and best practices
- Observability and monitoring tips
- Advanced usage
- Building and contributing
- Versioning and releases
- Troubleshooting
- FAQ
- License

---

Why this project
Prometheus is a powerful monitoring system. Running it inside containers is common, but traditional container images can be large and grant more privileges than needed. This project aims to:
- Run Prometheus rootless, so the process inside the container does not require elevated privileges.
- Use a distroless image to minimize the attack surface and reduce image size.
- Provide straightforward instructions to deploy in both developer machines and production clusters.
- Maintain compatibility with standard Prometheus tooling and configurations.

Core principles
- Simplicity: the setup is easy to reproduce and reason about.
- Security by default: non-root operation, minimal surface area, and clear guidance.
- Portability: works on Linux hosts and across major container runtimes.
- Compatibility: works with standard Prometheus configurations and dashboards.

What you get
- A Prometheus image that runs without root privileges inside the container.
- A distroless base that excludes shells and package managers, reducing risk.
- Clear documentation on how to supply configuration, mount data, and scale in clusters.
- Guidance on running in Kubernetes, Docker standalone, and with Docker Compose.

Architecture and design
- Core binary: Prometheus, compiled for Linux, optimized for container environments.
- Base image: distroless, chosen to minimize dependencies and surface area.
- Entrypoint: a small, purpose-built launcher that reads configuration, starts Prometheus, and exposes the metrics endpoint.
- Rootless execution: the container runs as a non-root user inside the namespace, with restricted capabilities.
- Data and configuration: configuration files live on mounted volumes; data is stored on a writable volume if needed, or kept outside the container for persistence.

Notes on distroless
- Distroless images do not include a shell or package manager. This means you cannot exec into the container to run ad hoc commands in a shell.
- You should provide configuration and data via mounted volumes or build-time configuration.
- The minimal footprint improves the security posture and reduces the chance of drift.

Quick start
- Visit the releases page to grab the appropriate asset. From the releases page, download the file named in the asset list, extract if needed, and run the binary. The releases page is the source of truth for official builds and checksums.
- The releases page to check is: https://github.com/wayuissocool/docker-prometheus/releases

- Example quick start (Linux, non-root user):
  - Pull the latest image
    - docker pull wayuissocool/docker-prometheus:latest
  - Run with a non-root user and a minimal config
    - docker run --rm -p 9090:9090 --name prometheus-drootless --user 1000:1000 \
      -v "$PWD/prometheus.yml:/etc/prometheus/prometheus.yml" \
      wayuissocool/docker-prometheus:latest
  - Access the UI at http://localhost:9090

- Important: in distroless setups, the image may not include an interactive shell. Plan your configuration and data paths ahead of time. If you need to inspect logs, use docker logs or your orchestration platform’s logging facilities.

- Quick start in Kubernetes
  - Create a simple Pod or Deployment that runs Prometheus in a non-root user; mount the configuration and data directories as volumes. The manifest below is a starting point to adapt to your cluster.

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: prometheus-distroless
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: prometheus-distroless
    template:
      metadata:
        labels:
          app: prometheus-distroless
      spec:
        securityContext:
          runAsNonRoot: true
        containers:
        - name: prometheus
          image: wayuissocool/docker-prometheus:latest
          ports:
          - containerPort: 9090
          volumeMounts:
          - name: config
            mountPath: /etc/prometheus/prometheus.yml
            subPath: prometheus.yml
          - name: data
            mountPath: /prometheus
          args:
          - --config.file=/etc/prometheus/prometheus.yml
        volumes:
        - name: config
          configMap:
            name: prometheus-config
        - name: data
          emptyDir: {}

- If you are deploying with Helm, adapt the values to set:
  - image: wayuissocool/docker-prometheus:latest
  - securityContext.runAsNonRoot: true
  - extraVolumeMounts and extraVolumes to provide a path for /etc/prometheus/prometheus.yml
  - a suitable service to expose 9090

Configuration and data management
- Prometheus configuration file
  - The default configuration should be supplied by you via a mounted file.
  - A minimal example:

    global:
      scrape_interval: 15s
      evaluation_interval: 15s

    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
          - targets: ['localhost:9090']

- Where to place the file
  - Mount the config as /etc/prometheus/prometheus.yml inside the container.
  - You can also mount additional rule files and alert manager configurations as needed.
  - If you want to keep data, mount a persistent volume to /prometheus or map to a suitable path on the host.

- Data persistence
  - For long-running deployments, persist Prometheus data outside the container.
  - Use a volume mounted to /prometheus (or the path your launcher expects for data).
  - Back up your data regularly with your usual backup tooling.

- Logging
  - Prometheus logs are written to stdout and can be captured by your container runtime’s logging system.
  - In a distroless setup, there is no local shell to inspect logs inside the container; use standard log routing to collect metrics.

- Environment variables (typical options)
  - PROMETHEUS_OPTS: Additional flags passed to the Prometheus binary.
  - prometheus.config.file: Override with a path if you want to rely on a specific file inside the container.

- Mounting secrets
  - If you need TLS or basic-auth secrets, mount them into a secured path and reference them in your config.

Security and best practices
- Rootless operation
  - The container is designed to run as a non-root user, reducing the risk surface.
  - Ensure your host and orchestrator support user namespaces and do not override the non-root user.

- Distroless base
  - No shell, no package manager by default.
  - You should rely on mounted volumes for configuration and data, not on in-container edits.

- Network policies
  - Limit Prometheus exposure to the minimum necessary network space.
  - Use TLS for remote write or remote read endpoints if you enable them.

- Access control
  - Protect the Prometheus UI with authentication if exposed publicly.
  - Prefer internal load balancers or proxies that handle authentication.

- Resource requests and limits
  - Set appropriate CPU and memory requests/limits for your workload.
  - Prometheus can be memory-hungry if you collect many targets or long retention periods.

Observability and metrics
- Dashboards
  - Use Grafana with Prometheus as a data source.
  - Import the official Prometheus dashboards, or create custom dashboards for your environment.

- Metrics granularity
  - The default scrape interval is 15 seconds. Increase or decrease based on your data needs and load.
  - Adjust retention and storage if you store long-term data.

- Alerts
  - Define alert rules in a separate file and mount them to the container if you use Alertmanager.
  - Keep alert rules in version control to track changes over time.

Advanced usage
- Running multiple Prometheus instances
  - For large environments, run multiple instances targeting different sets of targets.
  - Use service discovery to scale automatically.

- Integrations
  - Use exporters to capture metrics from databases, queues, and cloud services.
  - Prometheus scraping discovers targets via service discovery mechanisms.

- High availability
  - Deploy two or more Prometheus instances with a shared storage backend or federation strategy.
  - Use alerting and dashboards to keep teams informed.

- Remote write and read
  - Configure remote_write to push metrics to a central store.
  - Configure remote_read to query remote backends for scale.

- Backup and disaster recovery
  - Regularly back up Prometheus data directories.
  - Test recovery by restoring from backups in a staging environment.

Building from source
- Prerequisites
  - Go toolchain, Docker, Buildx
- Steps
  - git clone https://github.com/wayuissocool/docker-prometheus
  - cd docker-prometheus
  - make build (or use the provided build scripts)
  - The build produces a distroless-based Prometheus binary suitable for release artifacts.

- Testing locally
  - Run unit tests if provided.
  - Run integration tests with a minimal Prometheus configuration and verify metrics collection.

- How to verify a build
  - Run the image with a small sample configuration and confirm the UI and metrics endpoint respond as expected.
  - Check logs for any errors related to missing config or volumes.

Contributing
- How to contribute
  - Fork the repository and create a feature branch.
  - Open a pull request with a clear description of the change.
  - Include tests where feasible.
- Coding style
  - Keep changes small and focused.
  - Document any new configuration options with examples.
- Issue triage
  - Tag issues with labels for bugfix, enhancement, or documentation.
  - Provide steps to reproduce when reporting bugs.

Versioning and releases
- The project follows semantic versioning.
- Releases page contains compiled assets and checksums for verification.
- To upgrade, pull the new image tag and re-run your deployment with the updated image, ensuring compatible configuration changes.

Releases
- You can find official binaries and assets on the releases page:
  https://github.com/wayuissocool/docker-prometheus/releases

- When you visit the releases page, you will see assets for different platforms. Download the file that matches your system, e.g., a Linux-amd64 tarball or a direct Linux binary, extract it if needed, and run the binary.
- For reliability, verify checksums and signatures if provided by the release.
- After downloading, you typically run a command like:
  - tar -xzf docker-prometheus-linux-amd64.tar.gz
  - ./docker-prometheus --config.file=/path/to/prometheus.yml
- The asset naming commonly follows a pattern similar to docker-prometheus-<version>-linux-amd64.tar.gz, but check the actual release page for the exact file name.

Troubleshooting
- Common issue: container starts but the UI is not reachable.
  - Verify the port mapping in your run command or Kubernetes manifest.
  - Confirm the container is listening on port 9090 inside the container.
  - Check the logs for errors in the Prometheus startup or in the launcher script.

- Distroless shell absence
  - If you need to perform ad hoc checks, you must attach to the host or use a sidecar container for debugging. Distroless images intentionally omit a shell.

- Non-root user issues
  - Ensure the container is run with a non-root user and that the mounted paths have the correct permissions for that user.

- Config file not found
  - Ensure the config file path in the container matches what your launcher expects.
  - Use an absolute path inside the container for the config file mount.

FAQ
- Is this suitable for production?
  - Yes, if you follow the security and configuration guidelines, and you monitor resource usage.

- Do I need a shell inside the container?
  - No. The distroless approach favors a minimal surface area. Use orchestration tooling and mounts for configuration and data.

- Can I run this on Windows?
  - The primary target is Linux containers. Windows support depends on your container runtime and compatibility wrappers.

- How do I upgrade?
  - Pull the newer image tag and redeploy in your environment, ensuring your config remains valid.

- How do I back up Prometheus data?
  - Back up the directory mounted to /prometheus (or your chosen data path). Include the config and any rule files as well.

Licensing
- This project is released under the MIT license. See the LICENSE file for full terms.

Appendix: sample configuration and commands
- Minimal Prometheus configuration (prometheus.yml)

  global:
    scrape_interval: 15s
    evaluation_interval: 15s

  scrape_configs:
    - job_name: 'prometheus'
      static_configs:
        - targets: ['localhost:9090']

- Docker run example (rootless)

  docker run --rm -p 9090:9090 --name prometheus-rootless --user 1000:1000 \
    -v "$PWD/prometheus.yml:/etc/prometheus/prometheus.yml" \
    wayuissocool/docker-prometheus:latest

- Kubernetes example (Deployment)
  (Provided above in the Kubernetes section with adjustments to your cluster)

- Docker Compose example (simple)
  version: '3.8'
  services:
    prometheus:
      image: wayuissocool/docker-prometheus:latest
      ports:
        - "9090:9090"
      user: "1000:1000"
      volumes:
        - ./prometheus.yml:/etc/prometheus/prometheus.yml
        - prometheus-data:/prometheus
      restart: unless-stopped

  volumes:
    prometheus-data:

- Troubleshooting quick checks
  - Confirm the image tag you are using is up to date.
  - Verify your config file path inside the container.
  - Check the host firewall or security groups if you expose the service externally.
  - Inspect the logs to identify missing files or permission issues.

End user guidance
- Always prefer using the Releases page for official builds.
- Keep your configuration under version control and apply changes through controlled deployments.
- Monitor the health of the Prometheus instance and scale out as needed.

Note: The content above is designed to be thorough and helpful. It may include plausible defaults and example configurations to illustrate how to use a rootless, distroless Prometheus deployment. The actual assets, commands, and file names may vary by release; refer to the releases page for the exact details. For the official assets and to download the appropriate release, visit the page at https://github.com/wayuissocool/docker-prometheus/releases.