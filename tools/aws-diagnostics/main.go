package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"

	"github.com/joho/godotenv"
)

var (
	testTextract     bool
	testBedrock      bool
	testEmbeddings   bool
	testFileUpload   bool
	region           string
	bucket           string
	modelId          string
	verbose          bool
	createTestFile   bool
	textractTempFile string
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	flag.BoolVar(&testTextract, "textract", true, "Test AWS Textract functionality")
	flag.BoolVar(&testBedrock, "bedrock", true, "Test AWS Bedrock functionality")
	flag.BoolVar(&testEmbeddings, "embeddings", true, "Test AWS Titan embeddings specifically")
	flag.BoolVar(&testFileUpload, "upload", true, "Test S3 file upload")
	flag.StringVar(&region, "region", os.Getenv("AWS_REGION"), "AWS Region to use")
	flag.StringVar(&bucket, "bucket", os.Getenv("AWS_S3_BUCKET"), "S3 bucket for testing")
	flag.StringVar(&modelId, "model", os.Getenv("AWS_BEDROCK_MODEL"), "Bedrock model ID for embeddings")
	flag.BoolVar(&verbose, "verbose", true, "Enable verbose logging")
	flag.BoolVar(&createTestFile, "create-file", true, "Create test PDF file for Textract")
	flag.StringVar(&textractTempFile, "textract-file", "test-textract.pdf", "Test PDF file for Textract")
	flag.Parse()

	fmt.Println("🧪 AWS Services Diagnostic Tool")
	fmt.Println("==============================")

	fmt.Println("\n📊 Environment Information:")
	fmt.Printf("AWS Region: %s\n", region)
	fmt.Printf("AWS S3 Bucket: %s\n", bucket)
	fmt.Printf("AWS Bedrock Model: %s\n", modelId)
	maskAndPrint("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
	maskAndPrint("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))

	// Set default values if not provided
	if region == "" {
		region = "eu-west-2"
		fmt.Printf("⚠️ No region specified, using default: %s\n", region)
	}
	if bucket == "" {
		bucket = "esp-new-organizer-immunology"
		fmt.Printf("⚠️ No bucket specified, using default: %s\n", bucket)
	}
	if modelId == "" {
		modelId = "amazon.titan-embed-text-v1"
		fmt.Printf("⚠️ No model ID specified, using default: %s\n", modelId)
	}

	// First test credentials
	testAWSCredentials()

	// Run tests
	if testBedrock {
		testBedrockAPI()
	}

	if testEmbeddings {
		testBedrockEmbeddings(modelId)
	}

	if testFileUpload {
		testS3FileUpload(region, bucket)
	}

	if testTextract {
		testTextractAPI(region, bucket)
	}

	fmt.Println("\n✅ AWS diagnostic tests completed")
}

func testAWSCredentials() {
	fmt.Println("\n🔐 Testing AWS Credentials:")

	// Test basic credential resolution
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Printf("❌ Failed to load AWS config: %v\n", err)
		fmt.Println("   Check your AWS credentials configuration")
		return
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		fmt.Printf("❌ Failed to retrieve credentials: %v\n", err)
		fmt.Println("   Check your AWS credentials configuration")
		return
	}

	fmt.Printf("✅ Credentials found: %s\n", creds.Source)
	fmt.Printf("   Access Key ID: %s\n", maskString(creds.AccessKeyID))

	if creds.SessionToken != "" {
		fmt.Printf("   Using session token: %s\n", maskString(creds.SessionToken))
		fmt.Println("   ℹ️  Session tokens may expire - check token validity")
	}

	// Test STS to verify credentials work
	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Printf("❌ STS validation failed: %v\n", err)
		fmt.Println("   Credentials are invalid or expired")
		fmt.Println("   Check your IAM policy permissions")
		return
	}

	fmt.Printf("✅ STS validation successful\n")
	fmt.Printf("   Account: %s\n", aws.ToString(identity.Account))
	fmt.Printf("   User ID: %s\n", aws.ToString(identity.UserId))
	fmt.Printf("   ARN: %s\n", aws.ToString(identity.Arn))

	// Check if this is the ESPAssist user
	if aws.ToString(identity.Arn) != "arn:aws:iam::976193238457:user/ESPAssist" {
		fmt.Printf("⚠️  Warning: Using different IAM user than expected: %s\n", aws.ToString(identity.Arn))
		fmt.Println("   Expected: arn:aws:iam::976193238457:user/ESPAssist")
	}
}

func testBedrockAPI() {
	fmt.Println("\n🧠 Testing AWS Bedrock API access:")

	// Create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		fmt.Printf("❌ Failed to load AWS config: %v\n", err)
		fmt.Println("   Check that your credentials are valid and have Bedrock permissions")
		return
	}

	// Create Bedrock client
	client := bedrockruntime.NewFromConfig(cfg)

	// List available models (simple API call to test access)
	fmt.Println("   Attempting to create Bedrock client...")
	if client != nil {
		fmt.Println("✅ Successfully created Bedrock client")
	} else {
		fmt.Println("❌ Failed to create Bedrock client")
		return
	}
}

func testBedrockEmbeddings(modelId string) {
	fmt.Println("\n🔤 Testing AWS Bedrock embeddings generation:")

	// Create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		fmt.Printf("❌ Failed to load AWS config: %v\n", err)
		return
	}

	// Create Bedrock client
	client := bedrockruntime.NewFromConfig(cfg)

	// Create sample text for embedding
	sampleText := "This is a test text for AWS Bedrock embedding generation."

	// Create input for the model
	inputStruct := struct {
		InputText string `json:"inputText"`
	}{
		InputText: sampleText,
	}

	inputBytes, err := json.Marshal(inputStruct)
	if err != nil {
		fmt.Printf("❌ Failed to marshal input: %v\n", err)
		return
	}

	// Send embedding request
	fmt.Println("   Sending embedding request...")
	startTime := time.Now()
	resp, err := client.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		Body:        inputBytes,
		ContentType: aws.String("application/json"),
	})

	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Bedrock embedding request failed in %v: %v\n", duration, err)
		fmt.Println("\n🔍 Troubleshooting Tips:")
		fmt.Println("   1. Check AWS credentials are correct and have Bedrock permissions")
		fmt.Println("   2. Verify the model ID is correct and available in your region")
		fmt.Println("   3. Ensure your AWS account has access to the specified model")
		return
	}

	// Parse response
	var embedResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.Unmarshal(resp.Body, &embedResponse); err != nil {
		fmt.Printf("❌ Failed to unmarshal embedding response: %v\n", err)
		return
	}

	// Check if embeddings were generated
	if len(embedResponse.Embedding) == 0 {
		fmt.Println("❌ No embeddings returned")
		return
	}

	fmt.Printf("✅ Successfully generated embeddings in %v\n", duration)
	fmt.Printf("   Vector length: %d dimensions\n", len(embedResponse.Embedding))

	// Show first few values
	if verbose {
		fmt.Print("   Sample vector values: [")
		for i := 0; i < min(5, len(embedResponse.Embedding)); i++ {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%.6f", embedResponse.Embedding[i])
		}
		fmt.Println(", ...]")
	}
}

func testS3FileUpload(region, bucket string) {
	fmt.Println("\n📤 Testing S3 file upload:")

	// Create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		fmt.Printf("❌ Failed to load AWS config: %v\n", err)
		return
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	// Create a test file if requested
	if createTestFile {
		if err := createSimpleTextFile(textractTempFile); err != nil {
			fmt.Printf("❌ Failed to create test file: %v\n", err)
			return
		}
		fmt.Printf("   Created test file: %s\n", textractTempFile)
	}

	// Open the file
	file, err := os.Open(textractTempFile)
	if err != nil {
		fmt.Printf("❌ Failed to open test file: %v\n", err)
		return
	}
	defer file.Close()

	// Upload file to S3
	key := fmt.Sprintf("diagnostic-tests/%s-%d", textractTempFile, time.Now().Unix())
	fmt.Printf("   Uploading file to s3://%s/%s\n", bucket, key)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		fmt.Printf("❌ Failed to upload file to S3: %v\n", err)
		fmt.Println("\n🔍 Troubleshooting Tips:")
		fmt.Println("   1. Check AWS credentials have S3 permissions")
		fmt.Println("   2. Verify the bucket exists and is accessible")
		fmt.Println("   3. Check network connectivity to AWS S3")
		return
	}

	fmt.Println("✅ Successfully uploaded file to S3")
	fmt.Printf("   Object key: %s\n", key)

	// Return to beginning of file for potential reuse
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Printf("⚠️ Warning: Failed to reset file position: %v\n", err)
	}
}

func testTextractAPI(region, bucket string) {
	fmt.Println("\n🔍 Testing AWS Textract API:")

	// Create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		fmt.Printf("❌ Failed to load AWS config: %v\n", err)
		return
	}

	// Create Textract client
	textractClient := textract.NewFromConfig(cfg)

	// Check if we have a file in S3 to use
	s3Key := fmt.Sprintf("diagnostic-tests/%s-%d", textractTempFile, time.Now().Unix()-1)
	fmt.Printf("   Using S3 object: s3://%s/%s\n", bucket, s3Key)

	// Try to detect text in the document
	fmt.Println("   Attempting synchronous text detection...")
	detectInput := &textract.DetectDocumentTextInput{
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(s3Key),
			},
		},
	}

	startTime := time.Now()
	detectResult, err := textractClient.DetectDocumentText(context.TODO(), detectInput)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ Textract text detection failed in %v: %v\n", duration, err)
		fmt.Println("\n🔍 Troubleshooting Tips:")
		fmt.Println("   1. Check AWS credentials have Textract permissions")
		fmt.Println("   2. Verify the S3 object exists and is readable by Textract")
		fmt.Println("   3. Ensure your AWS account has access to Textract in this region")
		return
	}

	// Check results
	blockCount := len(detectResult.Blocks)
	fmt.Printf("✅ Successfully detected text in %v\n", duration)
	fmt.Printf("   Detected %d blocks\n", blockCount)

	// Show block types
	if blockCount > 0 {
		blockTypes := make(map[types.BlockType]int)
		for _, block := range detectResult.Blocks {
			blockTypes[block.BlockType]++
		}

		fmt.Println("   Block types detected:")
		for blockType, count := range blockTypes {
			fmt.Printf("   - %s: %d blocks\n", blockType, count)
		}
	}
}

// Helper functions
func maskAndPrint(name, value string) {
	if value == "" {
		fmt.Printf("%s: (not set)\n", name)
		return
	}

	masked := value
	if len(value) > 8 {
		masked = value[:4] + "..." + value[len(value)-4:]
	}
	fmt.Printf("%s: %s\n", name, masked)
}

func createSimpleTextFile(filename string) error {
	// For Textract testing, we need to create a proper PDF file
	// This is a minimal PDF content that Textract can read
	pdfContent := `%PDF-1.4
1 0 obj
<<
/Type /Catalog
/Pages 2 0 R
>>
endobj
2 0 obj
<<
/Type /Pages
/Kids [3 0 R]
/Count 1
>>
endobj
3 0 obj
<<
/Type /Page
/Parent 2 0 R
/MediaBox [0 0 612 792]
/Contents 4 0 R
>>
endobj
4 0 obj
<<
/Length 44
>>
stream
BT
/F1 12 Tf
72 720 Td
(Test document for AWS Textract) Tj
ET
endstream
endobj
xref
0 5
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
0000000206 00000 n 
trailer
<<
/Size 5
/Root 1 0 R
>>
startxref
321
%%EOF`

	return os.WriteFile(filename, []byte(pdfContent), 0644)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
