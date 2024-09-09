# ProjectDataDTO

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BuildConfig** | Pointer to [**BuildConfig**](BuildConfig.md) |  | [optional] 
**EnvVars** | **map[string]string** |  | 
**Image** | Pointer to **string** |  | [optional] 
**Name** | **string** |  | 
**Source** | [**CreateProjectSourceDTO**](CreateProjectSourceDTO.md) |  | 
**User** | Pointer to **string** |  | [optional] 

## Methods

### NewProjectDataDTO

`func NewProjectDataDTO(envVars map[string]string, name string, source CreateProjectSourceDTO, ) *ProjectDataDTO`

NewProjectDataDTO instantiates a new ProjectDataDTO object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProjectDataDTOWithDefaults

`func NewProjectDataDTOWithDefaults() *ProjectDataDTO`

NewProjectDataDTOWithDefaults instantiates a new ProjectDataDTO object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBuildConfig

`func (o *ProjectDataDTO) GetBuildConfig() BuildConfig`

GetBuildConfig returns the BuildConfig field if non-nil, zero value otherwise.

### GetBuildConfigOk

`func (o *ProjectDataDTO) GetBuildConfigOk() (*BuildConfig, bool)`

GetBuildConfigOk returns a tuple with the BuildConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuildConfig

`func (o *ProjectDataDTO) SetBuildConfig(v BuildConfig)`

SetBuildConfig sets BuildConfig field to given value.

### HasBuildConfig

`func (o *ProjectDataDTO) HasBuildConfig() bool`

HasBuildConfig returns a boolean if a field has been set.

### GetEnvVars

`func (o *ProjectDataDTO) GetEnvVars() map[string]string`

GetEnvVars returns the EnvVars field if non-nil, zero value otherwise.

### GetEnvVarsOk

`func (o *ProjectDataDTO) GetEnvVarsOk() (*map[string]string, bool)`

GetEnvVarsOk returns a tuple with the EnvVars field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnvVars

`func (o *ProjectDataDTO) SetEnvVars(v map[string]string)`

SetEnvVars sets EnvVars field to given value.


### GetImage

`func (o *ProjectDataDTO) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ProjectDataDTO) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ProjectDataDTO) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ProjectDataDTO) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetName

`func (o *ProjectDataDTO) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ProjectDataDTO) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ProjectDataDTO) SetName(v string)`

SetName sets Name field to given value.


### GetSource

`func (o *ProjectDataDTO) GetSource() CreateProjectSourceDTO`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *ProjectDataDTO) GetSourceOk() (*CreateProjectSourceDTO, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *ProjectDataDTO) SetSource(v CreateProjectSourceDTO)`

SetSource sets Source field to given value.


### GetUser

`func (o *ProjectDataDTO) GetUser() string`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *ProjectDataDTO) GetUserOk() (*string, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *ProjectDataDTO) SetUser(v string)`

SetUser sets User field to given value.

### HasUser

`func (o *ProjectDataDTO) HasUser() bool`

HasUser returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


