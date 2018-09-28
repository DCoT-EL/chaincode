
package main

type Event struct{
	Caller    string `json:"caller"`
	Role      string `json:"role"`
	Operation       string `json:"operation"`
	Moment string `json:"moment"`

}

type ChainOfCustody struct {
	Id                       string `json:"id"`
	TrackingId               string `json:"trackingId"`
	DocumentId               string `json:"documentId"`
	WeightOfParcel           float64    `json:"weightOfParcel"`
	SortingCenterDestination string `json:"sortingCenterDestination"`
	DistributionOfficeCode   string `json:"distributionOfficeCode"`
	DistributionZone         string `json:"distributionZone"`
	DeliveryMan              string `json:"deliveryMan"`
	CodeOwner                string `json:"codeOwner"`
	Text                     string `json:"text"`
	Status                   string `json:"status"`
	Event   `json:"event"`   
}

