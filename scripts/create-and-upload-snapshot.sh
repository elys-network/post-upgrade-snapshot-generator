#!/bin/bash

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

# Get snapshot file size
file_size=$(ls -lh $SNAPSHOT_PATH | awk '{print $5}')

# Generate info.json
echo -e "${YELLOW}Generating info.json...${NC}"
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
    --arg fileName "snapshot.tar.lz4" \
    --arg fileSize "$file_size" \
    --arg createdAt "$created_at" \
    --arg version "$version" \
    '{blockHeight: $blockHeight, fileName: $fileName, fileSize: $fileSize, createdAt: $createdAt, version: $version}' > $INFO_JSON_PATH

if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to generate info.json${NC}"
    # Cleanup and restart service before exiting
    rm -f $SNAPSHOT_PATH
    sudo systemctl start $SERVICE_NAME
    exit 1
fi
echo -e "${GREEN}Successfully generated info.json${NC}"

# Upload the snapshot
echo -e "${YELLOW}Uploading snapshot...${NC}"
$BINARY_NAME upload-snapshot $SNAPSHOT_PATH
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to upload snapshot${NC}"
    # Cleanup and restart service before exiting
    rm -f $SNAPSHOT_PATH $INFO_JSON_PATH
    sudo systemctl start $SERVICE_NAME
    exit 1
fi
echo -e "${GREEN}Successfully uploaded snapshot${NC}"

# Upload info.json
echo -e "${YELLOW}Uploading info.json...${NC}"
$BINARY_NAME upload-snapshot $INFO_JSON_PATH
if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to upload info.json${NC}"
    # Cleanup and restart service before exiting
    rm -f $SNAPSHOT_PATH $INFO_JSON_PATH
    sudo systemctl start $SERVICE_NAME
    exit 1
fi
echo -e "${GREEN}Successfully uploaded info.json${NC}"

# Clean up the temporary files
echo -e "${YELLOW}Cleaning up temporary files...${NC}"
rm -f $SNAPSHOT_PATH $INFO_JSON_PATH
if [ $? -ne 0 ]; then
    echo -e "${RED}Warning: Failed to clean up temporary files${NC}"
fi

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

echo -e "${GREEN}Process completed successfully!${NC}"
echo -e "${GREEN}- Snapshot created and uploaded${NC}"
echo -e "${GREEN}- info.json created and uploaded${NC}"
echo -e "${GREEN}- Temporary files cleaned up${NC}"
echo -e "${GREEN}- $SERVICE_NAME service restarted${NC}" 