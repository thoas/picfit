package mknote

import "github.com/rwcarlsen/goexif/exif"

// Useful resources used in creating these tables:
//    http://www.exiv2.org/makernote.html
//    http://www.exiv2.org/tags-canon.html
//    http://www.exiv2.org/tags-nikon.html

// Known Maker Note fields
const (
	// common fields
	ISOSpeed                   exif.FieldName = "ISOSpeed"
	ColorMode                  exif.FieldName = "ColorMode"
	Quality                    exif.FieldName = "Quality"
	Sharpening                 exif.FieldName = "Sharpening"
	Focus                      exif.FieldName = "Focus"
	FlashSetting               exif.FieldName = "FlashSetting"
	FlashDevice                exif.FieldName = "FlashDevice"
	WhiteBalanceBias           exif.FieldName = "WhiteBalanceBias"
	WB_RBLevels                exif.FieldName = "WB_RBLevels"
	ProgramShift               exif.FieldName = "ProgramShift"
	ExposureDiff               exif.FieldName = "ExposureDiff"
	ISOSelection               exif.FieldName = "ISOSelection"
	DataDump                   exif.FieldName = "DataDump"
	Preview                    exif.FieldName = "Preview"
	FlashComp                  exif.FieldName = "FlashComp"
	ISOSettings                exif.FieldName = "ISOSettings"
	ImageBoundary              exif.FieldName = "ImageBoundary"
	FlashExposureComp          exif.FieldName = "FlashExposureComp"
	FlashBracketComp           exif.FieldName = "FlashBracketComp"
	ExposureBracketComp        exif.FieldName = "ExposureBracketComp"
	ImageProcessing            exif.FieldName = "ImageProcessing"
	CropHiSpeed                exif.FieldName = "CropHiSpeed"
	ExposureTuning             exif.FieldName = "ExposureTuning"
	SerialNumber               exif.FieldName = "SerialNumber"
	ImageAuthentication        exif.FieldName = "ImageAuthentication"
	ActiveDLighting            exif.FieldName = "ActiveDLighting"
	VignetteControl            exif.FieldName = "VignetteControl"
	ImageAdjustment            exif.FieldName = "ImageAdjustment"
	ToneComp                   exif.FieldName = "ToneComp"
	AuxiliaryLens              exif.FieldName = "AuxiliaryLens"
	LensType                   exif.FieldName = "LensType"
	Lens                       exif.FieldName = "Lens"
	FocusDistance              exif.FieldName = "FocusDistance"
	DigitalZoom                exif.FieldName = "DigitalZoom"
	FlashMode                  exif.FieldName = "FlashMode"
	ShootingMode               exif.FieldName = "ShootingMode"
	AutoBracketRelease         exif.FieldName = "AutoBracketRelease"
	LensFStops                 exif.FieldName = "LensFStops"
	ContrastCurve              exif.FieldName = "ContrastCurve"
	ColorHue                   exif.FieldName = "ColorHue"
	SceneMode                  exif.FieldName = "SceneMode"
	HueAdjustment              exif.FieldName = "HueAdjustment"
	NEFCompression             exif.FieldName = "NEFCompression"
	NoiseReduction             exif.FieldName = "NoiseReduction"
	LinearizationTable         exif.FieldName = "LinearizationTable"
	RawImageCenter             exif.FieldName = "RawImageCenter"
	SensorPixelSize            exif.FieldName = "SensorPixelSize"
	SceneAssist                exif.FieldName = "SceneAssist"
	RetouchHistory             exif.FieldName = "RetouchHistory"
	ImageDataSize              exif.FieldName = "ImageDataSize"
	ImageCount                 exif.FieldName = "ImageCount"
	DeletedImageCount          exif.FieldName = "DeletedImageCount"
	ShutterCount               exif.FieldName = "ShutterCount"
	ImageOptimization          exif.FieldName = "ImageOptimization"
	SaturationText             exif.FieldName = "SaturationText"
	VariProgram                exif.FieldName = "VariProgram"
	ImageStabilization         exif.FieldName = "ImageStabilization"
	AFResponse                 exif.FieldName = "AFResponse"
	HighISONoiseReduction      exif.FieldName = "HighISONoiseReduction"
	ToningEffect               exif.FieldName = "ToningEffect"
	PrintIM                    exif.FieldName = "PrintIM"
	CaptureData                exif.FieldName = "CaptureData"
	CaptureVersion             exif.FieldName = "CaptureVersion"
	CaptureOffsets             exif.FieldName = "CaptureOffsets"
	ScanIFD                    exif.FieldName = "ScanIFD"
	ICCProfile                 exif.FieldName = "ICCProfile"
	CaptureOutput              exif.FieldName = "CaptureOutput"
	Panorama                   exif.FieldName = "Panorama"
	ImageType                  exif.FieldName = "ImageType"
	FirmwareVersion            exif.FieldName = "FirmwareVersion"
	FileNumber                 exif.FieldName = "FileNumber"
	OwnerName                  exif.FieldName = "OwnerName"
	CameraInfo                 exif.FieldName = "CameraInfo"
	CustomFunctions            exif.FieldName = "CustomFunctions"
	ModelID                    exif.FieldName = "ModelID"
	PictureInfo                exif.FieldName = "PictureInfo"
	ThumbnailImageValidArea    exif.FieldName = "ThumbnailImageValidArea"
	SerialNumberFormat         exif.FieldName = "SerialNumberFormat"
	SuperMacro                 exif.FieldName = "SuperMacro"
	OriginalDecisionDataOffset exif.FieldName = "OriginalDecisionDataOffset"
	WhiteBalanceTable          exif.FieldName = "WhiteBalanceTable"
	LensModel                  exif.FieldName = "LensModel"
	InternalSerialNumber       exif.FieldName = "InternalSerialNumber"
	DustRemovalData            exif.FieldName = "DustRemovalData"
	ProcessingInfo             exif.FieldName = "ProcessingInfo"
	MeasuredColor              exif.FieldName = "MeasuredColor"
	VRDOffset                  exif.FieldName = "VRDOffset"
	SensorInfo                 exif.FieldName = "SensorInfo"
	ColorData                  exif.FieldName = "ColorData"

	// Nikon-specific fields
	Nikon_Version        exif.FieldName = "Nikon.Version"
	Nikon_WhiteBalance   exif.FieldName = "Nikon.WhiteBalance"
	Nikon_ColorSpace     exif.FieldName = "Nikon.ColorSpace"
	Nikon_LightSource    exif.FieldName = "Nikon.LightSource"
	Nikon_Saturation     exif.FieldName = "Nikon_Saturation"
	Nikon_ShotInfo       exif.FieldName = "Nikon.ShotInfo"       // A sub-IFD
	Nikon_VRInfo         exif.FieldName = "Nikon.VRInfo"         // A sub-IFD
	Nikon_PictureControl exif.FieldName = "Nikon.PictureControl" // A sub-IFD
	Nikon_WorldTime      exif.FieldName = "Nikon.WorldTime"      // A sub-IFD
	Nikon_ISOInfo        exif.FieldName = "Nikon.ISOInfo"        // A sub-IFD
	Nikon_AFInfo         exif.FieldName = "Nikon.AFInfo"         // A sub-IFD
	Nikon_ColorBalance   exif.FieldName = "Nikon.ColorBalance"   // A sub-IFD
	Nikon_LensData       exif.FieldName = "Nikon.LensData"       // A sub-IFD
	Nikon_SerialNO       exif.FieldName = "Nikon.SerialNO"       // usually starts with "NO="
	Nikon_FlashInfo      exif.FieldName = "Nikon.FlashInfo"      // A sub-IFD
	Nikon_MultiExposure  exif.FieldName = "Nikon.MultiExposure"  // A sub-IFD
	Nikon_AFInfo2        exif.FieldName = "Nikon.AFInfo2"        // A sub-IFD
	Nikon_FileInfo       exif.FieldName = "Nikon.FileInfo"       // A sub-IFD
	Nikon_AFTune         exif.FieldName = "Nikon.AFTune"         // A sub-IFD
	Nikon3_0x000a        exif.FieldName = "Nikon3.0x000a"
	Nikon3_0x009b        exif.FieldName = "Nikon3.0x009b"
	Nikon3_0x009f        exif.FieldName = "Nikon3.0x009f"
	Nikon3_0x00a3        exif.FieldName = "Nikon3.0x00a3"

	// Canon-specific fiends
	Canon_CameraSettings exif.FieldName = "Canon.CameraSettings" // A sub-IFD
	Canon_ShotInfo       exif.FieldName = "Canon.ShotInfo"       // A sub-IFD
	Canon_AFInfo         exif.FieldName = "Canon.AFInfo"
	Canon_TimeInfo       exif.FieldName = "Canon.TimeInfo"
	Canon_0x0000         exif.FieldName = "Canon.0x0000"
	Canon_0x0003         exif.FieldName = "Canon.0x0003"
	Canon_0x00b5         exif.FieldName = "Canon.0x00b5"
	Canon_0x00c0         exif.FieldName = "Canon.0x00c0"
	Canon_0x00c1         exif.FieldName = "Canon.0x00c1"
)

var makerNoteCanonFields = map[uint16]exif.FieldName{
	0x0000: Canon_0x0000,
	0x0001: Canon_CameraSettings,
	0x0002: exif.FocalLength,
	0x0003: Canon_0x0003,
	0x0004: Canon_ShotInfo,
	0x0005: Panorama,
	0x0006: ImageType,
	0x0007: FirmwareVersion,
	0x0008: FileNumber,
	0x0009: OwnerName,
	0x000c: SerialNumber,
	0x000d: CameraInfo,
	0x000f: CustomFunctions,
	0x0010: ModelID,
	0x0012: PictureInfo,
	0x0013: ThumbnailImageValidArea,
	0x0015: SerialNumberFormat,
	0x001a: SuperMacro,
	0x0026: Canon_AFInfo,
	0x0035: Canon_TimeInfo,
	0x0083: OriginalDecisionDataOffset,
	0x00a4: WhiteBalanceTable,
	0x0095: LensModel,
	0x0096: InternalSerialNumber,
	0x0097: DustRemovalData,
	0x0099: CustomFunctions,
	0x00a0: ProcessingInfo,
	0x00aa: MeasuredColor,
	0x00b4: exif.ColorSpace,
	0x00b5: Canon_0x00b5,
	0x00c0: Canon_0x00c0,
	0x00c1: Canon_0x00c1,
	0x00d0: VRDOffset,
	0x00e0: SensorInfo,
	0x4001: ColorData,
}

// Nikon version 3 Maker Notes fields (used by E5400, SQ, D2H, D70, and newer)
var makerNoteNikon3Fields = map[uint16]exif.FieldName{
	0x0001: Nikon_Version,
	0x0002: ISOSpeed,
	0x0003: ColorMode,
	0x0004: Quality,
	0x0005: Nikon_WhiteBalance,
	0x0006: Sharpening,
	0x0007: Focus,
	0x0008: FlashSetting,
	0x0009: FlashDevice,
	0x000a: Nikon3_0x000a,
	0x000b: WhiteBalanceBias,
	0x000c: WB_RBLevels,
	0x000d: ProgramShift,
	0x000e: ExposureDiff,
	0x000f: ISOSelection,
	0x0010: DataDump,
	0x0011: Preview,
	0x0012: FlashComp,
	0x0013: ISOSettings,
	0x0016: ImageBoundary,
	0x0017: FlashExposureComp,
	0x0018: FlashBracketComp,
	0x0019: ExposureBracketComp,
	0x001a: ImageProcessing,
	0x001b: CropHiSpeed,
	0x001c: ExposureTuning,
	0x001d: SerialNumber,
	0x001e: Nikon_ColorSpace,
	0x001f: Nikon_VRInfo,
	0x0020: ImageAuthentication,
	0x0022: ActiveDLighting,
	0x0023: Nikon_PictureControl,
	0x0024: Nikon_WorldTime,
	0x0025: Nikon_ISOInfo,
	0x002a: VignetteControl,
	0x0080: ImageAdjustment,
	0x0081: ToneComp,
	0x0082: AuxiliaryLens,
	0x0083: LensType,
	0x0084: Lens,
	0x0085: FocusDistance,
	0x0086: DigitalZoom,
	0x0087: FlashMode,
	0x0088: Nikon_AFInfo,
	0x0089: ShootingMode,
	0x008a: AutoBracketRelease,
	0x008b: LensFStops,
	0x008c: ContrastCurve,
	0x008d: ColorHue,
	0x008f: SceneMode,
	0x0090: Nikon_LightSource,
	0x0091: Nikon_ShotInfo,
	0x0092: HueAdjustment,
	0x0093: NEFCompression,
	0x0094: Nikon_Saturation,
	0x0095: NoiseReduction,
	0x0096: LinearizationTable,
	0x0097: Nikon_ColorBalance,
	0x0098: Nikon_LensData,
	0x0099: RawImageCenter,
	0x009a: SensorPixelSize,
	0x009b: Nikon3_0x009b,
	0x009c: SceneAssist,
	0x009e: RetouchHistory,
	0x009f: Nikon3_0x009f,
	0x00a0: Nikon_SerialNO,
	0x00a2: ImageDataSize,
	0x00a3: Nikon3_0x00a3,
	0x00a5: ImageCount,
	0x00a6: DeletedImageCount,
	0x00a7: ShutterCount,
	0x00a8: Nikon_FlashInfo,
	0x00a9: ImageOptimization,
	0x00aa: SaturationText,
	0x00ab: VariProgram,
	0x00ac: ImageStabilization,
	0x00ad: AFResponse,
	0x00b0: Nikon_MultiExposure,
	0x00b1: HighISONoiseReduction,
	0x00b3: ToningEffect,
	0x00b7: Nikon_AFInfo2,
	0x00b8: Nikon_FileInfo,
	0x00b9: Nikon_AFTune,
	0x0e00: PrintIM,
	0x0e01: CaptureData,
	0x0e09: CaptureVersion,
	0x0e0e: CaptureOffsets,
	0x0e10: ScanIFD,
	0x0e1d: ICCProfile,
	0x0e1e: CaptureOutput,
}
