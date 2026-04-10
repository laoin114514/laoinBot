podman stop laoinbot
podman rm laoinbot
podman build -t laoinbot .
podman run -d --name laoinbot laoinbot