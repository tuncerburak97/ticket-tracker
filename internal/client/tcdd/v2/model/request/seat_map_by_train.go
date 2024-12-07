package request

type SeatMapByTrainRequest struct {
	ToStationId   int64 `json:"toStationId"`
	FromStationID int64 `json:"fromStationId"`
	TrainId       int64 `json:"trainId"`
	LegIndex      int64 `json:"legIndex"`
}
