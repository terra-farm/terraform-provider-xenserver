//
// This file is generated. To change the content of this file, please do not
// apply the change to this file because it will get overwritten. Instead,
// change xenapi.go and execute 'go generate'.
//

package xenAPI

import (
	"fmt"
	"github.com/amfranz/go-xmlrpc-client"
	"reflect"
	"strconv"
	"time"
)

var _ = fmt.Errorf
var _ = xmlrpc.NewClient
var _ = reflect.TypeOf
var _ = strconv.Atoi
var _ = time.UTC

type VMGuestMetricsRecord struct {
  // Unique identifier/object reference
	UUID string
  // version of the OS
	OsVersion map[string]string
  // version of the PV drivers
	PVDriversVersion map[string]string
  // Logical AND of network_paths_optimized and storage_paths_optimized
	PVDriversUpToDate bool
  // True if the network paths are optimized with PV driver
	NetworkPathsOptimized bool
  // True if the storage paths are optimized with PV driver
	StoragePathsOptimized bool
  // This field exists but has no data. Use the memory and memory_internal_free RRD data-sources instead.
	Memory map[string]string
  // This field exists but has no data.
	Disks map[string]string
  // network configuration
	Networks map[string]string
  // anything else
	Other map[string]string
  // Time at which this information was last updated
	LastUpdated time.Time
  // additional configuration
	OtherConfig map[string]string
  // True if the guest is sending heartbeat messages via the guest agent
	Live bool
}

type VMGuestMetricsRef string

// The metrics reported by the guest (as opposed to inferred from outside)
type VMGuestMetricsClass struct {
	client *Client
}

// Return a map of VM_guest_metrics references to VM_guest_metrics records for all VM_guest_metrics instances known to the system.
func (_class VMGuestMetricsClass) GetAllRecords(sessionID SessionRef) (_retval map[VMGuestMetricsRef]VMGuestMetricsRecord, _err error) {
	_method := "VM_guest_metrics.get_all_records"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg)
	if _err != nil {
		return
	}
	_retval, _err = convertVMGuestMetricsRefToVMGuestMetricsRecordMapToGo(_method + " -> ", _result.Value)
	return
}

// Return a list of all the VM_guest_metrics instances known to the system.
func (_class VMGuestMetricsClass) GetAll(sessionID SessionRef) (_retval []VMGuestMetricsRef, _err error) {
	_method := "VM_guest_metrics.get_all"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg)
	if _err != nil {
		return
	}
	_retval, _err = convertVMGuestMetricsRefSetToGo(_method + " -> ", _result.Value)
	return
}

// Remove the given key and its corresponding value from the other_config field of the given VM_guest_metrics.  If the key is not in that Map, then do nothing.
func (_class VMGuestMetricsClass) RemoveFromOtherConfig(sessionID SessionRef, self VMGuestMetricsRef, key string) (_err error) {
	_method := "VM_guest_metrics.remove_from_other_config"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_keyArg, _err := convertStringToXen(fmt.Sprintf("%s(%s)", _method, "key"), key)
	if _err != nil {
		return
	}
	_, _err =  _class.client.APICall(_method, _sessionIDArg, _selfArg, _keyArg)
	return
}

// Add the given key-value pair to the other_config field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) AddToOtherConfig(sessionID SessionRef, self VMGuestMetricsRef, key string, value string) (_err error) {
	_method := "VM_guest_metrics.add_to_other_config"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_keyArg, _err := convertStringToXen(fmt.Sprintf("%s(%s)", _method, "key"), key)
	if _err != nil {
		return
	}
	_valueArg, _err := convertStringToXen(fmt.Sprintf("%s(%s)", _method, "value"), value)
	if _err != nil {
		return
	}
	_, _err =  _class.client.APICall(_method, _sessionIDArg, _selfArg, _keyArg, _valueArg)
	return
}

// Set the other_config field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) SetOtherConfig(sessionID SessionRef, self VMGuestMetricsRef, value map[string]string) (_err error) {
	_method := "VM_guest_metrics.set_other_config"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_valueArg, _err := convertStringToStringMapToXen(fmt.Sprintf("%s(%s)", _method, "value"), value)
	if _err != nil {
		return
	}
	_, _err =  _class.client.APICall(_method, _sessionIDArg, _selfArg, _valueArg)
	return
}

// Get the live field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetLive(sessionID SessionRef, self VMGuestMetricsRef) (_retval bool, _err error) {
	_method := "VM_guest_metrics.get_live"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertBoolToGo(_method + " -> ", _result.Value)
	return
}

// Get the other_config field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetOtherConfig(sessionID SessionRef, self VMGuestMetricsRef) (_retval map[string]string, _err error) {
	_method := "VM_guest_metrics.get_other_config"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToStringMapToGo(_method + " -> ", _result.Value)
	return
}

// Get the last_updated field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetLastUpdated(sessionID SessionRef, self VMGuestMetricsRef) (_retval time.Time, _err error) {
	_method := "VM_guest_metrics.get_last_updated"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertTimeToGo(_method + " -> ", _result.Value)
	return
}

// Get the other field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetOther(sessionID SessionRef, self VMGuestMetricsRef) (_retval map[string]string, _err error) {
	_method := "VM_guest_metrics.get_other"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToStringMapToGo(_method + " -> ", _result.Value)
	return
}

// Get the networks field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetNetworks(sessionID SessionRef, self VMGuestMetricsRef) (_retval map[string]string, _err error) {
	_method := "VM_guest_metrics.get_networks"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToStringMapToGo(_method + " -> ", _result.Value)
	return
}

// Get the disks field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetDisks(sessionID SessionRef, self VMGuestMetricsRef) (_retval map[string]string, _err error) {
	_method := "VM_guest_metrics.get_disks"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToStringMapToGo(_method + " -> ", _result.Value)
	return
}

// Get the memory field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetMemory(sessionID SessionRef, self VMGuestMetricsRef) (_retval map[string]string, _err error) {
	_method := "VM_guest_metrics.get_memory"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToStringMapToGo(_method + " -> ", _result.Value)
	return
}

// Get the storage_paths_optimized field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetStoragePathsOptimized(sessionID SessionRef, self VMGuestMetricsRef) (_retval bool, _err error) {
	_method := "VM_guest_metrics.get_storage_paths_optimized"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertBoolToGo(_method + " -> ", _result.Value)
	return
}

// Get the network_paths_optimized field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetNetworkPathsOptimized(sessionID SessionRef, self VMGuestMetricsRef) (_retval bool, _err error) {
	_method := "VM_guest_metrics.get_network_paths_optimized"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertBoolToGo(_method + " -> ", _result.Value)
	return
}

// Get the PV_drivers_up_to_date field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetPVDriversUpToDate(sessionID SessionRef, self VMGuestMetricsRef) (_retval bool, _err error) {
	_method := "VM_guest_metrics.get_PV_drivers_up_to_date"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertBoolToGo(_method + " -> ", _result.Value)
	return
}

// Get the PV_drivers_version field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetPVDriversVersion(sessionID SessionRef, self VMGuestMetricsRef) (_retval map[string]string, _err error) {
	_method := "VM_guest_metrics.get_PV_drivers_version"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToStringMapToGo(_method + " -> ", _result.Value)
	return
}

// Get the os_version field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetOsVersion(sessionID SessionRef, self VMGuestMetricsRef) (_retval map[string]string, _err error) {
	_method := "VM_guest_metrics.get_os_version"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToStringMapToGo(_method + " -> ", _result.Value)
	return
}

// Get the uuid field of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetUUID(sessionID SessionRef, self VMGuestMetricsRef) (_retval string, _err error) {
	_method := "VM_guest_metrics.get_uuid"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertStringToGo(_method + " -> ", _result.Value)
	return
}

// Get a reference to the VM_guest_metrics instance with the specified UUID.
func (_class VMGuestMetricsClass) GetByUUID(sessionID SessionRef, uuid string) (_retval VMGuestMetricsRef, _err error) {
	_method := "VM_guest_metrics.get_by_uuid"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_uuidArg, _err := convertStringToXen(fmt.Sprintf("%s(%s)", _method, "uuid"), uuid)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _uuidArg)
	if _err != nil {
		return
	}
	_retval, _err = convertVMGuestMetricsRefToGo(_method + " -> ", _result.Value)
	return
}

// Get a record containing the current state of the given VM_guest_metrics.
func (_class VMGuestMetricsClass) GetRecord(sessionID SessionRef, self VMGuestMetricsRef) (_retval VMGuestMetricsRecord, _err error) {
	_method := "VM_guest_metrics.get_record"
	_sessionIDArg, _err := convertSessionRefToXen(fmt.Sprintf("%s(%s)", _method, "session_id"), sessionID)
	if _err != nil {
		return
	}
	_selfArg, _err := convertVMGuestMetricsRefToXen(fmt.Sprintf("%s(%s)", _method, "self"), self)
	if _err != nil {
		return
	}
	_result, _err := _class.client.APICall(_method, _sessionIDArg, _selfArg)
	if _err != nil {
		return
	}
	_retval, _err = convertVMGuestMetricsRecordToGo(_method + " -> ", _result.Value)
	return
}
