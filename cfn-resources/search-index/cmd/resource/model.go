// Code generated by 'cfn generate', changes will be undone by the next invocation. DO NOT EDIT.
// Updates to this type are made my editing the schema file and executing the 'generate' command.
package resource

// Model is autogenerated from the json schema
type Model struct {
	Analyzer       *string                                   `json:",omitempty"`
	Analyzers      []ApiAtlasFTSAnalyzersViewManual          `json:",omitempty"`
	ApiKeys        *ApiKeyDefinition                         `json:",omitempty"`
	ClusterName    *string                                   `json:",omitempty"`
	CollectionName *string                                   `json:",omitempty"`
	Database       *string                                   `json:",omitempty"`
	GroupId        *string                                   `json:",omitempty"`
	IndexId        *string                                   `json:",omitempty"`
	Mappings       *ApiAtlasFTSMappingsViewManual            `json:",omitempty"`
	Name           *string                                   `json:",omitempty"`
	SearchAnalyzer *string                                   `json:",omitempty"`
	Status         *string                                   `json:",omitempty"`
	Synonyms       []ApiAtlasFTSSynonymMappingDefinitionView `json:",omitempty"`
}

// ApiAtlasFTSAnalyzersViewManual is autogenerated from the json schema
type ApiAtlasFTSAnalyzersViewManual struct {
	CharFilters  []map[string]interface{} `json:",omitempty"`
	Name         *string                  `json:",omitempty"`
	TokenFilters []map[string]interface{} `json:",omitempty"`
	Tokenizer    map[string]interface{}   `json:",omitempty"`
}

// ApiKeyDefinition is autogenerated from the json schema
type ApiKeyDefinition struct {
	PrivateKey *string `json:",omitempty"`
	PublicKey  *string `json:",omitempty"`
}

// ApiAtlasFTSMappingsViewManual is autogenerated from the json schema
type ApiAtlasFTSMappingsViewManual struct {
	Dynamic *bool   `json:",omitempty"`
	Fields  *string `json:",omitempty"`
}

// ApiAtlasFTSSynonymMappingDefinitionView is autogenerated from the json schema
type ApiAtlasFTSSynonymMappingDefinitionView struct {
	Analyzer *string        `json:",omitempty"`
	Name     *string        `json:",omitempty"`
	Source   *SynonymSource `json:",omitempty"`
}

// SynonymSource is autogenerated from the json schema
type SynonymSource struct {
	Collection *string `json:",omitempty"`
}