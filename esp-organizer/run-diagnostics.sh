#!/bin/bash
# filepath: /Users/michaelstewart/Coding/assist/esp-organizer/run-diagnostics.sh

# Set up colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== ESP Organizer AWS Service Diagnostics ===${NC}"

# Load environment variables if .env exists
if [ -f .env ]; then
  echo -e "${GREEN}Loading environment variables from .env${NC}"
  set -a
  source .env
  set +a
else
  echo -e "${RED}No .env file found, using environment variables${NC}"
fi

# Check for AWS credentials
if [ -z "$AWS_ACCESS_KEY_ID" ] || [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
  echo -e "${RED}AWS credentials not found in environment${NC}"
  echo "Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY"
  exit 1
fi

# Run the diagnostic tool
echo -e "\n${YELLOW}1. Basic AWS configuration test${NC}"
go run cmd/aws-diagnostics/main.go -textract=false -bedrock=true -embeddings=false -upload=false

echo -e "\n${YELLOW}2. Testing S3 file upload${NC}"
go run cmd/aws-diagnostics/main.go -textract=false -bedrock=false -embeddings=false -upload=true

echo -e "\n${YELLOW}3. Testing AWS Bedrock embeddings${NC}"
go run cmd/aws-diagnostics/main.go -textract=false -bedrock=false -embeddings=true -upload=false

echo -e "\n${YELLOW}4. Testing AWS Textract${NC}"
go run cmd/aws-diagnostics/main.go -textract=true -bedrock=false -embeddings=false -upload=false

# Run a full test at the end
echo -e "\n${YELLOW}5. Running complete integration test${NC}"
go run cmd/aws-diagnostics/main.go -textract=true -bedrock=true -embeddings=true -upload=true

echo -e "\n${GREEN}✓ Diagnostics completed${NC}"
echo "Review the output above to identify any AWS service connectivity issues."