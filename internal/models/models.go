package models

import "time"

type VehicleType string

const (
	VehicleTypeNormal    VehicleType = "normal"
	VehicleTypeSevere    VehicleType = "severe"
	VehicleTypeNeonatal  VehicleType = "neonatal"
)

type VehicleStatus string

const (
	VehicleStatusStandby  VehicleStatus = "standby"
	VehicleStatusDispatched VehicleStatus = "dispatched"
	VehicleStatusOnScene  VehicleStatus = "on_scene"
	VehicleStatusInTransit VehicleStatus = "in_transit"
	VehicleStatusReturning VehicleStatus = "returning"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
)

type MaintenanceType string

const (
	MaintenanceTypeRepair MaintenanceType = "repair"
	MaintenanceTypeDisinfect MaintenanceType = "disinfect"
)

type Vehicle struct {
	ID              string        `json:"id"`
	PlateNumber     string        `json:"plate_number"`
	Type            VehicleType   `json:"type"`
	Crew            []string      `json:"crew"`
	Longitude       float64       `json:"longitude"`
	Latitude        float64       `json:"latitude"`
	Status          VehicleStatus `json:"status"`
	MaintenanceType MaintenanceType `json:"maintenance_type,omitempty"`
	CurrentDispatchID string      `json:"current_dispatch_id,omitempty"`
	LastStatusUpdate time.Time    `json:"last_status_update"`
	OnlineSince     time.Time     `json:"online_since"`
	TotalDrivingSeconds int64     `json:"total_driving_seconds"`
}

type HospitalDepartment string

const (
	DeptInternalMedicine HospitalDepartment = "internal_medicine"
	DeptSurgery          HospitalDepartment = "surgery"
	DeptPediatrics       HospitalDepartment = "pediatrics"
	DeptObstetrics       HospitalDepartment = "obstetrics"
	DeptStroke           HospitalDepartment = "stroke"
	DeptChestPain        HospitalDepartment = "chest_pain"
)

type Hospital struct {
	ID        string                        `json:"id"`
	Name      string                        `json:"name"`
	Address   string                        `json:"address"`
	Longitude float64                       `json:"longitude"`
	Latitude  float64                       `json:"latitude"`
	Beds      map[HospitalDepartment]int    `json:"beds"`
}

type IllnessType string

const (
	IllnessNormal   IllnessType = "normal"
	IllnessSevere   IllnessType = "severe"
	IllnessStroke   IllnessType = "stroke"
	IllnessHeart    IllnessType = "heart_attack"
	IllnessNeonatal IllnessType = "neonatal"
	IllnessTrauma   IllnessType = "trauma"
)

type DispatchStatus string

const (
	DispatchStatusQueued    DispatchStatus = "queued"
	DispatchStatusDispatched DispatchStatus = "dispatched"
	DispatchStatusOnScene    DispatchStatus = "on_scene"
	DispatchStatusInTransit  DispatchStatus = "in_transit"
	DispatchStatusCompleted  DispatchStatus = "completed"
	DispatchStatusCancelled  DispatchStatus = "cancelled"
	DispatchStatusReassigned DispatchStatus = "reassigned"
)

type ChangeLogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Operator  string    `json:"operator"`
	Action    string    `json:"action"`
	Reason    string    `json:"reason,omitempty"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
}

type Dispatch struct {
	ID                 string         `json:"id"`
	CallerID           string         `json:"caller_id"`
	CallerName         string         `json:"caller_name"`
	CallerPhone        string         `json:"caller_phone"`
	IsVIP              bool           `json:"is_vip"`
	QueueJumpReason    string         `json:"queue_jump_reason,omitempty"`
	VehicleID          string         `json:"vehicle_id,omitempty"`
	HospitalID         string         `json:"hospital_id,omitempty"`
	IllnessType        IllnessType    `json:"illness_type"`
	IllnessDescription string         `json:"illness_description"`
	Address            string         `json:"address"`
	Longitude          float64        `json:"longitude"`
	Latitude           float64        `json:"latitude"`
	Status             DispatchStatus `json:"status"`
	IsEmptyRun         bool           `json:"is_empty_run"`
	EmptyRunReason     string         `json:"empty_run_reason,omitempty"`
	CallTime           time.Time      `json:"call_time"`
	DispatchTime       *time.Time     `json:"dispatch_time,omitempty"`
	ArrivalTime        *time.Time     `json:"arrival_time,omitempty"`
	DepartSceneTime    *time.Time     `json:"depart_scene_time,omitempty"`
	ArriveHospitalTime *time.Time     `json:"arrive_hospital_time,omitempty"`
	ReturnTime         *time.Time     `json:"return_time,omitempty"`
	CompleteTime       *time.Time     `json:"complete_time,omitempty"`
	ChangeLog          []ChangeLogEntry `json:"change_log"`
	OriginalDispatchID string         `json:"original_dispatch_id,omitempty"`
	ReassignedFromID   string         `json:"reassigned_from_id,omitempty"`
	Priority           int            `json:"priority"`
}

type POI struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Type      string  `json:"type"`
}

type Caller struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Phone       string   `json:"phone"`
	IsVIP       bool     `json:"is_vip"`
	Address     string   `json:"address,omitempty"`
	DispatchIDs []string `json:"dispatch_ids"`
	CreatedAt   time.Time `json:"created_at"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type MonthlyStats struct {
	Month               string             `json:"month"`
	AvgResponseTimeSec  float64            `json:"avg_response_time_sec"`
	TotalDispatches     int                `json:"total_dispatches"`
	EmptyRunCount       int                `json:"empty_run_count"`
	EmptyRunRate        float64            `json:"empty_run_rate"`
	HospitalReceiveStats map[string]int    `json:"hospital_receive_stats"`
	VehicleUtilization  map[string]float64 `json:"vehicle_utilization"`
}
