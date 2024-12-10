package v2

import apiModel "ticket-tracker/internal/controller/dtos/tcdd"

type TccServiceInterfaceV2 interface {
	LoadAllStationV2() (*apiModel.StationInformation, error)
	QueryTrainV2(request *apiModel.QueryTrainRequest) (*apiModel.QueryTrainResponse, error)
	AddSearchRequest(request *apiModel.SearchTrainRequest) (*apiModel.SearchTrainResponse, error)
}
