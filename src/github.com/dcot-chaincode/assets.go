/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

