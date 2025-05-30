// this file was automatically generated using xgen from https://github.com/xuri/xgen
// it has been manually modified to fix some issues with the generated code
package l5x

import (
	"encoding/xml"
	"strings"
)

// RSLogix5000Content ...
type RSLogix5000Content *ProjectType

// ProjectType ...
type ProjectType struct {
	SchemaRevisionAttr    string                          `xml:"SchemaRevision,attr"`
	SoftwareRevisionAttr  string                          `xml:"SoftwareRevision,attr"`
	ArchitectureIDAttr    string                          `xml:"ArchitectureID,attr,omitempty"`
	TargetNameAttr        string                          `xml:"TargetName,attr,omitempty"`
	TargetTypeAttr        string                          `xml:"TargetType,attr,omitempty"`
	TargetSubTypeAttr     string                          `xml:"TargetSubType,attr,omitempty"`
	TargetClassAttr       string                          `xml:"TargetClass,attr,omitempty"`
	TargetCountAttr       string                          `xml:"TargetCount,attr,omitempty"`
	TargetRevisionAttr    string                          `xml:"TargetRevision,attr,omitempty"`
	TargetLastEditedAttr  string                          `xml:"TargetLastEdited,attr,omitempty"`
	TargetSignatureAttr   string                          `xml:"TargetSignature,attr,omitempty"`
	TargetIsEncryptedAttr string                          `xml:"TargetIsEncrypted,attr,omitempty"`
	CurrentLanguageAttr   string                          `xml:"CurrentLanguage,attr,omitempty"`
	ContainsContextAttr   string                          `xml:"ContainsContext,attr,omitempty"`
	OwnerAttr             string                          `xml:"Owner,attr,omitempty"`
	ExportDateAttr        string                          `xml:"ExportDate,attr,omitempty"`
	ExportOptionsAttr     string                          `xml:"ExportOptions,attr,omitempty"`
	ForceMemorySavingAttr string                          `xml:"ForceMemorySaving,attr,omitempty"`
	UseAttr               string                          `xml:"Use,attr,omitempty"`
	CustomProperties      *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Controller            *ControllerType                 `xml:"Controller"`
}

// DeviceObjType ...
type DeviceObjType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ExtDeviceObjType ...
type ExtDeviceObjType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ControllerType ...
type ControllerType struct {
	NameAttr                                              string                                 `xml:"Name,attr"`
	UIdAttr                                               string                                 `xml:"UId,attr,omitempty"`
	ParentUIdAttr                                         string                                 `xml:"ParentUId,attr,omitempty"`
	ProcessorTypeAttr                                     string                                 `xml:"ProcessorType,attr,omitempty"`
	MajorRevAttr                                          string                                 `xml:"MajorRev,attr,omitempty"`
	MinorRevAttr                                          string                                 `xml:"MinorRev,attr,omitempty"`
	TimeSliceAttr                                         string                                 `xml:"TimeSlice,attr,omitempty"`
	ShareUnusedTimeSliceAttr                              string                                 `xml:"ShareUnusedTimeSlice,attr,omitempty"`
	PowerLossProgramAttr                                  string                                 `xml:"PowerLossProgram,attr,omitempty"`
	PowerLossProgramUIdAttr                               string                                 `xml:"PowerLossProgramUId,attr,omitempty"`
	MajorFaultProgramAttr                                 string                                 `xml:"MajorFaultProgram,attr,omitempty"`
	MajorFaultProgramUIdAttr                              string                                 `xml:"MajorFaultProgramUId,attr,omitempty"`
	ProjectCreationDateAttr                               string                                 `xml:"ProjectCreationDate,attr,omitempty"`
	LastModifiedDateAttr                                  string                                 `xml:"LastModifiedDate,attr,omitempty"`
	SFCExecutionControlAttr                               string                                 `xml:"SFCExecutionControl,attr,omitempty"`
	SFCRestartPositionAttr                                string                                 `xml:"SFCRestartPosition,attr,omitempty"`
	SFCLastScanAttr                                       string                                 `xml:"SFCLastScan,attr,omitempty"`
	CommDriverAttr                                        string                                 `xml:"CommDriver,attr,omitempty"`
	CommPathAttr                                          string                                 `xml:"CommPath,attr,omitempty"`
	ProjectSNAttr                                         string                                 `xml:"ProjectSN,attr,omitempty"`
	OnlineSNAttr                                          string                                 `xml:"OnlineSN,attr,omitempty"`
	MatchProjectToControllerAttr                          string                                 `xml:"MatchProjectToController,attr,omitempty"`
	CanUseRPIFromProducerAttr                             string                                 `xml:"CanUseRPIFromProducer,attr,omitempty"`
	InhibitAutomaticFirmwareUpdateAttr                    string                                 `xml:"InhibitAutomaticFirmwareUpdate,attr,omitempty"`
	CurrentProjectLanguageAttr                            string                                 `xml:"CurrentProjectLanguage,attr,omitempty"`
	DefaultProjectLanguageAttr                            string                                 `xml:"DefaultProjectLanguage,attr,omitempty"`
	ControllerLanguageAttr                                string                                 `xml:"ControllerLanguage,attr,omitempty"`
	PassThroughConfigurationAttr                          string                                 `xml:"PassThroughConfiguration,attr,omitempty"`
	DownloadProjectDocumentationAndExtendedPropertiesAttr string                                 `xml:"DownloadProjectDocumentationAndExtendedProperties,attr,omitempty"`
	DownloadProjectCustomPropertiesAttr                   string                                 `xml:"DownloadProjectCustomProperties,attr,omitempty"`
	ReportMinorOverflowAttr                               string                                 `xml:"ReportMinorOverflow,attr,omitempty"`
	IOForcesEnabledAttr                                   string                                 `xml:"IOForcesEnabled,attr,omitempty"`
	SFCForcesEnabledAttr                                  string                                 `xml:"SFCForcesEnabled,attr,omitempty"`
	RedundancyUIdAttr                                     string                                 `xml:"RedundancyUId,attr,omitempty"`
	SafetyUIdAttr                                         string                                 `xml:"SafetyUId,attr,omitempty"`
	DataTypesUIdAttr                                      string                                 `xml:"DataTypesUId,attr,omitempty"`
	ModulesUIdAttr                                        string                                 `xml:"ModulesUId,attr,omitempty"`
	AddOnInstructionDefinitionsUIdAttr                    string                                 `xml:"AddOnInstructionDefinitionsUId,attr,omitempty"`
	TagsUIdAttr                                           string                                 `xml:"TagsUId,attr,omitempty"`
	ProgramsUIdAttr                                       string                                 `xml:"ProgramsUId,attr,omitempty"`
	TasksUIdAttr                                          string                                 `xml:"TasksUId,attr,omitempty"`
	AxesUIdAttr                                           string                                 `xml:"AxesUId,attr,omitempty"`
	CoordinateSystemsUIdAttr                              string                                 `xml:"CoordinateSystemsUId,attr,omitempty"`
	MotionGroupsUIdAttr                                   string                                 `xml:"MotionGroupsUId,attr,omitempty"`
	CSTUIdAttr                                            string                                 `xml:"CSTUId,attr,omitempty"`
	WallClockTimeUIdAttr                                  string                                 `xml:"WallClockTimeUId,attr,omitempty"`
	TrendGroupCollectionUIdAttr                           string                                 `xml:"TrendGroupCollectionUId,attr,omitempty"`
	DataLogCollectionUIdAttr                              string                                 `xml:"DataLogCollectionUId,attr,omitempty"`
	TimeSynchronizeUIdAttr                                string                                 `xml:"TimeSynchronizeUId,attr,omitempty"`
	FilePathAttr                                          string                                 `xml:"FilePath,attr,omitempty"`
	ProductCodeAttr                                       string                                 `xml:"ProductCode,attr,omitempty"`
	VerifiedAttr                                          string                                 `xml:"Verified,attr,omitempty"`
	ControllerEditsExistAttr                              string                                 `xml:"ControllerEditsExist,attr,omitempty"`
	PendingEditsExistAttr                                 string                                 `xml:"PendingEditsExist,attr,omitempty"`
	ForcesExistAttr                                       string                                 `xml:"ForcesExist,attr,omitempty"`
	EditResourceAttr                                      string                                 `xml:"EditResource,attr,omitempty"`
	FaultedAttr                                           string                                 `xml:"Faulted,attr,omitempty"`
	GeneralStatusAttr                                     string                                 `xml:"GeneralStatus,attr,omitempty"`
	IOMapLEDStatusAttr                                    string                                 `xml:"IOMapLEDStatus,attr,omitempty"`
	KeySwitchPositionAttr                                 string                                 `xml:"KeySwitchPosition,attr,omitempty"`
	ModeAttr                                              string                                 `xml:"Mode,attr,omitempty"`
	ConnectedStateAttr                                    string                                 `xml:"ConnectedState,attr,omitempty"`
	IsProjectDirtyAttr                                    string                                 `xml:"IsProjectDirty,attr,omitempty"`
	IsProjectDirtyFromAnotherWorkstationAttr              string                                 `xml:"IsProjectDirtyFromAnotherWorkstation,attr,omitempty"`
	EtherNetIPModeAttr                                    string                                 `xml:"EtherNetIPMode,attr,omitempty"`
	InstructionUsageConfiguredAOIsAttr                    string                                 `xml:"InstructionUsageConfiguredAOIs,attr,omitempty"`
	InstructionUsageConfiguredInstructionsAttr            string                                 `xml:"InstructionUsageConfiguredInstructions,attr,omitempty"`
	IsPlantPAxTaskingModelEnabledAttr                     string                                 `xml:"IsPlantPAxTaskingModelEnabled,attr,omitempty"`
	AutoDiagsEnabledAttr                                  string                                 `xml:"AutoDiagsEnabled,attr,omitempty"`
	WebServerEnabledAttr                                  string                                 `xml:"WebServerEnabled,attr,omitempty"`
	UseAttr                                               string                                 `xml:"Use,attr,omitempty"`
	CustomProperties                                      *CustomPropertiesCollectionType        `xml:"CustomProperties"`
	Description                                           *DescriptionType                       `xml:"Description"`
	RedundancyInfo                                        *RedundancyInfoType                    `xml:"RedundancyInfo"`
	Security                                              *SecurityInfoType                      `xml:"Security"`
	SafetyInfo                                            *SafetyInfoType                        `xml:"SafetyInfo"`
	DataTypes                                             *DataTypeCollectionType                `xml:"DataTypes"`
	DiagnosticMessages                                    *DiagnosticsMessageCollectionType      `xml:"DiagnosticMessages"`
	DeviceDiagnosticProfiles                              *DeviceDiagnosticProfileCollectionType `xml:"DeviceDiagnosticProfiles"`
	Modules                                               *MapDeviceCollectionType               `xml:"Modules"`
	AddOnInstructionDefinitions                           *UDIDefinitionCollectionType           `xml:"AddOnInstructionDefinitions"`
	AlarmDefinitions                                      *AlarmDefinitionCollectionAdaptorType  `xml:"AlarmDefinitions"`
	Tags                                                  *TagCollectionType                     `xml:"Tags"`
	Programs                                              *ProgramCollectionType                 `xml:"Programs"`
	Tasks                                                 *TaskCollectionType                    `xml:"Tasks"`
	ParameterConnections                                  *ParameterConnectionCollectionType     `xml:"ParameterConnections"`
	CommPorts                                             *CommPortCollectionType                `xml:"CommPorts"`
	CST                                                   *CSTType                               `xml:"CST"`
	WallClockTime                                         *WallClockTimeType                     `xml:"WallClockTime"`
	Trends                                                *TrendGroupCollectionType              `xml:"Trends"`
	DataLogs                                              *DataLogCollectionType                 `xml:"DataLogs"`
	QuickWatchLists                                       *QuickWatchCollectionAdaptorType       `xml:"QuickWatchLists"`
	TimeSynchronize                                       *TimeSynchronizeType                   `xml:"TimeSynchronize"`
	InternetProtocol                                      *TCPIPType                             `xml:"InternetProtocol"`
	EthernetPorts                                         *EthernetLinkCollectionType            `xml:"EthernetPorts"`
	EthernetNetwork                                       *DeviceLevelRingType                   `xml:"EthernetNetwork"`
}

// DescriptionType ...
type DescriptionType struct {
	UseAttr              string                          `xml:"Use,attr,omitempty"`
	CustomProperties     *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedDescription []*DescriptionTextType          `xml:"LocalizedDescription"`
	InnerValue           []byte                          `xml:",innerxml"`
}

func (d DescriptionType) CData() string {
	return ParseCData(d.InnerValue)
}

// DescriptionTextType ...
type DescriptionTextType struct {
	LangAttr         string                          `xml:"Lang,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// RedundancyInfoType ...
type RedundancyInfoType struct {
	UIdAttr                       string                          `xml:"UId,attr,omitempty"`
	EnabledAttr                   string                          `xml:"Enabled,attr,omitempty"`
	KeepTestEditsOnSwitchOverAttr string                          `xml:"KeepTestEditsOnSwitchOver,attr,omitempty"`
	IOMemoryPadPercentageAttr     string                          `xml:"IOMemoryPadPercentage,attr,omitempty"`
	DataTablePadPercentageAttr    string                          `xml:"DataTablePadPercentage,attr,omitempty"`
	UseAttr                       string                          `xml:"Use,attr,omitempty"`
	CustomProperties              *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// SecurityInfoType ...
type SecurityInfoType struct {
	CodeAttr                          string                          `xml:"Code,attr,omitempty"`
	SecurityAuthorityIDAttr           string                          `xml:"SecurityAuthorityID,attr,omitempty"`
	SecurityAuthorityURIAttr          string                          `xml:"SecurityAuthorityURI,attr,omitempty"`
	SecondarySecurityAuthorityIDAttr  string                          `xml:"SecondarySecurityAuthorityID,attr,omitempty"`
	SecondarySecurityAuthorityURIAttr string                          `xml:"SecondarySecurityAuthorityURI,attr,omitempty"`
	ChangesToDetectAttr               string                          `xml:"ChangesToDetect,attr,omitempty"`
	TrustedSlotsAttr                  string                          `xml:"TrustedSlots,attr,omitempty"`
	PermissionSetAttr                 string                          `xml:"PermissionSet,attr,omitempty"`
	UseAttr                           string                          `xml:"Use,attr,omitempty"`
	CustomProperties                  *CustomPropertiesCollectionType `xml:"CustomProperties"`
	PrimaryActionSets                 *PrimaryActionSetCollectionType `xml:"PrimaryActionSets"`
}

// SafetyTagMap ...
type SafetyTagMap struct {
	Value string `xml:",chardata"`
}

// SafetyInfoType ...
type SafetyInfoType struct {
	UIdAttr                     string                          `xml:"UId,attr,omitempty"`
	SafetySignatureAttr         string                          `xml:"SafetySignature,attr,omitempty"`
	SafetyLockedAttr            string                          `xml:"SafetyLocked,attr,omitempty"`
	SafetyLockPasswordAttr      string                          `xml:"SafetyLockPassword,attr,omitempty"`
	SafetyUnlockPasswordAttr    string                          `xml:"SafetyUnlockPassword,attr,omitempty"`
	SignatureRunModeProtectAttr string                          `xml:"SignatureRunModeProtect,attr,omitempty"`
	ConfigureSafetyIOAlwaysAttr string                          `xml:"ConfigureSafetyIOAlways,attr,omitempty"`
	SafetySignedAttr            string                          `xml:"SafetySigned,attr,omitempty"`
	SafetyLevelAttr             string                          `xml:"SafetyLevel,attr,omitempty"`
	UseAttr                     string                          `xml:"Use,attr,omitempty"`
	CustomProperties            *CustomPropertiesCollectionType `xml:"CustomProperties"`
	SafetyTaskFaultString       string                          `xml:"SafetyTaskFaultString"`
	SafetyTagMap                *SafetyTagMap                   `xml:"SafetyTagMap"`
}

// DataTypeCollectionType ...
type DataTypeCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	DataType         []*DataTypeType                 `xml:"DataType"`
}

// DataTypeType ...
type DataTypeType struct {
	NameAttr               string                          `xml:"Name,attr"`
	FamilyAttr             string                          `xml:"Family,attr,omitempty"`
	ClassAttr              string                          `xml:"Class,attr,omitempty"`
	UIdAttr                string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr          string                          `xml:"ParentUId,attr,omitempty"`
	SizeAttr               string                          `xml:"Size,attr,omitempty"`
	VerifiedAttr           string                          `xml:"Verified,attr,omitempty"`
	DeletionProhibitedAttr string                          `xml:"DeletionProhibited,attr,omitempty"`
	UseAttr                string                          `xml:"Use,attr,omitempty"`
	CustomProperties       *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description            *DescriptionType                `xml:"Description"`
	EngineeringUnit        *EngineeringUnitType            `xml:"EngineeringUnit"`
	Members                *TypeMemberCollectionType       `xml:"Members"`
	Dependencies           []*DependenciesWideType         `xml:"Dependencies"`
}

// TypeMemberCollectionType ...
type TypeMemberCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Member           []*DataTypeMemberType           `xml:"Member"`
}

// DataTypeMemberType ...
type DataTypeMemberType struct {
	NameAttr           string                          `xml:"Name,attr"`
	DataTypeAttr       string                          `xml:"DataType,attr"`
	DataTypeUIdAttr    string                          `xml:"DataTypeUId,attr,omitempty"`
	DimensionAttr      string                          `xml:"Dimension,attr,omitempty"`
	RadixAttr          string                          `xml:"Radix,attr,omitempty"`
	HiddenAttr         string                          `xml:"Hidden,attr,omitempty"`
	TargetAttr         string                          `xml:"Target,attr,omitempty"`
	BitNumberAttr      int                             `xml:"BitNumber,attr,omitempty"`
	SizeAttr           string                          `xml:"Size,attr,omitempty"`
	OffsetAttr         string                          `xml:"Offset,attr,omitempty"`
	MaxAttr            string                          `xml:"Max,attr,omitempty"`
	MinAttr            string                          `xml:"Min,attr,omitempty"`
	DefaultValueAttr   string                          `xml:"DefaultValue,attr,omitempty"`
	EditIDAttr         string                          `xml:"EditID,attr,omitempty"`
	ExternalAccessAttr string                          `xml:"ExternalAccess,attr,omitempty"`
	UIdAttr            string                          `xml:"UId,attr,omitempty"`
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description        *DescriptionType                `xml:"Description"`
	EngineeringUnit    *EngineeringUnitType            `xml:"EngineeringUnit"`
	State0             *State0Type                     `xml:"State0"`
	State1             *State1Type                     `xml:"State1"`
}

// DataTypeDependencyCollectionType ...
type DataTypeDependencyCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Dependency       []*DependencyWideType           `xml:"Dependency"`
}

// DataTypeDependencyType ...
type DataTypeDependencyType struct {
	TypeAttr         string                          `xml:"Type,attr,omitempty"`
	NameAttr         string                          `xml:"Name,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// DiagnosticsMessageCollectionType ...
type DiagnosticsMessageCollectionType struct {
	UseAttr           string                           `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType  `xml:"CustomProperties"`
	DiagnosticMessage []*DiagnosticsMessageAdaptorType `xml:"DiagnosticMessage"`
}

// DiagnosticsMessageAdaptorType ...
type DiagnosticsMessageAdaptorType struct {
	IDAttr           string                          `xml:"ID,attr"`
	VersionAttr      string                          `xml:"Version,attr"`
	DisplayCodeAttr  string                          `xml:"DisplayCode,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	DisplayText      []*DisplayTextType              `xml:"DisplayText"`
}

// DisplayTextType ...
type DisplayTextType struct {
	LangAttr         string                          `xml:"Lang,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// DeviceDiagnosticProfileCollectionType ...
type DeviceDiagnosticProfileCollectionType struct {
	UseAttr                 string                          `xml:"Use,attr,omitempty"`
	CustomProperties        *CustomPropertiesCollectionType `xml:"CustomProperties"`
	DeviceDiagnosticProfile []*DeviceDiagnosticProfileType  `xml:"DeviceDiagnosticProfile"`
}

// ProfileDefinition ...
type ProfileDefinition struct {
}

// DeviceDiagnosticProfileType ...
type DeviceDiagnosticProfileType struct {
	IDAttr                string                          `xml:"ID,attr"`
	VersionAttr           string                          `xml:"Version,attr"`
	ProfileReferencesAttr string                          `xml:"ProfileReferences,attr,omitempty"`
	MessageReferencesAttr string                          `xml:"MessageReferences,attr,omitempty"`
	UseAttr               string                          `xml:"Use,attr,omitempty"`
	CustomProperties      *CustomPropertiesCollectionType `xml:"CustomProperties"`
	ProfileDefinition     *ProfileDefinition              `xml:"ProfileDefinition"`
}

// MapDeviceCollectionType ...
type MapDeviceCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Module           []*MapDeviceType                `xml:"Module"`
}

// ExtendedProperties ...
type ExtendedProperties struct {
}

// MapDeviceType ...
type MapDeviceType struct {
	NameAttr                      string                             `xml:"Name,attr,omitempty"`
	UIdAttr                       string                             `xml:"UId,attr,omitempty"`
	CatalogNumberAttr             string                             `xml:"CatalogNumber,attr,omitempty"`
	VendorAttr                    string                             `xml:"Vendor,attr,omitempty"`
	ProductTypeAttr               string                             `xml:"ProductType,attr,omitempty"`
	ProductCodeAttr               string                             `xml:"ProductCode,attr,omitempty"`
	MajorAttr                     string                             `xml:"Major,attr,omitempty"`
	MinorAttr                     string                             `xml:"Minor,attr,omitempty"`
	UserDefinedVendorAttr         string                             `xml:"UserDefinedVendor,attr,omitempty"`
	UserDefinedProductTypeAttr    string                             `xml:"UserDefinedProductType,attr,omitempty"`
	UserDefinedProductCodeAttr    string                             `xml:"UserDefinedProductCode,attr,omitempty"`
	UserDefinedMajorAttr          string                             `xml:"UserDefinedMajor,attr,omitempty"`
	UserDefinedMinorAttr          string                             `xml:"UserDefinedMinor,attr,omitempty"`
	ParentUIdAttr                 string                             `xml:"ParentUId,attr,omitempty"`
	ParentModuleAttr              string                             `xml:"ParentModule,attr,omitempty"`
	ParentModuleUIdAttr           string                             `xml:"ParentModuleUId,attr,omitempty"`
	ParentModPortIdAttr           string                             `xml:"ParentModPortId,attr,omitempty"`
	InhibitedAttr                 string                             `xml:"Inhibited,attr,omitempty"`
	MajorFaultAttr                string                             `xml:"MajorFault,attr,omitempty"`
	DriverTypeAttr                string                             `xml:"DriverType,attr,omitempty"`
	ThrottleAttr                  string                             `xml:"Throttle,attr,omitempty"`
	ControlNetSignatureAttr       string                             `xml:"ControlNetSignature,attr,omitempty"`
	SafetyNetworkAttr             string                             `xml:"SafetyNetwork,attr,omitempty"`
	ConfigRevSupportedAttr        string                             `xml:"ConfigRevSupported,attr,omitempty"`
	BumplessReconfigSupportedAttr string                             `xml:"BumplessReconfigSupported,attr,omitempty"`
	IncludePortSegmentAttr        string                             `xml:"IncludePortSegment,attr,omitempty"`
	IncludeKeySegmentAttr         string                             `xml:"IncludeKeySegment,attr,omitempty"`
	ShutdownParentOnFaultAttr     string                             `xml:"ShutdownParentOnFault,attr,omitempty"`
	DrivesADCModeAttr             string                             `xml:"DrivesADCMode,attr,omitempty"`
	DrivesADCEnabledAttr          string                             `xml:"DrivesADCEnabled,attr,omitempty"`
	VerifiedAttr                  string                             `xml:"Verified,attr,omitempty"`
	FaultCodeAttr                 string                             `xml:"FaultCode,attr,omitempty"`
	FaultInfoAttr                 string                             `xml:"FaultInfo,attr,omitempty"`
	StatusAttr                    string                             `xml:"Status,attr,omitempty"`
	StatusStringAttr              string                             `xml:"StatusString,attr,omitempty"`
	UserDefinedCatalogNumberAttr  string                             `xml:"UserDefinedCatalogNumber,attr,omitempty"`
	SafetyEnabledAttr             string                             `xml:"SafetyEnabled,attr,omitempty"`
	TrackingGroupsAttr            string                             `xml:"TrackingGroups,attr,omitempty"`
	ImportCreateModeAttr          string                             `xml:"ImportCreateMode,attr,omitempty"`
	AutoDiagsEnabledAttr          string                             `xml:"AutoDiagsEnabled,attr,omitempty"`
	DiagnosticProfileIDAttr       string                             `xml:"DiagnosticProfileID,attr,omitempty"`
	DuplexEnabledAttr             string                             `xml:"DuplexEnabled,attr,omitempty"`
	UseAttr                       string                             `xml:"Use,attr,omitempty"`
	CustomProperties              *CustomPropertiesCollectionType    `xml:"CustomProperties"`
	Description                   *DescriptionType                   `xml:"Description"`
	FaultString                   string                             `xml:"FaultString"`
	EKey                          *ModuleEKeyType                    `xml:"EKey"`
	Ports                         *PortCollectionType                `xml:"Ports"`
	Communications                *CommunicationsType                `xml:"Communications"`
	ExtendedProperties            *ExtendedProperties                `xml:"ExtendedProperties"`
	PingMethod                    *PingMethodType                    `xml:"PingMethod"`
	SiblingDependencies           *SiblingDependenciesCollectionType `xml:"SiblingDependencies"`
}

// ModuleEKeyType ...
type ModuleEKeyType struct {
	StateAttr        string                          `xml:"State,attr"`
	RelaxedAttr      string                          `xml:"Relaxed,attr,omitempty"`
	VendorAttr       string                          `xml:"Vendor,attr,omitempty"`
	ProductCodeAttr  string                          `xml:"ProductCode,attr,omitempty"`
	ProductTypeAttr  string                          `xml:"ProductType,attr,omitempty"`
	MajorAttr        string                          `xml:"Major,attr,omitempty"`
	MinorAttr        string                          `xml:"Minor,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// PortCollectionType ...
type PortCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Port             []*PortType                     `xml:"Port"`
}

// PortType ...
type PortType struct {
	IdAttr                 string                          `xml:"Id,attr"`
	UIdAttr                string                          `xml:"UId,attr,omitempty"`
	AddressAttr            string                          `xml:"Address,attr,omitempty"`
	AddressBAttr           string                          `xml:"AddressB,attr,omitempty"`
	RcpGatewayAddress1Attr string                          `xml:"RcpGatewayAddress1,attr,omitempty"`
	RcpGatewayAddress2Attr string                          `xml:"RcpGatewayAddress2,attr,omitempty"`
	TypeAttr               string                          `xml:"Type,attr,omitempty"`
	UpstreamAttr           string                          `xml:"Upstream,attr,omitempty"`
	ConnectorOffsetAttr    string                          `xml:"ConnectorOffset,attr,omitempty"`
	WidthAttr              string                          `xml:"Width,attr,omitempty"`
	NATActualAddressAttr   string                          `xml:"NATActualAddress,attr,omitempty"`
	ValidationAttr         string                          `xml:"Validation,attr,omitempty"`
	SafetyNetworkAttr      string                          `xml:"SafetyNetwork,attr,omitempty"`
	UseAttr                string                          `xml:"Use,attr,omitempty"`
	CustomProperties       *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Bus                    *ABusType                       `xml:"Bus"`
}

// ABusType ...
type ABusType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	SizeAttr         string                          `xml:"Size,attr,omitempty"`
	BaudAttr         string                          `xml:"Baud,attr,omitempty"`
	AddressMaskAttr  string                          `xml:"AddressMask,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// CommunicationsType ...
type CommunicationsType struct {
	CommMethodAttr        string                          `xml:"CommMethod,attr,omitempty"`
	ConfigTagUIdAttr      string                          `xml:"ConfigTagUId,attr,omitempty"`
	PrimCxnInputSizeAttr  string                          `xml:"PrimCxnInputSize,attr,omitempty"`
	PrimCxnOutputSizeAttr string                          `xml:"PrimCxnOutputSize,attr,omitempty"`
	SecCxnInputSizeAttr   string                          `xml:"SecCxnInputSize,attr,omitempty"`
	SecCxnOutputSizeAttr  string                          `xml:"SecCxnOutputSize,attr,omitempty"`
	UseAttr               string                          `xml:"Use,attr,omitempty"`
	CustomProperties      *CustomPropertiesCollectionType `xml:"CustomProperties"`
	ConfigData            *ConfigDataType                 `xml:"ConfigData"`
	ConfigTag             *ConfigTagType                  `xml:"ConfigTag"`
	ConfigScript          *ConfigScriptType               `xml:"ConfigScript"`
	SafetyScript          *SafetyScriptType               `xml:"SafetyScript"`
	Connections           *MapConnectionCollectionType    `xml:"Connections"`
}

// ConfigTagType ...
type ConfigTagType struct {
	SuffixAttr         string                            `xml:"Suffix,attr,omitempty"`
	UIdAttr            string                            `xml:"UId,attr,omitempty"`
	ParentUIdAttr      string                            `xml:"ParentUId,attr,omitempty"`
	ConfigSizeAttr     string                            `xml:"ConfigSize,attr,omitempty"`
	DataTypeAttr       string                            `xml:"DataType,attr,omitempty"`
	DimensionsAttr     string                            `xml:"Dimensions,attr,omitempty"`
	RadixAttr          string                            `xml:"Radix,attr,omitempty"`
	MaxAttr            string                            `xml:"Max,attr,omitempty"`
	MinAttr            string                            `xml:"Min,attr,omitempty"`
	ExternalAccessAttr string                            `xml:"ExternalAccess,attr,omitempty"`
	ConstantAttr       string                            `xml:"Constant,attr,omitempty"`
	PermissionSetAttr  string                            `xml:"PermissionSet,attr,omitempty"`
	VerifiedAttr       string                            `xml:"Verified,attr,omitempty"`
	TrackingGroupsAttr string                            `xml:"TrackingGroups,attr,omitempty"`
	UseAttr            string                            `xml:"Use,attr,omitempty"`
	CustomProperties   []*CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description        []*DescriptionType                `xml:"Description"`
	Data               []*DataWideType                   `xml:"Data"`
	Comments           []*CommentCollectionType          `xml:"Comments"`
	EngineeringUnits   []*EngineeringUnitCollectionType  `xml:"EngineeringUnits"`
	Maxes              []*MaxLimitCollectionType         `xml:"Maxes"`
	Mins               []*MinLimitCollectionType         `xml:"Mins"`
	State0s            []*State0CollectionType           `xml:"State0s"`
	State1s            []*State1CollectionType           `xml:"State1s"`
}

// ConfigDataType ...
type ConfigDataType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ConfigSizeAttr   string                          `xml:"ConfigSize,attr,omitempty"`
	DataTypeAttr     string                          `xml:"DataType,attr,omitempty"`
	DimensionAttr    string                          `xml:"Dimension,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Data             []*DataWideType                 `xml:"Data"`
}

// ConfigScriptType ...
type ConfigScriptType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	SizeAttr         string                          `xml:"Size,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Data             []*DataWideType                 `xml:"Data"`
}

// SafetyScriptType ...
type SafetyScriptType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	SizeAttr         string                          `xml:"Size,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Data             []*DataWideType                 `xml:"Data"`
}

// MapConnectionCollectionType ...
type MapConnectionCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Connection       []*MapConnectionType            `xml:"Connection"`
	RackConnection   *RackConnectionType             `xml:"RackConnection"`
}

// MapConnectionType ...
type MapConnectionType struct {
	NameAttr                             string                          `xml:"Name,attr"`
	UIdAttr                              string                          `xml:"UId,attr,omitempty"`
	RPIAttr                              string                          `xml:"RPI,attr"`
	TypeAttr                             string                          `xml:"Type,attr,omitempty"`
	InputTagUIdAttr                      string                          `xml:"InputTagUId,attr,omitempty"`
	OutputTagUIdAttr                     string                          `xml:"OutputTagUId,attr,omitempty"`
	InputCxnPointAttr                    string                          `xml:"InputCxnPoint,attr,omitempty"`
	OutputCxnPointAttr                   string                          `xml:"OutputCxnPoint,attr,omitempty"`
	InputSizeAttr                        string                          `xml:"InputSize,attr,omitempty"`
	OutputSizeAttr                       string                          `xml:"OutputSize,attr,omitempty"`
	EventIDAttr                          string                          `xml:"EventID,attr,omitempty"`
	ProgrammaticallySendEventTriggerAttr string                          `xml:"ProgrammaticallySendEventTrigger,attr,omitempty"`
	TimeoutMultiplierAttr                string                          `xml:"TimeoutMultiplier,attr,omitempty"`
	NetworkDelayMultiplierAttr           string                          `xml:"NetworkDelayMultiplier,attr,omitempty"`
	ReactionTimeLimitAttr                string                          `xml:"ReactionTimeLimit,attr,omitempty"`
	MaxObservedNetworkDelayAttr          string                          `xml:"MaxObservedNetworkDelay,attr,omitempty"`
	ControlNetScheduledAttr              string                          `xml:"ControlNetScheduled,attr,omitempty"`
	UnicastAttr                          string                          `xml:"Unicast,attr,omitempty"`
	OpenOrderAttr                        string                          `xml:"OpenOrder,attr,omitempty"`
	InputRealTimeFormatAttr              string                          `xml:"InputRealTimeFormat,attr,omitempty"`
	OutputRealTimeFormatAttr             string                          `xml:"OutputRealTimeFormat,attr,omitempty"`
	ConnMultiplierAttr                   string                          `xml:"ConnMultiplier,attr,omitempty"`
	OutputSizeVariableAttr               string                          `xml:"OutputSizeVariable,attr,omitempty"`
	InputSizeVariableAttr                string                          `xml:"InputSizeVariable,attr,omitempty"`
	OutputPriorityAttr                   string                          `xml:"OutputPriority,attr,omitempty"`
	PriorityAttr                         string                          `xml:"Priority,attr,omitempty"`
	OutputConnectionTypeAttr             string                          `xml:"OutputConnectionType,attr,omitempty"`
	InputConnectionTypeAttr              string                          `xml:"InputConnectionType,attr,omitempty"`
	OutputRedundantOwnerAttr             string                          `xml:"OutputRedundantOwner,attr,omitempty"`
	OutputRPIAttr                        string                          `xml:"OutputRPI,attr,omitempty"`
	TransportClassAttr                   string                          `xml:"TransportClass,attr,omitempty"`
	InputProductionTriggerAttr           string                          `xml:"InputProductionTrigger,attr,omitempty"`
	DirectionAttr                        string                          `xml:"Direction,attr,omitempty"`
	ConnectionPathAttr                   string                          `xml:"ConnectionPath,attr,omitempty"`
	IncludeDataSegmentAttr               string                          `xml:"IncludeDataSegment,attr,omitempty"`
	GenerateFaultMemberAttr              string                          `xml:"GenerateFaultMember,attr,omitempty"`
	InputTagSuffixAttr                   string                          `xml:"InputTagSuffix,attr,omitempty"`
	OutputTagSuffixAttr                  string                          `xml:"OutputTagSuffix,attr,omitempty"`
	DisplayNameAttr                      string                          `xml:"DisplayName,attr,omitempty"`
	MaxConsumerNumberAttr                string                          `xml:"MaxConsumerNumber,attr,omitempty"`
	ExcludeFromAggregationAttr           string                          `xml:"ExcludeFromAggregation,attr,omitempty"`
	IsConcurrentAttr                     string                          `xml:"IsConcurrent,attr,omitempty"`
	UseAttr                              string                          `xml:"Use,attr,omitempty"`
	CustomProperties                     *CustomPropertiesCollectionType `xml:"CustomProperties"`
	InputTag                             *InputTagType                   `xml:"InputTag"`
	InputData                            *InputDataType                  `xml:"InputData"`
	OutputTag                            *OutputTagType                  `xml:"OutputTag"`
	OutputData                           *OutputDataType                 `xml:"OutputData"`
	InputImageInit                       *InputImageInitType             `xml:"InputImageInit"`
}

// InputTagType ...
type InputTagType struct {
	SuffixAttr         string                            `xml:"Suffix,attr,omitempty"`
	UIdAttr            string                            `xml:"UId,attr,omitempty"`
	ParentUIdAttr      string                            `xml:"ParentUId,attr,omitempty"`
	DataTypeAttr       string                            `xml:"DataType,attr,omitempty"`
	MaxAttr            string                            `xml:"Max,attr,omitempty"`
	MinAttr            string                            `xml:"Min,attr,omitempty"`
	ExternalAccessAttr string                            `xml:"ExternalAccess,attr,omitempty"`
	ConstantAttr       string                            `xml:"Constant,attr,omitempty"`
	PermissionSetAttr  string                            `xml:"PermissionSet,attr,omitempty"`
	VerifiedAttr       string                            `xml:"Verified,attr,omitempty"`
	TrackingGroupsAttr string                            `xml:"TrackingGroups,attr,omitempty"`
	UseAttr            string                            `xml:"Use,attr,omitempty"`
	CustomProperties   []*CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description        []*DescriptionType                `xml:"Description"`
	Comments           []*CommentCollectionType          `xml:"Comments"`
	EngineeringUnits   []*EngineeringUnitCollectionType  `xml:"EngineeringUnits"`
	Maxes              []*MaxLimitCollectionType         `xml:"Maxes"`
	Mins               []*MinLimitCollectionType         `xml:"Mins"`
	State0s            []*State0CollectionType           `xml:"State0s"`
	State1s            []*State1CollectionType           `xml:"State1s"`
	Data               []*DataWideType                   `xml:"Data"`
	ForceData          []*ForceDataWideType              `xml:"ForceData"`
}

// InputDataType ...
type InputDataType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	DataTypeAttr     string                          `xml:"DataType,attr,omitempty"`
	DimensionAttr    string                          `xml:"Dimension,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Data             []*DataWideType                 `xml:"Data"`
}

// OutputTagType ...
type OutputTagType struct {
	SuffixAttr         string                            `xml:"Suffix,attr,omitempty"`
	UIdAttr            string                            `xml:"UId,attr,omitempty"`
	ParentUIdAttr      string                            `xml:"ParentUId,attr,omitempty"`
	DataTypeAttr       string                            `xml:"DataType,attr,omitempty"`
	MaxAttr            string                            `xml:"Max,attr,omitempty"`
	MinAttr            string                            `xml:"Min,attr,omitempty"`
	ExternalAccessAttr string                            `xml:"ExternalAccess,attr,omitempty"`
	ConstantAttr       string                            `xml:"Constant,attr,omitempty"`
	PermissionSetAttr  string                            `xml:"PermissionSet,attr,omitempty"`
	VerifiedAttr       string                            `xml:"Verified,attr,omitempty"`
	TrackingGroupsAttr string                            `xml:"TrackingGroups,attr,omitempty"`
	UseAttr            string                            `xml:"Use,attr,omitempty"`
	CustomProperties   []*CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description        []*DescriptionType                `xml:"Description"`
	Comments           []*CommentCollectionType          `xml:"Comments"`
	EngineeringUnits   []*EngineeringUnitCollectionType  `xml:"EngineeringUnits"`
	Maxes              []*MaxLimitCollectionType         `xml:"Maxes"`
	Mins               []*MinLimitCollectionType         `xml:"Mins"`
	State0s            []*State0CollectionType           `xml:"State0s"`
	State1s            []*State1CollectionType           `xml:"State1s"`
	Data               []*DataWideType                   `xml:"Data"`
	ForceData          []*ForceDataWideType              `xml:"ForceData"`
}

// OutputDataType ...
type OutputDataType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	DataTypeAttr     string                          `xml:"DataType,attr,omitempty"`
	DimensionAttr    string                          `xml:"Dimension,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Data             []*DataWideType                 `xml:"Data"`
}

// InputImageInitType ...
type InputImageInitType struct {
	UIdAttr           string                          `xml:"UId,attr,omitempty"`
	TypeAttr          string                          `xml:"Type,attr"`
	SizeAttr          string                          `xml:"Size,attr,omitempty"`
	MaskValueSizeAttr string                          `xml:"MaskValueSize,attr,omitempty"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Data              []*DataWideType                 `xml:"Data"`
}

// RackConnectionType ...
type RackConnectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	InAliasTag       *InAliasTagType                 `xml:"InAliasTag"`
	OutAliasTag      *OutAliasTagType                `xml:"OutAliasTag"`
}

// InAliasTagType ...
type InAliasTagType struct {
	SuffixAttr       string                          `xml:"Suffix,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	MaxAttr          string                          `xml:"Max,attr,omitempty"`
	MinAttr          string                          `xml:"Min,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description      *DescriptionType                `xml:"Description"`
	Comments         *CommentCollectionType          `xml:"Comments"`
	EngineeringUnits *EngineeringUnitCollectionType  `xml:"EngineeringUnits"`
	Maxes            *MaxLimitCollectionType         `xml:"Maxes"`
	Mins             *MinLimitCollectionType         `xml:"Mins"`
	State0s          *State0CollectionType           `xml:"State0s"`
	State1s          *State1CollectionType           `xml:"State1s"`
}

// OutAliasTagType ...
type OutAliasTagType struct {
	SuffixAttr       string                          `xml:"Suffix,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	MaxAttr          string                          `xml:"Max,attr,omitempty"`
	MinAttr          string                          `xml:"Min,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description      *DescriptionType                `xml:"Description"`
	Comments         *CommentCollectionType          `xml:"Comments"`
	EngineeringUnits *EngineeringUnitCollectionType  `xml:"EngineeringUnits"`
	Maxes            *MaxLimitCollectionType         `xml:"Maxes"`
	Mins             *MinLimitCollectionType         `xml:"Mins"`
	State0s          *State0CollectionType           `xml:"State0s"`
	State1s          *State1CollectionType           `xml:"State1s"`
}

// PingMethodType ...
type PingMethodType struct {
	TypeAttr         string                          `xml:"Type,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// SiblingDependenciesCollectionType ...
type SiblingDependenciesCollectionType struct {
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
	SiblingDependency []*SiblingDependencyAdaptorType `xml:"SiblingDependency"`
}

// SiblingDependencyAdaptorType ...
type SiblingDependencyAdaptorType struct {
	ParentModPortIdAttr string                          `xml:"ParentModPortId,attr"`
	AddressAttr         string                          `xml:"Address,attr"`
	UseAttr             string                          `xml:"Use,attr,omitempty"`
	CustomProperties    *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// UDIDefinitionCollectionType ...
type UDIDefinitionCollectionType struct {
	UIdAttr                    string                            `xml:"UId,attr,omitempty"`
	ParentUIdAttr              string                            `xml:"ParentUId,attr,omitempty"`
	UseAttr                    string                            `xml:"Use,attr,omitempty"`
	EncodedData                []*EncodedUDIDefinitionType       `xml:"EncodedData"`
	CustomProperties           []*CustomPropertiesCollectionType `xml:"CustomProperties"`
	AddOnInstructionDefinition []*UDIDefinitionType              `xml:"AddOnInstructionDefinition"`
}

// EncryptedAOIContent ...
type EncryptedAOIContent struct {
	Value string `xml:",chardata"`
}

// UDIDefinitionType ...
type UDIDefinitionType struct {
	NameAttr                 string                          `xml:"Name,attr"`
	UIdAttr                  string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr            string                          `xml:"ParentUId,attr,omitempty"`
	ClassAttr                string                          `xml:"Class,attr,omitempty"`
	RevisionAttr             string                          `xml:"Revision,attr,omitempty"`
	RevisionExtensionAttr    string                          `xml:"RevisionExtension,attr,omitempty"`
	VendorAttr               string                          `xml:"Vendor,attr,omitempty"`
	ExecutePrescanAttr       string                          `xml:"ExecutePrescan,attr,omitempty"`
	ExecutePostscanAttr      string                          `xml:"ExecutePostscan,attr,omitempty"`
	ExecuteEnableInFalseAttr string                          `xml:"ExecuteEnableInFalse,attr,omitempty"`
	SignatureIDAttr          string                          `xml:"SignatureID,attr,omitempty"`
	SignatureTimestampAttr   string                          `xml:"SignatureTimestamp,attr,omitempty"`
	SafetySignatureIDAttr    string                          `xml:"SafetySignatureID,attr,omitempty"`
	CreatedDateAttr          string                          `xml:"CreatedDate,attr,omitempty"`
	CreatedByAttr            string                          `xml:"CreatedBy,attr,omitempty"`
	EditedDateAttr           string                          `xml:"EditedDate,attr,omitempty"`
	EditedByAttr             string                          `xml:"EditedBy,attr,omitempty"`
	SoftwareRevisionAttr     string                          `xml:"SoftwareRevision,attr,omitempty"`
	OriginalLibraryAttr      string                          `xml:"OriginalLibrary,attr,omitempty"`
	OriginalNameAttr         string                          `xml:"OriginalName,attr,omitempty"`
	OriginalRevisionAttr     string                          `xml:"OriginalRevision,attr,omitempty"`
	OriginalVendorAttr       string                          `xml:"OriginalVendor,attr,omitempty"`
	AddedToProjectAttr       string                          `xml:"AddedToProject,attr,omitempty"`
	SourceKeyAttr            string                          `xml:"SourceKey,attr,omitempty"`
	EncodedSourceKeyAttr     string                          `xml:"EncodedSourceKey,attr,omitempty"`
	SourceProtectionTypeAttr string                          `xml:"SourceProtectionType,attr,omitempty"`
	IsEncryptedAttr          string                          `xml:"IsEncrypted,attr,omitempty"`
	PermissionSetAttr        string                          `xml:"PermissionSet,attr,omitempty"`
	TrackingGroupsAttr       string                          `xml:"TrackingGroups,attr,omitempty"`
	LogicHashAttr            string                          `xml:"LogicHash,attr,omitempty"`
	DescriptionHashAttr      string                          `xml:"DescriptionHash,attr,omitempty"`
	UseAttr                  string                          `xml:"Use,attr,omitempty"`
	CustomProperties         *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description              *DescriptionType                `xml:"Description"`
	RevisionNote             *RevisionNoteType               `xml:"RevisionNote"`
	SignatureHistory         *SignatureHistoryType           `xml:"SignatureHistory"`
	AdditionalHelpText       *AdditionalHelpTextType         `xml:"AdditionalHelpText"`
	EncryptionInfo           *EncryptionInfoType             `xml:"EncryptionInfo"`
	Parameters               *UDIParameterCollectionType     `xml:"Parameters"`
	LocalTags                *UDILocalTagCollectionType      `xml:"LocalTags"`
	Tags                     *TagCollectionType              `xml:"Tags"`
	Routines                 *RoutineCollectionType          `xml:"Routines"`
	Dependencies             []*DependenciesWideType         `xml:"Dependencies"`
	EncryptedAOIContent      *EncryptedAOIContent            `xml:"EncryptedAOIContent"`
}

// EncodedUDIDefinitionType ...
type EncodedUDIDefinitionType struct {
	EncodedTypeAttr        string                          `xml:"EncodedType,attr,omitempty"`
	NameAttr               string                          `xml:"Name,attr,omitempty"`
	UIdAttr                string                          `xml:"UId,attr,omitempty"`
	ClassAttr              string                          `xml:"Class,attr,omitempty"`
	RevisionAttr           string                          `xml:"Revision,attr,omitempty"`
	RevisionExtensionAttr  string                          `xml:"RevisionExtension,attr,omitempty"`
	VendorAttr             string                          `xml:"Vendor,attr,omitempty"`
	SignatureIDAttr        string                          `xml:"SignatureID,attr,omitempty"`
	SignatureTimestampAttr string                          `xml:"SignatureTimestamp,attr,omitempty"`
	SafetySignatureIDAttr  string                          `xml:"SafetySignatureID,attr,omitempty"`
	EditedDateAttr         string                          `xml:"EditedDate,attr,omitempty"`
	SoftwareRevisionAttr   string                          `xml:"SoftwareRevision,attr,omitempty"`
	OriginalLibraryAttr    string                          `xml:"OriginalLibrary,attr,omitempty"`
	OriginalNameAttr       string                          `xml:"OriginalName,attr,omitempty"`
	OriginalRevisionAttr   string                          `xml:"OriginalRevision,attr,omitempty"`
	OriginalVendorAttr     string                          `xml:"OriginalVendor,attr,omitempty"`
	PermissionSetAttr      string                          `xml:"PermissionSet,attr,omitempty"`
	EncryptionConfigAttr   string                          `xml:"EncryptionConfig,attr,omitempty"`
	IsEncryptedAttr        string                          `xml:"IsEncrypted,attr,omitempty"`
	LogicHashAttr          string                          `xml:"LogicHash,attr,omitempty"`
	DescriptionHashAttr    string                          `xml:"DescriptionHash,attr,omitempty"`
	UseAttr                string                          `xml:"Use,attr,omitempty"`
	CustomProperties       *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description            *DescriptionType                `xml:"Description"`
	RevisionNote           *RevisionNoteType               `xml:"RevisionNote"`
	SignatureHistory       *SignatureHistoryType           `xml:"SignatureHistory"`
	AdditionalHelpText     *AdditionalHelpTextType         `xml:"AdditionalHelpText"`
	Parameters             *UDIParameterCollectionType     `xml:"Parameters"`
	Dependencies           []*DependenciesWideType         `xml:"Dependencies"`
}

// RevisionNoteType ...
type RevisionNoteType struct {
	UseAttr               string                          `xml:"Use,attr,omitempty"`
	CustomProperties      *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedRevisionNote []*RevisionNoteTextType         `xml:"LocalizedRevisionNote"`
}

// RevisionNoteTextType ...
type RevisionNoteTextType struct {
	LangAttr         string                          `xml:"Lang,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// SignatureHistoryType ...
type SignatureHistoryType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	HistoryEntry     []*HistoryEntryType             `xml:"HistoryEntry"`
}

// HistoryEntryType ...
type HistoryEntryType struct {
	UserAttr         string                          `xml:"User,attr"`
	TimestampAttr    string                          `xml:"Timestamp,attr"`
	SignatureIDAttr  string                          `xml:"SignatureID,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description      *DescriptionType                `xml:"Description"`
}

// AdditionalHelpTextType ...
type AdditionalHelpTextType struct {
	UseAttr                     string                             `xml:"Use,attr,omitempty"`
	CustomProperties            *CustomPropertiesCollectionType    `xml:"CustomProperties"`
	LocalizedAdditionalHelpText []*LocalizedAdditionalHelpTextType `xml:"LocalizedAdditionalHelpText"`
}

// LocalizedAdditionalHelpTextType ...
type LocalizedAdditionalHelpTextType struct {
	LangAttr         string                          `xml:"Lang,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// UDIParameterCollectionType ...
type UDIParameterCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Parameter        []*UDIParameterType             `xml:"Parameter"`
}

// UDIParameterType ...
type UDIParameterType struct {
	NameAttr           string                          `xml:"Name,attr"`
	UIdAttr            string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr      string                          `xml:"ParentUId,attr,omitempty"`
	TagTypeAttr        string                          `xml:"TagType,attr,omitempty"`
	DataTypeAttr       string                          `xml:"DataType,attr,omitempty"`
	DataTypeUIdAttr    string                          `xml:"DataTypeUId,attr,omitempty"`
	DimensionsAttr     string                          `xml:"Dimensions,attr,omitempty"`
	UsageAttr          string                          `xml:"Usage,attr"`
	RadixAttr          string                          `xml:"Radix,attr,omitempty"`
	AliasForAttr       string                          `xml:"AliasFor,attr,omitempty"`
	AliasBaseAttr      string                          `xml:"AliasBase,attr,omitempty"`
	RequiredAttr       string                          `xml:"Required,attr,omitempty"`
	VisibleAttr        string                          `xml:"Visible,attr,omitempty"`
	ConstantAttr       string                          `xml:"Constant,attr,omitempty"`
	ExternalAccessAttr string                          `xml:"ExternalAccess,attr,omitempty"`
	MaxAttr            string                          `xml:"Max,attr,omitempty"`
	MinAttr            string                          `xml:"Min,attr,omitempty"`
	VerifiedAttr       string                          `xml:"Verified,attr,omitempty"`
	CanBeNullAttr      string                          `xml:"CanBeNull,attr,omitempty"`
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description        *DescriptionType                `xml:"Description"`
	Comments           *CommentCollectionType          `xml:"Comments"`
	EngineeringUnits   *EngineeringUnitCollectionType  `xml:"EngineeringUnits"`
	Maxes              *MaxLimitCollectionType         `xml:"Maxes"`
	Mins               *MinLimitCollectionType         `xml:"Mins"`
	State0s            *State0CollectionType           `xml:"State0s"`
	State1s            *State1CollectionType           `xml:"State1s"`
	DefaultData        []*DefaultDataWideType          `xml:"DefaultData"`
}

// UDILocalTagCollectionType ...
type UDILocalTagCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalTag         []*UDILocalTagType              `xml:"LocalTag"`
}

// UDILocalTagType ...
type UDILocalTagType struct {
	NameAttr           string                          `xml:"Name,attr"`
	UIdAttr            string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr      string                          `xml:"ParentUId,attr,omitempty"`
	DataTypeAttr       string                          `xml:"DataType,attr"`
	DataTypeUIdAttr    string                          `xml:"DataTypeUId,attr,omitempty"`
	DimensionsAttr     string                          `xml:"Dimensions,attr,omitempty"`
	RadixAttr          string                          `xml:"Radix,attr,omitempty"`
	ExternalAccessAttr string                          `xml:"ExternalAccess,attr,omitempty"`
	MaxAttr            string                          `xml:"Max,attr,omitempty"`
	MinAttr            string                          `xml:"Min,attr,omitempty"`
	VerifiedAttr       string                          `xml:"Verified,attr,omitempty"`
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description        *DescriptionType                `xml:"Description"`
	Comments           *CommentCollectionType          `xml:"Comments"`
	EngineeringUnits   *EngineeringUnitCollectionType  `xml:"EngineeringUnits"`
	Maxes              *MaxLimitCollectionType         `xml:"Maxes"`
	Mins               *MinLimitCollectionType         `xml:"Mins"`
	State0s            *State0CollectionType           `xml:"State0s"`
	State1s            *State1CollectionType           `xml:"State1s"`
	DefaultData        []*DefaultDataWideType          `xml:"DefaultData"`
}

// UDIDependencyCollectionType ...
type UDIDependencyCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Dependency       []*DependencyWideType           `xml:"Dependency"`
}

// UDIDependencyType ...
type UDIDependencyType struct {
	TypeAttr         string                          `xml:"Type,attr,omitempty"`
	NameAttr         string                          `xml:"Name,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// RawDefaultDataType ...
type RawDefaultDataType struct {
	FormatAttr       string                          `xml:"Format,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// L5KDefaultDataType ...
type L5KDefaultDataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// StringDefaultDataType ...
type StringDefaultDataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	LengthAttr       string                          `xml:"Length,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// DecoratedDefaultDataType ...
type DecoratedDefaultDataType struct {
}

// AlarmDefinitionCollectionAdaptorType ...
type AlarmDefinitionCollectionAdaptorType struct {
	UseAttr                 string                          `xml:"Use,attr,omitempty"`
	CustomProperties        *CustomPropertiesCollectionType `xml:"CustomProperties"`
	DatatypeAlarmDefinition []*DatatypeAlarmDefinitionType  `xml:"DatatypeAlarmDefinition"`
}

// DatatypeAlarmDefinitionType ...
type DatatypeAlarmDefinitionType struct {
	NameAttr              string                          `xml:"Name,attr"`
	UseAttr               string                          `xml:"Use,attr,omitempty"`
	CustomProperties      *CustomPropertiesCollectionType `xml:"CustomProperties"`
	MemberAlarmDefinition []*MemberAlarmDefinitionType    `xml:"MemberAlarmDefinition"`
}

// MemberAlarmDefinitionType ...
type MemberAlarmDefinitionType struct {
	AckRequiredAttr            string                          `xml:"AckRequired,attr,omitempty"`
	LatchedAttr                string                          `xml:"Latched,attr,omitempty"`
	RequiredAttr               string                          `xml:"Required,attr,omitempty"`
	AlarmSetOperIncludedAttr   string                          `xml:"AlarmSetOperIncluded,attr,omitempty"`
	AlarmSetRollupIncludedAttr string                          `xml:"AlarmSetRollupIncluded,attr,omitempty"`
	NameAttr                   string                          `xml:"Name,attr"`
	InputAttr                  string                          `xml:"Input,attr"`
	ConditionTypeAttr          string                          `xml:"ConditionType,attr,omitempty"`
	LimitAttr                  string                          `xml:"Limit,attr,omitempty"`
	SeverityAttr               string                          `xml:"Severity,attr,omitempty"`
	OnDelayAttr                string                          `xml:"OnDelay,attr,omitempty"`
	OffDelayAttr               string                          `xml:"OffDelay,attr,omitempty"`
	ShelveDurationAttr         string                          `xml:"ShelveDuration,attr,omitempty"`
	MaxShelveDurationAttr      string                          `xml:"MaxShelveDuration,attr,omitempty"`
	DeadbandAttr               string                          `xml:"Deadband,attr,omitempty"`
	EvaluationPeriodAttr       string                          `xml:"EvaluationPeriod,attr,omitempty"`
	ExpressionAttr             string                          `xml:"Expression,attr,omitempty"`
	TargetTagAttr              string                          `xml:"TargetTag,attr,omitempty"`
	AssocTag1Attr              string                          `xml:"AssocTag1,attr,omitempty"`
	AssocTag2Attr              string                          `xml:"AssocTag2,attr,omitempty"`
	AssocTag3Attr              string                          `xml:"AssocTag3,attr,omitempty"`
	AssocTag4Attr              string                          `xml:"AssocTag4,attr,omitempty"`
	UIdAttr                    string                          `xml:"UId,attr,omitempty"`
	IsTemplateAttr             string                          `xml:"IsTemplate,attr,omitempty"`
	UseAttr                    string                          `xml:"Use,attr,omitempty"`
	CustomProperties           *CustomPropertiesCollectionType `xml:"CustomProperties"`
	AlarmConfig                *AlarmConfigType                `xml:"AlarmConfig"`
}

// TagCollectionType ...
type TagCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Tag              []*TagType                      `xml:"Tag"`
}

// TagType ...
type TagType struct {
	NameAttr           string                                  `xml:"Name,attr"`
	UIdAttr            string                                  `xml:"UId,attr,omitempty"`
	ParentUIdAttr      string                                  `xml:"ParentUId,attr,omitempty"`
	ClassAttr          string                                  `xml:"Class,attr,omitempty"`
	TagTypeAttr        string                                  `xml:"TagType,attr,omitempty"`
	DataTypeAttr       string                                  `xml:"DataType,attr,omitempty"`
	DataTypeUIdAttr    string                                  `xml:"DataTypeUId,attr,omitempty"`
	DimensionsAttr     string                                  `xml:"Dimensions,attr,omitempty"`
	RadixAttr          string                                  `xml:"Radix,attr,omitempty"`
	AliasForAttr       string                                  `xml:"AliasFor,attr,omitempty"`
	AliasBaseAttr      string                                  `xml:"AliasBase,attr,omitempty"`
	UsageAttr          string                                  `xml:"Usage,attr,omitempty"`
	SequencingAttr     string                                  `xml:"Sequencing,attr,omitempty"`
	ConstantAttr       string                                  `xml:"Constant,attr,omitempty"`
	RequiredAttr       string                                  `xml:"Required,attr,omitempty"`
	VisibleAttr        string                                  `xml:"Visible,attr,omitempty"`
	ExternalAccessAttr string                                  `xml:"ExternalAccess,attr,omitempty"`
	MaxAttr            string                                  `xml:"Max,attr,omitempty"`
	MinAttr            string                                  `xml:"Min,attr,omitempty"`
	IOAttr             string                                  `xml:"IO,attr,omitempty"`
	PermissionSetAttr  string                                  `xml:"PermissionSet,attr,omitempty"`
	ScopeAttr          string                                  `xml:"Scope,attr,omitempty"`
	CanForceAttr       string                                  `xml:"CanForce,attr,omitempty"`
	FullNameAttr       string                                  `xml:"FullName,attr,omitempty"`
	VerifiedAttr       string                                  `xml:"Verified,attr,omitempty"`
	TrackingGroupsAttr string                                  `xml:"TrackingGroups,attr,omitempty"`
	CanBeNullAttr      string                                  `xml:"CanBeNull,attr,omitempty"`
	UseAttr            string                                  `xml:"Use,attr,omitempty"`
	CustomProperties   []*CustomPropertiesCollectionType       `xml:"CustomProperties"`
	ConsumeInfo        []*ConsumeTagInfoType                   `xml:"ConsumeInfo"`
	ProduceInfo        []*ProduceTagInfoType                   `xml:"ProduceInfo"`
	TagConfiguration   []*TagConfigurationAdaptorType          `xml:"TagConfiguration"`
	AlarmConditions    []*ConfiguredAlarmCollectionAdaptorType `xml:"AlarmConditions"`
	Description        []*DescriptionType                      `xml:"Description"`
	BaseDescription    []string                                `xml:"BaseDescription"`
	Comments           []*CommentCollectionType                `xml:"Comments"`
	EngineeringUnits   []*EngineeringUnitCollectionType        `xml:"EngineeringUnits"`
	Maxes              []*MaxLimitCollectionType               `xml:"Maxes"`
	Mins               []*MinLimitCollectionType               `xml:"Mins"`
	State0s            []*State0CollectionType                 `xml:"State0s"`
	State1s            []*State1CollectionType                 `xml:"State1s"`
	Librarys           []*LibraryAdaptorCollectionType         `xml:"Librarys"`
	Labels             []*LabelsWideType                       `xml:"Labels"`
	Instructions       []*InstructionAdaptorCollectionType     `xml:"Instructions"`
	Areas              []*AreaAdaptorCollectionType            `xml:"Areas"`
	URLs               []*URLAdaptorCollectionType             `xml:"URLs"`
	Navigations        []*NavigationCollectionType             `xml:"Navigations"`
	Data               []*DataWideType                         `xml:"Data"`
	ForceData          []*ForceDataWideType                    `xml:"ForceData"`
}

func (t *TagType) Descriptions() string {
	comment := ""
	for d := range t.Description {
		comment = comment + t.Description[d].CData()
	}

	return comment
}

// CommentCollectionType ...
type CommentCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Comment          []*CommentWideType              `xml:"Comment"`
}

// CommentAdaptorType ...
type CommentAdaptorType struct {
	OperandAttr      string                          `xml:"Operand,attr"`
	UnusedAttr       string                          `xml:"Unused,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedComment []*LocalizedCommentWideType     `xml:"LocalizedComment"`
}

// CommentAdaptorTextType ...
type CommentAdaptorTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// EngineeringUnitCollectionType ...
type EngineeringUnitCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	EngineeringUnit  []*EngineeringUnitType          `xml:"EngineeringUnit"`
}

// EngineeringUnitType ...
type EngineeringUnitType struct {
	OperandAttr              string                          `xml:"Operand,attr,omitempty"`
	UseAttr                  string                          `xml:"Use,attr,omitempty"`
	CustomProperties         *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedEngineeringUnit []*EngineeringUnitTextType      `xml:"LocalizedEngineeringUnit"`
}

// EngineeringUnitTextType ...
type EngineeringUnitTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// MaxLimitCollectionType ...
type MaxLimitCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Max              []*MaxLimitType                 `xml:"Max"`
}

// MaxLimitType ...
type MaxLimitType struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// MinLimitCollectionType ...
type MinLimitCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Min              []*MinLimitType                 `xml:"Min"`
}

// MinLimitType ...
type MinLimitType struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// State0CollectionType ...
type State0CollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	State0           []*State0Type                   `xml:"State0"`
}

// State0Type ...
type State0Type struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedState0  []*State0TextType               `xml:"LocalizedState0"`
}

// State0TextType ...
type State0TextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// State1CollectionType ...
type State1CollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	State1           []*State1Type                   `xml:"State1"`
}

// State1Type ...
type State1Type struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedState1  []*State1TextType               `xml:"LocalizedState1"`
}

// State1TextType ...
type State1TextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// CustomPropertiesCollectionType ...
type CustomPropertiesCollectionType struct {
	UseAttr  string                         `xml:"Use,attr,omitempty"`
	Provider []*CustomPropertiesAdaptorType `xml:"Provider"`
}

// CustomPropertiesAdaptorType ...
type CustomPropertiesAdaptorType struct {
}

// CustomPropertiesTextType ...
type CustomPropertiesTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// PrimaryActionSetCollectionType ...
type PrimaryActionSetCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	PrimaryActionSet []*PrimaryActionSetAdaptorType  `xml:"PrimaryActionSet"`
}

// PrimaryActionSetAdaptorType ...
type PrimaryActionSetAdaptorType struct {
	PermissionSetAttr         string                          `xml:"PermissionSet,attr"`
	IsPermissionSetAttr       string                          `xml:"IsPermissionSet,attr"`
	UseAttr                   string                          `xml:"Use,attr,omitempty"`
	CustomProperties          *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedPrimaryActionSet []*PrimaryActionSetTextType     `xml:"LocalizedPrimaryActionSet"`
}

// PrimaryActionSetTextType ...
type PrimaryActionSetTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// LibraryAdaptorCollectionType ...
type LibraryAdaptorCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Library          []*LibraryAdaptorType           `xml:"Library"`
}

// LibraryAdaptorType ...
type LibraryAdaptorType struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// LabelAdaptorCollectionType ...
type LabelAdaptorCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Label            []*LabelWideType                `xml:"Label"`
}

// LabelAdaptorType ...
type LabelAdaptorType struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedLabel   []*LabelTextType                `xml:"LocalizedLabel"`
}

// LabelTextType ...
type LabelTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// InstructionAdaptorCollectionType ...
type InstructionAdaptorCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Instruction      []*InstructionAdaptorType       `xml:"Instruction"`
}

// InstructionAdaptorType ...
type InstructionAdaptorType struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// AreaAdaptorCollectionType ...
type AreaAdaptorCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Area             []*AreaAdaptorType              `xml:"Area"`
}

// AreaAdaptorType ...
type AreaAdaptorType struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// URLAdaptorCollectionType ...
type URLAdaptorCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	URL              []*URLAdaptorType               `xml:"URL"`
}

// URLAdaptorType ...
type URLAdaptorType struct {
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// NavigationCollectionType ...
type NavigationCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Navigation       []*NavigationAdaptorType        `xml:"Navigation"`
}

// NavigationAdaptorType ...
type NavigationAdaptorType struct {
	OperandAttr      string                          `xml:"Operand,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ConsumeTagInfoType ...
type ConsumeTagInfoType struct {
	ProducerAttr                string                          `xml:"Producer,attr"`
	RemoteTagAttr               string                          `xml:"RemoteTag,attr,omitempty"`
	RemoteInstanceAttr          string                          `xml:"RemoteInstance,attr,omitempty"`
	RPIAttr                     string                          `xml:"RPI,attr"`
	TimeoutMultiplierAttr       string                          `xml:"TimeoutMultiplier,attr,omitempty"`
	NetworkDelayMultiplierAttr  string                          `xml:"NetworkDelayMultiplier,attr,omitempty"`
	ReactionTimeLimitAttr       string                          `xml:"ReactionTimeLimit,attr,omitempty"`
	MaxObservedNetworkDelayAttr string                          `xml:"MaxObservedNetworkDelay,attr,omitempty"`
	UnicastAttr                 string                          `xml:"Unicast,attr,omitempty"`
	UseAttr                     string                          `xml:"Use,attr,omitempty"`
	CustomProperties            *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ProduceTagInfoType ...
type ProduceTagInfoType struct {
	ProduceCountAttr                     string                          `xml:"ProduceCount,attr,omitempty"`
	PLCMappingFileAttr                   string                          `xml:"PLCMappingFile,attr,omitempty"`
	PLC2MappingAttr                      string                          `xml:"PLC2Mapping,attr,omitempty"`
	ProgrammaticallySendEventTriggerAttr string                          `xml:"ProgrammaticallySendEventTrigger,attr,omitempty"`
	UnicastPermittedAttr                 string                          `xml:"UnicastPermitted,attr,omitempty"`
	MinimumRPIAttr                       string                          `xml:"MinimumRPI,attr,omitempty"`
	MaximumRPIAttr                       string                          `xml:"MaximumRPI,attr,omitempty"`
	DefaultRPIAttr                       string                          `xml:"DefaultRPI,attr,omitempty"`
	UseAttr                              string                          `xml:"Use,attr,omitempty"`
	CustomProperties                     *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ConfiguredAlarmCollectionAdaptorType ...
type ConfiguredAlarmCollectionAdaptorType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	AlarmCondition   []*ConfiguredAlarmAdaptorType   `xml:"AlarmCondition"`
}

// ConfiguredAlarmAdaptorType ...
type ConfiguredAlarmAdaptorType struct {
	InFaultAttr                  string                          `xml:"InFault,attr,omitempty"`
	AckRequiredAttr              string                          `xml:"AckRequired,attr,omitempty"`
	LatchedAttr                  string                          `xml:"Latched,attr,omitempty"`
	ProgAckAttr                  string                          `xml:"ProgAck,attr,omitempty"`
	OperAckAttr                  string                          `xml:"OperAck,attr,omitempty"`
	ProgResetAttr                string                          `xml:"ProgReset,attr,omitempty"`
	OperResetAttr                string                          `xml:"OperReset,attr,omitempty"`
	ProgSuppressAttr             string                          `xml:"ProgSuppress,attr,omitempty"`
	OperSuppressAttr             string                          `xml:"OperSuppress,attr,omitempty"`
	ProgUnsuppressAttr           string                          `xml:"ProgUnsuppress,attr,omitempty"`
	OperUnsuppressAttr           string                          `xml:"OperUnsuppress,attr,omitempty"`
	OperShelveAttr               string                          `xml:"OperShelve,attr,omitempty"`
	ProgUnshelveAttr             string                          `xml:"ProgUnshelve,attr,omitempty"`
	OperUnshelveAttr             string                          `xml:"OperUnshelve,attr,omitempty"`
	ProgDisableAttr              string                          `xml:"ProgDisable,attr,omitempty"`
	OperDisableAttr              string                          `xml:"OperDisable,attr,omitempty"`
	ProgEnableAttr               string                          `xml:"ProgEnable,attr,omitempty"`
	OperEnableAttr               string                          `xml:"OperEnable,attr,omitempty"`
	AlarmCountResetAttr          string                          `xml:"AlarmCountReset,attr,omitempty"`
	UsedAttr                     string                          `xml:"Used,attr,omitempty"`
	AlarmSetOperIncludedAttr     string                          `xml:"AlarmSetOperIncluded,attr,omitempty"`
	AlarmSetRollupIncludedAttr   string                          `xml:"AlarmSetRollupIncluded,attr,omitempty"`
	NameAttr                     string                          `xml:"Name,attr"`
	InputAttr                    string                          `xml:"Input,attr"`
	ConditionTypeAttr            string                          `xml:"ConditionType,attr,omitempty"`
	AlarmConditionDefinitionAttr string                          `xml:"AlarmConditionDefinition,attr,omitempty"`
	LimitAttr                    string                          `xml:"Limit,attr,omitempty"`
	SeverityAttr                 string                          `xml:"Severity,attr,omitempty"`
	OnDelayAttr                  string                          `xml:"OnDelay,attr,omitempty"`
	OffDelayAttr                 string                          `xml:"OffDelay,attr,omitempty"`
	ShelveDurationAttr           string                          `xml:"ShelveDuration,attr,omitempty"`
	MaxShelveDurationAttr        string                          `xml:"MaxShelveDuration,attr,omitempty"`
	DeadbandAttr                 string                          `xml:"Deadband,attr,omitempty"`
	EvaluationPeriodAttr         string                          `xml:"EvaluationPeriod,attr,omitempty"`
	ExpressionAttr               string                          `xml:"Expression,attr,omitempty"`
	TargetTagAttr                string                          `xml:"TargetTag,attr,omitempty"`
	AssocTag1Attr                string                          `xml:"AssocTag1,attr,omitempty"`
	AssocTag2Attr                string                          `xml:"AssocTag2,attr,omitempty"`
	AssocTag3Attr                string                          `xml:"AssocTag3,attr,omitempty"`
	AssocTag4Attr                string                          `xml:"AssocTag4,attr,omitempty"`
	UIdAttr                      string                          `xml:"UId,attr,omitempty"`
	UseAttr                      string                          `xml:"Use,attr,omitempty"`
	CustomProperties             *CustomPropertiesCollectionType `xml:"CustomProperties"`
	AlarmConfig                  *AlarmConfigType                `xml:"AlarmConfig"`
}

// StringDataType ...
type StringDataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	LengthAttr       string                          `xml:"Length,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// RawDataType ...
type RawDataType struct {
	FormatAttr       string                          `xml:"Format,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// L5KDataType ...
type L5KDataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// DecoratedDataType ...
type DecoratedDataType struct {
}

// RawForceDataType ...
type RawForceDataType struct {
	FormatAttr       string                          `xml:"Format,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// L5KForceDataType ...
type L5KForceDataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// MessageDataType ...
type MessageDataType struct {
	FormatAttr        string                          `xml:"Format,attr"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
	MessageParameters *MsgType                        `xml:"MessageParameters"`
}

// MotionGroupDataType ...
type MotionGroupDataType struct {
	FormatAttr            string                          `xml:"Format,attr"`
	UseAttr               string                          `xml:"Use,attr,omitempty"`
	CustomProperties      *CustomPropertiesCollectionType `xml:"CustomProperties"`
	MotionGroupParameters *MotionGroupType                `xml:"MotionGroupParameters"`
}

// CoordinateSystemDataType ...
type CoordinateSystemDataType struct {
	FormatAttr                 string                          `xml:"Format,attr"`
	UseAttr                    string                          `xml:"Use,attr,omitempty"`
	CustomProperties           *CustomPropertiesCollectionType `xml:"CustomProperties"`
	CoordinateSystemParameters *CoordinateSystemType           `xml:"CoordinateSystemParameters"`
}

// AxisDataType ...
type AxisDataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	AxisParameters   *AxisType                       `xml:"AxisParameters"`
}

// AlarmDataType ...
type AlarmDataType struct {
	FormatAttr             string                          `xml:"Format,attr"`
	UseAttr                string                          `xml:"Use,attr,omitempty"`
	CustomProperties       *CustomPropertiesCollectionType `xml:"CustomProperties"`
	AlarmAnalogParameters  *AlarmAnalogType                `xml:"AlarmAnalogParameters"`
	AlarmDigitalParameters *AlarmDigitalType               `xml:"AlarmDigitalParameters"`
	AlarmConfig            *AlarmConfigType                `xml:"AlarmConfig"`
}

// HMIBCDataType ...
type HMIBCDataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	HMIBCParameters  *HMIBCType                      `xml:"HMIBCParameters"`
}

// BEODataType ...
type BEODataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	BEOParameters    *BEOType                        `xml:"BEOParameters"`
}

// EEODataType ...
type EEODataType struct {
	FormatAttr       string                          `xml:"Format,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	EEOParameters    *EEOType                        `xml:"EEOParameters"`
}

// AlarmConfigType ...
type AlarmConfigType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Messages         *AlarmMessageCollectionType     `xml:"Messages"`
	AlarmClass       string                          `xml:"AlarmClass"`
	HMICmd           string                          `xml:"HMICmd"`
	HMIGroup         string                          `xml:"HMIGroup"`
}

// AlarmMessageCollectionType ...
type AlarmMessageCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Message          []*AlarmMessageType             `xml:"Message"`
}

// AlarmMessageType ...
type AlarmMessageType struct {
	TypeAttr         string                          `xml:"Type,attr"`
	IDAttr           string                          `xml:"ID,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Text             []*TextWideType                 `xml:"Text"`
}

// AlarmMessageTextType ...
type AlarmMessageTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ProgramCollectionType ...
type ProgramCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Program          []*AProgramType                 `xml:"Program"`
}

// AProgramType ...
type AProgramType struct {
	NameAttr                                    string                          `xml:"Name,attr"`
	UIdAttr                                     string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr                               string                          `xml:"ParentUId,attr,omitempty"`
	TypeAttr                                    string                          `xml:"Type,attr,omitempty"`
	TestEditsAttr                               string                          `xml:"TestEdits,attr,omitempty"`
	MainRoutineNameAttr                         string                          `xml:"MainRoutineName,attr,omitempty"`
	PreStateRoutineNameAttr                     string                          `xml:"PreStateRoutineName,attr,omitempty"`
	FaultRoutineNameAttr                        string                          `xml:"FaultRoutineName,attr,omitempty"`
	ExecutingTaskNameAttr                       string                          `xml:"ExecutingTaskName,attr,omitempty"`
	VerifiedAttr                                string                          `xml:"Verified,attr,omitempty"`
	EditsExistAttr                              string                          `xml:"EditsExist,attr,omitempty"`
	DisabledAttr                                string                          `xml:"Disabled,attr,omitempty"`
	InitialStepIndexAttr                        string                          `xml:"InitialStepIndex,attr,omitempty"`
	InitialStateAttr                            string                          `xml:"InitialState,attr,omitempty"`
	CompleteStateIfNotImplAttr                  string                          `xml:"CompleteStateIfNotImpl,attr,omitempty"`
	LossOfCommCmdAttr                           string                          `xml:"LossOfCommCmd,attr,omitempty"`
	ExternalRequestActionAttr                   string                          `xml:"ExternalRequestAction,attr,omitempty"`
	EquipmentIdAttr                             string                          `xml:"EquipmentId,attr,omitempty"`
	RecipePhaseNamesAttr                        string                          `xml:"RecipePhaseNames,attr,omitempty"`
	LastScanTimeAttr                            string                          `xml:"LastScanTime,attr,omitempty"`
	MaxScanTimeAttr                             string                          `xml:"MaxScanTime,attr,omitempty"`
	TagsUIdAttr                                 string                          `xml:"TagsUId,attr,omitempty"`
	RoutinesUIdAttr                             string                          `xml:"RoutinesUId,attr,omitempty"`
	ClassAttr                                   string                          `xml:"Class,attr,omitempty"`
	SynchronizeRedundancyDataAfterExecutionAttr string                          `xml:"SynchronizeRedundancyDataAfterExecution,attr,omitempty"`
	UseAsFolderAttr                             string                          `xml:"UseAsFolder,attr,omitempty"`
	AutoValueAssignStepToPhaseAttr              string                          `xml:"AutoValueAssignStepToPhase,attr,omitempty"`
	AutoValueAssignPhaseToStepOnCompleteAttr    string                          `xml:"AutoValueAssignPhaseToStepOnComplete,attr,omitempty"`
	AutoValueAssignPhaseToStepOnStoppedAttr     string                          `xml:"AutoValueAssignPhaseToStepOnStopped,attr,omitempty"`
	AutoValueAssignPhaseToStepOnAbortedAttr     string                          `xml:"AutoValueAssignPhaseToStepOnAborted,attr,omitempty"`
	RevisionAttr                                string                          `xml:"Revision,attr,omitempty"`
	RevisionExtensionAttr                       string                          `xml:"RevisionExtension,attr,omitempty"`
	UnitIDAttr                                  string                          `xml:"UnitID,attr,omitempty"`
	RetainSequenceIDOnResetAttr                 string                          `xml:"RetainSequenceIDOnReset,attr,omitempty"`
	GenerateSequenceEventsAttr                  string                          `xml:"GenerateSequenceEvents,attr,omitempty"`
	ValuesToUseOnStartAttr                      string                          `xml:"ValuesToUseOnStart,attr,omitempty"`
	ValuesToUseOnResetAttr                      string                          `xml:"ValuesToUseOnReset,attr,omitempty"`
	UseAttr                                     string                          `xml:"Use,attr,omitempty"`
	CustomProperties                            *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description                                 *DescriptionType                `xml:"Description"`
	RevisionNote                                *RevisionNoteType               `xml:"RevisionNote"`
	Tags                                        *TagCollectionType              `xml:"Tags"`
	Parameters                                  *UDIParameterCollectionType     `xml:"Parameters"`
	LocalTags                                   *UDILocalTagCollectionType      `xml:"LocalTags"`
	Routines                                    *RoutineCollectionType          `xml:"Routines"`
	ChildPrograms                               *ChildProgramCollectionType     `xml:"ChildPrograms"`
}

// RoutineCollectionType ...
type RoutineCollectionType struct {
	UIdAttr          string                            `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                            `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                            `xml:"Use,attr,omitempty"`
	EncodedData      []*EncodedRoutineType             `xml:"EncodedData"`
	CustomProperties []*CustomPropertiesCollectionType `xml:"CustomProperties"`
	Routine          []*RoutineType                    `xml:"Routine"`
}

// EncryptedSegments ...
type EncryptedSegments struct {
	Value string `xml:",chardata"`
}

// RoutineType ...
type RoutineType struct {
	NameAttr                 string                          `xml:"Name,attr"`
	TypeAttr                 string                          `xml:"Type,attr"`
	SourceKeyAttr            string                          `xml:"SourceKey,attr,omitempty"`
	EncodedSourceKeyAttr     string                          `xml:"EncodedSourceKey,attr,omitempty"`
	SourceProtectionTypeAttr string                          `xml:"SourceProtectionType,attr,omitempty"`
	UIdAttr                  string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr            string                          `xml:"ParentUId,attr,omitempty"`
	PermissionSetAttr        string                          `xml:"PermissionSet,attr,omitempty"`
	VerifiedAttr             string                          `xml:"Verified,attr,omitempty"`
	TrackingGroupsAttr       string                          `xml:"TrackingGroups,attr,omitempty"`
	LogicHashAttr            string                          `xml:"LogicHash,attr,omitempty"`
	DescriptionHashAttr      string                          `xml:"DescriptionHash,attr,omitempty"`
	UseAttr                  string                          `xml:"Use,attr,omitempty"`
	CustomProperties         *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description              *DescriptionType                `xml:"Description"`
	EncryptionInfo           *EncryptionInfoType             `xml:"EncryptionInfo"`
	RLLContent               []*RungCollectionType           `xml:"RLLContent"`
	FBDContent               []*FBDContentType               `xml:"FBDContent"`
	SFCContent               []*SFCContentType               `xml:"SFCContent"`
	STContent                []*STContentType                `xml:"STContent"`
	ExternalContent          []*ExternalContentType          `xml:"ExternalContent"`
	SEQContent               []*SEQContentType               `xml:"SEQContent"`
	EncryptedContent         []*EncryptedContentType         `xml:"EncryptedContent"`
	EncryptedSegments        *EncryptedSegments              `xml:"EncryptedSegments"`
}

// EncodedRoutineType ...
type EncodedRoutineType struct {
	EncodedTypeAttr      string                          `xml:"EncodedType,attr,omitempty"`
	NameAttr             string                          `xml:"Name,attr,omitempty"`
	UIdAttr              string                          `xml:"UId,attr,omitempty"`
	TypeAttr             string                          `xml:"Type,attr,omitempty"`
	PermissionSetAttr    string                          `xml:"PermissionSet,attr,omitempty"`
	EncryptionConfigAttr string                          `xml:"EncryptionConfig,attr,omitempty"`
	LogicHashAttr        string                          `xml:"LogicHash,attr,omitempty"`
	DescriptionHashAttr  string                          `xml:"DescriptionHash,attr,omitempty"`
	UseAttr              string                          `xml:"Use,attr,omitempty"`
	CustomProperties     *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description          *DescriptionType                `xml:"Description"`
}

// RungCollectionType ...
type RungCollectionType struct {
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	StartAttr          string                          `xml:"Start,attr,omitempty"`
	CountAttr          string                          `xml:"Count,attr,omitempty"`
	OnlineEditTypeAttr string                          `xml:"OnlineEditType,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Labels             []*LabelsWideType               `xml:"Labels"`
	Rung               []*RungType                     `xml:"Rung"`
}

// RungType ...
type RungType struct {
	NumberAttr       string                            `xml:"Number,attr,omitempty"`
	TypeAttr         string                            `xml:"Type,attr,omitempty"`
	UIdAttr          string                            `xml:"UId,attr,omitempty"`
	RungIdAttr       string                            `xml:"RungId,attr,omitempty"`
	VerifiedAttr     string                            `xml:"Verified,attr,omitempty"`
	RegionIdAttr     string                            `xml:"RegionId,attr,omitempty"`
	UseAttr          string                            `xml:"Use,attr,omitempty"`
	CustomProperties []*CustomPropertiesCollectionType `xml:"CustomProperties"`
	Comment          []*CommentWideType                `xml:"Comment"`
	Text             []*TextWideType                   `xml:"Text"`
	TextIOI          []string                          `xml:"TextIOI"`
	TextNTT          []string                          `xml:"TextNTT"`
}

// RungCommentType ...
type RungCommentType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedComment []*LocalizedCommentWideType     `xml:"LocalizedComment"`
}

// RungCommentTextType ...
type RungCommentTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// LabelCollectionType ...
type LabelCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Label            []*LabelWideType                `xml:"Label"`
}

// LabelType ...
type LabelType struct {
	NameAttr         string                          `xml:"Name,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDContentType ...
type FBDContentType struct {
	UseAttr              string                          `xml:"Use,attr,omitempty"`
	StartAttr            string                          `xml:"Start,attr,omitempty"`
	CountAttr            string                          `xml:"Count,attr,omitempty"`
	SheetSizeAttr        string                          `xml:"SheetSize,attr,omitempty"`
	SheetOrientationAttr string                          `xml:"SheetOrientation,attr,omitempty"`
	OnlineEditTypeAttr   string                          `xml:"OnlineEditType,attr,omitempty"`
	UIdAttr              string                          `xml:"UId,attr,omitempty"`
	CustomProperties     *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Sheet                []*SheetType                    `xml:"Sheet"`
}

// SheetType ...
type SheetType struct {
	NumberAttr       string                          `xml:"Number,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description      *DescriptionType                `xml:"Description"`
	LanguageElements []string                        `xml:"LanguageElements"`
	IRef             []*FBDInputRefType              `xml:"IRef"`
	ORef             []*FBDOutputRefType             `xml:"ORef"`
	ICon             []*FBDInputWireConnectorType    `xml:"ICon"`
	OCon             []*FBDOutputWireConnectorType   `xml:"OCon"`
	Block            []*AFBDBlockType                `xml:"Block"`
	Function         []*AFBDFunctionType             `xml:"Function"`
	AddOnInstruction []*FBDUDIBlockType              `xml:"AddOnInstruction"`
	GSV              []*FBDGsvBlockType              `xml:"GSV"`
	SSV              []*FBDSsvBlockType              `xml:"SSV"`
	JSR              []*FBDJSRType                   `xml:"JSR"`
	SBR              []*FBDSBRType                   `xml:"SBR"`
	RET              []*FBDRETType                   `xml:"RET"`
	Wire             []*FBDWireType                  `xml:"Wire"`
	FeedbackWire     []*FBDFeedbackWireType          `xml:"FeedbackWire"`
	TextBox          []*TextBoxType                  `xml:"TextBox"`
	Attachment       []*AAttachmentType              `xml:"Attachment"`
}

// FBDInputRefType ...
type FBDInputRefType struct {
	XMLName              xml.Name                        `xml:"IRef"`
	IDAttr               string                          `xml:"ID,attr"`
	XAttr                string                          `xml:"X,attr"`
	YAttr                string                          `xml:"Y,attr"`
	OperandAttr          string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr       string                          `xml:"OperandIOI,attr,omitempty"`
	OperandIOINTTAttr    string                          `xml:"OperandIOINTT,attr,omitempty"`
	HideDescAttr         string                          `xml:"HideDesc,attr,omitempty"`
	VerifiedAttr         string                          `xml:"Verified,attr,omitempty"`
	LatchedTagIOIAttr    string                          `xml:"LatchedTagIOI,attr,omitempty"`
	LatchedTagIOINTTAttr string                          `xml:"LatchedTagIOINTT,attr,omitempty"`
	UIdAttr              string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr       string                          `xml:"MetadataId,attr,omitempty"`
	RegionIdAttr         string                          `xml:"RegionId,attr,omitempty"`
	UseAttr              string                          `xml:"Use,attr,omitempty"`
	CustomProperties     *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDOutputRefType ...
type FBDOutputRefType struct {
	XMLName           xml.Name                        `xml:"ORef"`
	IDAttr            string                          `xml:"ID,attr"`
	XAttr             string                          `xml:"X,attr"`
	YAttr             string                          `xml:"Y,attr"`
	OperandAttr       string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr    string                          `xml:"OperandIOI,attr,omitempty"`
	OperandIOINTTAttr string                          `xml:"OperandIOINTT,attr,omitempty"`
	HideDescAttr      string                          `xml:"HideDesc,attr,omitempty"`
	VerifiedAttr      string                          `xml:"Verified,attr,omitempty"`
	RegionIdAttr      string                          `xml:"RegionId,attr,omitempty"`
	UIdAttr           string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr    string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDInputWireConnectorType ...
type FBDInputWireConnectorType struct {
	XMLName          xml.Name                        `xml:"ICon"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	NameAttr         string                          `xml:"Name,attr"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDOutputWireConnectorType ...
type FBDOutputWireConnectorType struct {
	XMLName          xml.Name                        `xml:"OCon"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	NameAttr         string                          `xml:"Name,attr"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDUDIBlockType ...
type FBDUDIBlockType struct {
	XMLName           xml.Name                        `xml:"AddOnInstruction"`
	NameAttr          string                          `xml:"Name,attr"`
	IDAttr            string                          `xml:"ID,attr"`
	XAttr             string                          `xml:"X,attr"`
	YAttr             string                          `xml:"Y,attr"`
	OperandAttr       string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr    string                          `xml:"OperandIOI,attr,omitempty"`
	OperandIOINTTAttr string                          `xml:"OperandIOINTT,attr,omitempty"`
	VisiblePinsAttr   string                          `xml:"VisiblePins,attr,omitempty"`
	VerifiedAttr      string                          `xml:"Verified,attr,omitempty"`
	RegionIdAttr      string                          `xml:"RegionId,attr,omitempty"`
	UIdAttr           string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr    string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
	InOutParameter    []*FBDUDIArgumentType           `xml:"InOutParameter"`
}

// FBDUDIArgumentType ...
type FBDUDIArgumentType struct {
	XMLName            xml.Name                        `xml:"FBD_UDIArgumentType"`
	NameAttr           string                          `xml:"Name,attr"`
	ArgumentAttr       string                          `xml:"Argument,attr,omitempty"`
	ArgumentIOIAttr    string                          `xml:"ArgumentIOI,attr,omitempty"`
	ArgumentIOINTTAttr string                          `xml:"ArgumentIOINTT,attr,omitempty"`
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// AFBDBlockType ...
type AFBDBlockType struct {
	XMLName               xml.Name                        `xml:"Block"`
	TypeAttr              string                          `xml:"Type,attr"`
	IDAttr                string                          `xml:"ID,attr"`
	XAttr                 string                          `xml:"X,attr"`
	YAttr                 string                          `xml:"Y,attr"`
	OperandAttr           string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr        string                          `xml:"OperandIOI,attr,omitempty"`
	OperandIOINTTAttr     string                          `xml:"OperandIOINTT,attr,omitempty"`
	VisiblePinsAttr       string                          `xml:"VisiblePins,attr,omitempty"`
	HideDescAttr          string                          `xml:"HideDesc,attr,omitempty"`
	AutotuneTagAttr       string                          `xml:"AutotuneTag,attr,omitempty"`
	AutotuneTagIOIAttr    string                          `xml:"AutotuneTagIOI,attr,omitempty"`
	AutotuneTagIOINTTAttr string                          `xml:"AutotuneTagIOINTT,attr,omitempty"`
	VerifiedAttr          string                          `xml:"Verified,attr,omitempty"`
	UIdAttr               string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr        string                          `xml:"MetadataId,attr,omitempty"`
	RegionIdAttr          string                          `xml:"RegionId,attr,omitempty"`
	UseAttr               string                          `xml:"Use,attr,omitempty"`
	CustomProperties      *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Array                 []*FBDSpecialArrayType          `xml:"Array"`
}

// AFBDFunctionType ...
type AFBDFunctionType struct {
	XMLName                 xml.Name                        `xml:"Function"`
	TypeAttr                string                          `xml:"Type,attr"`
	IDAttr                  string                          `xml:"ID,attr"`
	XAttr                   string                          `xml:"X,attr"`
	YAttr                   string                          `xml:"Y,attr"`
	OutputTagIOIAttr        string                          `xml:"OutputTagIOI,attr,omitempty"`
	OutputTagIOINTTAttr     string                          `xml:"OutputTagIOINTT,attr,omitempty"`
	RegionIdAttr            string                          `xml:"RegionId,attr,omitempty"`
	OutputTagArrayIndexAttr string                          `xml:"OutputTagArrayIndex,attr,omitempty"`
	UIdAttr                 string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr          string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr                 string                          `xml:"Use,attr,omitempty"`
	CustomProperties        *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDSpecialArrayType ...
type FBDSpecialArrayType struct {
	XMLName           xml.Name                        `xml:"Array"`
	NameAttr          string                          `xml:"Name,attr"`
	OperandAttr       string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr    string                          `xml:"OperandIOI,attr,omitempty"`
	OperandIOINTTAttr string                          `xml:"OperandIOINTT,attr,omitempty"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDGsvBlockType ...
type FBDGsvBlockType struct {
	XMLName          xml.Name                        `xml:"FBD_GsvBlockType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	ObjectAttr       string                          `xml:"Object,attr"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDSsvBlockType ...
type FBDSsvBlockType struct {
	XMLName          xml.Name                        `xml:"FBD_SsvBlockType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	ObjectAttr       string                          `xml:"Object,attr"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDJSRType ...
type FBDJSRType struct {
	XMLName           xml.Name                        `xml:"FBD_JSRType"`
	IDAttr            string                          `xml:"ID,attr"`
	XAttr             string                          `xml:"X,attr"`
	YAttr             string                          `xml:"Y,attr"`
	RoutineAttr       string                          `xml:"Routine,attr,omitempty"`
	RoutineIOIAttr    string                          `xml:"RoutineIOI,attr,omitempty"`
	RoutineIOINTTAttr string                          `xml:"RoutineIOINTT,attr,omitempty"`
	InAttr            string                          `xml:"In,attr,omitempty"`
	InIOIAttr         string                          `xml:"InIOI,attr,omitempty"`
	InIOINTTAttr      string                          `xml:"InIOINTT,attr,omitempty"`
	RetAttr           string                          `xml:"Ret,attr,omitempty"`
	RetIOIAttr        string                          `xml:"RetIOI,attr,omitempty"`
	RetIOINTTAttr     string                          `xml:"RetIOINTT,attr,omitempty"`
	UIdAttr           string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr    string                          `xml:"MetadataId,attr,omitempty"`
	VerifiedAttr      string                          `xml:"Verified,attr,omitempty"`
	RegionIdAttr      string                          `xml:"RegionId,attr,omitempty"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDSBRType ...
type FBDSBRType struct {
	XMLName          xml.Name                        `xml:"FBD_SBRType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	RetAttr          string                          `xml:"Ret,attr,omitempty"`
	RetIOIAttr       string                          `xml:"RetIOI,attr,omitempty"`
	RetIOINTTAttr    string                          `xml:"RetIOINTT,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	RegionIdAttr     string                          `xml:"RegionId,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDRETType ...
type FBDRETType struct {
	XMLName          xml.Name                        `xml:"FBD_RETType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	InAttr           string                          `xml:"In,attr,omitempty"`
	InIOIAttr        string                          `xml:"InIOI,attr,omitempty"`
	InIOINTTAttr     string                          `xml:"InIOINTT,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	RegionIdAttr     string                          `xml:"RegionId,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDWireType ...
type FBDWireType struct {
	XMLName          xml.Name                        `xml:"Wire"`
	FromIDAttr       string                          `xml:"FromID,attr"`
	FromParamAttr    string                          `xml:"FromParam,attr,omitempty"`
	ToIDAttr         string                          `xml:"ToID,attr"`
	ToParamAttr      string                          `xml:"ToParam,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FBDFeedbackWireType ...
type FBDFeedbackWireType struct {
	XMLName          xml.Name                        `xml:"FeedbackWire"`
	FromIDAttr       string                          `xml:"FromID,attr"`
	FromParamAttr    string                          `xml:"FromParam,attr,omitempty"`
	ToIDAttr         string                          `xml:"ToID,attr"`
	ToParamAttr      string                          `xml:"ToParam,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// TextBoxType ...
type TextBoxType struct {
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	WidthAttr        string                          `xml:"Width,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Text             []*TextWideType                 `xml:"Text"`
}

// TextBoxAdaptorType ...
type TextBoxAdaptorType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LocalizedText    []*TextBoxAdaptorTextType       `xml:"LocalizedText"`
}

// TextBoxAdaptorTextType ...
type TextBoxAdaptorTextType struct {
	LangAttr         string                          `xml:"Lang,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// SFCContentType ...
type SFCContentType struct {
	SheetSizeAttr        string                          `xml:"SheetSize,attr,omitempty"`
	SheetOrientationAttr string                          `xml:"SheetOrientation,attr,omitempty"`
	StepNameAttr         string                          `xml:"StepName,attr,omitempty"`
	TransitionNameAttr   string                          `xml:"TransitionName,attr,omitempty"`
	ActionNameAttr       string                          `xml:"ActionName,attr,omitempty"`
	StopNameAttr         string                          `xml:"StopName,attr,omitempty"`
	OnlineEditTypeAttr   string                          `xml:"OnlineEditType,attr,omitempty"`
	UIdAttr              string                          `xml:"UId,attr,omitempty"`
	UseAttr              string                          `xml:"Use,attr,omitempty"`
	CustomProperties     *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LanguageElements     []string                        `xml:"LanguageElements"`
	Step                 []*ASFCStepElementType          `xml:"Step"`
	Transition           []*ASFCTransitionElementType    `xml:"Transition"`
	Branch               []*ASFCBranchElementType        `xml:"Branch"`
	SbrRet               []*SFCSBRRETType                `xml:"SbrRet"`
	Stop                 []*ASFCStopElementType          `xml:"Stop"`
	DirectedLink         []*ASFCLangElemWireType         `xml:"DirectedLink"`
	TextBox              []*TextBoxType                  `xml:"TextBox"`
	Attachment           []*AAttachmentType              `xml:"Attachment"`
}

// SEQContentType ...
type SEQContentType struct {
	SheetSizeAttr        string                          `xml:"SheetSize,attr,omitempty"`
	SheetOrientationAttr string                          `xml:"SheetOrientation,attr,omitempty"`
	UIdAttr              string                          `xml:"UId,attr,omitempty"`
	UseAttr              string                          `xml:"Use,attr,omitempty"`
	CustomProperties     *CustomPropertiesCollectionType `xml:"CustomProperties"`
	LanguageElements     []string                        `xml:"LanguageElements"`
	Step                 []*ASFCStepElementType          `xml:"Step"`
	Transition           []*ASFCTransitionElementType    `xml:"Transition"`
	Branch               []*ASFCBranchElementType        `xml:"Branch"`
	Stop                 []*ASFCStopElementType          `xml:"Stop"`
	DirectedLink         []*ASFCLangElemWireType         `xml:"DirectedLink"`
	TextBox              []*TextBoxType                  `xml:"TextBox"`
	Attachment           []*AAttachmentType              `xml:"Attachment"`
	TagConfigurations    *TagConfigurationsAdaptorType   `xml:"TagConfigurations"`
}

// TagConfigurationsAdaptorType ...
type TagConfigurationsAdaptorType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	TagConfiguration []*TagConfigurationAdaptorType  `xml:"TagConfiguration"`
}

// TagConfigurationAdaptorType ...
type TagConfigurationAdaptorType struct {
	TagNameAttr              string                          `xml:"TagName,attr,omitempty"`
	TagIOIAttr               string                          `xml:"TagIOI,attr,omitempty"`
	Action1FirstRegionIDAttr string                          `xml:"Action1FirstRegionID,attr,omitempty"`
	Action2FirstRegionIDAttr string                          `xml:"Action2FirstRegionID,attr,omitempty"`
	UseAttr                  string                          `xml:"Use,attr,omitempty"`
	CustomProperties         *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Expression               *EmbeddedLanguageBlockType      `xml:"Expression"`
}

// ASFCStepElementType ...
type ASFCStepElementType struct {
	XMLName                     xml.Name                        `xml:"ASFC_StepElementType"`
	IDAttr                      string                          `xml:"ID,attr"`
	XAttr                       string                          `xml:"X,attr"`
	YAttr                       string                          `xml:"Y,attr"`
	OperandAttr                 string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr              string                          `xml:"OperandIOI,attr,omitempty"`
	ComponentIOIAttr            string                          `xml:"ComponentIOI,attr,omitempty"`
	HideDescAttr                string                          `xml:"HideDesc,attr,omitempty"`
	DescXAttr                   string                          `xml:"DescX,attr,omitempty"`
	DescYAttr                   string                          `xml:"DescY,attr,omitempty"`
	DescWidthAttr               string                          `xml:"DescWidth,attr,omitempty"`
	InitialStepAttr             string                          `xml:"InitialStep,attr,omitempty"`
	PresetUsesExprAttr          string                          `xml:"PresetUsesExpr,attr,omitempty"`
	LimitHighUsesExprAttr       string                          `xml:"LimitHighUsesExpr,attr,omitempty"`
	LimitLowUsesExprAttr        string                          `xml:"LimitLowUsesExpr,attr,omitempty"`
	ShowActionsAttr             string                          `xml:"ShowActions,attr,omitempty"`
	NoPhaseStepAttr             string                          `xml:"NoPhaseStep,attr,omitempty"`
	PhaseNameAttr               string                          `xml:"PhaseName,attr,omitempty"`
	PhaseIOIAttr                string                          `xml:"PhaseIOI,attr,omitempty"`
	TransferOfControlSourceAttr string                          `xml:"TransferOfControlSource,attr,omitempty"`
	TransferOfControlTargetAttr string                          `xml:"TransferOfControlTarget,attr,omitempty"`
	UIdAttr                     string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr              string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr                     string                          `xml:"Use,attr,omitempty"`
	CustomProperties            *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Preset                      *EmbeddedLanguageBlockType      `xml:"Preset"`
	LimitHigh                   *EmbeddedLanguageBlockType      `xml:"LimitHigh"`
	LimitLow                    *EmbeddedLanguageBlockType      `xml:"LimitLow"`
	Action                      []*SFCActionElementType         `xml:"Action"`
}

// EmbeddedLanguageBlockType ...
type EmbeddedLanguageBlockType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	STContent        *STContentType                  `xml:"STContent"`
}

// SFCActionElementType ...
type SFCActionElementType struct {
	XMLName             xml.Name                        `xml:"SFC_ActionElementType"`
	IDAttr              string                          `xml:"ID,attr"`
	OperandAttr         string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr      string                          `xml:"OperandIOI,attr,omitempty"`
	ComponentIOIAttr    string                          `xml:"ComponentIOI,attr,omitempty"`
	QualifierAttr       string                          `xml:"Qualifier,attr,omitempty"`
	IsBooleanAttr       string                          `xml:"IsBoolean,attr,omitempty"`
	PresetUsesExprAttr  string                          `xml:"PresetUsesExpr,attr,omitempty"`
	IndicatorTagAttr    string                          `xml:"IndicatorTag,attr,omitempty"`
	IndicatorTagIOIAttr string                          `xml:"IndicatorTagIOI,attr,omitempty"`
	UIdAttr             string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr      string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr             string                          `xml:"Use,attr,omitempty"`
	CustomProperties    *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Preset              *EmbeddedLanguageBlockType      `xml:"Preset"`
	Body                *EmbeddedLanguageBlockType      `xml:"Body"`
}

// ASFCTransitionElementType ...
type ASFCTransitionElementType struct {
	XMLName          xml.Name                        `xml:"ASFC_TransitionElementType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr   string                          `xml:"OperandIOI,attr,omitempty"`
	ComponentIOIAttr string                          `xml:"ComponentIOI,attr,omitempty"`
	HideDescAttr     string                          `xml:"HideDesc,attr,omitempty"`
	DescXAttr        string                          `xml:"DescX,attr,omitempty"`
	DescYAttr        string                          `xml:"DescY,attr,omitempty"`
	DescWidthAttr    string                          `xml:"DescWidth,attr,omitempty"`
	ForceAttr        string                          `xml:"Force,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Condition        *EmbeddedLanguageBlockType      `xml:"Condition"`
}

// ASFCBranchElementType ...
type ASFCBranchElementType struct {
	XMLName          xml.Name                        `xml:"ASFC_BranchElementType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr,omitempty"`
	YAttr            string                          `xml:"Y,attr"`
	BranchTypeAttr   string                          `xml:"BranchType,attr"`
	BranchFlowAttr   string                          `xml:"BranchFlow,attr"`
	PriorityAttr     string                          `xml:"Priority,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Leg              []*ASFCLegElementType           `xml:"Leg"`
}

// ASFCLegElementType ...
type ASFCLegElementType struct {
	XMLName          xml.Name                        `xml:"ASFC_LegElementType"`
	IDAttr           string                          `xml:"ID,attr"`
	ForceAttr        string                          `xml:"Force,attr,omitempty"`
	ComponentIOIAttr string                          `xml:"ComponentIOI,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// SFCSBRRETType ...
type SFCSBRRETType struct {
	XMLName          xml.Name                        `xml:"SFC_SBR_RETType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	InAttr           string                          `xml:"In,attr,omitempty"`
	InIOIAttr        string                          `xml:"InIOI,attr,omitempty"`
	RetAttr          string                          `xml:"Ret,attr,omitempty"`
	RetIOIAttr       string                          `xml:"RetIOI,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	InRegionIdAttr   string                          `xml:"InRegionId,attr,omitempty"`
	RetRegionIdAttr  string                          `xml:"RetRegionId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ASFCStopElementType ...
type ASFCStopElementType struct {
	XMLName          xml.Name                        `xml:"ASFC_StopElementType"`
	IDAttr           string                          `xml:"ID,attr"`
	XAttr            string                          `xml:"X,attr"`
	YAttr            string                          `xml:"Y,attr"`
	OperandAttr      string                          `xml:"Operand,attr,omitempty"`
	OperandIOIAttr   string                          `xml:"OperandIOI,attr,omitempty"`
	ComponentIOIAttr string                          `xml:"ComponentIOI,attr,omitempty"`
	HideDescAttr     string                          `xml:"HideDesc,attr,omitempty"`
	DescXAttr        string                          `xml:"DescX,attr,omitempty"`
	DescYAttr        string                          `xml:"DescY,attr,omitempty"`
	DescWidthAttr    string                          `xml:"DescWidth,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ASFCLangElemWireType ...
type ASFCLangElemWireType struct {
	XMLName          xml.Name                        `xml:"ASFC_LangElemWireType"`
	FromIDAttr       string                          `xml:"FromID,attr"`
	ToIDAttr         string                          `xml:"ToID,attr"`
	ShowAttr         string                          `xml:"Show,attr"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// AAttachmentType ...
type AAttachmentType struct {
	FromIDAttr       string                          `xml:"FromID,attr"`
	ToIDAttr         string                          `xml:"ToID,attr"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// STContentType ...
type STContentType struct {
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	StartAttr          string                          `xml:"Start,attr,omitempty"`
	CountAttr          string                          `xml:"Count,attr,omitempty"`
	OnlineEditTypeAttr string                          `xml:"OnlineEditType,attr,omitempty"`
	UIdAttr            string                          `xml:"UId,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
	RegionIdList       string                          `xml:"RegionIdList"`
	CompactSTLines     string                          `xml:"CompactSTLines"`
	MetadataIdList     string                          `xml:"MetadataIdList"`
	Line               []*STLineType                   `xml:"Line"`
}

// STLineType ...
type STLineType struct {
	NumberAttr       string                          `xml:"Number,attr,omitempty"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	MetadataIdAttr   string                          `xml:"MetadataId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Value            string                          `xml:"Value"`
}

// DLL ...
type DLL struct {
	Value string `xml:",chardata"`
}

// EntryPoint ...
type EntryPoint struct {
	Value string `xml:",chardata"`
}

// ExternalContentType ...
type ExternalContentType struct {
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
	DLL                *DLL                            `xml:"DLL"`
	EntryPoint         *EntryPoint                     `xml:"EntryPoint"`
	ExternalRoutineXml string                          `xml:"ExternalRoutineXml"`
}

// ChildProgramCollectionType ...
type ChildProgramCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	ChildProgram     []*ChildProgramType             `xml:"ChildProgram"`
}

// ChildProgramType ...
type ChildProgramType struct {
	NameAttr         string                          `xml:"Name,attr"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// TaskCollectionType ...
type TaskCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Task             []*TaskType                     `xml:"Task"`
}

// TaskType ...
type TaskType struct {
	NameAttr                 string                          `xml:"Name,attr"`
	UIdAttr                  string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr            string                          `xml:"ParentUId,attr,omitempty"`
	TypeAttr                 string                          `xml:"Type,attr"`
	WatchdogAttr             string                          `xml:"Watchdog,attr,omitempty"`
	PriorityAttr             string                          `xml:"Priority,attr,omitempty"`
	RateAttr                 string                          `xml:"Rate,attr,omitempty"`
	DisableUpdateOutputsAttr string                          `xml:"DisableUpdateOutputs,attr,omitempty"`
	InhibitTaskAttr          string                          `xml:"InhibitTask,attr,omitempty"`
	VerifiedAttr             string                          `xml:"Verified,attr,omitempty"`
	LastScanTimeAttr         string                          `xml:"LastScanTime,attr,omitempty"`
	MaxScanTimeAttr          string                          `xml:"MaxScanTime,attr,omitempty"`
	MaxIntervalAttr          string                          `xml:"MaxInterval,attr,omitempty"`
	MinIntervalAttr          string                          `xml:"MinInterval,attr,omitempty"`
	StartTimeAttr            string                          `xml:"StartTime,attr,omitempty"`
	ClassAttr                string                          `xml:"Class,attr,omitempty"`
	UseAttr                  string                          `xml:"Use,attr,omitempty"`
	CustomProperties         *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description              *DescriptionType                `xml:"Description"`
	EventInfo                *TaskEventInfoType              `xml:"EventInfo"`
	ScheduledPrograms        *ScheduledProgramCollectionType `xml:"ScheduledPrograms"`
}

// TaskEventInfoType ...
type TaskEventInfoType struct {
	EventTriggerAttr  string                          `xml:"EventTrigger,attr,omitempty"`
	EventTagAttr      string                          `xml:"EventTag,attr,omitempty"`
	EnableTimeoutAttr string                          `xml:"EnableTimeout,attr,omitempty"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ScheduledProgramCollectionType ...
type ScheduledProgramCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	ScheduledProgram []*ScheduledProgramType         `xml:"ScheduledProgram"`
}

// ScheduledProgramType ...
type ScheduledProgramType struct {
	NameAttr         string                          `xml:"Name,attr"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// CSTType ...
type CSTType struct {
	MasterIDAttr     string                          `xml:"MasterID,attr"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	VerifiedAttr     string                          `xml:"Verified,attr,omitempty"`
	StatusAttr       string                          `xml:"Status,attr,omitempty"`
	MasterAttr       string                          `xml:"Master,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// WallClockTimeType ...
type WallClockTimeType struct {
	UIdAttr                 string                          `xml:"UId,attr,omitempty"`
	LocalTimeAdjustmentAttr string                          `xml:"LocalTimeAdjustment,attr,omitempty"`
	TimeZoneAttr            string                          `xml:"TimeZone,attr,omitempty"`
	VerifiedAttr            string                          `xml:"Verified,attr,omitempty"`
	UseAttr                 string                          `xml:"Use,attr,omitempty"`
	CustomProperties        *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// SafetyControllerType ...
type SafetyControllerType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// RedundancyType ...
type RedundancyType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// FaultLogType ...
type FaultLogType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// CommPortCollectionType ...
type CommPortCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	SerialPort       []*SerialPortType               `xml:"SerialPort"`
}

// SerialPortType ...
type SerialPortType struct {
	ChannelAttr                 string                          `xml:"Channel,attr,omitempty"`
	BaudRateAttr                string                          `xml:"BaudRate,attr,omitempty"`
	ParityAttr                  string                          `xml:"Parity,attr,omitempty"`
	DataBitsAttr                string                          `xml:"DataBits,attr,omitempty"`
	StopBitsAttr                string                          `xml:"StopBits,attr,omitempty"`
	ComDriverIdAttr             string                          `xml:"ComDriverId,attr,omitempty"`
	RTSOffDelayAttr             string                          `xml:"RTSOffDelay,attr,omitempty"`
	RTSSendDelayAttr            string                          `xml:"RTSSendDelay,attr,omitempty"`
	ControlLineAttr             string                          `xml:"ControlLine,attr,omitempty"`
	RemoteModeChangeFlagAttr    string                          `xml:"RemoteModeChangeFlag,attr,omitempty"`
	ModeChangeAttentionCharAttr string                          `xml:"ModeChangeAttentionChar,attr,omitempty"`
	SystemModeCharacterAttr     string                          `xml:"SystemModeCharacter,attr,omitempty"`
	UserModeCharacterAttr       string                          `xml:"UserModeCharacter,attr,omitempty"`
	DCDWaitDelayAttr            string                          `xml:"DCDWaitDelay,attr,omitempty"`
	VerifiedAttr                string                          `xml:"Verified,attr,omitempty"`
	UseAttr                     string                          `xml:"Use,attr,omitempty"`
	CustomProperties            *CustomPropertiesCollectionType `xml:"CustomProperties"`
	ASCII                       *ASCIIDriverType                `xml:"ASCII"`
	DF1                         *DF1DriverType                  `xml:"DF1"`
}

// DF1DriverType ...
type DF1DriverType struct {
	UIdAttr                    string                          `xml:"UId,attr,omitempty"`
	DuplicateDetectionAttr     string                          `xml:"DuplicateDetection,attr"`
	ErrorDetectionAttr         string                          `xml:"ErrorDetection,attr"`
	EmbeddedResponseEnableAttr string                          `xml:"EmbeddedResponseEnable,attr"`
	DF1ModeAttr                string                          `xml:"DF1Mode,attr"`
	ACKTimeoutAttr             string                          `xml:"ACKTimeout,attr"`
	NAKReceiveLimitAttr        string                          `xml:"NAKReceiveLimit,attr"`
	ENQTransmitLimitAttr       string                          `xml:"ENQTransmitLimit,attr"`
	TransmitRetriesAttr        string                          `xml:"TransmitRetries,attr"`
	StationAddressAttr         string                          `xml:"StationAddress,attr"`
	ReplyMessageWaitAttr       string                          `xml:"ReplyMessageWait,attr"`
	PollingModeAttr            string                          `xml:"PollingMode,attr"`
	MasterMessageTransmitAttr  string                          `xml:"MasterMessageTransmit,attr"`
	NormalPollNodeFileAttr     string                          `xml:"NormalPollNodeFile,attr"`
	NormalPollGroupSizeAttr    string                          `xml:"NormalPollGroupSize,attr"`
	PriorityPollNodeFileAttr   string                          `xml:"PriorityPollNodeFile,attr"`
	ActiveStationFileAttr      string                          `xml:"ActiveStationFile,attr"`
	SlavePollTimeoutAttr       string                          `xml:"SlavePollTimeout,attr"`
	EOTSuppressionAttr         string                          `xml:"EOTSuppression,attr"`
	MaxStationAddressAttr      string                          `xml:"MaxStationAddress,attr"`
	TokenHoldFactorAttr        string                          `xml:"TokenHoldFactor,attr"`
	StoreFwdFileAttr           string                          `xml:"StoreFwdFile,attr"`
	EnableStoreFwdAttr         string                          `xml:"EnableStoreFwd,attr"`
	VerifiedAttr               string                          `xml:"Verified,attr,omitempty"`
	UseAttr                    string                          `xml:"Use,attr,omitempty"`
	CustomProperties           *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// ASCIIDriverType ...
type ASCIIDriverType struct {
	UIdAttr              string                          `xml:"UId,attr,omitempty"`
	XONXOFFEnableAttr    string                          `xml:"XONXOFFEnable,attr"`
	DeleteModeAttr       string                          `xml:"DeleteMode,attr"`
	EchoModeAttr         string                          `xml:"EchoMode,attr"`
	TerminationCharsAttr string                          `xml:"TerminationChars,attr"`
	AppendCharsAttr      string                          `xml:"AppendChars,attr"`
	BufferSizeAttr       string                          `xml:"BufferSize,attr"`
	VerifiedAttr         string                          `xml:"Verified,attr,omitempty"`
	UseAttr              string                          `xml:"Use,attr,omitempty"`
	CustomProperties     *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// AlarmAnalogType ...
type AlarmAnalogType struct {
	EnableInAttr            string                          `xml:"EnableIn,attr,omitempty"`
	InFaultAttr             string                          `xml:"InFault,attr,omitempty"`
	HHEnabledAttr           string                          `xml:"HHEnabled,attr,omitempty"`
	HEnabledAttr            string                          `xml:"HEnabled,attr,omitempty"`
	LEnabledAttr            string                          `xml:"LEnabled,attr,omitempty"`
	LLEnabledAttr           string                          `xml:"LLEnabled,attr,omitempty"`
	AckRequiredAttr         string                          `xml:"AckRequired,attr,omitempty"`
	ProgAckAllAttr          string                          `xml:"ProgAckAll,attr,omitempty"`
	OperAckAllAttr          string                          `xml:"OperAckAll,attr,omitempty"`
	HHProgAckAttr           string                          `xml:"HHProgAck,attr,omitempty"`
	HHOperAckAttr           string                          `xml:"HHOperAck,attr,omitempty"`
	HProgAckAttr            string                          `xml:"HProgAck,attr,omitempty"`
	HOperAckAttr            string                          `xml:"HOperAck,attr,omitempty"`
	LProgAckAttr            string                          `xml:"LProgAck,attr,omitempty"`
	LOperAckAttr            string                          `xml:"LOperAck,attr,omitempty"`
	LLProgAckAttr           string                          `xml:"LLProgAck,attr,omitempty"`
	LLOperAckAttr           string                          `xml:"LLOperAck,attr,omitempty"`
	ROCPosProgAckAttr       string                          `xml:"ROCPosProgAck,attr,omitempty"`
	ROCPosOperAckAttr       string                          `xml:"ROCPosOperAck,attr,omitempty"`
	ROCNegProgAckAttr       string                          `xml:"ROCNegProgAck,attr,omitempty"`
	ROCNegOperAckAttr       string                          `xml:"ROCNegOperAck,attr,omitempty"`
	ProgSuppressAttr        string                          `xml:"ProgSuppress,attr,omitempty"`
	OperSuppressAttr        string                          `xml:"OperSuppress,attr,omitempty"`
	ProgUnsuppressAttr      string                          `xml:"ProgUnsuppress,attr,omitempty"`
	OperUnsuppressAttr      string                          `xml:"OperUnsuppress,attr,omitempty"`
	HHOperShelveAttr        string                          `xml:"HHOperShelve,attr,omitempty"`
	HOperShelveAttr         string                          `xml:"HOperShelve,attr,omitempty"`
	LOperShelveAttr         string                          `xml:"LOperShelve,attr,omitempty"`
	LLOperShelveAttr        string                          `xml:"LLOperShelve,attr,omitempty"`
	ROCPosOperShelveAttr    string                          `xml:"ROCPosOperShelve,attr,omitempty"`
	ROCNegOperShelveAttr    string                          `xml:"ROCNegOperShelve,attr,omitempty"`
	ProgUnshelveAllAttr     string                          `xml:"ProgUnshelveAll,attr,omitempty"`
	HHOperUnshelveAttr      string                          `xml:"HHOperUnshelve,attr,omitempty"`
	HOperUnshelveAttr       string                          `xml:"HOperUnshelve,attr,omitempty"`
	LOperUnshelveAttr       string                          `xml:"LOperUnshelve,attr,omitempty"`
	LLOperUnshelveAttr      string                          `xml:"LLOperUnshelve,attr,omitempty"`
	ROCPosOperUnshelveAttr  string                          `xml:"ROCPosOperUnshelve,attr,omitempty"`
	ROCNegOperUnshelveAttr  string                          `xml:"ROCNegOperUnshelve,attr,omitempty"`
	ProgDisableAttr         string                          `xml:"ProgDisable,attr,omitempty"`
	OperDisableAttr         string                          `xml:"OperDisable,attr,omitempty"`
	ProgEnableAttr          string                          `xml:"ProgEnable,attr,omitempty"`
	OperEnableAttr          string                          `xml:"OperEnable,attr,omitempty"`
	AlarmCountResetAttr     string                          `xml:"AlarmCountReset,attr,omitempty"`
	HHMinDurationEnableAttr string                          `xml:"HHMinDurationEnable,attr,omitempty"`
	HMinDurationEnableAttr  string                          `xml:"HMinDurationEnable,attr,omitempty"`
	LMinDurationEnableAttr  string                          `xml:"LMinDurationEnable,attr,omitempty"`
	LLMinDurationEnableAttr string                          `xml:"LLMinDurationEnable,attr,omitempty"`
	InAttr                  string                          `xml:"In,attr,omitempty"`
	HHLimitAttr             string                          `xml:"HHLimit,attr,omitempty"`
	HHSeverityAttr          string                          `xml:"HHSeverity,attr,omitempty"`
	HLimitAttr              string                          `xml:"HLimit,attr,omitempty"`
	HSeverityAttr           string                          `xml:"HSeverity,attr,omitempty"`
	LLimitAttr              string                          `xml:"LLimit,attr,omitempty"`
	LSeverityAttr           string                          `xml:"LSeverity,attr,omitempty"`
	LLLimitAttr             string                          `xml:"LLLimit,attr,omitempty"`
	LLSeverityAttr          string                          `xml:"LLSeverity,attr,omitempty"`
	MinDurationPREAttr      string                          `xml:"MinDurationPRE,attr,omitempty"`
	ShelveDurationAttr      string                          `xml:"ShelveDuration,attr,omitempty"`
	MaxShelveDurationAttr   string                          `xml:"MaxShelveDuration,attr,omitempty"`
	DeadbandAttr            string                          `xml:"Deadband,attr,omitempty"`
	ROCPosLimitAttr         string                          `xml:"ROCPosLimit,attr,omitempty"`
	ROCPosSeverityAttr      string                          `xml:"ROCPosSeverity,attr,omitempty"`
	ROCNegLimitAttr         string                          `xml:"ROCNegLimit,attr,omitempty"`
	ROCNegSeverityAttr      string                          `xml:"ROCNegSeverity,attr,omitempty"`
	ROCPeriodAttr           string                          `xml:"ROCPeriod,attr,omitempty"`
	AssocTag1Attr           string                          `xml:"AssocTag1,attr,omitempty"`
	AssocTag2Attr           string                          `xml:"AssocTag2,attr,omitempty"`
	AssocTag3Attr           string                          `xml:"AssocTag3,attr,omitempty"`
	AssocTag4Attr           string                          `xml:"AssocTag4,attr,omitempty"`
	EnableOutAttr           string                          `xml:"EnableOut,attr,omitempty"`
	InAlarmAttr             string                          `xml:"InAlarm,attr,omitempty"`
	AnyInAlarmUnackAttr     string                          `xml:"AnyInAlarmUnack,attr,omitempty"`
	HHInAlarmAttr           string                          `xml:"HHInAlarm,attr,omitempty"`
	HHAckedAttr             string                          `xml:"HHAcked,attr,omitempty"`
	HHInAlarmUnackAttr      string                          `xml:"HHInAlarmUnack,attr,omitempty"`
	HInAlarmAttr            string                          `xml:"HInAlarm,attr,omitempty"`
	HAckedAttr              string                          `xml:"HAcked,attr,omitempty"`
	HInAlarmUnackAttr       string                          `xml:"HInAlarmUnack,attr,omitempty"`
	LInAlarmAttr            string                          `xml:"LInAlarm,attr,omitempty"`
	LAckedAttr              string                          `xml:"LAcked,attr,omitempty"`
	LInAlarmUnackAttr       string                          `xml:"LInAlarmUnack,attr,omitempty"`
	LLInAlarmAttr           string                          `xml:"LLInAlarm,attr,omitempty"`
	LLAckedAttr             string                          `xml:"LLAcked,attr,omitempty"`
	LLInAlarmUnackAttr      string                          `xml:"LLInAlarmUnack,attr,omitempty"`
	ROCPosInAlarmAttr       string                          `xml:"ROCPosInAlarm,attr,omitempty"`
	ROCPosAckedAttr         string                          `xml:"ROCPosAcked,attr,omitempty"`
	ROCPosInAlarmUnackAttr  string                          `xml:"ROCPosInAlarmUnack,attr,omitempty"`
	ROCNegInAlarmAttr       string                          `xml:"ROCNegInAlarm,attr,omitempty"`
	ROCNegAckedAttr         string                          `xml:"ROCNegAcked,attr,omitempty"`
	ROCNegInAlarmUnackAttr  string                          `xml:"ROCNegInAlarmUnack,attr,omitempty"`
	SuppressedAttr          string                          `xml:"Suppressed,attr,omitempty"`
	HHShelvedAttr           string                          `xml:"HHShelved,attr,omitempty"`
	HShelvedAttr            string                          `xml:"HShelved,attr,omitempty"`
	LShelvedAttr            string                          `xml:"LShelved,attr,omitempty"`
	LLShelvedAttr           string                          `xml:"LLShelved,attr,omitempty"`
	ROCPosShelvedAttr       string                          `xml:"ROCPosShelved,attr,omitempty"`
	ROCNegShelvedAttr       string                          `xml:"ROCNegShelved,attr,omitempty"`
	DisabledAttr            string                          `xml:"Disabled,attr,omitempty"`
	HHInAlarmTimeAttr       string                          `xml:"HHInAlarmTime,attr,omitempty"`
	HHAlarmCountAttr        string                          `xml:"HHAlarmCount,attr,omitempty"`
	HInAlarmTimeAttr        string                          `xml:"HInAlarmTime,attr,omitempty"`
	HAlarmCountAttr         string                          `xml:"HAlarmCount,attr,omitempty"`
	LInAlarmTimeAttr        string                          `xml:"LInAlarmTime,attr,omitempty"`
	LAlarmCountAttr         string                          `xml:"LAlarmCount,attr,omitempty"`
	LLInAlarmTimeAttr       string                          `xml:"LLInAlarmTime,attr,omitempty"`
	LLAlarmCountAttr        string                          `xml:"LLAlarmCount,attr,omitempty"`
	ROCPosInAlarmTimeAttr   string                          `xml:"ROCPosInAlarmTime,attr,omitempty"`
	ROCPosAlarmCntAttr      string                          `xml:"ROCPosAlarmCnt,attr,omitempty"`
	ROCNegInAlarmTimeAttr   string                          `xml:"ROCNegInAlarmTime,attr,omitempty"`
	ROCNegAlarmCntAttr      string                          `xml:"ROCNegAlarmCnt,attr,omitempty"`
	AckTimeAttr             string                          `xml:"AckTime,attr,omitempty"`
	RetToNormalTimeAttr     string                          `xml:"RetToNormalTime,attr,omitempty"`
	AlarmCountResetTimeAttr string                          `xml:"AlarmCountResetTime,attr,omitempty"`
	ShelveTimeAttr          string                          `xml:"ShelveTime,attr,omitempty"`
	UnshelveTimeAttr        string                          `xml:"UnshelveTime,attr,omitempty"`
	InstructFaultAttr       string                          `xml:"InstructFault,attr,omitempty"`
	InFaultedAttr           string                          `xml:"InFaulted,attr,omitempty"`
	SeverityInvAttr         string                          `xml:"SeverityInv,attr,omitempty"`
	AlarmLimitsInvAttr      string                          `xml:"AlarmLimitsInv,attr,omitempty"`
	DeadbandInvAttr         string                          `xml:"DeadbandInv,attr,omitempty"`
	ROCPosLimitInvAttr      string                          `xml:"ROCPosLimitInv,attr,omitempty"`
	ROCNegLimitInvAttr      string                          `xml:"ROCNegLimitInv,attr,omitempty"`
	ROCPeriodInvAttr        string                          `xml:"ROCPeriodInv,attr,omitempty"`
	UseAttr                 string                          `xml:"Use,attr,omitempty"`
	CustomProperties        *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// AlarmDigitalType ...
type AlarmDigitalType struct {
	EnableInAttr            string                          `xml:"EnableIn,attr,omitempty"`
	InAttr                  string                          `xml:"In,attr,omitempty"`
	InFaultAttr             string                          `xml:"InFault,attr,omitempty"`
	ConditionAttr           string                          `xml:"Condition,attr,omitempty"`
	AckRequiredAttr         string                          `xml:"AckRequired,attr,omitempty"`
	LatchedAttr             string                          `xml:"Latched,attr,omitempty"`
	ProgAckAttr             string                          `xml:"ProgAck,attr,omitempty"`
	OperAckAttr             string                          `xml:"OperAck,attr,omitempty"`
	ProgResetAttr           string                          `xml:"ProgReset,attr,omitempty"`
	OperResetAttr           string                          `xml:"OperReset,attr,omitempty"`
	ProgSuppressAttr        string                          `xml:"ProgSuppress,attr,omitempty"`
	OperSuppressAttr        string                          `xml:"OperSuppress,attr,omitempty"`
	ProgUnsuppressAttr      string                          `xml:"ProgUnsuppress,attr,omitempty"`
	OperUnsuppressAttr      string                          `xml:"OperUnsuppress,attr,omitempty"`
	OperShelveAttr          string                          `xml:"OperShelve,attr,omitempty"`
	ProgUnshelveAttr        string                          `xml:"ProgUnshelve,attr,omitempty"`
	OperUnshelveAttr        string                          `xml:"OperUnshelve,attr,omitempty"`
	ProgDisableAttr         string                          `xml:"ProgDisable,attr,omitempty"`
	OperDisableAttr         string                          `xml:"OperDisable,attr,omitempty"`
	ProgEnableAttr          string                          `xml:"ProgEnable,attr,omitempty"`
	OperEnableAttr          string                          `xml:"OperEnable,attr,omitempty"`
	AlarmCountResetAttr     string                          `xml:"AlarmCountReset,attr,omitempty"`
	UseProgTimeAttr         string                          `xml:"UseProgTime,attr,omitempty"`
	SeverityAttr            string                          `xml:"Severity,attr,omitempty"`
	MinDurationPREAttr      string                          `xml:"MinDurationPRE,attr,omitempty"`
	ShelveDurationAttr      string                          `xml:"ShelveDuration,attr,omitempty"`
	MaxShelveDurationAttr   string                          `xml:"MaxShelveDuration,attr,omitempty"`
	ProgTimeAttr            string                          `xml:"ProgTime,attr,omitempty"`
	AssocTag1Attr           string                          `xml:"AssocTag1,attr,omitempty"`
	AssocTag2Attr           string                          `xml:"AssocTag2,attr,omitempty"`
	AssocTag3Attr           string                          `xml:"AssocTag3,attr,omitempty"`
	AssocTag4Attr           string                          `xml:"AssocTag4,attr,omitempty"`
	EnableOutAttr           string                          `xml:"EnableOut,attr,omitempty"`
	InAlarmAttr             string                          `xml:"InAlarm,attr,omitempty"`
	AckedAttr               string                          `xml:"Acked,attr,omitempty"`
	InAlarmUnackAttr        string                          `xml:"InAlarmUnack,attr,omitempty"`
	SuppressedAttr          string                          `xml:"Suppressed,attr,omitempty"`
	ShelvedAttr             string                          `xml:"Shelved,attr,omitempty"`
	DisabledAttr            string                          `xml:"Disabled,attr,omitempty"`
	AlarmCountAttr          string                          `xml:"AlarmCount,attr,omitempty"`
	InAlarmTimeAttr         string                          `xml:"InAlarmTime,attr,omitempty"`
	AckTimeAttr             string                          `xml:"AckTime,attr,omitempty"`
	RetToNormalTimeAttr     string                          `xml:"RetToNormalTime,attr,omitempty"`
	AlarmCountResetTimeAttr string                          `xml:"AlarmCountResetTime,attr,omitempty"`
	ShelveTimeAttr          string                          `xml:"ShelveTime,attr,omitempty"`
	UnshelveTimeAttr        string                          `xml:"UnshelveTime,attr,omitempty"`
	InstructFaultAttr       string                          `xml:"InstructFault,attr,omitempty"`
	InFaultedAttr           string                          `xml:"InFaulted,attr,omitempty"`
	SeverityInvAttr         string                          `xml:"SeverityInv,attr,omitempty"`
	UseAttr                 string                          `xml:"Use,attr,omitempty"`
	CustomProperties        *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// HMIBCType ...
type HMIBCType struct {
	EnableInAttr      string                          `xml:"EnableIn,attr,omitempty"`
	ProgFBAttr        string                          `xml:"ProgFB,attr,omitempty"`
	BitIndexAttr      string                          `xml:"BitIndex,attr,omitempty"`
	TerminalCountAttr string                          `xml:"TerminalCount,attr,omitempty"`
	EnableOutAttr     string                          `xml:"EnableOut,attr,omitempty"`
	ButtonStateAttr   string                          `xml:"ButtonState,attr,omitempty"`
	OutAttr           string                          `xml:"Out,attr,omitempty"`
	UseAttr           string                          `xml:"Use,attr,omitempty"`
	CustomProperties  *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// MsgType ...
type MsgType struct {
	MessageTypeAttr           string                          `xml:"MessageType,attr"`
	RemoteElementAttr         string                          `xml:"RemoteElement,attr,omitempty"`
	RequestedLengthAttr       string                          `xml:"RequestedLength,attr,omitempty"`
	ConnectedFlagAttr         string                          `xml:"ConnectedFlag,attr,omitempty"`
	ConnectionPathAttr        string                          `xml:"ConnectionPath,attr,omitempty"`
	CommTypeCodeAttr          string                          `xml:"CommTypeCode,attr,omitempty"`
	ServiceCodeAttr           string                          `xml:"ServiceCode,attr,omitempty"`
	ObjectTypeAttr            string                          `xml:"ObjectType,attr,omitempty"`
	TargetObjectAttr          string                          `xml:"TargetObject,attr,omitempty"`
	AttributeNumberAttr       string                          `xml:"AttributeNumber,attr,omitempty"`
	ChannelAttr               string                          `xml:"Channel,attr,omitempty"`
	DHPlusSourceLinkAttr      string                          `xml:"DHPlusSourceLink,attr,omitempty"`
	DHPlusDestinationLinkAttr string                          `xml:"DHPlusDestinationLink,attr,omitempty"`
	DHPlusDestinationNodeAttr string                          `xml:"DHPlusDestinationNode,attr,omitempty"`
	RackAttr                  string                          `xml:"Rack,attr,omitempty"`
	GroupAttr                 string                          `xml:"Group,attr,omitempty"`
	SlotAttr                  string                          `xml:"Slot,attr,omitempty"`
	LocalIndexAttr            string                          `xml:"LocalIndex,attr,omitempty"`
	RemoteIndexAttr           string                          `xml:"RemoteIndex,attr,omitempty"`
	LocalElementAttr          string                          `xml:"LocalElement,attr,omitempty"`
	DestinationTagAttr        string                          `xml:"DestinationTag,attr,omitempty"`
	CacheConnectionsAttr      string                          `xml:"CacheConnections,attr,omitempty"`
	LargePacketUsageAttr      string                          `xml:"LargePacketUsage,attr,omitempty"`
	UseAttr                   string                          `xml:"Use,attr,omitempty"`
	CustomProperties          *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// MotionGroupType ...
type MotionGroupType struct {
	CoarseUpdatePeriodAttr         string                          `xml:"CoarseUpdatePeriod,attr,omitempty"`
	PhaseShiftAttr                 string                          `xml:"PhaseShift,attr,omitempty"`
	GeneralFaultTypeAttr           string                          `xml:"GeneralFaultType,attr,omitempty"`
	AutoTagUpdateAttr              string                          `xml:"AutoTagUpdate,attr,omitempty"`
	Alternate1UpdateMultiplierAttr string                          `xml:"Alternate1UpdateMultiplier,attr,omitempty"`
	Alternate2UpdateMultiplierAttr string                          `xml:"Alternate2UpdateMultiplier,attr,omitempty"`
	UseAttr                        string                          `xml:"Use,attr,omitempty"`
	CustomProperties               *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// CoordinateSystemType ...
type CoordinateSystemType struct {
	MotionGroupInstanceAttr              string                          `xml:"MotionGroupInstance,attr,omitempty"`
	ApplicationCatalogNumberInstanceAttr string                          `xml:"ApplicationCatalogNumberInstance,attr,omitempty"`
	ApplicationCatalogNumberVersionAttr  string                          `xml:"ApplicationCatalogNumberVersion,attr,omitempty"`
	ApplicationCatalogNumberAttr         string                          `xml:"ApplicationCatalogNumber,attr,omitempty"`
	SystemTypeAttr                       string                          `xml:"SystemType,attr,omitempty"`
	CoordinateDefinitionAttr             string                          `xml:"CoordinateDefinition,attr,omitempty"`
	DimensionAttr                        string                          `xml:"Dimension,attr,omitempty"`
	AxesAttr                             string                          `xml:"Axes,attr,omitempty"`
	MaximumPendingMovesAttr              string                          `xml:"MaximumPendingMoves,attr,omitempty"`
	CoordinationModeAttr                 string                          `xml:"CoordinationMode,attr,omitempty"`
	CoordinationUnitsAttr                string                          `xml:"CoordinationUnits,attr,omitempty"`
	ConversionRatioNumeratorAttr         string                          `xml:"ConversionRatioNumerator,attr,omitempty"`
	ConversionRatioDenominatorAttr       string                          `xml:"ConversionRatioDenominator,attr,omitempty"`
	CoordinateSystemAutoTagUpdateAttr    string                          `xml:"CoordinateSystemAutoTagUpdate,attr,omitempty"`
	MaximumSpeedAttr                     string                          `xml:"MaximumSpeed,attr,omitempty"`
	MaximumAccelerationAttr              string                          `xml:"MaximumAcceleration,attr,omitempty"`
	MaximumDecelerationAttr              string                          `xml:"MaximumDeceleration,attr,omitempty"`
	ActualPositionToleranceAttr          string                          `xml:"ActualPositionTolerance,attr,omitempty"`
	CommandPositionToleranceAttr         string                          `xml:"CommandPositionTolerance,attr,omitempty"`
	TransformDimensionAttr               string                          `xml:"TransformDimension,attr,omitempty"`
	JointRatioNumeratorAttr              string                          `xml:"JointRatioNumerator,attr,omitempty"`
	JointRatioDenominatorAttr            string                          `xml:"JointRatioDenominator,attr,omitempty"`
	LinkLength1Attr                      string                          `xml:"LinkLength1,attr,omitempty"`
	LinkLength2Attr                      string                          `xml:"LinkLength2,attr,omitempty"`
	ZeroAngleOffset1Attr                 string                          `xml:"ZeroAngleOffset1,attr,omitempty"`
	ZeroAngleOffset2Attr                 string                          `xml:"ZeroAngleOffset2,attr,omitempty"`
	ZeroAngleOffset3Attr                 string                          `xml:"ZeroAngleOffset3,attr,omitempty"`
	BaseOffset1Attr                      string                          `xml:"BaseOffset1,attr,omitempty"`
	BaseOffset2Attr                      string                          `xml:"BaseOffset2,attr,omitempty"`
	BaseOffset3Attr                      string                          `xml:"BaseOffset3,attr,omitempty"`
	EndEffectorOffset1Attr               string                          `xml:"EndEffectorOffset1,attr,omitempty"`
	EndEffectorOffset2Attr               string                          `xml:"EndEffectorOffset2,attr,omitempty"`
	EndEffectorOffset3Attr               string                          `xml:"EndEffectorOffset3,attr,omitempty"`
	DynamicsConfigurationBitsAttr        string                          `xml:"DynamicsConfigurationBits,attr,omitempty"`
	MaximumAccelerationJerkAttr          string                          `xml:"MaximumAccelerationJerk,attr,omitempty"`
	MaximumDecelerationJerkAttr          string                          `xml:"MaximumDecelerationJerk,attr,omitempty"`
	MasterInputConfigurationBitsAttr     string                          `xml:"MasterInputConfigurationBits,attr,omitempty"`
	MasterPositionFilterBandwidthAttr    string                          `xml:"MasterPositionFilterBandwidth,attr,omitempty"`
	LinkLength3Attr                      string                          `xml:"LinkLength3,attr,omitempty"`
	BallScrewLeadAttr                    string                          `xml:"BallScrewLead,attr,omitempty"`
	ZeroAngleOffset4Attr                 string                          `xml:"ZeroAngleOffset4,attr,omitempty"`
	ZeroAngleOffset5Attr                 string                          `xml:"ZeroAngleOffset5,attr,omitempty"`
	ZeroAngleOffset6Attr                 string                          `xml:"ZeroAngleOffset6,attr,omitempty"`
	MaximumOrientationSpeedAttr          string                          `xml:"MaximumOrientationSpeed,attr,omitempty"`
	MaximumOrientationAccelerationAttr   string                          `xml:"MaximumOrientationAcceleration,attr,omitempty"`
	MaximumOrientationDecelerationAttr   string                          `xml:"MaximumOrientationDeceleration,attr,omitempty"`
	SwingArmA3Attr                       string                          `xml:"SwingArmA3,attr,omitempty"`
	SwingArmD3Attr                       string                          `xml:"SwingArmD3,attr,omitempty"`
	SwingArmA4Attr                       string                          `xml:"SwingArmA4,attr,omitempty"`
	SwingArmD4Attr                       string                          `xml:"SwingArmD4,attr,omitempty"`
	SwingArmD5Attr                       string                          `xml:"SwingArmD5,attr,omitempty"`
	SwingArmCouplingRatioNumeratorAttr   string                          `xml:"SwingArmCouplingRatioNumerator,attr,omitempty"`
	SwingArmCouplingRatioDenominatorAttr string                          `xml:"SwingArmCouplingRatioDenominator,attr,omitempty"`
	SwingArmCouplingDirectionAttr        string                          `xml:"SwingArmCouplingDirection,attr,omitempty"`
	RobotJointsDirectionSenseBitsAttr    string                          `xml:"RobotJointsDirectionSenseBits,attr,omitempty"`
	UseAttr                              string                          `xml:"Use,attr,omitempty"`
	CustomProperties                     *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// AxisType ...
type AxisType struct {
	MotionGroupAttr                           string                          `xml:"MotionGroup,attr,omitempty"`
	CIPDriveSetAttrUpdateBitsAttr             string                          `xml:"CIPDriveSetAttrUpdateBits,attr,omitempty"`
	CIPDriveGetAttrUpdateBitsAttr             string                          `xml:"CIPDriveGetAttrUpdateBits,attr,omitempty"`
	CIPControllerSetAttrUpdateBitsAttr        string                          `xml:"CIPControllerSetAttrUpdateBits,attr,omitempty"`
	CIPControllerGetAttrUpdateBitsAttr        string                          `xml:"CIPControllerGetAttrUpdateBits,attr,omitempty"`
	MotionModuleAttr                          string                          `xml:"MotionModule,attr,omitempty"`
	ApplicationCatalogNumberInstanceAttr      string                          `xml:"ApplicationCatalogNumberInstance,attr,omitempty"`
	ApplicationCatalogNumberVersionAttr       string                          `xml:"ApplicationCatalogNumberVersion,attr,omitempty"`
	ApplicationCatalogNumberAttr              string                          `xml:"ApplicationCatalogNumber,attr,omitempty"`
	AxisConfigurationAttr                     string                          `xml:"AxisConfiguration,attr,omitempty"`
	FeedbackConfigurationAttr                 string                          `xml:"FeedbackConfiguration,attr,omitempty"`
	MotorDataSourceAttr                       string                          `xml:"MotorDataSource,attr,omitempty"`
	MotorCatalogNumberAttr                    string                          `xml:"MotorCatalogNumber,attr,omitempty"`
	Feedback1TypeAttr                         string                          `xml:"Feedback1Type,attr,omitempty"`
	Feedback2TypeAttr                         string                          `xml:"Feedback2Type,attr,omitempty"`
	MotorTypeAttr                             string                          `xml:"MotorType,attr,omitempty"`
	MotionScalingConfigurationAttr            string                          `xml:"MotionScalingConfiguration,attr,omitempty"`
	RotationalPosResolutionAttr               string                          `xml:"RotationalPosResolution,attr,omitempty"`
	ConversionConstantAttr                    string                          `xml:"ConversionConstant,attr,omitempty"`
	OutputCamExecutionTargetsAttr             string                          `xml:"OutputCamExecutionTargets,attr,omitempty"`
	PositionUnitsAttr                         string                          `xml:"PositionUnits,attr,omitempty"`
	AverageVelocityTimebaseAttr               string                          `xml:"AverageVelocityTimebase,attr,omitempty"`
	RotaryAxisAttr                            string                          `xml:"RotaryAxis,attr,omitempty"`
	PositionUnwindAttr                        string                          `xml:"PositionUnwind,attr,omitempty"`
	HomeModeAttr                              string                          `xml:"HomeMode,attr,omitempty"`
	HomeDirectionAttr                         string                          `xml:"HomeDirection,attr,omitempty"`
	HomeSequenceAttr                          string                          `xml:"HomeSequence,attr,omitempty"`
	HomeConfigurationBitsAttr                 string                          `xml:"HomeConfigurationBits,attr,omitempty"`
	HomePositionAttr                          string                          `xml:"HomePosition,attr,omitempty"`
	HomeOffsetAttr                            string                          `xml:"HomeOffset,attr,omitempty"`
	HomeSpeedAttr                             string                          `xml:"HomeSpeed,attr,omitempty"`
	HomeReturnSpeedAttr                       string                          `xml:"HomeReturnSpeed,attr,omitempty"`
	MaximumSpeedAttr                          string                          `xml:"MaximumSpeed,attr,omitempty"`
	MaximumAccelerationAttr                   string                          `xml:"MaximumAcceleration,attr,omitempty"`
	MaximumDecelerationAttr                   string                          `xml:"MaximumDeceleration,attr,omitempty"`
	ProgrammedStopModeAttr                    string                          `xml:"ProgrammedStopMode,attr,omitempty"`
	MasterInputConfigurationBitsAttr          string                          `xml:"MasterInputConfigurationBits,attr,omitempty"`
	MasterPositionFilterBandwidthAttr         string                          `xml:"MasterPositionFilterBandwidth,attr,omitempty"`
	AxisTypeAttr                              string                          `xml:"AxisType,attr,omitempty"`
	ServoLoopConfigurationAttr                string                          `xml:"ServoLoopConfiguration,attr,omitempty"`
	ExternalDriveTypeAttr                     string                          `xml:"ExternalDriveType,attr,omitempty"`
	FaultConfigurationBitsAttr                string                          `xml:"FaultConfigurationBits,attr,omitempty"`
	AxisInfoSelect1Attr                       string                          `xml:"AxisInfoSelect1,attr,omitempty"`
	AxisInfoSelect2Attr                       string                          `xml:"AxisInfoSelect2,attr,omitempty"`
	LDTTypeAttr                               string                          `xml:"LDTType,attr,omitempty"`
	LDTRecirculationsAttr                     string                          `xml:"LDTRecirculations,attr,omitempty"`
	LDTCalibrationConstantAttr                string                          `xml:"LDTCalibrationConstant,attr,omitempty"`
	LDTCalibrationConstantUnitsAttr           string                          `xml:"LDTCalibrationConstantUnits,attr,omitempty"`
	LDTScalingAttr                            string                          `xml:"LDTScaling,attr,omitempty"`
	LDTScalingUnitsAttr                       string                          `xml:"LDTScalingUnits,attr,omitempty"`
	LDTLengthAttr                             string                          `xml:"LDTLength,attr,omitempty"`
	LDTLengthUnitsAttr                        string                          `xml:"LDTLengthUnits,attr,omitempty"`
	SSICodeTypeAttr                           string                          `xml:"SSICodeType,attr,omitempty"`
	SSIDataLengthAttr                         string                          `xml:"SSIDataLength,attr,omitempty"`
	SSIClockFrequencyAttr                     string                          `xml:"SSIClockFrequency,attr,omitempty"`
	AbsoluteFeedbackEnableAttr                string                          `xml:"AbsoluteFeedbackEnable,attr,omitempty"`
	AbsoluteFeedbackOffsetAttr                string                          `xml:"AbsoluteFeedbackOffset,attr,omitempty"`
	ServoFeedbackTypeAttr                     string                          `xml:"ServoFeedbackType,attr,omitempty"`
	ServoPolarityBitsAttr                     string                          `xml:"ServoPolarityBits,attr,omitempty"`
	VelocityFeedforwardGainAttr               string                          `xml:"VelocityFeedforwardGain,attr,omitempty"`
	AccelerationFeedforwardGainAttr           string                          `xml:"AccelerationFeedforwardGain,attr,omitempty"`
	PositionProportionalGainAttr              string                          `xml:"PositionProportionalGain,attr,omitempty"`
	PositionIntegralGainAttr                  string                          `xml:"PositionIntegralGain,attr,omitempty"`
	VelocityProportionalGainAttr              string                          `xml:"VelocityProportionalGain,attr,omitempty"`
	VelocityIntegralGainAttr                  string                          `xml:"VelocityIntegralGain,attr,omitempty"`
	VelocityScalingAttr                       string                          `xml:"VelocityScaling,attr,omitempty"`
	TorqueScalingAttr                         string                          `xml:"TorqueScaling,attr,omitempty"`
	OutputLPFilterBandwidthAttr               string                          `xml:"OutputLPFilterBandwidth,attr,omitempty"`
	IntegratorHoldEnableAttr                  string                          `xml:"IntegratorHoldEnable,attr,omitempty"`
	PositionDifferentialGainAttr              string                          `xml:"PositionDifferentialGain,attr,omitempty"`
	DirectionalScalingRatioAttr               string                          `xml:"DirectionalScalingRatio,attr,omitempty"`
	MaximumPositiveTravelAttr                 string                          `xml:"MaximumPositiveTravel,attr,omitempty"`
	MaximumNegativeTravelAttr                 string                          `xml:"MaximumNegativeTravel,attr,omitempty"`
	PositionErrorToleranceAttr                string                          `xml:"PositionErrorTolerance,attr,omitempty"`
	PositionLockToleranceAttr                 string                          `xml:"PositionLockTolerance,attr,omitempty"`
	OutputLimitAttr                           string                          `xml:"OutputLimit,attr,omitempty"`
	DirectDriveRampRateAttr                   string                          `xml:"DirectDriveRampRate,attr,omitempty"`
	OutputOffsetAttr                          string                          `xml:"OutputOffset,attr,omitempty"`
	VelocityOffsetAttr                        string                          `xml:"VelocityOffset,attr,omitempty"`
	TorqueOffsetAttr                          string                          `xml:"TorqueOffset,attr,omitempty"`
	FrictionCompensationAttr                  string                          `xml:"FrictionCompensation,attr,omitempty"`
	FrictionCompensationWindowAttr            string                          `xml:"FrictionCompensationWindow,attr,omitempty"`
	BacklashStabilizationWindowAttr           string                          `xml:"BacklashStabilizationWindow,attr,omitempty"`
	BacklashReversalOffsetAttr                string                          `xml:"BacklashReversalOffset,attr,omitempty"`
	HardOvertravelFaultActionAttr             string                          `xml:"HardOvertravelFaultAction,attr,omitempty"`
	SoftOvertravelFaultActionAttr             string                          `xml:"SoftOvertravelFaultAction,attr,omitempty"`
	PositionErrorFaultActionAttr              string                          `xml:"PositionErrorFaultAction,attr,omitempty"`
	FeedbackFaultActionAttr                   string                          `xml:"FeedbackFaultAction,attr,omitempty"`
	FeedbackNoiseFaultActionAttr              string                          `xml:"FeedbackNoiseFaultAction,attr,omitempty"`
	DriveFaultActionAttr                      string                          `xml:"DriveFaultAction,attr,omitempty"`
	TestIncrementAttr                         string                          `xml:"TestIncrement,attr,omitempty"`
	TuningTravelLimitAttr                     string                          `xml:"TuningTravelLimit,attr,omitempty"`
	TuningSpeedAttr                           string                          `xml:"TuningSpeed,attr,omitempty"`
	TuningTorqueAttr                          string                          `xml:"TuningTorque,attr,omitempty"`
	DampingFactorAttr                         string                          `xml:"DampingFactor,attr,omitempty"`
	DriveModelTimeConstantAttr                string                          `xml:"DriveModelTimeConstant,attr,omitempty"`
	PositionServoBandwidthAttr                string                          `xml:"PositionServoBandwidth,attr,omitempty"`
	VelocityServoBandwidthAttr                string                          `xml:"VelocityServoBandwidth,attr,omitempty"`
	TuningConfigurationBitsAttr               string                          `xml:"TuningConfigurationBits,attr,omitempty"`
	TorqueLimitSourceAttr                     string                          `xml:"TorqueLimitSource,attr,omitempty"`
	DriveUnitAttr                             string                          `xml:"DriveUnit,attr,omitempty"`
	PositionDataScalingAttr                   string                          `xml:"PositionDataScaling,attr,omitempty"`
	PositionDataScalingFactorAttr             string                          `xml:"PositionDataScalingFactor,attr,omitempty"`
	PositionDataScalingExpAttr                string                          `xml:"PositionDataScalingExp,attr,omitempty"`
	VelocityDataScalingAttr                   string                          `xml:"VelocityDataScaling,attr,omitempty"`
	VelocityDataScalingFactorAttr             string                          `xml:"VelocityDataScalingFactor,attr,omitempty"`
	VelocityDataScalingExpAttr                string                          `xml:"VelocityDataScalingExp,attr,omitempty"`
	AccelerationDataScalingAttr               string                          `xml:"AccelerationDataScaling,attr,omitempty"`
	AccelerationDataScalingFactorAttr         string                          `xml:"AccelerationDataScalingFactor,attr,omitempty"`
	AccelerationDataScalingExpAttr            string                          `xml:"AccelerationDataScalingExp,attr,omitempty"`
	TorqueDataScalingAttr                     string                          `xml:"TorqueDataScaling,attr,omitempty"`
	TorqueDataScalingFactorAttr               string                          `xml:"TorqueDataScalingFactor,attr,omitempty"`
	TorqueDataScalingExpAttr                  string                          `xml:"TorqueDataScalingExp,attr,omitempty"`
	DrivePolarityAttr                         string                          `xml:"DrivePolarity,attr,omitempty"`
	MotorFeedbackTypeAttr                     string                          `xml:"MotorFeedbackType,attr,omitempty"`
	MotorFeedbackResolutionAttr               string                          `xml:"MotorFeedbackResolution,attr,omitempty"`
	AuxFeedbackTypeAttr                       string                          `xml:"AuxFeedbackType,attr,omitempty"`
	AuxFeedbackResolutionAttr                 string                          `xml:"AuxFeedbackResolution,attr,omitempty"`
	MotorFeedbackUnitAttr                     string                          `xml:"MotorFeedbackUnit,attr,omitempty"`
	AuxFeedbackUnitAttr                       string                          `xml:"AuxFeedbackUnit,attr,omitempty"`
	OutputNotchFilterFrequencyAttr            string                          `xml:"OutputNotchFilterFrequency,attr,omitempty"`
	VelocityDroopAttr                         string                          `xml:"VelocityDroop,attr,omitempty"`
	VelocityLimitBipolarAttr                  string                          `xml:"VelocityLimitBipolar,attr,omitempty"`
	AccelerationLimitBipolarAttr              string                          `xml:"AccelerationLimitBipolar,attr,omitempty"`
	TorqueLimitBipolarAttr                    string                          `xml:"TorqueLimitBipolar,attr,omitempty"`
	VelocityLimitPositiveAttr                 string                          `xml:"VelocityLimitPositive,attr,omitempty"`
	VelocityLimitNegativeAttr                 string                          `xml:"VelocityLimitNegative,attr,omitempty"`
	VelocityThresholdAttr                     string                          `xml:"VelocityThreshold,attr,omitempty"`
	VelocityWindowAttr                        string                          `xml:"VelocityWindow,attr,omitempty"`
	VelocityStandstillWindowAttr              string                          `xml:"VelocityStandstillWindow,attr,omitempty"`
	AccelerationLimitPositiveAttr             string                          `xml:"AccelerationLimitPositive,attr,omitempty"`
	AccelerationLimitNegativeAttr             string                          `xml:"AccelerationLimitNegative,attr,omitempty"`
	TorqueLimitPositiveAttr                   string                          `xml:"TorqueLimitPositive,attr,omitempty"`
	TorqueLimitNegativeAttr                   string                          `xml:"TorqueLimitNegative,attr,omitempty"`
	TorqueThresholdAttr                       string                          `xml:"TorqueThreshold,attr,omitempty"`
	DriveThermalFaultActionAttr               string                          `xml:"DriveThermalFaultAction,attr,omitempty"`
	MotorThermalFaultActionAttr               string                          `xml:"MotorThermalFaultAction,attr,omitempty"`
	DriveEnableInputFaultActionAttr           string                          `xml:"DriveEnableInputFaultAction,attr,omitempty"`
	StoppingTorqueAttr                        string                          `xml:"StoppingTorque,attr,omitempty"`
	StoppingTimeLimitAttr                     string                          `xml:"StoppingTimeLimit,attr,omitempty"`
	BrakeEngageDelayTimeAttr                  string                          `xml:"BrakeEngageDelayTime,attr,omitempty"`
	BrakeReleaseDelayTimeAttr                 string                          `xml:"BrakeReleaseDelayTime,attr,omitempty"`
	PowerSupplyIDAttr                         string                          `xml:"PowerSupplyID,attr,omitempty"`
	BusRegulatorIDAttr                        string                          `xml:"BusRegulatorID,attr,omitempty"`
	PWMFrequencySelectAttr                    string                          `xml:"PWMFrequencySelect,attr,omitempty"`
	LoadInertiaRatioAttr                      string                          `xml:"LoadInertiaRatio,attr,omitempty"`
	AmplifierCatalogNumberAttr                string                          `xml:"AmplifierCatalogNumber,attr,omitempty"`
	AuxFeedbackRatioAttr                      string                          `xml:"AuxFeedbackRatio,attr,omitempty"`
	ContinuousTorqueLimitAttr                 string                          `xml:"ContinuousTorqueLimit,attr,omitempty"`
	ResistiveBrakeContactDelayAttr            string                          `xml:"ResistiveBrakeContactDelay,attr,omitempty"`
	ConfigurationProfileAttr                  string                          `xml:"ConfigurationProfile,attr,omitempty"`
	RegistrationInputsAttr                    string                          `xml:"RegistrationInputs,attr,omitempty"`
	MaximumAccelerationJerkAttr               string                          `xml:"MaximumAccelerationJerk,attr,omitempty"`
	MaximumDecelerationJerkAttr               string                          `xml:"MaximumDecelerationJerk,attr,omitempty"`
	DynamicsConfigurationBitsAttr             string                          `xml:"DynamicsConfigurationBits,attr,omitempty"`
	PhaseLossFaultActionAttr                  string                          `xml:"PhaseLossFaultAction,attr,omitempty"`
	HomeTorqueLevelAttr                       string                          `xml:"HomeTorqueLevel,attr,omitempty"`
	InputPowerPhaseAttr                       string                          `xml:"InputPowerPhase,attr,omitempty"`
	FeedbackUnitRatioAttr                     string                          `xml:"FeedbackUnitRatio,attr,omitempty"`
	AccelerationLimitAttr                     string                          `xml:"AccelerationLimit,attr,omitempty"`
	DecelerationLimitAttr                     string                          `xml:"DecelerationLimit,attr,omitempty"`
	RampJerkControlAttr                       string                          `xml:"RampJerkControl,attr,omitempty"`
	FlyingStartEnableAttr                     string                          `xml:"FlyingStartEnable,attr,omitempty"`
	SkipSpeed1Attr                            string                          `xml:"SkipSpeed1,attr,omitempty"`
	SkipSpeed2Attr                            string                          `xml:"SkipSpeed2,attr,omitempty"`
	SkipSpeed3Attr                            string                          `xml:"SkipSpeed3,attr,omitempty"`
	SkipSpeedBandAttr                         string                          `xml:"SkipSpeedBand,attr,omitempty"`
	CommandTorqueAttr                         string                          `xml:"CommandTorque,attr,omitempty"`
	PositionIntegratorBandwidthAttr           string                          `xml:"PositionIntegratorBandwidth,attr,omitempty"`
	PositionErrorToleranceTimeAttr            string                          `xml:"PositionErrorToleranceTime,attr,omitempty"`
	PositionIntegratorControlAttr             string                          `xml:"PositionIntegratorControl,attr,omitempty"`
	PositionIntegratorPreloadAttr             string                          `xml:"PositionIntegratorPreload,attr,omitempty"`
	VelocityErrorToleranceAttr                string                          `xml:"VelocityErrorTolerance,attr,omitempty"`
	VelocityErrorToleranceTimeAttr            string                          `xml:"VelocityErrorToleranceTime,attr,omitempty"`
	VelocityIntegratorControlAttr             string                          `xml:"VelocityIntegratorControl,attr,omitempty"`
	VelocityIntegratorPreloadAttr             string                          `xml:"VelocityIntegratorPreload,attr,omitempty"`
	VelocityLockToleranceAttr                 string                          `xml:"VelocityLockTolerance,attr,omitempty"`
	SystemInertiaAttr                         string                          `xml:"SystemInertia,attr,omitempty"`
	TorqueLowPassFilterBandwidthAttr          string                          `xml:"TorqueLowPassFilterBandwidth,attr,omitempty"`
	TorqueNotchFilterFrequencyAttr            string                          `xml:"TorqueNotchFilterFrequency,attr,omitempty"`
	TorqueRateLimitAttr                       string                          `xml:"TorqueRateLimit,attr,omitempty"`
	OvertorqueLimitAttr                       string                          `xml:"OvertorqueLimit,attr,omitempty"`
	OvertorqueLimitTimeAttr                   string                          `xml:"OvertorqueLimitTime,attr,omitempty"`
	UndertorqueLimitAttr                      string                          `xml:"UndertorqueLimit,attr,omitempty"`
	UndertorqueLimitTimeAttr                  string                          `xml:"UndertorqueLimitTime,attr,omitempty"`
	FluxCurrentReferenceAttr                  string                          `xml:"FluxCurrentReference,attr,omitempty"`
	CurrentErrorAttr                          string                          `xml:"CurrentError,attr,omitempty"`
	TorqueLoopBandwidthAttr                   string                          `xml:"TorqueLoopBandwidth,attr,omitempty"`
	TorqueIntegralTimeConstantAttr            string                          `xml:"TorqueIntegralTimeConstant,attr,omitempty"`
	FluxLoopBandwidthAttr                     string                          `xml:"FluxLoopBandwidth,attr,omitempty"`
	FluxIntegralTimeConstantAttr              string                          `xml:"FluxIntegralTimeConstant,attr,omitempty"`
	FluxUpControlAttr                         string                          `xml:"FluxUpControl,attr,omitempty"`
	FluxUpTimeAttr                            string                          `xml:"FluxUpTime,attr,omitempty"`
	FrequencyControlMethodAttr                string                          `xml:"FrequencyControlMethod,attr,omitempty"`
	MaximumVoltageAttr                        string                          `xml:"MaximumVoltage,attr,omitempty"`
	MaximumFrequencyAttr                      string                          `xml:"MaximumFrequency,attr,omitempty"`
	BreakVoltageAttr                          string                          `xml:"BreakVoltage,attr,omitempty"`
	BreakFrequencyAttr                        string                          `xml:"BreakFrequency,attr,omitempty"`
	StartBoostAttr                            string                          `xml:"StartBoost,attr,omitempty"`
	RunBoostAttr                              string                          `xml:"RunBoost,attr,omitempty"`
	StoppingActionAttr                        string                          `xml:"StoppingAction,attr,omitempty"`
	MechanicalBrakeControlAttr                string                          `xml:"MechanicalBrakeControl,attr,omitempty"`
	MechanicalBrakeReleaseDelayAttr           string                          `xml:"MechanicalBrakeReleaseDelay,attr,omitempty"`
	MechanicalBrakeEngageDelayAttr            string                          `xml:"MechanicalBrakeEngageDelay,attr,omitempty"`
	PowerLossActionAttr                       string                          `xml:"PowerLossAction,attr,omitempty"`
	PowerLossThresholdAttr                    string                          `xml:"PowerLossThreshold,attr,omitempty"`
	ShutdownActionAttr                        string                          `xml:"ShutdownAction,attr,omitempty"`
	InverterCapacityAttr                      string                          `xml:"InverterCapacity,attr,omitempty"`
	ConverterCapacityAttr                     string                          `xml:"ConverterCapacity,attr,omitempty"`
	InverterOverloadActionAttr                string                          `xml:"InverterOverloadAction,attr,omitempty"`
	MotorOverloadActionAttr                   string                          `xml:"MotorOverloadAction,attr,omitempty"`
	CIPAxisExceptionActionAttr                string                          `xml:"CIPAxisExceptionAction,attr,omitempty"`
	CIPAxisExceptionActionRAAttr              string                          `xml:"CIPAxisExceptionActionRA,attr,omitempty"`
	MotorOverspeedUserLimitAttr               string                          `xml:"MotorOverspeedUserLimit,attr,omitempty"`
	MotorThermalOverloadUserLimitAttr         string                          `xml:"MotorThermalOverloadUserLimit,attr,omitempty"`
	InverterThermalOverloadUserLimitAttr      string                          `xml:"InverterThermalOverloadUserLimit,attr,omitempty"`
	LocalControlAttr                          string                          `xml:"LocalControl,attr,omitempty"`
	PositionLeadLagFilterBandwidthAttr        string                          `xml:"PositionLeadLagFilterBandwidth,attr,omitempty"`
	PositionLeadLagFilterGainAttr             string                          `xml:"PositionLeadLagFilterGain,attr,omitempty"`
	PositionNotchFilterFrequencyAttr          string                          `xml:"PositionNotchFilterFrequency,attr,omitempty"`
	VelocityNegativeFeedforwardGainAttr       string                          `xml:"VelocityNegativeFeedforwardGain,attr,omitempty"`
	BacklashCompensationWindowAttr            string                          `xml:"BacklashCompensationWindow,attr,omitempty"`
	TorqueLeadLagFilterBandwidthAttr          string                          `xml:"TorqueLeadLagFilterBandwidth,attr,omitempty"`
	TorqueLeadLagFilterGainAttr               string                          `xml:"TorqueLeadLagFilterGain,attr,omitempty"`
	SLATConfigurationAttr                     string                          `xml:"SLATConfiguration,attr,omitempty"`
	SLATSetPointAttr                          string                          `xml:"SLATSetPoint,attr,omitempty"`
	SLATTimeDelayAttr                         string                          `xml:"SLATTimeDelay,attr,omitempty"`
	CurrentDisturbanceAttr                    string                          `xml:"CurrentDisturbance,attr,omitempty"`
	CommutationSelfSensingCurrentAttr         string                          `xml:"CommutationSelfSensingCurrent,attr,omitempty"`
	DCInjectionBrakeCurrentAttr               string                          `xml:"DCInjectionBrakeCurrent,attr,omitempty"`
	FluxBrakingEnableAttr                     string                          `xml:"FluxBrakingEnable,attr,omitempty"`
	DCInjectionBrakeTimeAttr                  string                          `xml:"DCInjectionBrakeTime,attr,omitempty"`
	MotorDeviceCodeAttr                       string                          `xml:"MotorDeviceCode,attr,omitempty"`
	MotorUnitAttr                             string                          `xml:"MotorUnit,attr,omitempty"`
	MotorPolarityAttr                         string                          `xml:"MotorPolarity,attr,omitempty"`
	MotorRatedVoltageAttr                     string                          `xml:"MotorRatedVoltage,attr,omitempty"`
	MotorRatedContinuousCurrentAttr           string                          `xml:"MotorRatedContinuousCurrent,attr,omitempty"`
	MotorRatedPeakCurrentAttr                 string                          `xml:"MotorRatedPeakCurrent,attr,omitempty"`
	MotorRatedOutputPowerAttr                 string                          `xml:"MotorRatedOutputPower,attr,omitempty"`
	MotorOverloadLimitAttr                    string                          `xml:"MotorOverloadLimit,attr,omitempty"`
	MotorIntegralThermalSwitchAttr            string                          `xml:"MotorIntegralThermalSwitch,attr,omitempty"`
	MotorMaxWindingTemperatureAttr            string                          `xml:"MotorMaxWindingTemperature,attr,omitempty"`
	MotorWindingToAmbientCapacitanceAttr      string                          `xml:"MotorWindingToAmbientCapacitance,attr,omitempty"`
	MotorWindingToAmbientResistanceAttr       string                          `xml:"MotorWindingToAmbientResistance,attr,omitempty"`
	PMMotorResistanceAttr                     string                          `xml:"PMMotorResistance,attr,omitempty"`
	PMMotorInductanceAttr                     string                          `xml:"PMMotorInductance,attr,omitempty"`
	RotaryMotorPolesAttr                      string                          `xml:"RotaryMotorPoles,attr,omitempty"`
	RotaryMotorInertiaAttr                    string                          `xml:"RotaryMotorInertia,attr,omitempty"`
	RotaryMotorRatedSpeedAttr                 string                          `xml:"RotaryMotorRatedSpeed,attr,omitempty"`
	RotaryMotorMaxSpeedAttr                   string                          `xml:"RotaryMotorMaxSpeed,attr,omitempty"`
	RotaryMotorDampingCoefficientAttr         string                          `xml:"RotaryMotorDampingCoefficient,attr,omitempty"`
	LinearMotorPolePitchAttr                  string                          `xml:"LinearMotorPolePitch,attr,omitempty"`
	LinearMotorRatedSpeedAttr                 string                          `xml:"LinearMotorRatedSpeed,attr,omitempty"`
	LinearMotorMassAttr                       string                          `xml:"LinearMotorMass,attr,omitempty"`
	LinearMotorMaxSpeedAttr                   string                          `xml:"LinearMotorMaxSpeed,attr,omitempty"`
	LinearMotorDampingCoefficientAttr         string                          `xml:"LinearMotorDampingCoefficient,attr,omitempty"`
	PMMotorRatedTorqueAttr                    string                          `xml:"PMMotorRatedTorque,attr,omitempty"`
	PMMotorTorqueConstantAttr                 string                          `xml:"PMMotorTorqueConstant,attr,omitempty"`
	PMMotorRotaryVoltageConstantAttr          string                          `xml:"PMMotorRotaryVoltageConstant,attr,omitempty"`
	PMMotorRatedForceAttr                     string                          `xml:"PMMotorRatedForce,attr,omitempty"`
	PMMotorForceConstantAttr                  string                          `xml:"PMMotorForceConstant,attr,omitempty"`
	PMMotorLinearVoltageConstantAttr          string                          `xml:"PMMotorLinearVoltageConstant,attr,omitempty"`
	InductionMotorRatedFrequencyAttr          string                          `xml:"InductionMotorRatedFrequency,attr,omitempty"`
	InductionMotorFluxCurrentAttr             string                          `xml:"InductionMotorFluxCurrent,attr,omitempty"`
	InductionMotorStatorResistanceAttr        string                          `xml:"InductionMotorStatorResistance,attr,omitempty"`
	InductionMotorStatorLeakageReactanceAttr  string                          `xml:"InductionMotorStatorLeakageReactance,attr,omitempty"`
	InductionMotorMagnetizationReactanceAttr  string                          `xml:"InductionMotorMagnetizationReactance,attr,omitempty"`
	InductionMotorRotorResistanceAttr         string                          `xml:"InductionMotorRotorResistance,attr,omitempty"`
	InductionMotorRotorLeakageReactanceAttr   string                          `xml:"InductionMotorRotorLeakageReactance,attr,omitempty"`
	Feedback1UnitAttr                         string                          `xml:"Feedback1Unit,attr,omitempty"`
	Feedback1PolarityAttr                     string                          `xml:"Feedback1Polarity,attr,omitempty"`
	Feedback1StartupMethodAttr                string                          `xml:"Feedback1StartupMethod,attr,omitempty"`
	Feedback1CycleResolutionAttr              string                          `xml:"Feedback1CycleResolution,attr,omitempty"`
	Feedback1CycleInterpolationAttr           string                          `xml:"Feedback1CycleInterpolation,attr,omitempty"`
	Feedback1TurnsAttr                        string                          `xml:"Feedback1Turns,attr,omitempty"`
	Feedback1LengthAttr                       string                          `xml:"Feedback1Length,attr,omitempty"`
	Feedback1DataLengthAttr                   string                          `xml:"Feedback1DataLength,attr,omitempty"`
	Feedback1DataCodeAttr                     string                          `xml:"Feedback1DataCode,attr,omitempty"`
	Feedback1ResolverTransformerRatioAttr     string                          `xml:"Feedback1ResolverTransformerRatio,attr,omitempty"`
	Feedback1ResolverExcitationVoltageAttr    string                          `xml:"Feedback1ResolverExcitationVoltage,attr,omitempty"`
	Feedback1ResolverExcitationFrequencyAttr  string                          `xml:"Feedback1ResolverExcitationFrequency,attr,omitempty"`
	Feedback1ResolverCableBalanceAttr         string                          `xml:"Feedback1ResolverCableBalance,attr,omitempty"`
	Feedback1VelocityFilterBandwidthAttr      string                          `xml:"Feedback1VelocityFilterBandwidth,attr,omitempty"`
	Feedback1AccelFilterBandwidthAttr         string                          `xml:"Feedback1AccelFilterBandwidth,attr,omitempty"`
	PMMotorFluxSaturationAttr                 string                          `xml:"PMMotorFluxSaturation,attr,omitempty"`
	RotaryMotorFanCoolingSpeedAttr            string                          `xml:"RotaryMotorFanCoolingSpeed,attr,omitempty"`
	RotaryMotorFanCoolingDeratingAttr         string                          `xml:"RotaryMotorFanCoolingDerating,attr,omitempty"`
	LinearMotorIntegralLimitSwitchAttr        string                          `xml:"LinearMotorIntegralLimitSwitch,attr,omitempty"`
	Feedback1LossActionAttr                   string                          `xml:"Feedback1LossAction,attr,omitempty"`
	Feedback1VelocityFilterTapsAttr           string                          `xml:"Feedback1VelocityFilterTaps,attr,omitempty"`
	Feedback1AccelFilterTapsAttr              string                          `xml:"Feedback1AccelFilterTaps,attr,omitempty"`
	CIPAxisFaultLogAttr                       string                          `xml:"CIPAxisFaultLog,attr,omitempty"`
	CyclicReadUpdateListAttr                  string                          `xml:"CyclicReadUpdateList,attr,omitempty"`
	CyclicWriteUpdateListAttr                 string                          `xml:"CyclicWriteUpdateList,attr,omitempty"`
	ScalingSourceAttr                         string                          `xml:"ScalingSource,attr,omitempty"`
	LoadTypeAttr                              string                          `xml:"LoadType,attr,omitempty"`
	ActuatorTypeAttr                          string                          `xml:"ActuatorType,attr,omitempty"`
	TravelModeAttr                            string                          `xml:"TravelMode,attr,omitempty"`
	PositionScalingNumeratorAttr              string                          `xml:"PositionScalingNumerator,attr,omitempty"`
	PositionScalingDenominatorAttr            string                          `xml:"PositionScalingDenominator,attr,omitempty"`
	PositionUnwindNumeratorAttr               string                          `xml:"PositionUnwindNumerator,attr,omitempty"`
	PositionUnwindDenominatorAttr             string                          `xml:"PositionUnwindDenominator,attr,omitempty"`
	TravelRangeAttr                           string                          `xml:"TravelRange,attr,omitempty"`
	MotionResolutionAttr                      string                          `xml:"MotionResolution,attr,omitempty"`
	MotionPolarityAttr                        string                          `xml:"MotionPolarity,attr,omitempty"`
	MotorTestResistanceAttr                   string                          `xml:"MotorTestResistance,attr,omitempty"`
	MotorTestInductanceAttr                   string                          `xml:"MotorTestInductance,attr,omitempty"`
	MotorTestFluxCurrentAttr                  string                          `xml:"MotorTestFluxCurrent,attr,omitempty"`
	MotorTestSlipSpeedAttr                    string                          `xml:"MotorTestSlipSpeed,attr,omitempty"`
	TuneFrictionAttr                          string                          `xml:"TuneFriction,attr,omitempty"`
	TuneLoadOffsetAttr                        string                          `xml:"TuneLoadOffset,attr,omitempty"`
	TotalInertiaAttr                          string                          `xml:"TotalInertia,attr,omitempty"`
	TuningSelectAttr                          string                          `xml:"TuningSelect,attr,omitempty"`
	TuningDirectionAttr                       string                          `xml:"TuningDirection,attr,omitempty"`
	ApplicationTypeAttr                       string                          `xml:"ApplicationType,attr,omitempty"`
	LoopResponseAttr                          string                          `xml:"LoopResponse,attr,omitempty"`
	FeedbackCommutationAlignedAttr            string                          `xml:"FeedbackCommutationAligned,attr,omitempty"`
	FrictionCompensationSlidingAttr           string                          `xml:"FrictionCompensationSliding,attr,omitempty"`
	FrictionCompensationStaticAttr            string                          `xml:"FrictionCompensationStatic,attr,omitempty"`
	FrictionCompensationViscousAttr           string                          `xml:"FrictionCompensationViscous,attr,omitempty"`
	PositionLoopBandwidthAttr                 string                          `xml:"PositionLoopBandwidth,attr,omitempty"`
	VelocityLoopBandwidthAttr                 string                          `xml:"VelocityLoopBandwidth,attr,omitempty"`
	VelocityIntegratorBandwidthAttr           string                          `xml:"VelocityIntegratorBandwidth,attr,omitempty"`
	FeedbackNoiseUserLimitAttr                string                          `xml:"FeedbackNoiseUserLimit,attr,omitempty"`
	FeedbackSignalLossUserLimitAttr           string                          `xml:"FeedbackSignalLossUserLimit,attr,omitempty"`
	FeedbackDataLossUserLimitAttr             string                          `xml:"FeedbackDataLossUserLimit,attr,omitempty"`
	Feedback1BatteryAbsoluteAttr              string                          `xml:"Feedback1BatteryAbsolute,attr,omitempty"`
	CIPAxisAlarmtLogAttr                      string                          `xml:"CIPAxisAlarmtLog,attr,omitempty"`
	MotionExceptionActionAttr                 string                          `xml:"MotionExceptionAction,attr,omitempty"`
	SoftTravelLimitCheckingAttr               string                          `xml:"SoftTravelLimitChecking,attr,omitempty"`
	MoveRepeatModeAttr                        string                          `xml:"MoveRepeatMode,attr,omitempty"`
	MoveRepeatDwellAttr                       string                          `xml:"MoveRepeatDwell,attr,omitempty"`
	Feedback2UnitAttr                         string                          `xml:"Feedback2Unit,attr,omitempty"`
	Feedback2PolarityAttr                     string                          `xml:"Feedback2Polarity,attr,omitempty"`
	Feedback2StartupMethodAttr                string                          `xml:"Feedback2StartupMethod,attr,omitempty"`
	Feedback2CycleResolutionAttr              string                          `xml:"Feedback2CycleResolution,attr,omitempty"`
	Feedback2CycleInterpolationAttr           string                          `xml:"Feedback2CycleInterpolation,attr,omitempty"`
	Feedback2TurnsAttr                        string                          `xml:"Feedback2Turns,attr,omitempty"`
	Feedback2LengthAttr                       string                          `xml:"Feedback2Length,attr,omitempty"`
	Feedback2DataLengthAttr                   string                          `xml:"Feedback2DataLength,attr,omitempty"`
	Feedback2DataCodeAttr                     string                          `xml:"Feedback2DataCode,attr,omitempty"`
	Feedback2ResolverTransformerRatioAttr     string                          `xml:"Feedback2ResolverTransformerRatio,attr,omitempty"`
	Feedback2ResolverExcitationVoltageAttr    string                          `xml:"Feedback2ResolverExcitationVoltage,attr,omitempty"`
	Feedback2ResolverExcitationFrequencyAttr  string                          `xml:"Feedback2ResolverExcitationFrequency,attr,omitempty"`
	Feedback2ResolverCableBalanceAttr         string                          `xml:"Feedback2ResolverCableBalance,attr,omitempty"`
	Feedback2VelocityFilterBandwidthAttr      string                          `xml:"Feedback2VelocityFilterBandwidth,attr,omitempty"`
	Feedback2AccelFilterBandwidthAttr         string                          `xml:"Feedback2AccelFilterBandwidth,attr,omitempty"`
	Feedback2LossActionAttr                   string                          `xml:"Feedback2LossAction,attr,omitempty"`
	Feedback2VelocityFilterTapsAttr           string                          `xml:"Feedback2VelocityFilterTaps,attr,omitempty"`
	Feedback2AccelFilterTapsAttr              string                          `xml:"Feedback2AccelFilterTaps,attr,omitempty"`
	Feedback2BatteryAbsoluteAttr              string                          `xml:"Feedback2BatteryAbsolute,attr,omitempty"`
	LoadRatioAttr                             string                          `xml:"LoadRatio,attr,omitempty"`
	TotalMassAttr                             string                          `xml:"TotalMass,attr,omitempty"`
	TuneInertiaMassAttr                       string                          `xml:"TuneInertiaMass,attr,omitempty"`
	SoftTravelLimitPositiveAttr               string                          `xml:"SoftTravelLimitPositive,attr,omitempty"`
	SoftTravelLimitNegativeAttr               string                          `xml:"SoftTravelLimitNegative,attr,omitempty"`
	GainTuningConfigurationBitsAttr           string                          `xml:"GainTuningConfigurationBits,attr,omitempty"`
	CommutationOffsetAttr                     string                          `xml:"CommutationOffset,attr,omitempty"`
	PowerLossTimeAttr                         string                          `xml:"PowerLossTime,attr,omitempty"`
	SystemBandwidthAttr                       string                          `xml:"SystemBandwidth,attr,omitempty"`
	VelocityLowPassFilterBandwidthAttr        string                          `xml:"VelocityLowPassFilterBandwidth,attr,omitempty"`
	FeedbackMasterSelectAttr                  string                          `xml:"FeedbackMasterSelect,attr,omitempty"`
	TransmissionRatioInputAttr                string                          `xml:"TransmissionRatioInput,attr,omitempty"`
	TransmissionRatioOutputAttr               string                          `xml:"TransmissionRatioOutput,attr,omitempty"`
	ActuatorLeadAttr                          string                          `xml:"ActuatorLead,attr,omitempty"`
	ActuatorLeadUnitAttr                      string                          `xml:"ActuatorLeadUnit,attr,omitempty"`
	ActuatorDiameterAttr                      string                          `xml:"ActuatorDiameter,attr,omitempty"`
	ActuatorDiameterUnitAttr                  string                          `xml:"ActuatorDiameterUnit,attr,omitempty"`
	SystemAccelerationBaseAttr                string                          `xml:"SystemAccelerationBase,attr,omitempty"`
	DriveModelTimeConstantBaseAttr            string                          `xml:"DriveModelTimeConstantBase,attr,omitempty"`
	DriveRatedPeakCurrentAttr                 string                          `xml:"DriveRatedPeakCurrent,attr,omitempty"`
	HookupTestDistanceAttr                    string                          `xml:"HookupTestDistance,attr,omitempty"`
	HookupTestTimeAttr                        string                          `xml:"HookupTestTime,attr,omitempty"`
	HookupTestFeedbackChannelAttr             string                          `xml:"HookupTestFeedbackChannel,attr,omitempty"`
	LoadCouplingAttr                          string                          `xml:"LoadCoupling,attr,omitempty"`
	SystemDampingAttr                         string                          `xml:"SystemDamping,attr,omitempty"`
	CurrentVectorLimitAttr                    string                          `xml:"CurrentVectorLimit,attr,omitempty"`
	InductionMotorRatedSlipSpeedAttr          string                          `xml:"InductionMotorRatedSlipSpeed,attr,omitempty"`
	CommutationPolarityAttr                   string                          `xml:"CommutationPolarity,attr,omitempty"`
	LoadObserverConfigurationAttr             string                          `xml:"LoadObserverConfiguration,attr,omitempty"`
	LoadObserverBandwidthAttr                 string                          `xml:"LoadObserverBandwidth,attr,omitempty"`
	LoadObserverIntegratorBandwidthAttr       string                          `xml:"LoadObserverIntegratorBandwidth,attr,omitempty"`
	LoadObserverFeedbackGainAttr              string                          `xml:"LoadObserverFeedbackGain,attr,omitempty"`
	AxisIDAttr                                string                          `xml:"AxisID,attr,omitempty"`
	MotorDataAttr                             string                          `xml:"MotorData,attr,omitempty"`
	AdditionalBusCapacitanceAttr              string                          `xml:"AdditionalBusCapacitance,attr,omitempty"`
	InterpolatedPositionConfigurationAttr     string                          `xml:"InterpolatedPositionConfiguration,attr,omitempty"`
	AxisUpdateScheduleAttr                    string                          `xml:"AxisUpdateSchedule,attr,omitempty"`
	ProvingConfigurationAttr                  string                          `xml:"ProvingConfiguration,attr,omitempty"`
	TorqueProveCurrentAttr                    string                          `xml:"TorqueProveCurrent,attr,omitempty"`
	BrakeTestTorqueAttr                       string                          `xml:"BrakeTestTorque,attr,omitempty"`
	BrakeProveRampTimeAttr                    string                          `xml:"BrakeProveRampTime,attr,omitempty"`
	BrakeSlipToleranceAttr                    string                          `xml:"BrakeSlipTolerance,attr,omitempty"`
	ZeroSpeedAttr                             string                          `xml:"ZeroSpeed,attr,omitempty"`
	ZeroSpeedTimeAttr                         string                          `xml:"ZeroSpeedTime,attr,omitempty"`
	MotorPhaseLossLimitAttr                   string                          `xml:"MotorPhaseLossLimit,attr,omitempty"`
	AdaptiveTuningConfigurationAttr           string                          `xml:"AdaptiveTuningConfiguration,attr,omitempty"`
	TorqueNotchFilterHighFrequencyLimitAttr   string                          `xml:"TorqueNotchFilterHighFrequencyLimit,attr,omitempty"`
	TorqueNotchFilterLowFrequencyLimitAttr    string                          `xml:"TorqueNotchFilterLowFrequencyLimit,attr,omitempty"`
	TorqueNotchFilterTuningThresholdAttr      string                          `xml:"TorqueNotchFilterTuningThreshold,attr,omitempty"`
	AutoSagConfigurationAttr                  string                          `xml:"AutoSagConfiguration,attr,omitempty"`
	AutoSagSlipIncrementAttr                  string                          `xml:"AutoSagSlipIncrement,attr,omitempty"`
	AutoSagSlipTimeLimitAttr                  string                          `xml:"AutoSagSlipTimeLimit,attr,omitempty"`
	AutoSagStartAttr                          string                          `xml:"AutoSagStart,attr,omitempty"`
	CoastingTimeLimitAttr                     string                          `xml:"CoastingTimeLimit,attr,omitempty"`
	SafeTorqueOffActionAttr                   string                          `xml:"SafeTorqueOffAction,attr,omitempty"`
	FlyingStartMethodAttr                     string                          `xml:"FlyingStartMethod,attr,omitempty"`
	CommutationOffsetCompensationAttr         string                          `xml:"CommutationOffsetCompensation,attr,omitempty"`
	PMMotorLqInductanceAttr                   string                          `xml:"PMMotorLqInductance,attr,omitempty"`
	PMMotorLdInductanceAttr                   string                          `xml:"PMMotorLdInductance,attr,omitempty"`
	PMMotorLqFluxSaturationAttr               string                          `xml:"PMMotorLqFluxSaturation,attr,omitempty"`
	PMMotorLdFluxSaturationAttr               string                          `xml:"PMMotorLdFluxSaturation,attr,omitempty"`
	PMMotorRotaryBusOvervoltageSpeedAttr      string                          `xml:"PMMotorRotaryBusOvervoltageSpeed,attr,omitempty"`
	PMMotorLinearBusOvervoltageSpeedAttr      string                          `xml:"PMMotorLinearBusOvervoltageSpeed,attr,omitempty"`
	PMMotorRotaryMaxExtendedSpeedAttr         string                          `xml:"PMMotorRotaryMaxExtendedSpeed,attr,omitempty"`
	PMMotorLinearMaxExtendedSpeedAttr         string                          `xml:"PMMotorLinearMaxExtendedSpeed,attr,omitempty"`
	PMMotorExtendedSpeedPermissiveAttr        string                          `xml:"PMMotorExtendedSpeedPermissive,attr,omitempty"`
	BusOvervoltageOperationalLimitAttr        string                          `xml:"BusOvervoltageOperationalLimit,attr,omitempty"`
	ConnectionLossStoppingActionAttr          string                          `xml:"ConnectionLossStoppingAction,attr,omitempty"`
	SafeTorqueOffActionSourceAttr             string                          `xml:"SafeTorqueOffActionSource,attr,omitempty"`
	SafeStoppingActionAttr                    string                          `xml:"SafeStoppingAction,attr,omitempty"`
	SafeStoppingActionSourceAttr              string                          `xml:"SafeStoppingActionSource,attr,omitempty"`
	VerticalLoadControlAttr                   string                          `xml:"VerticalLoadControl,attr,omitempty"`
	BusVoltageSetPointAttr                    string                          `xml:"BusVoltageSetPoint,attr,omitempty"`
	ConverterMotoringPowerLimitAttr           string                          `xml:"ConverterMotoringPowerLimit,attr,omitempty"`
	ConverterRegenerativePowerLimitAttr       string                          `xml:"ConverterRegenerativePowerLimit,attr,omitempty"`
	ConverterOvertemperatureUserLimitAttr     string                          `xml:"ConverterOvertemperatureUserLimit,attr,omitempty"`
	ConverterThermalOverloadUserLimitAttr     string                          `xml:"ConverterThermalOverloadUserLimit,attr,omitempty"`
	ConverterGroundCurrentUserLimitAttr       string                          `xml:"ConverterGroundCurrentUserLimit,attr,omitempty"`
	ConverterPreChargeOverloadUserLimitAttr   string                          `xml:"ConverterPreChargeOverloadUserLimit,attr,omitempty"`
	TotalDCBusCapacitanceAttr                 string                          `xml:"TotalDCBusCapacitance,attr,omitempty"`
	ExternalDCBusCapacitanceAttr              string                          `xml:"ExternalDCBusCapacitance,attr,omitempty"`
	ConverterModelTimeConstantBaseAttr        string                          `xml:"ConverterModelTimeConstantBase,attr,omitempty"`
	ConverterCurrentLoopBandwidthBaseAttr     string                          `xml:"ConverterCurrentLoopBandwidthBase,attr,omitempty"`
	ConverterRatedCurrentAttr                 string                          `xml:"ConverterRatedCurrent,attr,omitempty"`
	ConverterRatedPeakCurrentAttr             string                          `xml:"ConverterRatedPeakCurrent,attr,omitempty"`
	ConverterRatedVoltageAttr                 string                          `xml:"ConverterRatedVoltage,attr,omitempty"`
	ConverterDCBusCapacitanceAttr             string                          `xml:"ConverterDCBusCapacitance,attr,omitempty"`
	ACLineVoltageTimeConstantAttr             string                          `xml:"ACLineVoltageTimeConstant,attr,omitempty"`
	BusVoltageErrorToleranceAttr              string                          `xml:"BusVoltageErrorTolerance,attr,omitempty"`
	BusVoltageErrorToleranceTimeAttr          string                          `xml:"BusVoltageErrorToleranceTime,attr,omitempty"`
	ReactivePowerSetPointAttr                 string                          `xml:"ReactivePowerSetPoint,attr,omitempty"`
	ReactivePowerRateLimitAttr                string                          `xml:"ReactivePowerRateLimit,attr,omitempty"`
	SystemCapacitanceAttr                     string                          `xml:"SystemCapacitance,attr,omitempty"`
	ConverterCurrentLoopBandwidthAttr         string                          `xml:"ConverterCurrentLoopBandwidth,attr,omitempty"`
	ConverterCurrentIntegratorBandwidthAttr   string                          `xml:"ConverterCurrentIntegratorBandwidth,attr,omitempty"`
	ConverterCurrentVectorLimitAttr           string                          `xml:"ConverterCurrentVectorLimit,attr,omitempty"`
	ACLineVoltageSagThresholdAttr             string                          `xml:"ACLineVoltageSagThreshold,attr,omitempty"`
	ACLineVoltageSagTimeAttr                  string                          `xml:"ACLineVoltageSagTime,attr,omitempty"`
	ConverterInputPhaseLossTimeAttr           string                          `xml:"ConverterInputPhaseLossTime,attr,omitempty"`
	ACLineFrequencyChangeThresholdAttr        string                          `xml:"ACLineFrequencyChangeThreshold,attr,omitempty"`
	ACLineFrequencyChangeTimeAttr             string                          `xml:"ACLineFrequencyChangeTime,attr,omitempty"`
	ACLineSyncLossTimeAttr                    string                          `xml:"ACLineSyncLossTime,attr,omitempty"`
	ACLineOvervoltageUserLimitAttr            string                          `xml:"ACLineOvervoltageUserLimit,attr,omitempty"`
	ACLineUndervoltageUserLimitAttr           string                          `xml:"ACLineUndervoltageUserLimit,attr,omitempty"`
	ACLineOvervoltageUserLimitAlternateAttr   string                          `xml:"ACLineOvervoltageUserLimitAlternate,attr,omitempty"`
	ACLineUndervoltageUserLimitAlternateAttr  string                          `xml:"ACLineUndervoltageUserLimitAlternate,attr,omitempty"`
	ACLineHighFreqUserLimitAttr               string                          `xml:"ACLineHighFreqUserLimit,attr,omitempty"`
	ACLineLowFreqUserLimitAttr                string                          `xml:"ACLineLowFreqUserLimit,attr,omitempty"`
	ACLineHighFreqUserLimitAlternateAttr      string                          `xml:"ACLineHighFreqUserLimitAlternate,attr,omitempty"`
	ACLineLowFreqUserLimitAlternateAttr       string                          `xml:"ACLineLowFreqUserLimitAlternate,attr,omitempty"`
	ConverterHeatsinkOvertempUserLimitAttr    string                          `xml:"ConverterHeatsinkOvertempUserLimit,attr,omitempty"`
	ACLineOverloadUserLimitAttr               string                          `xml:"ACLineOverloadUserLimit,attr,omitempty"`
	ConverterCurrentLoopDampingAttr           string                          `xml:"ConverterCurrentLoopDamping,attr,omitempty"`
	BusObserverBandwidthAttr                  string                          `xml:"BusObserverBandwidth,attr,omitempty"`
	BusObserverIntegratorBandwidthAttr        string                          `xml:"BusObserverIntegratorBandwidth,attr,omitempty"`
	ACLineVoltageUnbalanceLimitAttr           string                          `xml:"ACLineVoltageUnbalanceLimit,attr,omitempty"`
	ACLineCurrentUnbalanceLimitAttr           string                          `xml:"ACLineCurrentUnbalanceLimit,attr,omitempty"`
	ACLineSyncErrorToleranceAttr              string                          `xml:"ACLineSyncErrorTolerance,attr,omitempty"`
	ACLineSourceImpedanceAttr                 string                          `xml:"ACLineSourceImpedance,attr,omitempty"`
	ACLineSourcePowerAttr                     string                          `xml:"ACLineSourcePower,attr,omitempty"`
	ACLineSourceImpedanceAlternateAttr        string                          `xml:"ACLineSourceImpedanceAlternate,attr,omitempty"`
	ACLineSourcePowerAlternateAttr            string                          `xml:"ACLineSourcePowerAlternate,attr,omitempty"`
	BusVoltageLoopBandwidthAttr               string                          `xml:"BusVoltageLoopBandwidth,attr,omitempty"`
	BusVoltageIntegratorBandwidthAttr         string                          `xml:"BusVoltageIntegratorBandwidth,attr,omitempty"`
	BusVoltageRateLimitAttr                   string                          `xml:"BusVoltageRateLimit,attr,omitempty"`
	ActiveCurrentCommandAttr                  string                          `xml:"ActiveCurrentCommand,attr,omitempty"`
	ReactiveCurrentCommandAttr                string                          `xml:"ReactiveCurrentCommand,attr,omitempty"`
	ActiveCurrentTrimAttr                     string                          `xml:"ActiveCurrentTrim,attr,omitempty"`
	ActiveCurrentLowPassFilterBandwidthAttr   string                          `xml:"ActiveCurrentLowPassFilterBandwidth,attr,omitempty"`
	ActiveCurrentNotchFilterFrequencyAttr     string                          `xml:"ActiveCurrentNotchFilterFrequency,attr,omitempty"`
	ActiveCurrentRateLimitAttr                string                          `xml:"ActiveCurrentRateLimit,attr,omitempty"`
	ReactiveCurrentRateLimitAttr              string                          `xml:"ReactiveCurrentRateLimit,attr,omitempty"`
	ConverterControlModeAttr                  string                          `xml:"ConverterControlMode,attr,omitempty"`
	ReactivePowerControlAttr                  string                          `xml:"ReactivePowerControl,attr,omitempty"`
	ConverterStartupMethodAttr                string                          `xml:"ConverterStartupMethod,attr,omitempty"`
	ACLineVoltageSagActionAttr                string                          `xml:"ACLineVoltageSagAction,attr,omitempty"`
	ConverterInputPhaseLossActionAttr         string                          `xml:"ConverterInputPhaseLossAction,attr,omitempty"`
	ACLineFrequencyChangeActionAttr           string                          `xml:"ACLineFrequencyChangeAction,attr,omitempty"`
	ACLineSyncLossActionAttr                  string                          `xml:"ACLineSyncLossAction,attr,omitempty"`
	ConverterOverloadActionAttr               string                          `xml:"ConverterOverloadAction,attr,omitempty"`
	ConverterCurrentLoopTuningMethodAttr      string                          `xml:"ConverterCurrentLoopTuningMethod,attr,omitempty"`
	ConverterACInputPhasingAttr               string                          `xml:"ConverterACInputPhasing,attr,omitempty"`
	ConverterACInputFrequencyAttr             string                          `xml:"ConverterACInputFrequency,attr,omitempty"`
	ACLineSourceSelectAttr                    string                          `xml:"ACLineSourceSelect,attr,omitempty"`
	BusVoltageReferenceSourceAttr             string                          `xml:"BusVoltageReferenceSource,attr,omitempty"`
	BusObserverConfigurationAttr              string                          `xml:"BusObserverConfiguration,attr,omitempty"`
	ConverterACInputVoltageAttr               string                          `xml:"ConverterACInputVoltage,attr,omitempty"`
	CIPAxisExceptionAction2Attr               string                          `xml:"CIPAxisExceptionAction2,attr,omitempty"`
	CIPAxisExceptionAction2RAAttr             string                          `xml:"CIPAxisExceptionAction2RA,attr,omitempty"`
	ConverterConfigurationAttr                string                          `xml:"ConverterConfiguration,attr,omitempty"`
	ConverterCurrentLimitSourceAttr           string                          `xml:"ConverterCurrentLimitSource,attr,omitempty"`
	SafetyFaultActionAttr                     string                          `xml:"SafetyFaultAction,attr,omitempty"`
	ACLineContactorInputCheckingAttr          string                          `xml:"ACLineContactorInputChecking,attr,omitempty"`
	ConverterModelTimeConstantAttr            string                          `xml:"ConverterModelTimeConstant,attr,omitempty"`
	ConverterRatedPowerAttr                   string                          `xml:"ConverterRatedPower,attr,omitempty"`
	CurrentLoopBandwidthScalingFactorAttr     string                          `xml:"CurrentLoopBandwidthScalingFactor,attr,omitempty"`
	CurrentLoopBandwidthAttr                  string                          `xml:"CurrentLoopBandwidth,attr,omitempty"`
	DriveRatedVoltageAttr                     string                          `xml:"DriveRatedVoltage,attr,omitempty"`
	MaxOutputFrequencyAttr                    string                          `xml:"MaxOutputFrequency,attr,omitempty"`
	ACLineResonanceUserLimitAttr              string                          `xml:"ACLineResonanceUserLimit,attr,omitempty"`
	MotorTestDataValidAttr                    string                          `xml:"MotorTestDataValid,attr,omitempty"`
	TrackMoverLengthAttr                      string                          `xml:"TrackMoverLength,attr,omitempty"`
	TrackMoverCenterofMassOffsetAttr          string                          `xml:"TrackMoverCenterofMassOffset,attr,omitempty"`
	Feedback1CalibrationOffsetAttr            string                          `xml:"Feedback1CalibrationOffset,attr,omitempty"`
	Feedback2CalibrationOffsetAttr            string                          `xml:"Feedback2CalibrationOffset,attr,omitempty"`
	CommandNotchFilterFrequencyAttr           string                          `xml:"CommandNotchFilterFrequency,attr,omitempty"`
	CommandNotchFilterWidthAttr               string                          `xml:"CommandNotchFilterWidth,attr,omitempty"`
	CommandNotchFilterDepthAttr               string                          `xml:"CommandNotchFilterDepth,attr,omitempty"`
	CommandNotchFilterGainAttr                string                          `xml:"CommandNotchFilterGain,attr,omitempty"`
	CommandNotchFilter2FrequencyAttr          string                          `xml:"CommandNotchFilter2Frequency,attr,omitempty"`
	CommandNotchFilter2WidthAttr              string                          `xml:"CommandNotchFilter2Width,attr,omitempty"`
	CommandNotchFilter2DepthAttr              string                          `xml:"CommandNotchFilter2Depth,attr,omitempty"`
	CommandNotchFilter2GainAttr               string                          `xml:"CommandNotchFilter2Gain,attr,omitempty"`
	TorqueNotchFilter2FrequencyAttr           string                          `xml:"TorqueNotchFilter2Frequency,attr,omitempty"`
	TorqueNotchFilter3FrequencyAttr           string                          `xml:"TorqueNotchFilter3Frequency,attr,omitempty"`
	TorqueNotchFilter4FrequencyAttr           string                          `xml:"TorqueNotchFilter4Frequency,attr,omitempty"`
	TorqueNotchFilterWidthAttr                string                          `xml:"TorqueNotchFilterWidth,attr,omitempty"`
	TorqueNotchFilterDepthAttr                string                          `xml:"TorqueNotchFilterDepth,attr,omitempty"`
	TorqueNotchFilterGainAttr                 string                          `xml:"TorqueNotchFilterGain,attr,omitempty"`
	TorqueNotchFilter2WidthAttr               string                          `xml:"TorqueNotchFilter2Width,attr,omitempty"`
	TorqueNotchFilter2DepthAttr               string                          `xml:"TorqueNotchFilter2Depth,attr,omitempty"`
	TorqueNotchFilter2GainAttr                string                          `xml:"TorqueNotchFilter2Gain,attr,omitempty"`
	TorqueNotchFilter3WidthAttr               string                          `xml:"TorqueNotchFilter3Width,attr,omitempty"`
	TorqueNotchFilter3DepthAttr               string                          `xml:"TorqueNotchFilter3Depth,attr,omitempty"`
	TorqueNotchFilter3GainAttr                string                          `xml:"TorqueNotchFilter3Gain,attr,omitempty"`
	TorqueNotchFilter4WidthAttr               string                          `xml:"TorqueNotchFilter4Width,attr,omitempty"`
	TorqueNotchFilter4DepthAttr               string                          `xml:"TorqueNotchFilter4Depth,attr,omitempty"`
	TorqueNotchFilter4GainAttr                string                          `xml:"TorqueNotchFilter4Gain,attr,omitempty"`
	AdaptiveTuningTrackingNotchFiltersAttr    string                          `xml:"AdaptiveTuningTrackingNotchFilters,attr,omitempty"`
	AdaptiveTuningGainScalingFactorMinAttr    string                          `xml:"AdaptiveTuningGainScalingFactorMin,attr,omitempty"`
	TorqueLowPassFilterBandwidthMinAttr       string                          `xml:"TorqueLowPassFilterBandwidthMin,attr,omitempty"`
	TorqueNotchFilterWidthMinAttr             string                          `xml:"TorqueNotchFilterWidthMin,attr,omitempty"`
	TorqueNotchFilterWidthMaxAttr             string                          `xml:"TorqueNotchFilterWidthMax,attr,omitempty"`
	MotorTestTravelLimitAttr                  string                          `xml:"MotorTestTravelLimit,attr,omitempty"`
	MotorTestSpeedAttr                        string                          `xml:"MotorTestSpeed,attr,omitempty"`
	MotorTestTorqueAttr                       string                          `xml:"MotorTestTorque,attr,omitempty"`
	HookupTestSpeedAttr                       string                          `xml:"HookupTestSpeed,attr,omitempty"`
	TorqueEstimateCrossoverSpeedAttr          string                          `xml:"TorqueEstimateCrossoverSpeed,attr,omitempty"`
	TorqueEstimateNotch1FrequencyAttr         string                          `xml:"TorqueEstimateNotch1Frequency,attr,omitempty"`
	TorqueEstimateNotch1WidthAttr             string                          `xml:"TorqueEstimateNotch1Width,attr,omitempty"`
	TorqueEstimateNotch1DepthAttr             string                          `xml:"TorqueEstimateNotch1Depth,attr,omitempty"`
	TorqueEstimateNotch1GainAttr              string                          `xml:"TorqueEstimateNotch1Gain,attr,omitempty"`
	TorqueEstimateNotch2FrequencyAttr         string                          `xml:"TorqueEstimateNotch2Frequency,attr,omitempty"`
	TorqueEstimateNotch2WidthAttr             string                          `xml:"TorqueEstimateNotch2Width,attr,omitempty"`
	TorqueEstimateNotch2DepthAttr             string                          `xml:"TorqueEstimateNotch2Depth,attr,omitempty"`
	TorqueEstimateNotch2GainAttr              string                          `xml:"TorqueEstimateNotch2Gain,attr,omitempty"`
	FlyingStartCEMFReconnectDelayAttr         string                          `xml:"FlyingStartCEMFReconnectDelay,attr,omitempty"`
	FlyingStartCEMFCurrentRegKpAttr           string                          `xml:"FlyingStartCEMFCurrentRegKp,attr,omitempty"`
	FlyingStartCEMFCurrentRegKiAttr           string                          `xml:"FlyingStartCEMFCurrentRegKi,attr,omitempty"`
	FlyingStartCEMFVelocityRegKpAttr          string                          `xml:"FlyingStartCEMFVelocityRegKp,attr,omitempty"`
	FlyingStartCEMFVelocityRegKiAttr          string                          `xml:"FlyingStartCEMFVelocityRegKi,attr,omitempty"`
	FlyingStartCEMFExcitationRegKpAttr        string                          `xml:"FlyingStartCEMFExcitationRegKp,attr,omitempty"`
	FlyingStartCEMFExcitationRegKiAttr        string                          `xml:"FlyingStartCEMFExcitationRegKi,attr,omitempty"`
	FlyingStartCEMFBrakeLevelAttr             string                          `xml:"FlyingStartCEMFBrakeLevel,attr,omitempty"`
	FlyingStartCEMFBrakeTimeAttr              string                          `xml:"FlyingStartCEMFBrakeTime,attr,omitempty"`
	FlyingStartCEMFZeroSpeedThresholdAttr     string                          `xml:"FlyingStartCEMFZeroSpeedThreshold,attr,omitempty"`
	FlyingStartSweepReconnectDelayAttr        string                          `xml:"FlyingStartSweepReconnectDelay,attr,omitempty"`
	FlyingStartSweepInitialVoltageRegKpAttr   string                          `xml:"FlyingStartSweepInitialVoltageRegKp,attr,omitempty"`
	FlyingStartSweepInitialVoltageRegKiAttr   string                          `xml:"FlyingStartSweepInitialVoltageRegKi,attr,omitempty"`
	FlyingStartSweepTimeAttr                  string                          `xml:"FlyingStartSweepTime,attr,omitempty"`
	FlyingStartSweepVHzDCBoostAdjustAttr      string                          `xml:"FlyingStartSweepVHzDCBoostAdjust,attr,omitempty"`
	FlyingStartSweepVHzRatioAttr              string                          `xml:"FlyingStartSweepVHzRatio,attr,omitempty"`
	FlyingStartSweepSpeedDetectLevelAttr      string                          `xml:"FlyingStartSweepSpeedDetectLevel,attr,omitempty"`
	FlyingStartSweepSpeedDetectTimeAttr       string                          `xml:"FlyingStartSweepSpeedDetectTime,attr,omitempty"`
	FlyingStartSweepRecoveryCurrentRegKiAttr  string                          `xml:"FlyingStartSweepRecoveryCurrentRegKi,attr,omitempty"`
	FlyingStartSweepVelocityRegKpAttr         string                          `xml:"FlyingStartSweepVelocityRegKp,attr,omitempty"`
	FlyingStartSweepVelocityRegKiAttr         string                          `xml:"FlyingStartSweepVelocityRegKi,attr,omitempty"`
	FlyingStartSweepBrakeLevelAttr            string                          `xml:"FlyingStartSweepBrakeLevel,attr,omitempty"`
	FlyingStartSweepBrakeTimeAttr             string                          `xml:"FlyingStartSweepBrakeTime,attr,omitempty"`
	FlyingStartSweepZeroSpeedThresholdAttr    string                          `xml:"FlyingStartSweepZeroSpeedThreshold,attr,omitempty"`
	FrictionCompensationMethodAttr            string                          `xml:"FrictionCompensationMethod,attr,omitempty"`
	FrictionCompensationStartSpeedAttr        string                          `xml:"FrictionCompensationStartSpeed,attr,omitempty"`
	FrictionCompensationHysteresisAttr        string                          `xml:"FrictionCompensationHysteresis,attr,omitempty"`
	FrictionCompensationBreakawayTimeAttr     string                          `xml:"FrictionCompensationBreakawayTime,attr,omitempty"`
	MotorAdaptionSpeedAttr                    string                          `xml:"MotorAdaptionSpeed,attr,omitempty"`
	TorqueAdaptionEnableAttr                  string                          `xml:"TorqueAdaptionEnable,attr,omitempty"`
	SlipAdaptionEnableAttr                    string                          `xml:"SlipAdaptionEnable,attr,omitempty"`
	FluxAdaptionEnableAttr                    string                          `xml:"FluxAdaptionEnable,attr,omitempty"`
	TorqueAdaptionRegulatorKpAttr             string                          `xml:"TorqueAdaptionRegulatorKp,attr,omitempty"`
	TorqueAdaptionRegulatorKiAttr             string                          `xml:"TorqueAdaptionRegulatorKi,attr,omitempty"`
	TorqueAdaptionRegulatorLimitPositiveAttr  string                          `xml:"TorqueAdaptionRegulatorLimitPositive,attr,omitempty"`
	TorqueAdaptionRegulatorLimitNegativeAttr  string                          `xml:"TorqueAdaptionRegulatorLimitNegative,attr,omitempty"`
	SlipandFluxRegulatorSlewTimeAttr          string                          `xml:"SlipandFluxRegulatorSlewTime,attr,omitempty"`
	SlipandFluxRegulatorSlewRateAttr          string                          `xml:"SlipandFluxRegulatorSlewRate,attr,omitempty"`
	SlipandFluxRegulatorConvergenceLevelAttr  string                          `xml:"SlipandFluxRegulatorConvergenceLevel,attr,omitempty"`
	SlipandFluxRegulatorConvergenceTimeAttr   string                          `xml:"SlipandFluxRegulatorConvergenceTime,attr,omitempty"`
	SlipAdaptionRegulatorIqThresholdAttr      string                          `xml:"SlipAdaptionRegulatorIqThreshold,attr,omitempty"`
	SlipAdaptionRegulatorKpAttr               string                          `xml:"SlipAdaptionRegulatorKp,attr,omitempty"`
	SlipAdaptionRegulatorKiAttr               string                          `xml:"SlipAdaptionRegulatorKi,attr,omitempty"`
	FluxAdaptionRegulatorKpAttr               string                          `xml:"FluxAdaptionRegulatorKp,attr,omitempty"`
	FluxAdaptionRegulatorKiAttr               string                          `xml:"FluxAdaptionRegulatorKi,attr,omitempty"`
	VqIdDecouplingGainAttr                    string                          `xml:"VqIdDecouplingGain,attr,omitempty"`
	VdIqDecouplingGainAttr                    string                          `xml:"VdIqDecouplingGain,attr,omitempty"`
	LqIqFeedbackFilterBandwidthAttr           string                          `xml:"LqIqFeedbackFilterBandwidth,attr,omitempty"`
	CurrentLimitRegulatorKpAttr               string                          `xml:"CurrentLimitRegulatorKp,attr,omitempty"`
	CurrentLimitRegulatorKiAttr               string                          `xml:"CurrentLimitRegulatorKi,attr,omitempty"`
	CurrentLimitRegulatorKdAttr               string                          `xml:"CurrentLimitRegulatorKd,attr,omitempty"`
	LowFrequencyIdCurrentLimitRegulatorKpAttr string                          `xml:"LowFrequencyIdCurrentLimitRegulatorKp,attr,omitempty"`
	LowFrequencyIqCurrentLimitRegulatorKpAttr string                          `xml:"LowFrequencyIqCurrentLimitRegulatorKp,attr,omitempty"`
	PreChargeHoldControlAttr                  string                          `xml:"PreChargeHoldControl,attr,omitempty"`
	PreChargeHoldDelayAttr                    string                          `xml:"PreChargeHoldDelay,attr,omitempty"`
	MotorOvertemperatureUserLimitAttr         string                          `xml:"MotorOvertemperatureUserLimit,attr,omitempty"`
	InverterOvertemperatureUserLimitAttr      string                          `xml:"InverterOvertemperatureUserLimit,attr,omitempty"`
	SensorlessVectorEconomyAccelDecelKpAttr   string                          `xml:"SensorlessVectorEconomyAccelDecelKp,attr,omitempty"`
	SensorlessVectorEconomyAccelDecelKiAttr   string                          `xml:"SensorlessVectorEconomyAccelDecelKi,attr,omitempty"`
	SensorlessVectorEconomyAtSpeedKiAttr      string                          `xml:"SensorlessVectorEconomyAtSpeedKi,attr,omitempty"`
	SensorlessVectorBoostFilterBandwidthAttr  string                          `xml:"SensorlessVectorBoostFilterBandwidth,attr,omitempty"`
	VelocityFeedbackDelayCompensationAttr     string                          `xml:"VelocityFeedbackDelayCompensation,attr,omitempty"`
	TorqueCalibrationFactorMotoringAttr       string                          `xml:"TorqueCalibrationFactorMotoring,attr,omitempty"`
	TorqueCalibrationFactorRegeneratingAttr   string                          `xml:"TorqueCalibrationFactorRegenerating,attr,omitempty"`
	FluxVectorFrequencyRegulatorKpAttr        string                          `xml:"FluxVectorFrequencyRegulatorKp,attr,omitempty"`
	FluxVectorFrequencyRegulatorKiAttr        string                          `xml:"FluxVectorFrequencyRegulatorKi,attr,omitempty"`
	MotorStabilityControlEnableAttr           string                          `xml:"MotorStabilityControlEnable,attr,omitempty"`
	MotorStabilityControlFilterBandwidthAttr  string                          `xml:"MotorStabilityControlFilterBandwidth,attr,omitempty"`
	MotorStabilityControlVoltageGainAttr      string                          `xml:"MotorStabilityControlVoltageGain,attr,omitempty"`
	MotorStabilityControlFrequencyGainAttr    string                          `xml:"MotorStabilityControlFrequencyGain,attr,omitempty"`
	PowerDeviceCompensationEnableAttr         string                          `xml:"PowerDeviceCompensationEnable,attr,omitempty"`
	PowerDeviceDeadTimeCompensationAttr       string                          `xml:"PowerDeviceDeadTimeCompensation,attr,omitempty"`
	DCInjectionCurrentRegulatorKpAttr         string                          `xml:"DCInjectionCurrentRegulatorKp,attr,omitempty"`
	DCInjectionCurrentRegulatorKiAttr         string                          `xml:"DCInjectionCurrentRegulatorKi,attr,omitempty"`
	FluxBrakingRegulatorKpAttr                string                          `xml:"FluxBrakingRegulatorKp,attr,omitempty"`
	FluxBrakingRegulatorKiAttr                string                          `xml:"FluxBrakingRegulatorKi,attr,omitempty"`
	FluxBrakingVoltageLimitAttr               string                          `xml:"FluxBrakingVoltageLimit,attr,omitempty"`
	FluxDownRegulatorKpAttr                   string                          `xml:"FluxDownRegulatorKp,attr,omitempty"`
	FluxDownRegulatorKiAttr                   string                          `xml:"FluxDownRegulatorKi,attr,omitempty"`
	ACInjectionBrakeRegulatorKpAttr           string                          `xml:"ACInjectionBrakeRegulatorKp,attr,omitempty"`
	ACInjectionBrakeRegulatorKiAttr           string                          `xml:"ACInjectionBrakeRegulatorKi,attr,omitempty"`
	ACInjectionBrakePowerThresholdAttr        string                          `xml:"ACInjectionBrakePowerThreshold,attr,omitempty"`
	ACInjectionBrakeFrequencyThresholdAttr    string                          `xml:"ACInjectionBrakeFrequencyThreshold,attr,omitempty"`
	BusRegulatorVoltageLevelAttr              string                          `xml:"BusRegulatorVoltageLevel,attr,omitempty"`
	BusRegulatorSetPointSourceAttr            string                          `xml:"BusRegulatorSetPointSource,attr,omitempty"`
	BusRegulatorKpAttr                        string                          `xml:"BusRegulatorKp,attr,omitempty"`
	BusRegulatorKiAttr                        string                          `xml:"BusRegulatorKi,attr,omitempty"`
	BusLimitRegulatorKpAttr                   string                          `xml:"BusLimitRegulatorKp,attr,omitempty"`
	BusLimitRegulatorKdAttr                   string                          `xml:"BusLimitRegulatorKd,attr,omitempty"`
	BusLimitActiveCurrentRegulatorKpAttr      string                          `xml:"BusLimitActiveCurrentRegulatorKp,attr,omitempty"`
	BusLimitActiveCurrentRegulatorKiAttr      string                          `xml:"BusLimitActiveCurrentRegulatorKi,attr,omitempty"`
	InverterGroundCurrentUserLimitAttr        string                          `xml:"InverterGroundCurrentUserLimit,attr,omitempty"`
	TestModeConfigurationAttr                 string                          `xml:"TestModeConfiguration,attr,omitempty"`
	TestModeEnableAttr                        string                          `xml:"TestModeEnable,attr,omitempty"`
	UseAttr                                   string                          `xml:"Use,attr,omitempty"`
	CustomProperties                          *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// TrendGroupCollectionType ...
type TrendGroupCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Trend            []*TrendGroupType               `xml:"Trend"`
}

// Template ...
type Template struct {
	Value string `xml:",chardata"`
}

// TrendGroupType ...
type TrendGroupType struct {
	NameAttr                         string                          `xml:"Name,attr"`
	UIdAttr                          string                          `xml:"UId,attr,omitempty"`
	SamplePeriodAttr                 string                          `xml:"SamplePeriod,attr,omitempty"`
	NumberOfCapturesAttr             string                          `xml:"NumberOfCaptures,attr,omitempty"`
	CaptureSizeTypeAttr              string                          `xml:"CaptureSizeType,attr,omitempty"`
	CaptureSizeAttr                  string                          `xml:"CaptureSize,attr,omitempty"`
	StartTriggerTypeAttr             string                          `xml:"StartTriggerType,attr,omitempty"`
	StartTriggerTag1Attr             string                          `xml:"StartTriggerTag1,attr,omitempty"`
	StartTriggerOperation1Attr       string                          `xml:"StartTriggerOperation1,attr,omitempty"`
	StartTriggerTargetType1Attr      string                          `xml:"StartTriggerTargetType1,attr,omitempty"`
	StartTriggerTargetValue1Attr     string                          `xml:"StartTriggerTargetValue1,attr,omitempty"`
	StartTriggerTargetTag1Attr       string                          `xml:"StartTriggerTargetTag1,attr,omitempty"`
	StartTriggerLogicalOperationAttr string                          `xml:"StartTriggerLogicalOperation,attr,omitempty"`
	StartTriggerTag2Attr             string                          `xml:"StartTriggerTag2,attr,omitempty"`
	StartTriggerOperation2Attr       string                          `xml:"StartTriggerOperation2,attr,omitempty"`
	StartTriggerTargetType2Attr      string                          `xml:"StartTriggerTargetType2,attr,omitempty"`
	StartTriggerTargetValue2Attr     string                          `xml:"StartTriggerTargetValue2,attr,omitempty"`
	StartTriggerTargetTag2Attr       string                          `xml:"StartTriggerTargetTag2,attr,omitempty"`
	PreSampleTypeAttr                string                          `xml:"PreSampleType,attr,omitempty"`
	PreSamplesAttr                   string                          `xml:"PreSamples,attr,omitempty"`
	StopTriggerTypeAttr              string                          `xml:"StopTriggerType,attr,omitempty"`
	StopTriggerTag1Attr              string                          `xml:"StopTriggerTag1,attr,omitempty"`
	StopTriggerOperation1Attr        string                          `xml:"StopTriggerOperation1,attr,omitempty"`
	StopTriggerTargetType1Attr       string                          `xml:"StopTriggerTargetType1,attr,omitempty"`
	StopTriggerTargetValue1Attr      string                          `xml:"StopTriggerTargetValue1,attr,omitempty"`
	StopTriggerTargetTag1Attr        string                          `xml:"StopTriggerTargetTag1,attr,omitempty"`
	StopTriggerLogicalOperationAttr  string                          `xml:"StopTriggerLogicalOperation,attr,omitempty"`
	StopTriggerTag2Attr              string                          `xml:"StopTriggerTag2,attr,omitempty"`
	StopTriggerOperation2Attr        string                          `xml:"StopTriggerOperation2,attr,omitempty"`
	StopTriggerTargetType2Attr       string                          `xml:"StopTriggerTargetType2,attr,omitempty"`
	StopTriggerTargetValue2Attr      string                          `xml:"StopTriggerTargetValue2,attr,omitempty"`
	StopTriggerTargetTag2Attr        string                          `xml:"StopTriggerTargetTag2,attr,omitempty"`
	PostSampleTypeAttr               string                          `xml:"PostSampleType,attr,omitempty"`
	PostSamplesAttr                  string                          `xml:"PostSamples,attr,omitempty"`
	TrendxVersionAttr                string                          `xml:"TrendxVersion,attr,omitempty"`
	UseAttr                          string                          `xml:"Use,attr,omitempty"`
	CustomProperties                 *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description                      *DescriptionType                `xml:"Description"`
	Template                         *Template                       `xml:"Template"`
	Pens                             *PenCollectionType              `xml:"Pens"`
}

// PenCollectionType ...
type PenCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Pen              []*PenType                      `xml:"Pen"`
}

// PenType ...
type PenType struct {
	NameAttr         string                          `xml:"Name,attr"`
	ColorAttr        string                          `xml:"Color,attr,omitempty"`
	VisibleAttr      string                          `xml:"Visible,attr,omitempty"`
	WidthAttr        string                          `xml:"Width,attr,omitempty"`
	TypeAttr         string                          `xml:"Type,attr,omitempty"`
	StyleAttr        string                          `xml:"Style,attr,omitempty"`
	MarkerAttr       string                          `xml:"Marker,attr,omitempty"`
	MinAttr          string                          `xml:"Min,attr,omitempty"`
	MaxAttr          string                          `xml:"Max,attr,omitempty"`
	EngUnitsAttr     string                          `xml:"EngUnits,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description      *DescriptionType                `xml:"Description"`
}

// ParameterConnectionCollectionType ...
type ParameterConnectionCollectionType struct {
	UseAttr             string                          `xml:"Use,attr,omitempty"`
	CustomProperties    *CustomPropertiesCollectionType `xml:"CustomProperties"`
	ParameterConnection []*ParameterConnectionType      `xml:"ParameterConnection"`
}

// ParameterConnectionType ...
type ParameterConnectionType struct {
	EndPoint1Attr    string                          `xml:"EndPoint1,attr"`
	EndPoint2Attr    string                          `xml:"EndPoint2,attr"`
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// QuickWatchCollectionAdaptorType ...
type QuickWatchCollectionAdaptorType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	QuickWatchList   []*QuickWatchAdaptorType        `xml:"QuickWatchList"`
}

// QuickWatchAdaptorType ...
type QuickWatchAdaptorType struct {
	NameAttr         string                          `xml:"Name,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	WatchTag         []*QuickWatchTagAdaptorType     `xml:"WatchTag"`
}

// QuickWatchTagAdaptorType ...
type QuickWatchTagAdaptorType struct {
	SpecifierAttr    string                          `xml:"Specifier,attr"`
	ScopeAttr        string                          `xml:"Scope,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// TimeSynchronizeType ...
type TimeSynchronizeType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	Priority1Attr    string                          `xml:"Priority1,attr,omitempty"`
	Priority2Attr    string                          `xml:"Priority2,attr,omitempty"`
	PTPEnableAttr    string                          `xml:"PTPEnable,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// TCPIPType ...
type TCPIPType struct {
	ConfigTypeAttr   string                          `xml:"ConfigType,attr,omitempty"`
	IPAddressAttr    string                          `xml:"IPAddress,attr,omitempty"`
	SubnetMaskAttr   string                          `xml:"SubnetMask,attr,omitempty"`
	GatewayAttr      string                          `xml:"Gateway,attr,omitempty"`
	PrimaryDNSAttr   string                          `xml:"PrimaryDNS,attr,omitempty"`
	SecondaryDNSAttr string                          `xml:"SecondaryDNS,attr,omitempty"`
	DomainNameAttr   string                          `xml:"DomainName,attr,omitempty"`
	HostNameAttr     string                          `xml:"HostName,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// EthernetLinkCollectionType ...
type EthernetLinkCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	EthernetPort     []*EthernetLinkType             `xml:"EthernetPort"`
}

// EthernetLinkType ...
type EthernetLinkType struct {
	PortAttr                 string                          `xml:"Port,attr,omitempty"`
	LabelAttr                string                          `xml:"Label,attr,omitempty"`
	PortEnabledAttr          string                          `xml:"PortEnabled,attr,omitempty"`
	AutoNegotiateEnabledAttr string                          `xml:"AutoNegotiateEnabled,attr,omitempty"`
	InterfaceSpeedAttr       string                          `xml:"InterfaceSpeed,attr,omitempty"`
	DuplexModeAttr           string                          `xml:"DuplexMode,attr,omitempty"`
	UseAttr                  string                          `xml:"Use,attr,omitempty"`
	CustomProperties         *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// DeviceLevelRingType ...
type DeviceLevelRingType struct {
	SupervisorModeEnabledAttr string                          `xml:"SupervisorModeEnabled,attr,omitempty"`
	SupervisorPrecedenceAttr  string                          `xml:"SupervisorPrecedence,attr,omitempty"`
	BeaconIntervalAttr        string                          `xml:"BeaconInterval,attr,omitempty"`
	BeaconTimeoutAttr         string                          `xml:"BeaconTimeout,attr,omitempty"`
	VLANIDAttr                string                          `xml:"VLANID,attr,omitempty"`
	UseAttr                   string                          `xml:"Use,attr,omitempty"`
	CustomProperties          *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// DataLogCollectionType ...
type DataLogCollectionType struct {
	UIdAttr          string                          `xml:"UId,attr,omitempty"`
	ParentUIdAttr    string                          `xml:"ParentUId,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	DataLog          []*DataLogType                  `xml:"DataLog"`
}

// DataLogType ...
type DataLogType struct {
	NameAttr                        string                          `xml:"Name,attr"`
	UIdAttr                         string                          `xml:"UId,attr,omitempty"`
	PeriodAttr                      string                          `xml:"Period,attr,omitempty"`
	PeriodUnitAttr                  string                          `xml:"PeriodUnit,attr,omitempty"`
	EnableDataLoggingOnDownloadAttr string                          `xml:"EnableDataLoggingOnDownload,attr,omitempty"`
	CollectDataOnlyInRunModeAttr    string                          `xml:"CollectDataOnlyInRunMode,attr,omitempty"`
	DataCollectionTypeAttr          string                          `xml:"DataCollectionType,attr,omitempty"`
	SampleEventAttr                 string                          `xml:"SampleEvent,attr,omitempty"`
	StartTriggerConditionAttr       string                          `xml:"StartTriggerCondition,attr,omitempty"`
	TriggerConditionAttr            string                          `xml:"TriggerCondition,attr,omitempty"`
	PreSamplesAttr                  string                          `xml:"PreSamples,attr,omitempty"`
	PreSampleUnitAttr               string                          `xml:"PreSampleUnit,attr,omitempty"`
	StopTriggerConditionAttr        string                          `xml:"StopTriggerCondition,attr,omitempty"`
	PostSamplesAttr                 string                          `xml:"PostSamples,attr,omitempty"`
	PostSampleUnitAttr              string                          `xml:"PostSampleUnit,attr,omitempty"`
	CaptureSizeAttr                 string                          `xml:"CaptureSize,attr,omitempty"`
	CaptureSizeUnitAttr             string                          `xml:"CaptureSizeUnit,attr,omitempty"`
	CaptureSizeExceededAttr         string                          `xml:"CaptureSizeExceeded,attr,omitempty"`
	DataCapturesToKeepAttr          string                          `xml:"DataCapturesToKeep,attr,omitempty"`
	DataLogCapturesExceededAttr     string                          `xml:"DataLogCapturesExceeded,attr,omitempty"`
	NumberOfDataLogTagsAttr         string                          `xml:"NumberOfDataLogTags,attr,omitempty"`
	DataLogSizeAttr                 string                          `xml:"DataLogSize,attr,omitempty"`
	DataLogSizeUnitAttr             string                          `xml:"DataLogSizeUnit,attr,omitempty"`
	DataLogSizeExceededAttr         string                          `xml:"DataLogSizeExceeded,attr,omitempty"`
	UseAttr                         string                          `xml:"Use,attr,omitempty"`
	CustomProperties                *CustomPropertiesCollectionType `xml:"CustomProperties"`
	Description                     *DescriptionType                `xml:"Description"`
	StartTriggerExpression          string                          `xml:"StartTriggerExpression"`
	TriggerExpression               string                          `xml:"TriggerExpression"`
	StopTriggerExpression           string                          `xml:"StopTriggerExpression"`
	DataLogTags                     *DataLogTagCollectionType       `xml:"DataLogTags"`
}

// DataLogTagCollectionType ...
type DataLogTagCollectionType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	DataLogTag       []*DataLogTagType               `xml:"DataLogTag"`
}

// DataLogTagType ...
type DataLogTagType struct {
	NameAttr         string                          `xml:"Name,attr"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// EncryptionInfoType ...
type EncryptionInfoType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
	EncryptionKey    []*EncryptionKeyType            `xml:"EncryptionKey"`
}

// EncryptionKeyType ...
type EncryptionKeyType struct {
	NameAttr         string                          `xml:"Name,attr,omitempty"`
	IDAttr           string                          `xml:"ID,attr,omitempty"`
	DescriptionAttr  string                          `xml:"Description,attr,omitempty"`
	VendorAttr       string                          `xml:"Vendor,attr,omitempty"`
	PublicKeyAttr    string                          `xml:"PublicKey,attr,omitempty"`
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// EncryptedContentType ...
type EncryptedContentType struct {
	EncryptedTypeAttr  string                          `xml:"EncryptedType,attr"`
	OnlineEditTypeAttr string                          `xml:"OnlineEditType,attr,omitempty"`
	UseAttr            string                          `xml:"Use,attr,omitempty"`
	CustomProperties   *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// HardwareStatusType ...
type HardwareStatusType struct {
	UseAttr          string                          `xml:"Use,attr,omitempty"`
	CustomProperties *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// BEOType ...
type BEOType struct {
	EnergyResourceTypeAttr           string                          `xml:"EnergyResourceType,attr"`
	BaseEnergyObjectCapabilitiesAttr string                          `xml:"BaseEnergyObjectCapabilities,attr"`
	EnergyAccuracyAttr               string                          `xml:"EnergyAccuracy,attr,omitempty"`
	DataStatusAttr                   string                          `xml:"DataStatus,attr,omitempty"`
	ConsumedEnergyOdometerAttr       string                          `xml:"ConsumedEnergyOdometer,attr,omitempty"`
	GeneratedEnergyOdometerAttr      string                          `xml:"GeneratedEnergyOdometer,attr,omitempty"`
	NetEnergyOdometerAttr            string                          `xml:"NetEnergyOdometer,attr,omitempty"`
	EnergyTransferRateAttr           string                          `xml:"EnergyTransferRate,attr,omitempty"`
	ExtendedDataStatusAttr           string                          `xml:"ExtendedDataStatus,attr,omitempty"`
	UseAttr                          string                          `xml:"Use,attr,omitempty"`
	CustomProperties                 *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// EEOType ...
type EEOType struct {
	RealEnergyConsumedOdometerAttr  string                          `xml:"RealEnergyConsumedOdometer,attr,omitempty"`
	RealEnergyGeneratedOdometerAttr string                          `xml:"RealEnergyGeneratedOdometer,attr,omitempty"`
	RealEnergyNetOdometerAttr       string                          `xml:"RealEnergyNetOdometer,attr,omitempty"`
	TotalRealPowerAttr              string                          `xml:"TotalRealPower,attr,omitempty"`
	ThreePhaseTruePowerFactorAttr   string                          `xml:"ThreePhaseTruePowerFactor,attr,omitempty"`
	PhaseRotationAttr               string                          `xml:"PhaseRotation,attr,omitempty"`
	UseAttr                         string                          `xml:"Use,attr,omitempty"`
	CustomProperties                *CustomPropertiesCollectionType `xml:"CustomProperties"`
}

// DependenciesWideType ...
type DependenciesWideType struct {
	DecoratedDataElements []*DecoratedDataElements
	Dependency            []*DataTypeDependencyType `xml:"Dependency"`
}

// DependencyWideType ...
type DependencyWideType struct {
	TypeAttr              string `xml:"Type,attr,omitempty"`
	NameAttr              string `xml:"Name,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
}

// EncodedWideType ...
type EncodedWideType struct {
	EncodedTypeAttr        string `xml:"EncodedType,attr,omitempty"`
	NameAttr               string `xml:"Name,attr,omitempty"`
	UIdAttr                string `xml:"UId,attr,omitempty"`
	ClassAttr              string `xml:"Class,attr,omitempty"`
	RevisionAttr           string `xml:"Revision,attr,omitempty"`
	RevisionExtensionAttr  string `xml:"RevisionExtension,attr,omitempty"`
	VendorAttr             string `xml:"Vendor,attr,omitempty"`
	SignatureIDAttr        string `xml:"SignatureID,attr,omitempty"`
	SignatureTimestampAttr string `xml:"SignatureTimestamp,attr,omitempty"`
	SafetySignatureIDAttr  string `xml:"SafetySignatureID,attr,omitempty"`
	EditedDateAttr         string `xml:"EditedDate,attr,omitempty"`
	SoftwareRevisionAttr   string `xml:"SoftwareRevision,attr,omitempty"`
	OriginalLibraryAttr    string `xml:"OriginalLibrary,attr,omitempty"`
	OriginalNameAttr       string `xml:"OriginalName,attr,omitempty"`
	OriginalRevisionAttr   string `xml:"OriginalRevision,attr,omitempty"`
	OriginalVendorAttr     string `xml:"OriginalVendor,attr,omitempty"`
	PermissionSetAttr      string `xml:"PermissionSet,attr,omitempty"`
	EncryptionConfigAttr   string `xml:"EncryptionConfig,attr,omitempty"`
	IsEncryptedAttr        string `xml:"IsEncrypted,attr,omitempty"`
	LogicHashAttr          string `xml:"LogicHash,attr,omitempty"`
	DescriptionHashAttr    string `xml:"DescriptionHash,attr,omitempty"`
	TypeAttr               string `xml:"Type,attr,omitempty"`
	DecoratedDataElements  []*DecoratedDataElements
	Description            []*DescriptionType                  `xml:"Description"`
	RevisionNote           []*RevisionNoteType                 `xml:"RevisionNote"`
	SignatureHistory       []*SignatureHistoryType             `xml:"SignatureHistory"`
	AdditionalHelpText     []*AdditionalHelpTextType           `xml:"AdditionalHelpText"`
	Parameters             []*UDIParameterCollectionType       `xml:"Parameters"`
	Dependencies           []*DataTypeDependencyCollectionType `xml:"Dependencies"`
}

// Value ...
type Value struct {
}

// DefaultDataWideType ...
type DefaultDataWideType struct {
	FormatAttr            string `xml:"Format,attr,omitempty"`
	LengthAttr            string `xml:"Length,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
	Value                 []*Value `xml:"Value"`
}

// CommentWideType ...
type CommentWideType struct {
	OperandAttr           string `xml:"Operand,attr,omitempty"`
	UnusedAttr            string `xml:"Unused,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
	Value                 []string                  `xml:"Value"`
	LocalizedComment      []*CommentAdaptorTextType `xml:"LocalizedComment"`
	InnerValue            []byte                    `xml:",innerxml"`
}

func (d CommentWideType) CData() string {
	return ParseCData(d.InnerValue)
}

// LocalizedCommentWideType ...
type LocalizedCommentWideType struct {
	LangAttr              string `xml:"Lang,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
	Value                 []string `xml:"Value"`
}

// DataWideType ...
type DataWideType struct {
	FormatAttr                 string `xml:"Format,attr,omitempty"`
	LengthAttr                 string `xml:"Length,attr,omitempty"`
	DecoratedDataElements      []*DecoratedDataElements
	Value                      []*Value                `xml:"Value"`
	MessageParameters          []*MsgType              `xml:"MessageParameters"`
	MotionGroupParameters      []*MotionGroupType      `xml:"MotionGroupParameters"`
	HMIBCParameters            []*HMIBCType            `xml:"HMIBCParameters"`
	BEOParameters              []*BEOType              `xml:"BEOParameters"`
	EEOParameters              []*EEOType              `xml:"EEOParameters"`
	CoordinateSystemParameters []*CoordinateSystemType `xml:"CoordinateSystemParameters"`
	AxisParameters             []*AxisType             `xml:"AxisParameters"`
	AlarmAnalogParameters      []*AlarmAnalogType      `xml:"AlarmAnalogParameters"`
	AlarmDigitalParameters     []*AlarmDigitalType     `xml:"AlarmDigitalParameters"`
	AlarmConfig                []*AlarmConfigType      `xml:"AlarmConfig"`
	DataValue                  *DataValue              `xml:"DataValue"`
	Structure                  *DataStructure          `xml:"Structure"`
	Array                      *DataArray              `xml:"Array"`
	InnerValue                 []byte                  `xml:",innerxml"`
}

func (d DataWideType) CData() string {
	return ParseCData(d.InnerValue)
}

// ForceDataWideType ...
type ForceDataWideType struct {
	FormatAttr            string `xml:"Format,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
	Value                 []*Value `xml:"Value"`
}

// TextWideType ...
type TextWideType struct {
	LangAttr              string `xml:"Lang,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
	Value                 []string                  `xml:"Value"`
	LocalizedText         []*TextBoxAdaptorTextType `xml:"LocalizedText"`
}

// LabelsWideType ...
type LabelsWideType struct {
	UIdAttr               string `xml:"UId,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
	Label                 []*LabelType `xml:"Label"`
}

// LabelWideType ...
type LabelWideType struct {
	NameAttr              string `xml:"Name,attr,omitempty"`
	UIdAttr               string `xml:"UId,attr,omitempty"`
	OperandAttr           string `xml:"Operand,attr,omitempty"`
	DecoratedDataElements []*DecoratedDataElements
	Value                 []string         `xml:"Value"`
	LocalizedLabel        []*LabelTextType `xml:"LocalizedLabel"`
}

// DecoratedDataElements ...
type DecoratedDataElements struct {
	DataValue *DataValue
	Array     *DataArray
	Structure *DataStructure
}

// DataValue ...
type DataValue struct {
	NameAttr       string `xml:"Name,attr,omitempty"`
	DataTypeAttr   string `xml:"DataType,attr,omitempty"`
	RadixAttr      string `xml:"Radix,attr,omitempty"`
	ValueAttr      string `xml:"Value,attr,omitempty"`
	ForceValueAttr string `xml:"ForceValue,attr,omitempty"`
	InnerValue     []byte `xml:",innerxml"`
}

func (d DataValue) CData() string {
	return ParseCData(d.InnerValue)
}

// DataArray ...
type DataArray struct {
	NameAttr       string              `xml:"Name,attr,omitempty"`
	DataTypeAttr   string              `xml:"DataType,attr,omitempty"`
	DimensionsAttr string              `xml:"Dimensions,attr,omitempty"`
	RadixAttr      string              `xml:"Radix,attr,omitempty"`
	Element        []*DataArrayElement `xml:"Element"`
}

// DataArrayElement ...
type DataArrayElement struct {
	IndexAttr      string           `xml:"Index,attr,omitempty"`
	ValueAttr      string           `xml:"Value,attr,omitempty"`
	ForceValueAttr string           `xml:"ForceValue,attr,omitempty"`
	Structure      []*DataStructure `xml:"Structure"`
}

// DataStructure ...
type DataStructure struct {
	NameAttr        string           `xml:"Name,attr,omitempty"`
	DataTypeAttr    string           `xml:"DataType,attr,omitempty"`
	DataValueMember []*DataValue     `xml:"DataValueMember"`
	StructureMember []*DataStructure `xml:"StructureMember"`
	ArrayMember     []*DataArray     `xml:"ArrayMember"`
}

func ParseCData(iv []byte) string {
	cv := string(iv)
	cv = strings.TrimSpace(cv)
	// everything should have this prefix and suffix
	cv = strings.TrimPrefix(cv, "<![CDATA[")
	cv = strings.TrimSuffix(cv, "]]>")
	// strings have single quotes around them so remove them too
	// this does nothing if it's not a string
	cv = strings.TrimPrefix(cv, "'")
	cv = strings.TrimSuffix(cv, "'")
	return cv
}
