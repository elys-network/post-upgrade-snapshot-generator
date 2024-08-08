package uploadsnapshot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

func UploadSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-snapshot [file_path] [flags]",
		Short: "Upload snapshot file to R2",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// get args
			filePath := args[0]
			if filePath == "" {
				log.Fatalf(types.ColorRed + "file path is required")
			}

			// Fetch credentials and configuration from environment variables
			accessKey := os.Getenv("R2_ACCESS_KEY")
			secretKey := os.Getenv("R2_SECRET_KEY")
			accountId := os.Getenv("R2_ACCOUNT_ID")
			bucketName := os.Getenv("R2_BUCKET_NAME")

			// Ensure all required environment variables are set
			if accessKey == "" || secretKey == "" || accountId == "" || bucketName == "" {
				fmt.Println("Please set R2_ACCESS_KEY, R2_SECRET_KEY, R2_ACCOUNT_ID, and R2_BUCKET_NAME environment variables")
				os.Exit(1)
			}

			// Load AWS configuration with credentials
			cfg, err := config.LoadDefaultConfig(
				context.TODO(),
				config.WithCredentialsProvider(
					credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
				),
				config.WithRegion("auto"),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to load configuration, %v", err)
				os.Exit(1)
			}

			// Create an S3 client
			client := s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
			})

			// Create a presigner
			presignClient := s3.NewPresignClient(client)
			presigner := types.Presigner{PresignClient: presignClient}

			// Open the file to upload
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file %q, %v", filePath, err)
				os.Exit(1)
			}
			defer file.Close()

			// Get the file size
			fileInfo, err := file.Stat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get file stats %q, %v", filePath, err)
				os.Exit(1)
			}
			fileSize := fileInfo.Size()

			// Create a progress bar
			p := mpb.New(mpb.WithWidth(60), mpb.WithOutput(os.Stdout), mpb.WithAutoRefresh())
			bar := p.AddBar(fileSize,
				mpb.PrependDecorators(
					decor.Name("Upload progress:"),
					decor.Percentage(decor.WC{W: 5}),
				),
				mpb.AppendDecorators(
					decor.CountersKibiByte("% .2f / % .2f"),
					decor.Name("  "),
					decor.AverageSpeed(decor.SizeB1024(0), "% .2f", decor.WC{W: 7}),
					decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 12}),
				),
			)

			// Create a proxy reader
			proxyReader := bar.ProxyReader(file)
			defer proxyReader.Close()

			// Upload the file
			key := filepath.Base(filePath)

			presignedPutRequest, err := presigner.PutObject(bucketName, key, 600)
			if err != nil {
				panic(err)
			}

			uploadReq, err := http.NewRequest("PUT", presignedPutRequest.URL, proxyReader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create upload request, %v", err)
				os.Exit(1)
			}
			uploadReq.ContentLength = fileSize

			resp, err := http.DefaultClient.Do(uploadReq)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to upload file using presigned URL, %v", err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			// Wait for the bar to complete
			p.Wait()

			fmt.Printf("Successfully uploaded %q to %q\n", filePath, bucketName)
		},
	}

	return cmd
}
