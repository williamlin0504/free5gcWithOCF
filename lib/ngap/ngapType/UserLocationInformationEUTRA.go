package ngapType

// Need to import "free5gcWithOCF/lib/aper" if it uses "aper"

type UserLocationInformationEUTRA struct {
	EUTRACGI     EUTRACGI                                                      `aper:"valueExt"`
	TAI          TAI                                                           `aper:"valueExt"`
	TimeStamp    *TimeStamp                                                    `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUserLocationInformationEUTRAExtIEs `aper:"optional"`
}