package deletesnapshot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

// delete-snapshot command does the following:
// 1. delete an existing snapshot file from R2 bucket
func DeleteSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-snapshot [branch_name] [flags]",
		Short: "Delete an existing snapshot file from R2 bucket",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// get args
			branchName := args[0]
			if branchName == "" {
				log.Fatalf(types.ColorRed + "branch name is required")
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

			// Replace '/' with '_' in the branch name
			safeBranchName := strings.ReplaceAll(branchName, "/", "_")

			// Construct the key for the snapshot file
			key := fmt.Sprintf("elys-snapshot-%s.tar.lz4", safeBranchName)

			presignedDeleteRequest, err := presigner.DeleteObject(bucketName, key)
			if err != nil {
				panic(err)
			}

			delRequest, err := http.NewRequest("DELETE", presignedDeleteRequest.URL, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create delete request, %v", err)
				os.Exit(1)
			}

			resp, err := http.DefaultClient.Do(delRequest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to delete file %q from bucket %q using presigned URL, %v", key, bucketName, err)
				os.Exit(1)
			}
			defer resp.Body.Close()

			fmt.Printf("Successfully deleted %q from %q\n", key, bucketName)
		},
	}

	return cmd
}
