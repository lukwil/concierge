package hasura

import "time"

// SingleDeploymentPayload is the payload refering to the single_deployment table.
// Hasura/Postgres sends it if an event trigger is fired.
type SingleDeploymentPayload struct {
	Event struct {
		SessionVariables struct {
			XHasuraRole string `json:"x-hasura-role"`
		} `json:"session_variables"`
		Op   string `json:"op"`
		Data struct {
			Old struct {
				ID             int    `json:"id"`
				Name           string `json:"name"`
				NameK8S        string `json:"name_k8s"`
				ContainerImage string `json:"container_image"`
				CPU            int    `json:"cpu"`
				RAM            int    `json:"ram"`
				GPU            int    `json:"gpu"`
				URLPrefix      string `json:"url_prefix"`
				StatusID       int    `json:"status_id"`
				VolumeID       int    `json:"volume_id"`
			} `json:"old"`
			New struct {
				ID             int    `json:"id"`
				Name           string `json:"name"`
				NameK8S        string `json:"name_k8s"`
				ContainerImage string `json:"container_image"`
				CPU            int    `json:"cpu"`
				RAM            int    `json:"ram"`
				GPU            int    `json:"gpu"`
				URLPrefix      string `json:"url_prefix"`
				StatusID       int    `json:"status_id"`
				VolumeID       int    `json:"volume_id"`
			} `json:"new"`
		} `json:"data"`
	} `json:"event"`
	CreatedAt    time.Time `json:"created_at"`
	ID           string    `json:"id"`
	DeliveryInfo struct {
		MaxRetries   int `json:"max_retries"`
		CurrentRetry int `json:"current_retry"`
	} `json:"delivery_info"`
	Trigger struct {
		Name string `json:"name"`
	} `json:"trigger"`
	Table struct {
		Schema string `json:"schema"`
		Name   string `json:"name"`
	} `json:"table"`
}
