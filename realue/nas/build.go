// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package nas

import (
	"bytes"
	"fmt"

	realuectx "github.com/omec-project/gnbsim/realue/context"
	"github.com/omec-project/gnbsim/realue/util"
	"github.com/omec-project/gnbsim/util/nastestpacket"
	"github.com/omec-project/nas"

	"github.com/omec-project/nas/nasConvert"
	"github.com/omec-project/nas/nasMessage"
	"github.com/omec-project/nas/nasTestpacket"
	"github.com/omec-project/nas/nasType"
)

const (
	SWITCH_OFF uint8 = 0
)

func GetServiceRequest(ue *realuectx.RealUe) ([]byte, error) {

	nasMsg := nastestpacket.BuildServiceRequest(nasMessage.ServiceTypeData)
	serviceRequest := nasMsg.GmmMessage.ServiceRequest

	guti := nasConvert.GutiToNas(ue.Guti)
	serviceRequest.SetTypeOfIdentity(nasMessage.MobileIdentity5GSType5gSTmsi)
	serviceRequest.SetAMFSetID(guti.GetAMFSetID())
	serviceRequest.SetAMFPointer(guti.GetAMFPointer())
	serviceRequest.SetTMSI5G(guti.GetTMSI5G())
	serviceRequest.SetNasKeySetIdentifiler(uint8(ue.NgKsi.Ksi))

	data := new(bytes.Buffer)
	err := nasMsg.GmmMessageEncode(data)
	if err != nil {
		return nil, fmt.Errorf("encode failed:", err)
	}

	return data.Bytes(), nil
}

func GetRegisterRequest(ue *realuectx.RealUe) ([]byte, error) {
	var err error

	ueSecurityCapability := ue.GetUESecurityCapability()
	ue.Suci, err = util.SupiToSuci(ue.Supi, ue.Plmn)
	if err != nil {
		ue.Log.Errorln("SupiToSuci returned:", err)
		return nil, fmt.Errorf("failed to derive suci")
	}
	mobileId5GS := nasType.MobileIdentity5GS{
		Len:    uint16(len(ue.Suci)), // suci
		Buffer: ue.Suci,
	}
	snssaiBuf := nasConvert.SnssaiToNas(*ue.SNssai)
	requestedNSSAI := &nasType.RequestedNSSAI{
		Iei:    nasMessage.RegistrationRequestRequestedNSSAIType,
		Len:    uint8(len(snssaiBuf)),
		Buffer: snssaiBuf,
	}

	networkSlicingIndication := nasType.NewNetworkSlicingIndication(
		nasMessage.RegistrationRequestNetworkSlicingIndicationType)
	// LG: Useless, for memo
	networkSlicingIndication.SetDCNI(0)
	networkSlicingIndication.SetNSSCI(0)

	capability5GMM := &nasType.Capability5GMM{
		Iei: nasMessage.RegistrationRequestCapability5GMMType,
		Len: 1,
	}

	ue.Log.Traceln("Generating SUPI Registration Request Message")
	nasPdu := nasTestpacket.GetRegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration,
		mobileId5GS,
		requestedNSSAI,
		networkSlicingIndication,
		ueSecurityCapability,
		capability5GMM,
		nil,
		nil,
	)
	return nasPdu, nil
}

func GetDeregisterRequest(ue *realuectx.RealUe) ([]byte, error) {
	var err error

	if ue.Guti == "" {
		ue.Log.Errorln("guti not allocated")
		return nil, fmt.Errorf(
			"failed to create deregistration request: guti not allocated",
		)
	}
	gutiNas := nasConvert.GutiToNas(ue.Guti)
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    11, // 5g-guti
		Buffer: gutiNas.Octet[:],
	}

	nasPdu := nasTestpacket.GetDeregistrationRequest(nasMessage.AccessType3GPP,
		SWITCH_OFF, uint8(ue.NgKsi.Ksi), mobileIdentity5GS)
	nasPdu, err = EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		ue.Log.Errorln("EncodeNasPduWithSecurity() returned:", err)
		return nil, fmt.Errorf(
			"failed to encrypt deregistration request message",
		)
	}
	return nasPdu, nil
}

func GetRegistrationComplete(ue *realuectx.RealUe) ([]byte, error) {
	ue.Log.Traceln("Generating Registration Complete Message")
	nasPdu := nasTestpacket.GetRegistrationComplete(nil)
	nasPdu, err := EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		ue.Log.Errorln("EncodeNasPduWithSecurity() returned:", err)
		return nasPdu, fmt.Errorf(
			"failed to encrypt registration complete message",
		)
	}
	return nasPdu, nil
}

func NasGetTransferContent(
	ue *realuectx.RealUe,
	xNasTransport *nas.Message,
) (msgType uint8, nasMsg *nas.Message, err error) {

	msgType = xNasTransport.GmmHeader.GetMessageType()
	ue.Log.Infoln("Received Message Type:", msgType)

	if msgType == nas.MsgTypeDLNASTransport {
		ue.Log.Info(
			"Payload container type:",
			xNasTransport.GmmMessage.DLNASTransport.SpareHalfOctetAndPayloadContainerType,
		)
		payload := xNasTransport.GmmMessage.DLNASTransport.PayloadContainer
		if payload.Len == 0 {
			err = fmt.Errorf("payload container length is 0")
			return
		}
		buffer := payload.Buffer[:payload.Len]
		m := nas.NewMessage()
		err = m.PlainNasDecode(&buffer)
		if err != nil {
			ue.Log.Errorln("PlainNasDecode returned:", err)
			err = fmt.Errorf("failed to decode payload container")
			return
		}
		nasMsg = m
		msgType = nasMsg.GsmHeader.GetMessageType()
	}
	return
}
