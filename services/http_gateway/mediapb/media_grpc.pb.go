package mediapb

import (
	"context"

	grpc "google.golang.org/grpc"
)

// MediaServiceClient is a minimal client interface for calling MediaService.
type MediaServiceClient interface {
	UploadMedia(ctx context.Context, in *UploadMediaRequest, opts ...grpc.CallOption) (*UploadMediaResponse, error)
}

type mediaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMediaServiceClient(cc grpc.ClientConnInterface) MediaServiceClient {
	return &mediaServiceClient{cc}
}

const MediaService_UploadMedia_FullMethodName = "/media.MediaService/UploadMedia"

func (c *mediaServiceClient) UploadMedia(ctx context.Context, in *UploadMediaRequest, opts ...grpc.CallOption) (*UploadMediaResponse, error) {
	out := new(UploadMediaResponse)
	if err := c.cc.Invoke(ctx, MediaService_UploadMedia_FullMethodName, in, out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}
