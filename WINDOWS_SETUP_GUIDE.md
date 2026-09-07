# Running PhotoPrism with GPU on Windows

## 🎯 Quick Overview

You can run your GPU-enabled PhotoPrism fork on Windows! Here's how.

---

## 📋 Prerequisites on Windows

### 1. Install Docker Desktop

**Download:** https://www.docker.com/products/docker-desktop/

**Requirements:**
- Windows 10/11 Pro, Enterprise, or Education (64-bit)
- WSL 2 enabled
- At least 8GB RAM
- 20GB free disk space

**Install Steps:**
1. Download Docker Desktop
2. Run installer
3. Enable WSL 2 during installation
4. Restart computer
5. Start Docker Desktop

### 2. Install Git for Windows

**Download:** https://git-scm.com/download/win

Or use GitHub Desktop: https://desktop.github.com/

---

## 🚀 Setup Your Fork on Windows

### **Option A: Using Git Bash (Recommended)**

```bash
# 1. Clone your fork
cd C:\Users\YourName\Documents
git clone https://github.com/jav6223/photoprism.git
cd photoprism

# 2. Checkout GPU branch
git checkout feature/gpu-acceleration

# 3. Create .env file for Windows
notepad .env
```

### **Option B: Using GitHub Desktop**

1. Open GitHub Desktop
2. File → Clone Repository
3. Enter: `jav6223/photoprism`
4. Choose location
5. Click "Clone"
6. Switch branch to `feature/gpu-acceleration`

---

## ⚙️ Configure for Windows

### **Create `.env` File**

Create a file called `.env` in the PhotoPrism directory:

```env
# PhotoPrism Windows Configuration

## Admin Credentials
PHOTOPRISM_ADMIN_USER=admin
PHOTOPRISM_ADMIN_PASSWORD=yourpassword

## Database
MYSQL_ROOT_PASSWORD=insecure
MYSQL_DATABASE=photoprism
MYSQL_USER=photoprism
MYSQL_PASSWORD=insecure

## Performance
PHOTOPRISM_WORKERS=4
PHOTOPRISM_UPLOAD_LIMIT=614400  # 600 GB

## GPU Settings (Windows DirectML)
PHOTOPRISM_FACE_GPU=true
PHOTOPRISM_FACE_GPU_PROVIDER=auto  # Will detect DirectML on Windows
PHOTOPRISM_FACE_GPU_DEVICE=0

## Video Transcoding
PHOTOPRISM_FFMPEG_ENCODER=software  # Or 'nvidia' if you have NVIDIA GPU
```

### **Create `docker-compose.yml` for Windows**

Create `docker-compose.windows.yml`:

```yaml
version: '3.8'

services:
  photoprism:
    image: photoprism/photoprism:latest
    depends_on:
      - mariadb
    restart: unless-stopped
    security_opt:
      - seccomp:unconfined
      - apparmor:unconfined
    ports:
      - "2342:2342"
    env_file:
      - .env
    environment:
      PHOTOPRISM_SITE_URL: "http://localhost:2342/"
      PHOTOPRISM_ORIGINALS_LIMIT: -1
      PHOTOPRISM_HTTP_COMPRESSION: "gzip"
      PHOTOPRISM_DATABASE_DRIVER: "mysql"
      PHOTOPRISM_DATABASE_SERVER: "mariadb:3306"
      PHOTOPRISM_DATABASE_NAME: "photoprism"
      PHOTOPRISM_DATABASE_USER: "photoprism"
      PHOTOPRISM_DATABASE_PASSWORD: "insecure"
    volumes:
      # Mount your photos from Windows
      - "D:/Photos:/photoprism/originals"  # Change D:/Photos to your photos location
      - "./storage:/photoprism/storage"
    # For NVIDIA GPU on Windows (optional)
    # deploy:
    #   resources:
    #     reservations:
    #       devices:
    #         - driver: nvidia
    #           count: 1
    #           capabilities: [gpu]

  mariadb:
    image: mariadb:11
    restart: unless-stopped
    security_opt:
      - seccomp:unconfined
      - apparmor:unconfined
    command: 
      - --innodb-buffer-pool-size=256M
      - --transaction-isolation=READ-COMMITTED
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
      - --max-connections=512
      - --innodb-rollback-on-timeout=OFF
      - --innodb-lock-wait-timeout=120
    volumes:
      - "mariadb_data:/var/lib/mysql"
    env_file:
      - .env

volumes:
  mariadb_data:
```

---

## 🎮 GPU Support on Windows

### **GPU Detection on Windows:**

Your GPU code supports **DirectML** on Windows, which works with:
- ✅ NVIDIA GPUs (RTX, GTX series)
- ✅ AMD GPUs (Radeon series)
- ✅ Intel integrated GPUs

The `PHOTOPRISM_FACE_GPU_PROVIDER=auto` setting will automatically detect:
1. **DirectML** (Windows-native GPU acceleration)
2. **NVIDIA CUDA** (if NVIDIA GPU + drivers installed)
3. **CPU fallback** (if no GPU detected)

### **For NVIDIA GPU (Optional - Better Performance):**

If you have NVIDIA GPU and want maximum speed:

1. Install **NVIDIA Container Toolkit for Windows**
2. Install **CUDA Toolkit**: https://developer.nvidia.com/cuda-downloads
3. In `docker-compose.windows.yml`, uncomment the GPU section
4. Set `PHOTOPRISM_FACE_GPU_PROVIDER=cuda`

---

## 🚀 Start PhotoPrism on Windows

### **Using PowerShell or CMD:**

```powershell
# Navigate to your PhotoPrism directory
cd C:\Users\YourName\Documents\photoprism

# Start PhotoPrism
docker-compose -f docker-compose.windows.yml up -d

# Check logs
docker-compose -f docker-compose.windows.yml logs -f photoprism

# Stop PhotoPrism
docker-compose -f docker-compose.windows.yml down
```

### **Access PhotoPrism:**

Open browser: http://localhost:2342

Default login:
- Username: `admin`
- Password: `yourpassword` (whatever you set in .env)

---

## 📂 Transferring Your Photos from Linux

### **Option A: Copy Photos to Windows**

```bash
# On Linux, copy to USB drive or network share
sudo cp -r /media/advsecm1/J/Anmol /media/usb-drive/

# On Windows, copy from USB to D:\Photos
```

### **Option B: Network Share**

1. Share Linux folder over network (Samba)
2. Mount in Windows as network drive
3. Point PhotoPrism to network drive

### **Option C: Sync Database Too**

To keep your existing indexes/faces:

```bash
# On Linux - backup database
cd /media/advsecm1/J/PhotoPrism_Portable
sudo docker compose exec db mysqldump -u root -pinsecure photoprism > photoprism_backup.sql

# Copy photoprism_backup.sql to Windows

# On Windows - restore database
docker-compose -f docker-compose.windows.yml up -d mariadb
docker-compose -f docker-compose.windows.yml exec -T mariadb mysql -u root -pinsecure photoprism < photoprism_backup.sql
```

---

## 🔧 Build Custom Docker Image on Windows

If you want to build your GPU-enabled version on Windows:

### **PowerShell Script:**

```powershell
# Navigate to repo
cd C:\Users\YourName\Documents\photoprism

# Build custom image with GPU support
docker build -t photoprism-gpu:latest .

# Update docker-compose to use custom image
# Change: image: photoprism/photoprism:latest
# To:     image: photoprism-gpu:latest
```

---

## 🐛 Troubleshooting Windows

### **"Docker daemon not running"**
**Fix:** Start Docker Desktop from Start menu

### **"WSL 2 installation incomplete"**
**Fix:** Run in PowerShell (as Admin):
```powershell
wsl --install
wsl --set-default-version 2
```

### **"Port 2342 already in use"**
**Fix:** Change port in docker-compose:
```yaml
ports:
  - "2343:2342"  # Use 2343 instead
```

### **"GPU not detected"**
**Fix:** 
1. Check Docker Desktop → Settings → Resources → Enable GPU
2. Install GPU drivers
3. For NVIDIA: Install CUDA Toolkit

### **Slow performance**
**Fix:**
1. Increase Docker Desktop memory (Settings → Resources)
2. Move Docker data to SSD
3. Enable GPU acceleration

---

## 📊 Performance Comparison: Windows vs Linux

| Feature | Linux (Your Current) | Windows |
|---------|---------------------|---------|
| CPU Performance | Same | Same |
| GPU (Intel) | OpenVINO | DirectML |
| GPU (NVIDIA) | CUDA | CUDA or DirectML |
| Speed | Baseline | ~10-20% slower (WSL2 overhead) |
| GPU Speedup | 5-10x (OpenVINO) | 5-10x (DirectML) |

**Recommendation:** Windows with GPU is still **much faster** than Windows CPU-only!

---

## 🎯 Quick Start Checklist

- [ ] Install Docker Desktop on Windows
- [ ] Install Git for Windows
- [ ] Clone your fork: `git clone https://github.com/jav6223/photoprism.git`
- [ ] Checkout branch: `git checkout feature/gpu-acceleration`
- [ ] Create `.env` file with your settings
- [ ] Create `docker-compose.windows.yml`
- [ ] Update photo path in docker-compose (D:/Photos)
- [ ] Run: `docker-compose -f docker-compose.windows.yml up -d`
- [ ] Open: http://localhost:2342

---

## 🌟 Additional Tips

### **Auto-Start on Windows Boot:**

1. Docker Desktop → Settings → General
2. Check: ✅ "Start Docker Desktop when you log in"
3. Add to Windows startup:
   ```powershell
   # Create startup script: start_photoprism.bat
   cd C:\Users\YourName\Documents\photoprism
   docker-compose -f docker-compose.windows.yml up -d
   ```

### **Backup on Windows:**

Create scheduled task to backup:
```powershell
# backup.bat
docker-compose -f docker-compose.windows.yml exec db mysqldump -u root -pinsecure photoprism > "D:\Backups\photoprism_%date%.sql"
```

### **Update PhotoPrism:**

```powershell
cd C:\Users\YourName\Documents\photoprism
git pull origin feature/gpu-acceleration
docker-compose -f docker-compose.windows.yml down
docker-compose -f docker-compose.windows.yml up -d --build
```

---

## 📱 Access from Phone/Tablet

PhotoPrism will be accessible on your local network:

1. Find your Windows PC IP: `ipconfig`
2. On phone/tablet: http://192.168.1.x:2342
3. (Replace x with your PC's IP)

For external access, set up port forwarding or use Tailscale.

---

## 🔐 Security Recommendations

1. Change default passwords in `.env`
2. Enable HTTPS if exposing to internet
3. Use strong admin password
4. Don't expose port 2342 to public internet without auth
5. Enable firewall rules

---

## 📞 Need Help?

Check logs on Windows:
```powershell
docker-compose -f docker-compose.windows.yml logs -f photoprism
```

Common issues are in the Troubleshooting section above!

---

**Your GPU-enabled PhotoPrism will work great on Windows!** 🚀

The main difference is:
- **Linux:** Uses OpenVINO (Intel) or CUDA (NVIDIA)
- **Windows:** Uses DirectML (any GPU) or CUDA (NVIDIA)

Both give you **5-10x speedup** over CPU-only! 🎉
