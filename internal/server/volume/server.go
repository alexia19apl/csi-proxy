package volume

import (
	"context"
	"fmt"

	"k8s.io/klog"

	"github.com/kubernetes-csi/csi-proxy/client/apiversion"
	"github.com/kubernetes-csi/csi-proxy/internal/server/volume/internal"
)

type Server struct {
	hostAPI API
}

type API interface {
	ListVolumesOnDisk(diskID string) (volumeIDs []string, err error)
	// MountVolume mounts the volume at the requested global staging path
	MountVolume(volumeID, path string) error
	// DismountVolume gracefully dismounts a volume
	DismountVolume(volumeID, path string) error
	// IsVolumeFormatted checks if a volume is formatted with NTFS
	IsVolumeFormatted(volumeID string) (bool, error)
	// FormatVolume formats a volume with the provided file system
	FormatVolume(volumeID string) error
	// ResizeVolume performs resizing of the partition and file system for a block based volume
	ResizeVolume(volumeID string, size int64) error
}

func NewServer(hostAPI API) (*Server, error) {
	return &Server{
		hostAPI: hostAPI,
	}, nil
}

func (s *Server) ListVolumesOnDisk(context context.Context, request *internal.ListVolumesOnDiskRequest, version apiversion.Version) (*internal.ListVolumesOnDiskResponse, error) {
	klog.V(5).Infof("ListVolumesOnDisk: Request: %+v", request)
	response := &internal.ListVolumesOnDiskResponse{}

	diskID := request.DiskId
	if diskID == "" {
		return response, fmt.Errorf("disk id empty")
	}
	volumeIDs, err := s.hostAPI.ListVolumesOnDisk(diskID)
	if err != nil {
		return response, err
	}

	response.VolumeIds = volumeIDs
	return response, nil
}

func (s *Server) MountVolume(context context.Context, request *internal.MountVolumeRequest, version apiversion.Version) (*internal.MountVolumeResponse, error) {
	klog.V(5).Infof("MountVolume: Request: %+v", request)
	response := &internal.MountVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		return response, fmt.Errorf("volume id empty")
	}
	path := request.Path
	if path == "" {
		return response, fmt.Errorf("mount path empty")
	}

	err := s.hostAPI.MountVolume(volumeID, path)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *Server) DismountVolume(context context.Context, request *internal.DismountVolumeRequest, version apiversion.Version) (*internal.DismountVolumeResponse, error) {
	klog.V(5).Infof("DismountVolume: Request: %+v", request)
	response := &internal.DismountVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		return response, fmt.Errorf("volume id empty")
	}
	path := request.Path
	if path == "" {
		return response, fmt.Errorf("mount path empty")
	}
	err := s.hostAPI.DismountVolume(volumeID, path)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (s *Server) IsVolumeFormatted(context context.Context, request *internal.IsVolumeFormattedRequest, version apiversion.Version) (*internal.IsVolumeFormattedResponse, error) {
	klog.V(5).Infof("IsVolumeFormatted: Request: %+v", request)
	response := &internal.IsVolumeFormattedResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		return response, fmt.Errorf("volume id empty")
	}
	isFormatted, err := s.hostAPI.IsVolumeFormatted(volumeID)
	if err != nil {
		return response, err
	}
	klog.V(5).Infof("IsVolumeFormatted: return: %v", isFormatted)
	response.Formatted = isFormatted
	return response, nil
}

func (s *Server) FormatVolume(context context.Context, request *internal.FormatVolumeRequest, version apiversion.Version) (*internal.FormatVolumeResponse, error) {
	klog.V(5).Infof("FormatVolume: Request: %+v", request)
	response := &internal.FormatVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		return response, fmt.Errorf("volume id empty")
	}

	err := s.hostAPI.FormatVolume(volumeID)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *Server) ResizeVolume(context context.Context, request *internal.ResizeVolumeRequest, version apiversion.Version) (*internal.ResizeVolumeResponse, error) {
	klog.V(5).Infof("ResizeVolume: Request: %+v", request)
	response := &internal.ResizeVolumeResponse{}

	volumeID := request.VolumeId
	if volumeID == "" {
		return response, fmt.Errorf("volume id empty")
	}
	size := request.Size
	// TODO : Validate size param

	err := s.hostAPI.ResizeVolume(volumeID, size)
	if err != nil {
		return response, err
	}

	return response, nil
}
