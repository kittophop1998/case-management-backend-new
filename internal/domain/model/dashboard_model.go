package model

type GetCustInfoResponse struct {
	NationalID      string `json:"nationalId"`
	CustomerNameEng string `json:"customerNameEng"`
	CustomerNameTH  string `json:"customerNameTH"`
	MobileNO        string `json:"mobileNO"`
	MailToAddress   string `json:"mailToAddress"`
	MailTo          string `json:"mailTo"`
}

type GetCustProfileResponse struct {
	ErrorSystem                   string `json:"errorSystem"`
	LastCardApplyDate             string `json:"lastCardApplyDate"`
	CustomerSentiment             string `json:"customerSentiment"`
	PhoneNoLastUpdateDate         string `json:"phoneNoLastUpdateDate"`
	LastIncreaseCreditLimitUpdate string `json:"lastIncreaseCreditLimitUpdate"`
	LastReduceCreditLimitUpdate   string `json:"lastReduceCreditLimitUpdate"`
	LastIncomeUpdate              string `json:"lastIncomeUpdate"`
	SuggestedAction               string `json:"suggestedAction"`
	TypeOfJob                     string `json:"typeOfJob"`
	MaritalStatus                 string `json:"maritalStatus"`
	Gender                        string `json:"gender"`
	LastEStatementSentDate        string `json:"lastEStatementSentDate"`
	EStatementSentStatus          string `json:"eStatementSentStatus"`
	StatementChannel              string `json:"statementChannel"`
	ConsentForDisclose            string `json:"consentForDisclose"`
	BlockMedia                    string `json:"blockMedia"`
	ConsentForCollectUse          string `json:"consentForCollectUse"`
	PaymentStatus                 string `json:"paymentStatus"`
	DayPastDue                    string `json:"dayPastDue"`
	LastOverdueDate               string `json:"lastOverdueDate"`
	// MailTo                     string `json:"mailTo"`
}

type GetCustSegmentResponse struct {
	Sweetheart      string `json:"sweetheart"`
	ComplaintLevel  string `json:"complaintLevel"`
	CustomerGroup   string `json:"customerGroup"`
	ComplaintGroup  string `json:"complaintGroup"`
	CustomerType    string `json:"customerType"`
	MemberStatus    string `json:"memberStatus"`
	CustomerSegment string `json:"customerSegment"`
	UpdateData      string `json:"updateData"`
}

type GetCustSuggestionResponse struct {
	SuggestCards      []string                             `json:"suggest_cards"`
	SuggestPromotions []GetCustSuggestionPromotionResponse `json:"suggest_promotions"`
}

type GetCustSuggestionPromotionResponse struct {
	PromotionCode            string   `json:"promotionCode"`
	PromotionName            string   `json:"promotionName"`
	PromotionDetails         string   `json:"promotionDetails"`
	Action                   string   `json:"action"`
	PromotionResultTimestamp string   `json:"promotionResultTimestamp"`
	Period                   string   `json:"period"`
	EligibleCard             []string `json:"eligibleCard"`
}
