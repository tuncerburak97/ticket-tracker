package response

type SelectSeatResponse struct {
	AllocationID            string `json:"allocationId"`
	SelectedBookingClassID  int    `json:"selectedBookingClassId"`
	SelectedFareFamilyID    int    `json:"selectedFareFamilyId"`
	SelectedPassengerTypeID int    `json:"selectedPassengerTypeId"`
	LockFor                 int    `json:"lockFor"`
}
