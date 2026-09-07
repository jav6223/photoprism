import requests
import json
import random
import os
import datetime
from PIL import Image, ImageDraw, ImageFont, ImageFilter
import io

# --- CONFIGURATION ---
PHOTOPRISM_URL = "http://localhost:2342"
PHOTOPRISM_USER = "admin"
PHOTOPRISM_PASS = "YourSecurePassword" # From your docker-compose.yml
OLLAMA_URL = "http://localhost:11434/api/generate"
OLLAMA_MODEL = "llama3.2-vision:11b"
EXPORT_DIR = "./Memories_Export"
# ---------------------

def get_session():
    print("[*] Authenticating with PhotoPrism...")
    res = requests.post(f"{PHOTOPRISM_URL}/api/v1/session", json={
        "username": PHOTOPRISM_USER,
        "password": PHOTOPRISM_PASS
    })
    res.raise_for_status()
    return res.json()["id"]

def find_memory_photos(session_id):
    print("[*] Searching for 'On This Day' or Random Throwback photos...")
    headers = {"X-Session-ID": session_id}
    
    # Decide randomly whether to do "On This Day" or "Random Throwback"
    is_on_this_day = random.choice([True, False])
    
    now = datetime.datetime.now()
    month = now.month
    day = now.day
    
    query = f"month:{month} day:{day} before:{now.year}" if is_on_this_day else "quality:3"
    
    res = requests.get(f"{PHOTOPRISM_URL}/api/v1/photos", params={
        "count": 50,
        "q": query,
        "merged": "true"
    }, headers=headers)
    
    photos = res.json()
    if not photos or len(photos) < 4:
        # Fallback to random photos if not enough found
        res = requests.get(f"{PHOTOPRISM_URL}/api/v1/photos", params={"count": 50, "quality": 3, "merged": "true"}, headers=headers)
        photos = res.json()
    
    # Pick 4 random photos from the results
    selected = random.sample(photos, min(4, len(photos)))
    print(f"[*] Found {len(selected)} photos for the memory.")
    return selected

def download_image(session_id, hash_str):
    headers = {"X-Session-ID": session_id}
    res = requests.get(f"{PHOTOPRISM_URL}/api/v1/t/{hash_str}/fit_1280", headers=headers)
    return Image.open(io.BytesIO(res.content))

def generate_story(images_base64):
    print("[*] Asking Ollama to write a story...")
    prompt = """You are a nostalgic storyteller. Look at these photos from a past event.
1. Write a very short, catchy 2-4 word Title for the collage.
2. Write a short, nostalgic 2-sentence story about this memory.
Format your response exactly like this:
TITLE: <your title>
STORY: <your story>
"""
    # Just sending the first image to Ollama to keep it fast, or we could send all if the model supports it.
    res = requests.post(OLLAMA_URL, json={
        "model": OLLAMA_MODEL,
        "prompt": prompt,
        "images": [images_base64[0]], # Sending the best photo
        "stream": False
    })
    
    text = res.json().get("response", "")
    
    title = "A Beautiful Memory"
    story = "Looking back at these wonderful moments."
    
    for line in text.split("\n"):
        if line.startswith("TITLE:"): title = line.replace("TITLE:", "").strip()
        if line.startswith("STORY:"): story = line.replace("STORY:", "").strip()
        
    return title, story

def build_vertical_collage(images, title):
    print("[*] Building 9:16 vertical collage...")
    # 1080x1920 (9:16)
    canvas = Image.new("RGB", (1080, 1920), (20, 20, 20))
    
    # Simple 2x2 grid in the middle
    grid_y_offset = 400
    for i, img in enumerate(images):
        img = img.copy()
        # Crop to square
        min_dim = min(img.width, img.height)
        left = (img.width - min_dim) / 2
        top = (img.height - min_dim) / 2
        img = img.crop((left, top, left + min_dim, top + min_dim))
        img = img.resize((500, 500), Image.Resampling.LANCZOS)
        
        x = 30 if i % 2 == 0 else 550
        y = grid_y_offset if i < 2 else grid_y_offset + 520
        canvas.paste(img, (x, y))
        
    # Draw Title
    draw = ImageDraw.Draw(canvas)
    draw.text((50, 200), title, fill=(255, 255, 255), font=ImageFont.load_default())
    
    return canvas

def upload_to_photoprism(session_id, filepath, title, story):
    print("[*] Uploading to PhotoPrism...")
    import_dir = "./storage/import/Memories"
    os.makedirs(import_dir, exist_ok=True)
    
    filename = os.path.basename(filepath)
    import_path = os.path.join(import_dir, filename)
    
    with open(filepath, 'rb') as src, open(import_path, 'wb') as dst:
        dst.write(src.read())
        
    headers = {"X-Session-ID": session_id}
    requests.post(f"{PHOTOPRISM_URL}/api/v1/import", headers=headers, json={"path": "Memories", "move": True})
    print("[+] Memory imported into PhotoPrism!")

def main():
    os.makedirs(EXPORT_DIR, exist_ok=True)
    session_id = get_session()
    
    photos = find_memory_photos(session_id)
    if len(photos) < 4:
        print("[-] Not enough photos to create a collage.")
        return
        
    images = []
    import base64
    b64_images = []
    
    for p in photos:
        hash_str = p.get("Hash")
        img = download_image(session_id, hash_str)
        images.append(img)
        
        buffered = io.BytesIO()
        img.save(buffered, format="JPEG")
        b64_images.append(base64.b64encode(buffered.getvalue()).decode('utf-8'))
        
    title, story = generate_story(b64_images)
    print(f"Title: {title}\nStory: {story}")
    
    collage = build_vertical_collage(images, title)
    
    timestamp = datetime.datetime.now().strftime("%Y%m%d_%H%M%S")
    filepath = os.path.join(EXPORT_DIR, f"Memory_{timestamp}.jpg")
    
    collage.save(filepath, "JPEG")
    with open(filepath.replace(".jpg", ".txt"), "w") as f:
        f.write(f"{title}\n\n{story}")
        
    print(f"[+] Saved locally to {filepath}")
    upload_to_photoprism(session_id, filepath, title, story)

if __name__ == "__main__":
    main()
