package models

type TokenEndpointRequest struct {
	Permission string `json:"permission"`
	AppID      uint   `json:"appId"`
	DomainID   uint   `json:"domainId"`
}
