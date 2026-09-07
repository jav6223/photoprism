# Windows Quick Start Guide

## 🚀 Super Simple Setup (3 Steps)

Your PhotoPrism is on USB NVMe at: `/media/advsecm1/J/PhotoPrism_Portable`

---

## ✅ **Step 1: One-Time Installation**

**On your Windows PC, open PowerShell as Administrator:**

Right-click Start → "Windows PowerShell (Admin)"

**Run this ONE command:**

```powershell
# Copy windows_install.ps1 from USB to your PC first, then:
cd C:\Users\YourName\Downloads
.\windows_install.ps1
```

**OR manually install via Chocolatey:**

```powershell
# Install Chocolatey
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install everything
choco install docker-desktop git vscode googlechrome -y

# Install WSL2
wsl --install -d Ubuntu
```

**Then RESTART your PC!**

---

## ✅ **Step 2: After Restart**

**1. Start Docker Desktop**
   - Open from Start Menu
   - Wait for it to say "Docker Desktop is running"

**2. Setup PhotoPrism**

PowerShell as Administrator:

```powershell
cd C:\Users\YourName\Downloads
.\setup_photoprism.ps1
```

This will:
- Mount your USB NVMe
- Find PhotoPrism
- Start it

---

## ✅ **Step 3: Daily Use**

**Two options:**

### **Option A: Double-Click (Easiest)**

1. Copy `start_photoprism_windows.bat` to your Desktop
2. Edit it and change `PHYSICALDRIVE2` to your actual drive number
3. Right-click → "Run as Administrator"
4. PhotoPrism starts!

### **Option B: PowerShell**

```powershell
# Mount USB
wsl --mount \\.\PHYSICALDRIVE2 --partition 1

# Start PhotoPrism
wsl -e bash -c "cd /mnt/wsl/PHYSICALDRIVE2/media/advsecm1/J/PhotoPrism_Portable && docker compose up -d"

# Open browser
start http://localhost:2342
```

---

## 🔍 **How to Find Your Drive Number**

PowerShell as Administrator:

```powershell
wmic diskdrive list brief
```

Look for your USB NVMe (check the size), note the number:
- Example: `\\.\PHYSICALDRIVE2` → Use `2`
- Example: `\\.\PHYSICALDRIVE1` → Use `1`

---

## 📦 **What Gets Installed**

| Software | Purpose | Size |
|----------|---------|------|
| **Docker Desktop** | Run PhotoPrism containers | ~500 MB |
| **WSL2 + Ubuntu** | Linux environment | ~2 GB |
| **Git** | Version control | ~300 MB |
| **VS Code** | Code editor (optional) | ~200 MB |
| **Chrome** | Web browser | ~200 MB |

**Total:** ~3-4 GB

---

## 🎮 **GPU Acceleration on Windows**

Your GPU code will automatically work on Windows!

**Auto-detection priority:**
1. ✅ **NVIDIA CUDA** (if NVIDIA GPU + drivers)
2. ✅ **DirectML** (if any GPU - Intel/AMD/NVIDIA)
3. ✅ **CPU fallback** (if no GPU)

**To enable GPU in Docker Desktop:**
1. Open Docker Desktop
2. Settings → Resources → **Enable GPU support**
3. Restart Docker Desktop

Your `.env` already has:
```env
PHOTOPRISM_FACE_GPU=true
PHOTOPRISM_FACE_GPU_PROVIDER=auto
```

**It will just work!** 🎉

---

## 📂 **Your Data Locations**

**On USB NVMe (Linux ext4):**
```
/media/advsecm1/J/
├── PhotoPrism_Portable/    ← Main installation
│   ├── docker-compose.yml
│   ├── .env
│   └── storage/
└── Anmol/                  ← Your photos (8TB)
```

**In Windows (after mounting):**
```
\\wsl$\PHYSICALDRIVE2\media\advsecm1\J\PhotoPrism_Portable\
```

You can access files from Windows Explorer: `\\wsl$\`

---

## 🛠️ **Useful Commands**

### **View logs:**
```powershell
docker logs photoprism
docker logs photoprism-db
```

### **Stop PhotoPrism:**
```powershell
docker compose -f /mnt/wsl/PHYSICALDRIVE2/media/advsecm1/J/PhotoPrism_Portable/docker-compose.yml down
```

### **Restart PhotoPrism:**
```powershell
docker restart photoprism
```

### **Check if running:**
```powershell
docker ps
```

### **Import photos manually:**
```powershell
docker exec photoprism photoprism import
```

### **Check GPU is working:**
```powershell
docker logs photoprism | findstr "execution provider"
```

Should show: `onnx: using cuda execution provider` or `onnx: using directml execution provider`

---

## 🐛 **Troubleshooting**

### **"Cannot find PHYSICALDRIVE2"**

Find your drive:
```powershell
wmic diskdrive list brief
```

Look for your USB NVMe size, use that number.

### **"Docker daemon is not running"**

1. Start Docker Desktop from Start Menu
2. Wait for green icon in system tray
3. Try again

### **"WSL 2 installation is incomplete"**

Download WSL2 kernel update:
https://aka.ms/wsl2kernel

### **"Port 2342 already in use"**

Something else is using the port:
```powershell
netstat -ano | findstr :2342
# Kill that process or change PhotoPrism port
```

### **"Permission denied"**

Make sure you're running PowerShell/CMD as Administrator.

### **"Cannot access \\wsl$\"**

Make sure WSL2 is running:
```powershell
wsl --list --running
```

If not running:
```powershell
wsl
```

---

## 🎯 **Quick Reference**

### **First Time Setup:**
1. Run `windows_install.ps1` (installs everything)
2. Restart PC
3. Start Docker Desktop
4. Run `setup_photoprism.ps1` (mounts & starts)

### **Daily Use:**
1. Plug in USB NVMe
2. Double-click `start_photoprism_windows.bat`
3. Open http://localhost:2342

### **Stop:**
```powershell
docker compose down
```

---

## 📞 **Quick Help**

### **Is it running?**
```powershell
docker ps
# Should show 'photoprism' and 'photoprism-db'
```

### **Where are my photos?**
```
\\wsl$\PHYSICALDRIVE2\media\advsecm1\J\Anmol\
```

### **Why is it slow?**
Check if GPU is enabled:
```powershell
docker logs photoprism | findstr GPU
```

---

## 🎉 **Summary**

**Installation:**
```powershell
choco install docker-desktop git -y
wsl --install -d Ubuntu
# Restart PC
```

**Daily Start:**
```batch
start_photoprism_windows.bat
```

**Access:**
```
http://localhost:2342
```

**That's it!** Your entire PhotoPrism setup runs from the USB NVMe on Windows! 🚀

---

## 💡 **Pro Tips**

### **Create Desktop Shortcut:**

1. Right-click Desktop → New → Shortcut
2. Location: `C:\path\to\start_photoprism_windows.bat`
3. Name: "PhotoPrism"
4. Right-click shortcut → Properties → Advanced → ✅ Run as administrator

### **Auto-start on Windows Boot:**

1. Press `Win + R`
2. Type: `shell:startup`
3. Copy `start_photoprism_windows.bat` there
4. PhotoPrism starts automatically when Windows boots!

### **Access from Phone:**

1. Find Windows PC IP: `ipconfig`
2. On phone: `http://192.168.1.x:2342`
3. (Replace x with your PC's IP)

---

**Your PhotoPrism setup is now portable between Linux and Windows!** 🎊
