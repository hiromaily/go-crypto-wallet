## FAQ

### Is DevContainer required?

**No.** DevContainer is completely optional. You can continue with local development without any changes.

### Does DevContainer affect my host system?

**No.** Everything runs in an isolated container. Your host system remains unchanged.

### Can I use DevContainer with multiple projects?

**Yes.** Each project has its own container with its own configuration.

### What happens to my files?

Your files are **bind-mounted** from your host system. Changes in the container are reflected on your host and vice versa.

### Can I customize the DevContainer?

**Yes.** Edit `.devcontainer/devcontainer.json`:

- Change the base image
- Add more tools in `postCreateCommand`
- Add/remove VS Code extensions
- Modify settings

### How much disk space does it use?

- Base Go image: ~1GB
- Tools and dependencies: ~500MB
- **Total: ~1.5GB** for the first setup
- Subsequent projects share the base image

### Can I use DevContainer offline?

**Partially.** After the first build, you can work offline. However:

- Initial build requires internet (downloads image)
- `go get` requires internet
- `docker compose pull` requires internet

### Does DevContainer slow down my editor?

**No.** VS Code/Cursor connects to the container via a lightweight client-server architecture. Performance is similar to local development.

### Can I access localhost from the container?

**Yes.** Ports are automatically forwarded. If your app runs on `localhost:8080` in the container, access it at `localhost:8080` on your host.

### What about M1/M2/M3 Mac compatibility?

**Fully supported.** The DevContainer images support ARM64 architecture.

### Can I use Docker Compose in DevContainer?

**Yes.** Docker Compose commands work through the mounted host socket. Containers run on the host, accessible from DevContainer.

### How do I exit DevContainer?

```bash
# Close the VS Code/Cursor window
# OR
# F1 → "Dev Containers: Reopen Folder Locally"
```

### How do I delete the DevContainer?

```bash
# Remove container and image
docker ps -a | grep go-crypto-wallet | awk '{print $1}' | xargs docker rm
docker images | grep go-crypto-wallet | awk '{print $3}' | xargs docker rmi

# Or use Docker Desktop UI
# Containers → Delete
# Images → Delete
```

---
