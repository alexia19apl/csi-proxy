package volume

import (
	"fmt"
	"os/exec"
	"strings"

	"k8s.io/klog"
)

const formatFilesystem = "ntfs"

// VolAPIImplementor - struct for implementing the internal Volume APIs
type VolAPIImplementor struct{}

// New - Construct a new Volume API Implementation.
func New() VolAPIImplementor {
	return VolAPIImplementor{}
}

func runExec(cmd string) ([]byte, error) {
	klog.V(5).Infof("Running command: %s", cmd)
	out, err := exec.Command("powershell", "/c", cmd).CombinedOutput()
	klog.V(5).Infof("Result: %s. Len: %d. Error: %v.", string(out), len(string(out)), err)
	return out, err
}

// ListVolumesOnDisk - returns back list of volumes(volumeIDs) in the disk (requested in diskID).
func (VolAPIImplementor) ListVolumesOnDisk(diskID string) (volumeIDs []string, err error) {
	cmd := fmt.Sprintf("(Get-Disk -DeviceId %s |Get-Partition | Get-Volume).UniqueId", diskID)
	out, err := runExec(cmd)
	if err != nil {
		return []string{}, err
	}

	volumeIds := strings.Split(strings.TrimSpace(string(out)), "\r\n")
	return volumeIds, nil
}

// FormatVolume - Format a volume with a pre specified filesystem (typically ntfs)
func (VolAPIImplementor) FormatVolume(volumeID string) (err error) {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Format-Volume  -FileSystem %s -Confirm:$false", volumeID, formatFilesystem)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error formatting volume %s: %v, %v", volumeID, out, err)
	}
	// TODO: Do we need to handle anything for len(out) == 0
	return nil
}

// IsVolumeFormatted - Check if the volume is formatted with the pre specified filesystem(typically ntfs).
func (VolAPIImplementor) IsVolumeFormatted(volumeID string) (bool, error) {
	cmd := fmt.Sprintf("(Get-Volume -UniqueId \"%s\" -ErrorAction Stop).FileSystemType", volumeID)
	out, err := runExec(cmd)
	if err != nil {
		return false, fmt.Errorf("error checking if volume is formatted %s: %s, %v", volumeID, string(out), err)
	}
	if len(out) == 0 || !strings.EqualFold(strings.TrimSpace(string(out)), formatFilesystem) {
		return false, nil
	}
	return true, nil
}

// MountVolume - mounts a volume to a path. This is done using the Add-PartitionAccessPath for presenting the volume via a path.
func (VolAPIImplementor) MountVolume(volumeID, path string) error {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-Partition | Add-PartitionAccessPath -AccessPath %s", volumeID, path)
	_, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error mount volume %s to path %s. err: %v", volumeID, path, err)
	}
	return nil
}

// DismountVolume - unmounts the volume path by removing the patition access path
func (VolAPIImplementor) DismountVolume(volumeID, path string) error {
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-Partition | Remove-PartitionAccessPath -AccessPath %s", volumeID, path)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error getting driver letter to mount volume %s: %v, %v", volumeID, out, err)
	}
	return nil
}

// ResizeVolume - resize the volume to the size specified as parameter.
func (VolAPIImplementor) ResizeVolume(volumeID string, size int64) error {
	// TODO: Check the size of the resize
	// TODO: We have to get the right partition.
	cmd := fmt.Sprintf("Get-Volume -UniqueId \"%s\" | Get-partition | Resize-Partition -Size %d", volumeID, size)
	out, err := runExec(cmd)
	if err != nil {
		return fmt.Errorf("error resizing volume %s: %s, %v", volumeID, string(out), err)
	}
	return nil
}
