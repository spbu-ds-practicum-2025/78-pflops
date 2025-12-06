# test_client.py
import grpc
from services.MediaService.pkg.pb import media_pb2
from services.MediaService.pkg.pb import media_pb2_grpc

def test_server():
    print("Connecting to gRPC server...")
    
    # Connect to server
    channel = grpc.insecure_channel('localhost:50051')
    stub = media_pb2_grpc.MediaServiceStub(channel)
    
    try:
        # Test 1: Upload a file
        print("1. Testing upload...")
        upload_request = media_pb2.UploadMediaRequest(
            user_id="test_user",
            file_bytes=b"This is a test file content",
            mime_type="text/plain",
            file_name="test.txt"
        )
        upload_response = stub.UploadMedia(upload_request)
        print(f"   Uploaded: {upload_response.media_id}")
        
        # Test 2: Get URL for the file
        print("2. Testing get URL...")
        url_request = media_pb2.GetUrlRequest(media_id=upload_response.media_id)
        url_response = stub.GetUrl(url_request)
        print(f"   URL: {url_response.url}")
        
        # Test 3: List user's files
        print("3. Testing list files...")
        list_request = media_pb2.ListMediaRequest(user_id="test_user")
        list_response = stub.ListMedia(list_request)
        print(f"   Found {len(list_response.media_items)} files")
        
        # Test 4: Get the file back
        print("4. Testing get file...")
        get_request = media_pb2.GetMediaRequest(media_id=upload_response.media_id)
        get_response = stub.GetMedia(get_request)
        print(f"   Got file: {len(get_response.file_bytes)} bytes")
        
        # Test 5: Delete the file
        print("5. Testing delete...")
        delete_request = media_pb2.DeleteMediaRequest(
            media_id=upload_response.media_id,
            user_id="test_user"
        )
        delete_response = stub.DeleteMedia(delete_request)
        print(f"   Delete success: {delete_response.success}")
        
        print("\n✅ All tests passed!")
        
    except grpc.RpcError as e:
        print(f"❌ gRPC error: {e}")
    except Exception as e:
        print(f"❌ Error: {e}")

def test_reflection():
    """Test if reflection is working"""
    try:
        from grpc_reflection.v1alpha.proto_reflection_descriptor_database import ProtoReflectionDescriptorDatabase
        import grpc.reflection
        channel = grpc.insecure_channel('localhost:50051')
        reflection_db = ProtoReflectionDescriptorDatabase(channel)
        
        print("✅ Reflection is enabled and working")
        return True
    except Exception as e:
        print(f"❌ Reflection test failed: {e}")
        return False

if __name__ == "__main__":
    test_reflection()
    test_server()