#!/usr/bin/env bash
set -e

# ANSI Color Codes
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${CYAN}================================================${NC}"
echo -e "${CYAN}PhotoPrism Setup - Ubuntu Native Launcher${NC}"
echo -e "${CYAN}================================================${NC}"
echo ""

# ----------------------------------------------------
# Step 1: Check Docker Daemon Status
# ----------------------------------------------------
echo -e "${YELLOW}Step 1: Checking Docker status...${NC}"
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}ERROR: Docker is not running or current user does not have permission.${NC}"
    echo -e "${YELLOW}Try running: sudo systemctl start docker${NC}"
    echo -e "${YELLOW}Or ensure your user is in the docker group: sudo usermod -aG docker $USER${NC}"
    exit 1
fi
echo -e "${GREEN}Docker is running!${NC}"
echo ""

# ----------------------------------------------------
# Step 2: Locate PhotoPrism Directory
# ----------------------------------------------------
echo -e "${YELLOW}Step 2: Checking for PhotoPrism drive and path...${NC}"

# List of possible paths to check
CANDIDATE_PATHS=(
    "/run/media/jani/J/PhotoPrism_Portable"
    "/run/media/jani/J/photoprism/photoprism"
    "/media/jani/J/PhotoPrism_Portable"
    "/media/jani/J/photoprism/photoprism"
    "/media/$USER/J/PhotoPrism_Portable"
)

PHOTOPRISM_PATH=""

for path in "${CANDIDATE_PATHS[@]}"; do
    if [ -d "$path" ]; then
        PHOTOPRISM_PATH="$path"
        break
    fi
done

if [ -z "$PHOTOPRISM_PATH" ]; then
    echo -e "${RED}ERROR: Could not find PhotoPrism directory!${NC}"
    echo -e "${YELLOW}Please check if your USB drive 'J' is plugged in and mounted.${NC}"
    echo -e "Looked in:"
    for path in "${CANDIDATE_PATHS[@]}"; do
        echo -e "  - $path"
    done
    echo ""
    echo -e "${YELLOW}Currently mounted drives under /media and /run/media:${NC}"
    ls -la /media/$USER/ 2>/dev/null || ls -la /run/media/$USER/ 2>/dev/null || echo "No removable drives found."
    exit 1
fi

echo -e "${GREEN}Found PhotoPrism at: ${PHOTOPRISM_PATH}${NC}"

# Check for docker-compose.yml / compose.yaml
if [ ! -f "$PHOTOPRISM_PATH/docker-compose.yml" ] && [ ! -f "$PHOTOPRISM_PATH/compose.yaml" ]; then
    echo -e "${RED}ERROR: No docker-compose.yml or compose.yaml found inside $PHOTOPRISM_PATH${NC}"
    exit 1
fi

echo ""

# ----------------------------------------------------
# Step 3: Starting PhotoPrism
# ----------------------------------------------------
echo -e "${YELLOW}Step 3: Starting PhotoPrism with Docker Compose...${NC}"
cd "$PHOTOPRISM_PATH"

if docker compose up -d; then
    echo ""
    echo -e "${CYAN}================================================${NC}"
    echo -e "${GREEN}PhotoPrism is Starting!${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
    
    echo -e "${YELLOW}Waiting for services to initialize (15 seconds)...${NC}"
    sleep 15

    echo -e "${GREEN}PhotoPrism should now be running!${NC}"
    echo ""
    echo -e "${YELLOW}Access PhotoPrism at:${NC}"
    echo -e "${CYAN}  http://localhost:2342${NC}"
    echo ""

    # Open in default Linux browser if xdg-open exists
    if command -v xdg-open >/dev/null 2>&1; then
        echo -e "${YELLOW}Opening in your default browser...${NC}"
        xdg-open "http://localhost:2342" >/dev/null 2>&1 &
    fi

    echo ""
    echo -e "${YELLOW}Useful management commands:${NC}"
    echo -e "  View logs:  docker compose -f \"$PHOTOPRISM_PATH/docker-compose.yml\" logs -f"
    echo -e "  Stop:       docker compose -f \"$PHOTOPRISM_PATH/docker-compose.yml\" down"
    echo -e "  Restart:    docker compose -f \"$PHOTOPRISM_PATH/docker-compose.yml\" restart"
    echo ""
else
    echo -e "${RED}ERROR: Failed to start PhotoPrism!${NC}"
    echo -e "${YELLOW}Last 50 log lines:${NC}"
    docker compose logs --tail=50
    exit 1
fi