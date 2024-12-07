package response

type StationLoadResponse struct {
	Id                int64    `json:"id"`
	Name              string   `json:"name"`
	UnitId            int64    `json:"unitId"`
	StationCode       string   `json:"stationCode"`
	StationTrainTypes []string `json:"stationTrainTypes"`
	Pairs             []int64  `json:"pairs"`
}
