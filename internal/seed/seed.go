package seed

import (
	"time"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/store"
)

func SeedData() {
	seedVehicles()
	seedHospitals()
	seedPOIs()
	seedCallers()
}

func seedVehicles() {
	now := time.Now()

	vehicles := []*models.Vehicle{
		{
			ID:               "V001",
			PlateNumber:      "京A·12001",
			Type:             models.VehicleTypeNormal,
			Crew:             []string{"张医生", "李护士"},
			Longitude:        116.4074,
			Latitude:         39.9042,
			Status:           models.VehicleStatusStandby,
			LastStatusUpdate: now,
			OnlineSince:      now.Add(-2 * time.Hour),
		},
		{
			ID:               "V002",
			PlateNumber:      "京A·12002",
			Type:             models.VehicleTypeNormal,
			Crew:             []string{"王医生", "赵护士"},
			Longitude:        116.4174,
			Latitude:         39.9142,
			Status:           models.VehicleStatusStandby,
			LastStatusUpdate: now,
			OnlineSince:      now.Add(-3 * time.Hour),
		},
		{
			ID:               "V003",
			PlateNumber:      "京A·12003",
			Type:             models.VehicleTypeSevere,
			Crew:             []string{"刘主任", "陈护士", "孙急救员"},
			Longitude:        116.3974,
			Latitude:         39.9092,
			Status:           models.VehicleStatusStandby,
			LastStatusUpdate: now,
			OnlineSince:      now.Add(-1 * time.Hour),
		},
		{
			ID:               "V004",
			PlateNumber:      "京A·12004",
			Type:             models.VehicleTypeSevere,
			Crew:             []string{"周医生", "吴护士", "郑急救员"},
			Longitude:        116.4274,
			Latitude:         39.8942,
			Status:           models.VehicleStatusStandby,
			LastStatusUpdate: now,
			OnlineSince:      now.Add(-4 * time.Hour),
		},
		{
			ID:               "V005",
			PlateNumber:      "京A·12005",
			Type:             models.VehicleTypeNeonatal,
			Crew:             []string{"钱医生", "孙护士"},
			Longitude:        116.4124,
			Latitude:         39.9242,
			Status:           models.VehicleStatusStandby,
			LastStatusUpdate: now,
			OnlineSince:      now.Add(-2 * time.Hour),
		},
		{
			ID:               "V006",
			PlateNumber:      "京A·12006",
			Type:             models.VehicleTypeNormal,
			Crew:             []string{"冯医生", "褚护士"},
			Longitude:        116.3874,
			Latitude:         39.9242,
			Status:           models.VehicleStatusMaintenance,
			MaintenanceType:  models.MaintenanceTypeDisinfect,
			LastStatusUpdate: now,
			OnlineSince:      now.Add(-5 * time.Hour),
		},
	}

	for _, v := range vehicles {
		store.GlobalStore.AddVehicle(v)
	}
}

func seedHospitals() {
	hospitals := []*models.Hospital{
		{
			ID:        "H001",
			Name:      "北京协和医院",
			Address:   "东城区帅府园1号",
			Longitude: 116.4188,
			Latitude:  39.9199,
			Beds: map[models.HospitalDepartment]int{
				models.DeptInternalMedicine: 15,
				models.DeptSurgery:          12,
				models.DeptPediatrics:       8,
				models.DeptObstetrics:       6,
				models.DeptStroke:           10,
				models.DeptChestPain:        8,
			},
		},
		{
			ID:        "H002",
			Name:      "北京朝阳医院",
			Address:   "朝阳区工人体育场南路8号",
			Longitude: 116.4476,
			Latitude:  39.9326,
			Beds: map[models.HospitalDepartment]int{
				models.DeptInternalMedicine: 20,
				models.DeptSurgery:          18,
				models.DeptPediatrics:       10,
				models.DeptObstetrics:       5,
				models.DeptStroke:           15,
				models.DeptChestPain:        12,
			},
		},
		{
			ID:        "H003",
			Name:      "北京儿童医院",
			Address:   "西城区南礼士路56号",
			Longitude: 116.3622,
			Latitude:  39.9235,
			Beds: map[models.HospitalDepartment]int{
				models.DeptInternalMedicine: 0,
				models.DeptSurgery:          0,
				models.DeptPediatrics:       30,
				models.DeptObstetrics:       0,
				models.DeptStroke:           0,
				models.DeptChestPain:        0,
			},
		},
		{
			ID:        "H004",
			Name:      "天坛医院",
			Address:   "丰台区南四环西路119号",
			Longitude: 116.3356,
			Latitude:  39.8412,
			Beds: map[models.HospitalDepartment]int{
				models.DeptInternalMedicine: 10,
				models.DeptSurgery:          8,
				models.DeptPediatrics:       5,
				models.DeptObstetrics:       0,
				models.DeptStroke:           25,
				models.DeptChestPain:        5,
			},
		},
	}

	for _, h := range hospitals {
		store.GlobalStore.AddHospital(h)
	}
}

func seedPOIs() {
	pois := []*models.POI{
		{ID: "P001", Name: "天安门广场", Address: "东城区天安门广场", Longitude: 116.4039, Latitude: 39.9048, Type: "景点"},
		{ID: "P002", Name: "故宫博物院", Address: "东城区景山前街4号", Longitude: 116.3970, Latitude: 39.9163, Type: "景点"},
		{ID: "P003", Name: "王府井步行街", Address: "东城区王府井大街", Longitude: 116.4107, Latitude: 39.9149, Type: "商业区"},
		{ID: "P004", Name: "北京火车站", Address: "东城区毛家湾胡同甲13号", Longitude: 116.4273, Latitude: 39.9028, Type: "交通枢纽"},
		{ID: "P005", Name: "北京西站", Address: "丰台区莲花池东路118号", Longitude: 116.3218, Latitude: 39.8949, Type: "交通枢纽"},
		{ID: "P006", Name: "西单商场", Address: "西城区西单北大街120号", Longitude: 116.3745, Latitude: 39.9133, Type: "商业区"},
		{ID: "P007", Name: "三里屯太古里", Address: "朝阳区三里屯路19号", Longitude: 116.4556, Latitude: 39.9371, Type: "商业区"},
		{ID: "P008", Name: "国贸大厦", Address: "朝阳区建国门外大街1号", Longitude: 116.4604, Latitude: 39.9088, Type: "写字楼"},
		{ID: "P009", Name: "望京SOHO", Address: "朝阳区望京街10号", Longitude: 116.4782, Latitude: 39.9947, Type: "写字楼"},
		{ID: "P010", Name: "奥林匹克公园", Address: "朝阳区北辰东路15号", Longitude: 116.3972, Latitude: 40.0012, Type: "景点"},
		{ID: "P011", Name: "中关村", Address: "海淀区中关村大街", Longitude: 116.3176, Latitude: 39.9796, Type: "商业区"},
		{ID: "P012", Name: "北京大学", Address: "海淀区颐和园路5号", Longitude: 116.3055, Latitude: 39.9925, Type: "学校"},
		{ID: "P013", Name: "清华大学", Address: "海淀区清华园1号", Longitude: 116.3264, Latitude: 40.0030, Type: "学校"},
		{ID: "P014", Name: "北京南站", Address: "丰台区永外大街车站路12号", Longitude: 116.3776, Latitude: 39.8656, Type: "交通枢纽"},
		{ID: "P015", Name: "首都机场T3", Address: "顺义区机场西路", Longitude: 116.6056, Latitude: 40.0799, Type: "交通枢纽"},
		{ID: "P016", Name: "颐和园", Address: "海淀区新建宫门路19号", Longitude: 116.2721, Latitude: 39.9992, Type: "景点"},
		{ID: "P017", Name: "圆明园", Address: "海淀区清华西路28号", Longitude: 116.2987, Latitude: 40.0080, Type: "景点"},
		{ID: "P018", Name: "南锣鼓巷", Address: "东城区南锣鼓巷胡同", Longitude: 116.4038, Latitude: 39.9373, Type: "景点"},
		{ID: "P019", Name: "后海", Address: "西城区后海北沿", Longitude: 116.3898, Latitude: 39.9421, Type: "景点"},
		{ID: "P020", Name: "798艺术区", Address: "朝阳区酒仙桥路4号", Longitude: 116.4951, Latitude: 39.9848, Type: "景点"},
	}

	for _, p := range pois {
		store.GlobalStore.AddPOI(p)
	}
}

func seedCallers() {
	now := time.Now()
	callers := []*models.Caller{
		{
			ID:        "C001",
			Name:      "张三",
			Phone:     "13800138001",
			IsVIP:     false,
			Address:   "东城区王府井大街138号",
			CreatedAt: now.Add(-30 * 24 * time.Hour),
		},
		{
			ID:        "C002",
			Name:      "李四",
			Phone:     "13900139002",
			IsVIP:     true,
			Address:   "朝阳区建国门外大街1号",
			CreatedAt: now.Add(-60 * 24 * time.Hour),
		},
	}

	for _, c := range callers {
		store.GlobalStore.AddCaller(c)
	}
}
