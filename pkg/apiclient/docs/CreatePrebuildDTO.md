# CreatePrebuildDTO

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Branch** | **string** |  | 
**CommitInterval** | **int32** |  | 
**Id** | Pointer to **string** |  | [optional] 
**ProjectConfigName** | **string** |  | 
**RunAtInit** | **bool** |  | 
**TriggerFiles** | **[]string** |  | 

## Methods

### NewCreatePrebuildDTO

`func NewCreatePrebuildDTO(branch string, commitInterval int32, projectConfigName string, runAtInit bool, triggerFiles []string, ) *CreatePrebuildDTO`

NewCreatePrebuildDTO instantiates a new CreatePrebuildDTO object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreatePrebuildDTOWithDefaults

`func NewCreatePrebuildDTOWithDefaults() *CreatePrebuildDTO`

NewCreatePrebuildDTOWithDefaults instantiates a new CreatePrebuildDTO object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBranch

`func (o *CreatePrebuildDTO) GetBranch() string`

GetBranch returns the Branch field if non-nil, zero value otherwise.

### GetBranchOk

`func (o *CreatePrebuildDTO) GetBranchOk() (*string, bool)`

GetBranchOk returns a tuple with the Branch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBranch

`func (o *CreatePrebuildDTO) SetBranch(v string)`

SetBranch sets Branch field to given value.


### GetCommitInterval

`func (o *CreatePrebuildDTO) GetCommitInterval() int32`

GetCommitInterval returns the CommitInterval field if non-nil, zero value otherwise.

### GetCommitIntervalOk

`func (o *CreatePrebuildDTO) GetCommitIntervalOk() (*int32, bool)`

GetCommitIntervalOk returns a tuple with the CommitInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommitInterval

`func (o *CreatePrebuildDTO) SetCommitInterval(v int32)`

SetCommitInterval sets CommitInterval field to given value.


### GetId

`func (o *CreatePrebuildDTO) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CreatePrebuildDTO) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CreatePrebuildDTO) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CreatePrebuildDTO) HasId() bool`

HasId returns a boolean if a field has been set.

### GetProjectConfigName

`func (o *CreatePrebuildDTO) GetProjectConfigName() string`

GetProjectConfigName returns the ProjectConfigName field if non-nil, zero value otherwise.

### GetProjectConfigNameOk

`func (o *CreatePrebuildDTO) GetProjectConfigNameOk() (*string, bool)`

GetProjectConfigNameOk returns a tuple with the ProjectConfigName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectConfigName

`func (o *CreatePrebuildDTO) SetProjectConfigName(v string)`

SetProjectConfigName sets ProjectConfigName field to given value.


### GetRunAtInit

`func (o *CreatePrebuildDTO) GetRunAtInit() bool`

GetRunAtInit returns the RunAtInit field if non-nil, zero value otherwise.

### GetRunAtInitOk

`func (o *CreatePrebuildDTO) GetRunAtInitOk() (*bool, bool)`

GetRunAtInitOk returns a tuple with the RunAtInit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunAtInit

`func (o *CreatePrebuildDTO) SetRunAtInit(v bool)`

SetRunAtInit sets RunAtInit field to given value.


### GetTriggerFiles

`func (o *CreatePrebuildDTO) GetTriggerFiles() []string`

GetTriggerFiles returns the TriggerFiles field if non-nil, zero value otherwise.

### GetTriggerFilesOk

`func (o *CreatePrebuildDTO) GetTriggerFilesOk() (*[]string, bool)`

GetTriggerFilesOk returns a tuple with the TriggerFiles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTriggerFiles

`func (o *CreatePrebuildDTO) SetTriggerFiles(v []string)`

SetTriggerFiles sets TriggerFiles field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


