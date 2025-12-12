package mediapb

import proto "github.com/golang/protobuf/proto"

// Minimal protobuf message definitions for MediaService client.
// These rely on v1 protobuf reflection via struct tags.

type UploadMediaRequest struct {
	UserId    string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FileBytes []byte `protobuf:"bytes,2,opt,name=file_bytes,json=fileBytes,proto3" json:"file_bytes,omitempty"`
	MimeType  string `protobuf:"bytes,3,opt,name=mime_type,json=mimeType,proto3" json:"mime_type,omitempty"`
	FileName  string `protobuf:"bytes,4,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
}

func (m *UploadMediaRequest) Reset()         { *m = UploadMediaRequest{} }
func (m *UploadMediaRequest) String() string { return proto.CompactTextString(m) }
func (*UploadMediaRequest) ProtoMessage()    {}

type UploadMediaResponse struct {
	MediaId string `protobuf:"bytes,1,opt,name=media_id,json=mediaId,proto3" json:"media_id,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Url     string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
}

func (m *UploadMediaResponse) Reset()         { *m = UploadMediaResponse{} }
func (m *UploadMediaResponse) String() string { return proto.CompactTextString(m) }
func (*UploadMediaResponse) ProtoMessage()    {}
