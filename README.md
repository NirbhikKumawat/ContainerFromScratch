# 📦 Go-Container-From-Scratch
A educational implementation of Linux container primitives written in Go. This project demonstrates how tools like Docker use the Linux kernel's features—specifically Namespaces and Control Groups (Cgroups)—to isolate and limit processes.

## 🚀 Features
- __PID Isolation:__ The containerized process sees itself as PID 1.

- __Filesystem Isolation:__ Uses chroot to trap the process in a private root filesystem.

- __Mount Isolation:__ Private mount table via CLONE_NEWNS with a namespace-aware /proc.

- __UTS Isolation:__ Isolated hostname and domain name.

- __Network Isolation:__ A private network stack (initially empty).

- __Resource Limiting:__ Cgroup v2 integration to cap memory usage (set to 100MB).

- __The Re-exec Pattern:__ Uses the /proc/self/exe strategy to perform setup from inside the namespaces.

## 🛠️ Prerequisites
- __OS:__ Linux (Fedora, Ubuntu, etc.). This will not work on macOS or Windows.

- __Language:__ Go 1.18+

- __Privileges:__ sudo access is required to create namespaces.

## 📁 Setting up the Root Filesystem
Before running the container, you need a mini-Linux environment (a rootfs) for the container to "live" in.

### Create the directory structure:

```bash
mkdir -p my-rootfs/bin my-rootfs/lib64 my-rootfs/proc
```
### Copy the binaries:
You need to copy bash, ls, and ps (and their dependencies) into the folder.
```bash
cp /bin/bash /bin/ls /bin/ps my-rootfs/bin/
```

### Copy dependencies:
Use ldd /bin/bash to find the shared libraries and copy them into my-rootfs/lib64.

__Note:__ Ignore linux-vdso.so.1; it is a virtual library provided by the kernel.

## 🏃 How to Run
### Build the project:

```bash
go build -o main main.go
```

### Execute the container:

```bash
sudo ./main run /bin/bash
```

### Testing the Isolation
Once inside the container shell, try the following:

- __ps aux:__ You should see your shell as PID 1 and nothing from the host.

- __hostname:__ See the name change to my-container.

- `ls /`: You should only see the bin, lib64, and proc folders you created.

- `cat /proc/self/cgroup`: Verify you are inside the my-container control group.

## 🧠 How it Works: The Architecture
The Two-Stage Launch
To properly mount the `/proc` filesystem and set the hostname, we have to run code after the namespaces are created but before the user's command starts. We achieve this through Re-execution:

- __The Parent (run):__ Creates the namespaces using `syscall.SysProcAttr` and launches a copy of itself.

- __The Child (child):__ Now inside the namespaces, it sets the hostname, joins the Cgroup, chroots into the new filesystem, and mounts /proc.

- __The Exec:__ The child uses `syscall.Exec` to overwrite itself with the target command (e.g., `/bin/bash`).

## Resource Control
The program creates a cgroup at `/sys/fs/cgroup/my-container`. It writes 100M to memory.max and adds the child's PID to `cgroup.procs`. When the parent exits, it automatically cleans up the mount and the cgroup.
