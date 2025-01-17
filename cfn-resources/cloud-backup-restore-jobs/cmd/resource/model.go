// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	ProjectId                  *string                     `json:",omitempty"`
	InstanceType               *string                     `json:",omitempty"`
	InstanceName               *string                     `json:",omitempty"`
	Id                         *string                     `json:",omitempty"`
	DeliveryType               *string                     `json:",omitempty"`
	DeliveryUrl                []string                    `json:",omitempty"`
	Cancelled                  *bool                       `json:",omitempty"`
	CreatedAt                  *string                     `json:",omitempty"`
	Expired                    *bool                       `json:",omitempty"`
	ExpiresAt                  *string                     `json:",omitempty"`
	FinishedAt                 *string                     `json:",omitempty"`
	Timestamp                  *string                     `json:",omitempty"`
	SnapshotId                 *string                     `json:",omitempty"`
	Links                      []Links                     `json:",omitempty"`
	OpLogTs                    *string                     `json:",omitempty"`
	OpLogInc                   *string                     `json:",omitempty"`
	PointInTimeUtcSeconds      *int                        `json:",omitempty"`
	TargetProjectId            *string                     `json:",omitempty"`
	TargetClusterName          *string                     `json:",omitempty"`
	Profile                    *string                     `json:",omitempty"`
	EnableSynchronousCreation  *bool                       `json:",omitempty"`
	SynchronousCreationOptions *SynchronousCreationOptions `json:",omitempty"`
}

// Links is autogenerated from the json schema
type Links struct {
	Rel  *string `json:",omitempty"`
	Href *string `json:",omitempty"`
}

// SynchronousCreationOptions is autogenerated from the json schema
type SynchronousCreationOptions struct {
	TimeOutInSeconds       *int  `json:",omitempty"`
	CallbackDelaySeconds   *int  `json:",omitempty"`
	ReturnSuccessIfTimeOut *bool `json:",omitempty"`
}
