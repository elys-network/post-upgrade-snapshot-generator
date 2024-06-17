package deletesnapshot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"

	"github.com/elys-network/post-upgrade-snapshot-generator/types"
)

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
			s3URL := os.Getenv("R2_ENDPOINT")
			bucketName := os.Getenv("R2_BUCKET_NAME")

			// Ensure all required environment variables are set
			if accessKey == "" || secretKey == "" || s3URL == "" || bucketName == "" {
				fmt.Println("Please set R2_ACCESS_KEY, R2_SECRET_KEY, R2_ENDPOINT, and R2_BUCKET_NAME environment variables")
				os.Exit(1)
			}

			// Load AWS configuration with credentials
			cfg, err := config.LoadDefaultConfig(
				context.TODO(),
				config.WithCredentialsProvider(
					credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
				),
				config.WithRegion("auto"), // Ensure this region is appropriate or set it via environment variable if needed
				config.WithEndpointResolverWithOptions(
					aws.EndpointResolverWithOptionsFunc(
						func(service, region string, options ...interface{}) (aws.Endpoint, error) {
							return aws.Endpoint{
								URL: s3URL,
							}, nil
						},
					),
				),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to load configuration, %v", err)
				os.Exit(1)
			}

			// Create an S3 client
			client := s3.NewFromConfig(cfg)

			// Replace '/' with '_' in the branch name
			safeBranchName := strings.ReplaceAll(branchName, "/", "_")

			// Construct the key for the snapshot file
			key := fmt.Sprintf("elys-snapshot-%s.tar.lz4", safeBranchName)

			// Delete the file
			_, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
				Bucket: &bucketName,
				Key:    &key,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to delete file %q from bucket %q, %v", key, bucketName, err)
				os.Exit(1)
			}

			fmt.Printf("Successfully deleted %q from %q\n", key, bucketName)
		},
	}

	return cmd
}
