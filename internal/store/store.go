package store

import (
	"sync"
	"time"

	"emergency-dispatch/internal/models"
)

type Store struct {
	Vehicles      sync.Map
	Hospitals       sync.Map
	Dispatches      sync.Map
	POIs            sync.Map
	Callers         sync.Map
	CallerPhoneIdx  sync.Map
	DispatchQueue   []string
	QueueMu           sync.Mutex
}

var GlobalStore = &Store{}

func (s *Store) AddVehicle(v *models.Vehicle) {
	s.Vehicles.Store(v.ID, v)
}

func (s *Store) GetVehicle(id string) (*models.Vehicle, bool) {
	val, ok := s.Vehicles.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*models.Vehicle), true
}

func (s *Store) UpdateVehicle(v *models.Vehicle) {
	v.LastStatusUpdate = time.Now()
	s.Vehicles.Store(v.ID, v)
}

func (s *Store) DeleteVehicle(id string) {
	s.Vehicles.Delete(id)
}

func (s *Store) ListVehicles() []*models.Vehicle {
	var result []*models.Vehicle
	s.Vehicles.Range(func(key, value interface{}) bool {
		result = append(result, value.(*models.Vehicle))
		return true
	})
	return result
}

func (s *Store) ListVehiclesByStatus(status models.VehicleStatus) []*models.Vehicle {
	var result []*models.Vehicle
	s.Vehicles.Range(func(key, value interface{}) bool {
		v := value.(*models.Vehicle)
		if v.Status == status {
			result = append(result, v)
		}
		return true
	})
	return result
}

func (s *Store) ListVehiclesByTypeAndStatus(vType models.VehicleType, status models.VehicleStatus) []*models.Vehicle {
	var result []*models.Vehicle
	s.Vehicles.Range(func(key, value interface{}) bool {
		v := value.(*models.Vehicle)
		if v.Type == vType && v.Status == status {
			result = append(result, v)
		}
		return true
	})
	return result
}

func (s *Store) CompareAndSwapVehicleStatus(id string, oldStatus, newStatus models.VehicleStatus) (*models.Vehicle, bool) {
	val, ok := s.Vehicles.Load(id)
	if !ok {
		return nil, false
	}
	v := val.(*models.Vehicle)
	if v.Status != oldStatus {
		return v, false
	}
	newV := *v
	newV.Status = newStatus
	newV.LastStatusUpdate = time.Now()
	swapped := s.Vehicles.CompareAndSwap(id, val, &newV)
	if swapped {
		return &newV, true
	}
	return nil, false
}

func (s *Store) AssignVehicleToDispatch(vehicleID, dispatchID string) (*models.Vehicle, bool) {
	val, ok := s.Vehicles.Load(vehicleID)
	if !ok {
		return nil, false
	}
	v := val.(*models.Vehicle)
	if v.Status != models.VehicleStatusStandby {
		return v, false
	}
	newV := *v
	newV.Status = models.VehicleStatusDispatched
	newV.CurrentDispatchID = dispatchID
	newV.LastStatusUpdate = time.Now()
	swapped := s.Vehicles.CompareAndSwap(vehicleID, val, &newV)
	if swapped {
		return &newV, true
	}
	return nil, false
}

func (s *Store) ReleaseVehicleFromDispatch(vehicleID string) (*models.Vehicle, bool) {
	val, ok := s.Vehicles.Load(vehicleID)
	if !ok {
		return nil, false
	}
	v := val.(*models.Vehicle)
	if v.Status == models.VehicleStatusStandby {
		return v, true
	}
	newV := *v
	newV.Status = models.VehicleStatusStandby
	newV.CurrentDispatchID = ""
	newV.LastStatusUpdate = time.Now()
	swapped := s.Vehicles.CompareAndSwap(vehicleID, val, &newV)
	if swapped {
		return &newV, true
	}
	return nil, false
}

func (s *Store) AddHospital(h *models.Hospital) {
	s.Hospitals.Store(h.ID, h)
}

func (s *Store) GetHospital(id string) (*models.Hospital, bool) {
	val, ok := s.Hospitals.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*models.Hospital), true
}

func (s *Store) UpdateHospital(h *models.Hospital) {
	s.Hospitals.Store(h.ID, h)
}

func (s *Store) DeleteHospital(id string) {
	s.Hospitals.Delete(id)
}

func (s *Store) ListHospitals() []*models.Hospital {
	var result []*models.Hospital
	s.Hospitals.Range(func(key, value interface{}) bool {
		result = append(result, value.(*models.Hospital))
		return true
	})
	return result
}

func (s *Store) AddDispatch(d *models.Dispatch) {
	s.Dispatches.Store(d.ID, d)
}

func (s *Store) GetDispatch(id string) (*models.Dispatch, bool) {
	val, ok := s.Dispatches.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*models.Dispatch), true
}

func (s *Store) UpdateDispatch(d *models.Dispatch) {
	s.Dispatches.Store(d.ID, d)
}

func (s *Store) DeleteDispatch(id string) {
	s.Dispatches.Delete(id)
}

func (s *Store) ListDispatches() []*models.Dispatch {
	var result []*models.Dispatch
	s.Dispatches.Range(func(key, value interface{}) bool {
		result = append(result, value.(*models.Dispatch))
		return true
	})
	return result
}

func (s *Store) ListDispatchesByStatus(status models.DispatchStatus) []*models.Dispatch {
	var result []*models.Dispatch
	s.Dispatches.Range(func(key, value interface{}) bool {
		d := value.(*models.Dispatch)
		if d.Status == status {
			result = append(result, d)
		}
		return true
	})
	return result
}

func (s *Store) AddPOI(p *models.POI) {
	s.POIs.Store(p.ID, p)
}

func (s *Store) GetPOI(id string) (*models.POI, bool) {
	val, ok := s.POIs.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*models.POI), true
}

func (s *Store) UpdatePOI(p *models.POI) {
	s.POIs.Store(p.ID, p)
}

func (s *Store) DeletePOI(id string) {
	s.POIs.Delete(id)
}

func (s *Store) ListPOIs() []*models.POI {
	var result []*models.POI
	s.POIs.Range(func(key, value interface{}) bool {
		result = append(result, value.(*models.POI))
		return true
	})
	return result
}

func (s *Store) AddCaller(c *models.Caller) {
	s.Callers.Store(c.ID, c)
	s.CallerPhoneIdx.Store(c.Phone, c.ID)
}

func (s *Store) GetCaller(id string) (*models.Caller, bool) {
	val, ok := s.Callers.Load(id)
	if !ok {
		return nil, false
	}
	return val.(*models.Caller), true
}

func (s *Store) GetCallerByPhone(phone string) (*models.Caller, bool) {
	idVal, ok := s.CallerPhoneIdx.Load(phone)
	if !ok {
		return nil, false
	}
	return s.GetCaller(idVal.(string))
}

func (s *Store) UpdateCaller(c *models.Caller) {
	s.Callers.Store(c.ID, c)
}

func (s *Store) DeleteCaller(id string) {
	val, ok := s.Callers.Load(id)
	if ok {
		c := val.(*models.Caller)
		s.CallerPhoneIdx.Delete(c.Phone)
	}
	s.Callers.Delete(id)
}

func (s *Store) ListCallers() []*models.Caller {
	var result []*models.Caller
	s.Callers.Range(func(key, value interface{}) bool {
		result = append(result, value.(*models.Caller))
		return true
	})
	return result
}

func (s *Store) AddToQueue(dispatchID string, isVIP bool) {
	s.QueueMu.Lock()
	defer s.QueueMu.Unlock()
	if isVIP {
		idx := 0
		for i, id := range s.DispatchQueue {
			d, ok := s.GetDispatch(id)
			if ok && !d.IsVIP {
				idx = i
				break
			}
			idx = i + 1
		}
		s.DispatchQueue = append(s.DispatchQueue[:idx], append([]string{dispatchID}, s.DispatchQueue[idx:]...)
	} else {
		s.DispatchQueue = append(s.DispatchQueue, dispatchID)
	}
}

func (s *Store) RemoveFromQueue(dispatchID string) {
	s.QueueMu.Lock()
	defer s.QueueMu.Unlock()
	for i, id := range s.DispatchQueue {
		if id == dispatchID {
			s.DispatchQueue = append(s.DispatchQueue[:i], s.DispatchQueue[i+1:]...)
			break
		}
	}
}

func (s *Store) GetQueue() []string {
	s.QueueMu.Lock()
	defer s.QueueMu.Unlock()
	result := make([]string, len(s.DispatchQueue))
	copy(result, s.DispatchQueue)
	return result
}
