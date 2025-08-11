package model

type GetCustomerInfoResponse struct {
	NationalID      string `json:"national_id"`
	CustomerNameEng string `json:"customer_name_eng"`
	CustomerNameTH  string `json:"customer_name_th"`
	MobileNO        string `json:"mobile_no"`
	MailToAddress   string `json:"mail_to_address"`
	MailTo          string `json:"mail_to"`
}
