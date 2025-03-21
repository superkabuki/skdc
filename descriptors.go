package skdc

import (
	"encoding/json"
	"fmt"
)

// Tag, Length , Name and Identifier for Descriptors
type TagLenNameId struct {
	Tag        uint8
	Length     uint8
	Name       string
	Identifier string
}

// Avail Descriptor
type AvailDescriptor struct {
	TagLenNameId
	ProviderAvailID uint32
}

// DTMF Descriptor
type DTMFDescriptor struct {
	TagLenNameId
	PreRoll   uint8
	DTMFCount uint8
	DTMFChars string
}

// Segmentation Descriptor
type SegmentationDescriptor struct {
	TagLenNameId
	SegmentationEventID                    string
	SegmentationEventCancelIndicator       bool
	SegmentationEventIDComplianceIndicator bool
	ProgramSegmentationFlag                bool
	SegmentationDurationFlag               bool
	DeliveryNotRestrictedFlag              bool
	WebDeliveryAllowedFlag                 bool
	NoRegionalBlackoutFlag                 bool
	ArchiveAllowedFlag                     bool
	DeviceRestrictions                     string
	SegmentationDuration                   float64
	SegmentationMessage                    string
	SegmentationUpidType                   uint8
	SegmentationUpidLength                 uint8
	SegmentationUpid                       *Upid
	SegmentationTypeID                     uint8
	SegmentNum                             uint8
	SegmentsExpected                       uint8
	SubSegmentNum                          uint8
	SubSegmentsExpected                    uint8
}

// Time Descriptor
type TimeDescriptor struct {
	TagLenNameId
	TAISeconds uint64
	TAINano    uint32
	UTCOffset  uint16
}

/*
*

	Descriptor is the combination of all the descriptors,
	works kind of like a union, it's either an AvailDescriptor,
	or DTMFDescriptor, or SegmentationDescriptor or TimeDescriptor.
	

*
*/
type Descriptor struct {
	TagLenNameId
	AvailDescriptor
	DTMFDescriptor
	SegmentationDescriptor
	TimeDescriptor
}

/*
		 *
		    Custom MarshalJSON
		        Marshal a Descriptor into

	            0x0: AvailDescriptor,
			    0x1: DTMFDescriptor,
			    0x2: SegmentationDescriptor,
			    0x3: TimeDescriptor,
		        or just return the Descriptor

*
*/
func (dscptr *Descriptor) MarshalJSON() ([]byte, error) {
	switch dscptr.Tag {
	case 0x0:
		return json.Marshal(&dscptr.AvailDescriptor)

	case 0x1:
		return json.Marshal(&dscptr.DTMFDescriptor)

	case 0x2:
		return json.Marshal(&dscptr.SegmentationDescriptor)

	case 0x3:
		return json.Marshal(&dscptr.TimeDescriptor)

	}
	type Funk Descriptor
	return json.Marshal(&struct{ *Funk }{(*Funk)(dscptr)})
}

// Return Descriptor as JSON
func (dscptr *Descriptor) Json() string {
	stuff, err := dscptr.MarshalJSON()
	chk(err)
	return string(stuff)
}

// Print Descriptor as JSON
func (dscptr *Descriptor) Show() {
	fmt.Printf(dscptr.Json())
}

/*
*
Decode returns a Splice Descriptor by tag.

	The following Splice Descriptors are recognized.

	    0x0: Avail Descriptor,
	    0x1: DTMF Descriptor,
	    0x2: Segmentation Descriptor,
	    0x3: Time Descriptor,

*
*/
func (dscptr *Descriptor) decode(bd *bitDecoder, tag uint8, length uint8) {
	switch tag {
	case 0x0:
		dscptr.Tag = 0x0
		dscptr.decodeAvailDescriptor(bd, tag, length)
	case 0x1:
		dscptr.Tag = 0x1
		dscptr.decodeDTMFDescriptor(bd, tag, length)
	case 0x2:
		dscptr.Tag = 0x2
		dscptr.decodeSegmentationDescriptor(bd, tag, length)
	case 0x3:
		dscptr.Tag = 0x3
		dscptr.decodeTimeDescriptor(bd, tag, length)
	}
}


// Decode for  Avail Descriptors
func (dscptr *Descriptor) decodeAvailDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "Avail Descriptor"
	dscptr.ProviderAvailID = bd.uInt32(32)
	dscptr.AvailDescriptor.TagLenNameId = dscptr.TagLenNameId

}

// Decode for DTMF Splice Descriptor
func (dscptr *Descriptor) decodeDTMFDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "DTMF Descriptor"
	dscptr.PreRoll = bd.uInt8(8)
	dscptr.DTMFCount = bd.uInt8(8)>>5
	//bd.goForward(5)
	dscptr.DTMFChars = bd.asAscii(uint(8 * dscptr.DTMFCount))
	dscptr.DTMFDescriptor.TagLenNameId = dscptr.TagLenNameId

}

// Decode for the Time Descriptor
func (dscptr *Descriptor) decodeTimeDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "Time Descriptor"
	dscptr.TAISeconds = bd.uInt64(48)
	dscptr.TAINano = bd.uInt32(32)
	dscptr.UTCOffset = bd.uInt16(16)
	dscptr.TimeDescriptor.TagLenNameId = dscptr.TagLenNameId

}

// Decode for the Segmentation Descriptor
func (dscptr *Descriptor) decodeSegmentationDescriptor(bd *bitDecoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.asAscii(32)
	dscptr.Name = "Segmentation Descriptor"
	dscptr.SegmentationEventID = bd.asHex(32)
	dscptr.SegmentationEventCancelIndicator = bd.asFlag()
	dscptr.SegmentationEventIDComplianceIndicator = bd.asFlag()
	bd.goForward(6)
	if !dscptr.SegmentationEventCancelIndicator {
		dscptr.decodeSegFlags(bd)
		dscptr.decodeSegmentation(bd)
	}
	dscptr.SegmentationDescriptor.TagLenNameId = dscptr.TagLenNameId

}

func (dscptr *Descriptor) decodeSegFlags(bd *bitDecoder) {
	dscptr.ProgramSegmentationFlag = bd.asFlag()
	dscptr.SegmentationDurationFlag = bd.asFlag()
	dscptr.DeliveryNotRestrictedFlag = bd.asFlag()
	if !dscptr.DeliveryNotRestrictedFlag {
		dscptr.WebDeliveryAllowedFlag = bd.asFlag()
		dscptr.NoRegionalBlackoutFlag = bd.asFlag()
		dscptr.ArchiveAllowedFlag = bd.asFlag()
		dscptr.DeviceRestrictions = table20[bd.uInt8(2)] // 8
	} else {
		bd.goForward(5)
	}
}

func (dscptr *Descriptor) decodeSegmentation(bd *bitDecoder) {
	if dscptr.SegmentationDurationFlag {
		dscptr.SegmentationDuration = bd.as90k(40)
	}
	dscptr.SegmentationUpidType = bd.uInt8(8)
	dscptr.SegmentationUpidLength = bd.uInt8(8)
	if dscptr.SegmentationUpidLength > 0 {
		dscptr.SegmentationUpid = &Upid{}
		dscptr.SegmentationUpid.decode(bd, dscptr.SegmentationUpidType, dscptr.SegmentationUpidLength)
	}
	dscptr.SegmentationTypeID = bd.uInt8(8)
	mesg, ok := table22[dscptr.SegmentationTypeID]
	if ok {
		dscptr.SegmentationMessage = mesg
	}
	dscptr.SegmentNum = bd.uInt8(8)
	dscptr.SegmentsExpected = bd.uInt8(8)
	subSegIDs := []uint16{0x30, 0x32, 0x34, 0x36, 0x38, 0x3A, 0x44, 0x46}
	if IsIn(subSegIDs, uint16(dscptr.SegmentationTypeID)) {
		dscptr.SubSegmentNum = bd.uInt8(8)
		dscptr.SubSegmentsExpected = bd.uInt8(8)
	}
}

func (dscptr *Descriptor) encode(be *bitEncoder) {
	switch dscptr.Tag {
	case 0x0:
		dscptr.encodeAvailDescriptor(be)
	case 0x2:
		dscptr.encodeSegmentationDescriptor(be)
	}
}

// Encode for Avail Descriptors
func (dscptr *Descriptor) encodeAvailDescriptor(be *bitEncoder) {
	be.Add(uint32(dscptr.ProviderAvailID), 32)
}

// Encode a segmentation descriptor
func (dscptr *Descriptor) encodeSegmentationDescriptor(be *bitEncoder) {
	dscptr.SegmentationDescriptor.TagLenNameId = dscptr.TagLenNameId

	be.AddHex64(dscptr.SegmentationEventID, 32)
	be.Add(dscptr.SegmentationEventCancelIndicator, 1)
	be.Add(dscptr.SegmentationEventIDComplianceIndicator, 1)
	be.Reserve(6)
	if !dscptr.SegmentationEventCancelIndicator {
		dscptr.encodeFlags(be)
		dscptr.encodeSegmentation(be)
	}
}

func (dscptr *Descriptor) encodeFlags(be *bitEncoder) {
	be.Add(dscptr.ProgramSegmentationFlag, 1)
	be.Add(dscptr.SegmentationDurationFlag, 1)
	be.Add(dscptr.DeliveryNotRestrictedFlag, 1)
	if !dscptr.DeliveryNotRestrictedFlag {
		be.Add(dscptr.WebDeliveryAllowedFlag, 1)
		be.Add(dscptr.NoRegionalBlackoutFlag, 1)
		be.Add(dscptr.ArchiveAllowedFlag, 1)
		//   a_key = k_by_v(table20, dscptr.device_restrictions)
		//     nbin.add_int(a_key, 2)
		be.Add(3, 2) //  dscptr.device_restrictions
	} else {
		be.Reserve(5)
	}
}

func (dscptr *Descriptor) encodeSegmentation(be *bitEncoder) {
	if dscptr.SegmentationDurationFlag {
		be.Add(float64(dscptr.SegmentationDuration), 40)
	}
	be.Add(dscptr.SegmentationUpidType, 8)
	be.Add(dscptr.SegmentationUpidLength, 8)
	if dscptr.SegmentationUpidLength > 0 {
		dscptr.SegmentationUpid.encode(be, dscptr.SegmentationUpidType)
	}
	be.Add(dscptr.SegmentationTypeID, 8)
	dscptr.encodeSegments(be)
}

func (dscptr *Descriptor) encodeSegments(be *bitEncoder) {
	be.Add(dscptr.SegmentNum, 8)
	be.Add(dscptr.SegmentsExpected, 8)
	subSegIDs := []uint16{0x30, 0x32, 0x34, 0x36, 0x38, 0x3A, 0x44, 0x46}
	if IsIn(subSegIDs, uint16(dscptr.SegmentationTypeID)) {
		be.Add(dscptr.SubSegmentNum, 8)
		be.Add(dscptr.SubSegmentsExpected, 8)
	}

}
