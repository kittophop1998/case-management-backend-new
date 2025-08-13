package model

type ConnectorGetCustomerInfoRequest struct {
	UserRef string `json:"UserRef" validate:"required"`
	Mode    string `json:"Mode" validate:"required"`
}

type ConnectorGetCustomerInfoResponse struct {
	CustID          string `json:"IDCardNo"`
	CustomerNameEng string `json:"CustomerNameENG" validate:"required"`
	CustomerNameTH  string `json:"CustomerNameTH"`
	MobileNO        string `json:"MobileNo"`
	MailToAddress   string `json:"mail_to_address"`
	MailTo          string `json:"MailTo"`
	HomeAddress     string `json:"HomeAddress"`
	HomeZip         int    `json:"HomeZip"`
	OfficeName      string `json:"OfficeName"`
	OfficeAddress   string `json:"OfficeAddress"`
	OfficeZip       int    `json:"OfficeZip"`
}

type ConnectorErrorResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}
