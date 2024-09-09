# Template

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**ProjectConfigs** | **[]string** |  | 

## Methods

### NewTemplate

`func NewTemplate(name string, projectConfigs []string, ) *Template`

NewTemplate instantiates a new Template object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTemplateWithDefaults

`func NewTemplateWithDefaults() *Template`

NewTemplateWithDefaults instantiates a new Template object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Template) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Template) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Template) SetName(v string)`

SetName sets Name field to given value.


### GetProjectConfigs

`func (o *Template) GetProjectConfigs() []string`

GetProjectConfigs returns the ProjectConfigs field if non-nil, zero value otherwise.

### GetProjectConfigsOk

`func (o *Template) GetProjectConfigsOk() (*[]string, bool)`

GetProjectConfigsOk returns a tuple with the ProjectConfigs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProjectConfigs

`func (o *Template) SetProjectConfigs(v []string)`

SetProjectConfigs sets ProjectConfigs field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


