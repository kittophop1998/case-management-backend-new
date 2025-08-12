package model

type GetCustInfoResponse struct {
	NationalID      string `json:"national_id"`
	CustomerNameEng string `json:"customer_name_eng"`
	CustomerNameTH  string `json:"customer_name_th"`
	MobileNO        string `json:"mobile_no"`
	MailToAddress   string `json:"mail_to_address"`
	MailTo          string `json:"mail_to"`
}

type GetCustProfileResponse struct {
	ErrorSystem                   string `json:"error_system"`
	LastCardApplyDate             string `json:"last_card_apply_date"`
	CustomerSentiment             string `json:"customer_sentiment"`
	PhoneNoLastUpdateDate         string `json:"phone_no_last_update_date"`
	LastIncreaseCreditLimitUpdate string `json:"last_increase_credit_limit_update"`
	LastReduceCreditLimitUpdate   string `json:"last_reduce_credit_limit_update"`
	LastIncomeUpdate              string `json:"last_income_update"`
	SuggestedAction               string `json:"suggested_action"`
	TypeOfJob                     string `json:"type_of_job"`
	MaritalStatus                 string `json:"marital_status"`
	Gender                        string `json:"gender"`
	LastEStatementSentDate        string `json:"last_e_statement_sent_date"`
	EStatementSentStatus          string `json:"e_statement_sent_status"`
	StatementChannel              string `json:"statement_channel"`
	ConsentForDisclose            string `json:"consent_for_disclose"`
	BlockMedia                    string `json:"block_media"`
	ConsentForCollectUse          string `json:"consent_for_collect_use"`
	PaymentStatus                 string `json:"payment_status"`
	DayPastDue                    string `json:"day_past_due"`
	LastOverdueDate               string `json:"last_overdue_date"`
	// MailTo                        string `json:"mail_to"`
}

type GetCustSegmentResponse struct {
	Sweetheart      string `json:"sweetheart"`
	ComplaintLevel  string `json:"complaint_level"`
	CustomerGroup   string `json:"customer_group"`
	ComplaintGroup  string `json:"complaint_group"`
	CustomerType    string `json:"customer_type"`
	MemberStatus    string `json:"member_status"`
	CustomerSegment string `json:"customer_segment"`
	UpdateData      string `json:"update_data"`
}

type GetCustSuggestionResponse struct {
	SuggestCards      []string                             `json:"suggest_cards"`
	SuggestPromotions []GetCustSuggestionPromotionResponse `json:"suggest_promotions"`
}

type GetCustSuggestionPromotionResponse struct {
	PromotionCode            string   `json:"promotion_code"`
	PromotionName            string   `json:"promotion_name"`
	PromotionDetails         string   `json:"promotion_details"`
	Action                   string   `json:"action"`
	PromotionResultTimestamp string   `json:"promotion_result_timestamp"`
	Period                   string   `json:"period"`
	EligibleCard             []string `json:"eligible_card"`
}
