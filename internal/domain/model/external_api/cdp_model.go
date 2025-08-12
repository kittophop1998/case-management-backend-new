package model

import "encoding/json"

type Profile interface{}

type Profile1Response struct {
	LastCardApplyDate     string `json:"last_card_apply_date"`
	PhoneNoLastUpdateDate string `json:"phone_no_last_update_date" validate:"required"`
}

type Profile2Response struct {
	Gender        string `json:"gender" validate:"required"`
	MaritalStatus string `json:"marital_status" validate:"required"`
	TypeOfJob     string `json:"type_of_job" validate:"required"`
}

type Profile3Response struct {
	LastIncreaseCreditLimitUpdate string `json:"last_increase_credit_limit_update" validate:"required"`
	LastReduceCreditLimitUpdate   string `json:"last_reduce_credit_limit_update" validate:"required"`
	LastIncomeUpdate              string `json:"last_income_update" validate:"required"`
	SuggestedAction               string `json:"suggested_action" validate:"required"`
}

type Profile4Response struct {
	LastEStatementSentDate string `json:"last_e_statement_sent_date" validate:"required"`
	EStatementSentStatus   string `json:"e_statement_sent_status" validate:"required"`
	StatementChannel       string `json:"statement_channel" validate:"required"`
}

type Profile5Response struct {
	ConsentForCollectUse string `json:"consent_for_collect_use" validate:"required"`
	ConsentForDisclose   string `json:"consent_for_disclose" validate:"required"`
	BlockMedia           string `json:"block_media" validate:"required"`
}

type Item struct {
	Values     []string        `json:"values"`
	Attributes json.RawMessage `json:"attributes"`
	Key        Key             `json:"key"`
	AudienceId string          `json:"audienceId"`
}

type Key struct {
	AeonID string `json:"aeon_id"`
}

type Attributes struct {
	VvipCustomerGroupFlag       string     `json:"vvip_customer_group_flag" validate:"required"`
	VvipCustomerPosition        string     `json:"vvip_customer_position"`
	SweetheartCustomerGroupFlag string     `json:"sweetheart_customer_group_flag"`
	CustomerLevel               string     `json:"customer_level" validate:"required"`
	ComplaintGroup              string     `json:"complaint_group" validate:"required"`
	ComplaintTopic              string     `json:"complaint_topic" validate:"required"`
	CustomerType                string     `json:"customer_type"`
	MemberStatus                string     `json:"member_status" validate:"required"`
	CbaSegment                  string     `json:"cba_segment" validate:"required"`
	UpdateData                  string     `json:"update_data" validate:"required"`
	NameOfCards                 string     `json:"name_of_cards" validate:"required"`
	PromotionArray              [][]string `json:"promotion_array"`
}
