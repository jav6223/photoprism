# PhotoPrism Configuration Summary

## ✅ What We've Configured

### 1. Auto-Import & Upload Settings
Created `.env` file with:
- **Auto-import**: Enabled (300 seconds / 5 minutes delay)
- **Auto-index**: Enabled (120 seconds / 2 minutes delay)
- **Upload limit**: 600 GB (614,400 MB)
- **Upload restrictions**: None (all file types allowed)

### 2. Your Photo Collection
- **Mounted**: `/media/advsecm1/J/Anmol` → PhotoPrism import folder
- **Auto-folder albums**: Will create albums for each subfolder:
  - DCIM, Manali 11-24, Farewell, Food, etc.

---

## 🎮 GPU Acceleration Status

### Your Hardware
- **GPU**: Intel Arrow Lake-U (Integrated Graphics)
- **Devices**: `/dev/dri/renderD128` available ✅

### What Works with GPU

#### ✅ Video Transcoding (Available Now)
**Intel Quick Sync** or **VA-API** for hardware video encoding:

```bash
# Option 1: Intel Quick Sync (recommended for newer Intel)
PHOTOPRISM_FFMPEG_ENCODER=intel

# Option 2: VA-API (more portable)
PHOTOPRISM_FFMPEG_ENCODER=vaapi
```

Add to your `.env` file and restart to enable.

#### ❌ AI Models (Not Available - CPU Only)
**Current limitation**: AI models run on CPU only
- Face detection (ONNX)
- Face recognition (ONNX embeddings)
- Image classification/labels (TensorFlow)
- NSFW detection (TensorFlow)

**Why?** The code doesn't configure GPU execution providers for ONNX Runtime or TensorFlow.

**To add GPU support would require**:
1. Modifying `internal/ai/onnx/runtime.go` to add CUDA execution provider
2. Modifying `internal/ai/tensorflow/` to enable GPU sessions
3. Installing CUDA/cuDNN libraries in Docker image
4. Compiling TensorFlow/ONNX with GPU support

This is a significant code change and isn't currently implemented.

---

## 🔞 NSFW Detection (Yes, It Exists!)

PhotoPrism has a **built-in NSFW detection model** that prevents inappropriate content.

### How It Works
- **Model**: TensorFlow-based classifier
- **Categories**: Classifies images as:
  - **Drawing** - Artwork/illustrations
  - **Hentai** - Adult anime/manga
  - **Neutral** - Safe for work
  - **Porn** - Pornographic content
  - **Sexy** - Suggestive content

### Configuration Options

```bash
# Option 1: Flag uploads as private (won't upload to public albums)
PHOTOPRISM_DETECT_NSFW=true          # Auto-flag potential NSFW as private
PHOTOPRISM_UPLOAD_NSFW=false         # Block NSFW uploads entirely

# Option 2: Allow everything
PHOTOPRISM_DETECT_NSFW=false         # Don't check
PHOTOPRISM_UPLOAD_NSFW=true          # Allow all uploads
```

### Thresholds
- **Safe**: < 75% NSFW score
- **Medium**: 75-85% NSFW score  
- **High**: 85-98% NSFW score
- **Blocked**: > 98% in Porn/Sexy/Hentai categories

### Current Setting
Your `.env` has: `PHOTOPRISM_UPLOAD_NSFW=true` (allows all uploads)

**To enable NSFW blocking**, change in `.env`:
```bash
PHOTOPRISM_DETECT_NSFW=true          # Flag NSFW as private
PHOTOPRISM_UPLOAD_NSFW=false         # Reject NSFW uploads
```

---

## 🚀 Next Steps

### 1. Enable Intel GPU for Videos (Optional)
Edit `.env` and uncomment:
```bash
PHOTOPRISM_FFMPEG_ENCODER=intel
```

Add to `compose.yaml` under `photoprism` service:
```yaml
devices:
  - "/dev/dri:/dev/dri"
```

### 2. Start PhotoPrism
```bash
cd /media/advsecm1/J/photoprism/photoprism
docker compose up -d
```

### 3. Watch Auto-Import
```bash
docker compose logs -f photoprism
```

PhotoPrism will automatically:
1. Find new files in `/media/advsecm1/J/Anmol`
2. Import them every 5 minutes
3. Index them every 2 minutes
4. Create folder albums for each directory

### 4. Manual Import (Faster)
If you want to import immediately instead of waiting:
```bash
docker compose exec photoprism photoprism import
```

Or index directly:
```bash
docker compose exec photoprism photoprism index
```

---

## 📊 Your Photo Collection
Based on `/media/advsecm1/J/Anmol/`, you have folders like:
- DCIM (Camera photos)
- Manali 11-24 (Trip photos)
- Farewell, Food, Gujrat all, etc.
- Expert RAW (High-quality images)
- Heavy videos 8k,4k,2k (Large videos - will benefit from GPU transcoding!)

Each will become a separate album automatically.

---

## ⚙️ Configuration Files Modified
1. **`.env`** - Custom settings (auto-import, upload limits)
2. **`compose.yaml`** - Added Anmol folder mount

## 🔧 Useful Commands
```bash
# Restart PhotoPrism
docker compose restart photoprism

# View logs
docker compose logs -f photoprism

# Force re-index
docker compose exec photoprism photoprism index --cleanup

# Check storage
docker compose exec photoprism photoprism show storage

# Import specific folder
docker compose exec photoprism photoprism import --path Anmol/DCIM
```

---

## 💡 Pro Tips

1. **Large video collection**: Enable Intel GPU encoding to speed up 4K/8K video processing
2. **NSFW filtering**: Enable if you want automatic content filtering
3. **Face detection**: Will automatically detect faces in photos (CPU-only, may be slow)
4. **RAW photos**: Automatically converts RAW formats (may be slow without GPU)

---

## ❓ FAQ

**Q: Can I add GPU support for AI models?**  
A: Not without code modifications. Would require implementing ONNX CUDA execution provider and TensorFlow GPU sessions.

**Q: Will auto-import delete my original files?**  
A: No! It copies files to PhotoPrism storage, originals stay untouched.

**Q: How fast is upload?**  
A: Limited by network/disk speed, not by PhotoPrism. For bulk uploads, copy directly to `storage/import/` folder.

**Q: Do folder albums update automatically?**  
A: Yes! When you add photos to a folder, the album updates on next index.

---

Enjoy your PhotoPrism setup! 🎉
