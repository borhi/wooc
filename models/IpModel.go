package models

type IpModel struct {
	IpAddress string   `json:"ip_address" validate:"required,ip"`
	ASN       int32    `json:"ASN" validate:"required,numeric"`
	Domains   []string `json:"domains" validate:"required,max=100,dive,required,fqdn"`
}

type UnauthorizedIpModel struct {
	IpAddress string `json:"ip_address"`
}
