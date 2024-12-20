#!/bin/bash

# R2 credentials
export R2_ACCESS_KEY=
export R2_SECRET_KEY=
export R2_BUCKET_NAME=
export R2_ENDPOINT=

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="post-upgrade-snapshot-generator"
SNAPSHOT_PATH="/tmp/snapshot.tar.lz4"
INFO_JSON_PATH="/tmp/info.json"
SERVICE_NAME="elysd"
HOME_PATH="$HOME/.elys"

# Check if the binary exists
if ! command -v $BINARY_NAME &> /dev/null; then
    echo -e "${RED}Binary $BINARY_NAME not found. Please install post-upgrade-snapshot-generator first.${NC}"
    exit 1
fi

# Check if the service exists
if ! systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${RED}Service $SERVICE_NAME not found. Please install the service first.${NC}"
    exit 1
fi

# Check if the home path exists
if [ ! -d $HOME_PATH ]; then
    echo -e "${RED}Home path $HOME_PATH not found. Please set the home path first.${NC}"
    exit 1
fi

# Check if rclone is installed
if ! command -v rclone &> /dev/null; then
    echo -e "${RED}Rclone not found. Please install rclone first.${NC}"
    exit 1
fi

# Check if lz4 is installed
if ! command -v lz4 &> /dev/null; then
    echo -e "${RED}Lz4 not found. Please install lz4 first.${NC}"
    exit 1
fi

echo -e "${YELLOW}Starting snapshot creation and upload process...${NC}"

# Stop the systemd service
echo -e "${YELLOW}Stopping $SERVICE_NAME service...${NC}"
sudo systemctl stop $SERVICE_NAME
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to stop $SERVICE_NAME service${NC}"
    exit 1
fi
echo -e "${GREEN}Successfully stopped $SERVICE_NAME service${NC}"

# Wait a few seconds to ensure the service is fully stopped
sleep 5

# Create the snapshot
echo -e "${YELLOW}Creating snapshot...${NC}"
$BINARY_NAME create-snapshot $SNAPSHOT_PATH --home $HOME_PATH
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to create snapshot${NC}"
    # Restart service before exiting
    sudo systemctl start $SERVICE_NAME
    exit 1
fi
echo -e "${GREEN}Successfully created snapshot${NC}"

# Restart the service
echo -e "${YELLOW}Restarting $SERVICE_NAME service...${NC}"
sudo systemctl start $SERVICE_NAME
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to restart $SERVICE_NAME service${NC}"
    exit 1
fi

# Wait a few seconds and check if service is running
sleep 5
if ! systemctl is-active --quiet $SERVICE_NAME; then
    echo -e "${RED}Warning: $SERVICE_NAME service failed to start properly${NC}"
    exit 1
fi

# Get snapshot file size
file_size=$(ls -lh $SNAPSHOT_PATH | awk '{print $5}')

# Generate info.json
echo -e "${YELLOW}Generating ${INFO_JSON_PATH}...${NC}"
# Fetch the JSON data from the Elys testnet RPC endpoint
input_json=$(curl -s https://rpc.testnet.elys.network/abci_info?)

# Extract the relevant fields using jq
block_height=$(echo "$input_json" | jq -r '.result.response.last_block_height')
data=$(echo "$input_json" | jq -r '.result.response.data')
version=$(echo "$input_json" | jq -r '.result.response.version')

# Construct the desired output
created_at=$(date -Iseconds)
jq -n \
    --arg blockHeight "$block_height" \
    --arg fileName "$SNAPSHOT_PATH" \
    --arg fileSize "$file_size" \
    --arg createdAt "$created_at" \
    --arg version "$version" \
    '{blockHeight: $blockHeight, fileName: $fileName, fileSize: $fileSize, createdAt: $createdAt, version: $version}' > $INFO_JSON_PATH

if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to generate ${INFO_JSON_PATH}${NC}"
    # Cleanup before exiting
    rm -f $SNAPSHOT_PATH
    exit 1
fi
echo -e "${GREEN}Successfully generated ${INFO_JSON_PATH}${NC}"

# Set rclone environment variables
export RCLONE_CONFIG_R2_TYPE=s3
export RCLONE_CONFIG_R2_PROVIDER=Cloudflare
export RCLONE_CONFIG_R2_ACCESS_KEY_ID=$R2_ACCESS_KEY
export RCLONE_CONFIG_R2_SECRET_ACCESS_KEY=$R2_SECRET_KEY
export RCLONE_CONFIG_R2_REGION=enam
export RCLONE_CONFIG_R2_ENDPOINT=$R2_ENDPOINT

# Upload the snapshot
echo -e "${YELLOW}Uploading snapshot...${NC}"
rclone -vv copy $SNAPSHOT_PATH r2:${R2_BUCKET_NAME}/
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to upload snapshot${NC}"
    # Cleanup before exiting
    rm -f $SNAPSHOT_PATH $INFO_JSON_PATH
    exit 1
fi
echo -e "${GREEN}Successfully uploaded snapshot${NC}"

# Upload info.json
echo -e "${YELLOW}Uploading ${INFO_JSON_PATH}...${NC}"
rclone -vv copy $INFO_JSON_PATH r2:${R2_BUCKET_NAME}/
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to upload ${INFO_JSON_PATH}${NC}"
    # Cleanup before exiting
    rm -f $SNAPSHOT_PATH $INFO_JSON_PATH
    exit 1
fi
echo -e "${GREEN}Successfully uploaded ${INFO_JSON_PATH}${NC}"

# Clean up the temporary files
echo -e "${YELLOW}Cleaning up temporary files...${NC}"
rm -f $SNAPSHOT_PATH $INFO_JSON_PATH
if [ $? -ne 0 ]; then
    echo -e "${RED}Warning: Failed to clean up temporary files${NC}"
fi

echo -e "${GREEN}Process completed successfully!${NC}"
echo -e "${GREEN}- $SNAPSHOT_PATH created and uploaded${NC}"
echo -e "${GREEN}- $SERVICE_NAME service restarted${NC}"
echo -e "${GREEN}- $INFO_JSON_PATH created and uploaded${NC}"
echo -e "${GREEN}- Temporary files cleaned up${NC}"