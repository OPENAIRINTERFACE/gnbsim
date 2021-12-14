// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
// SPDX-License-Identifier: LicenseRef-ONF-Member-Only-1.0

package realue

import (
	"fmt"
	"gnbsim/common"
	"gnbsim/realue/context"
	"gnbsim/realue/util"
	"gnbsim/realue/worker/pdusessworker"
	"gnbsim/util/test"
	"net"

	realue_nas "gnbsim/realue/nas"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasTestpacket"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/openapi/models"
	"github.com/omec-project/nas"
	"github.com/omec-project/nas/nasConvert"
)

//TODO Remove the hardcoding
const (
	SN_NAME    string = "5G:mnc093.mcc208.3gppnetwork.org"
	SWITCH_OFF uint8  = 0
)

func HandleRegRequestEvent(ue *context.RealUe,
	msg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Registration Request Event")

	ueSecurityCapability := ue.GetUESecurityCapability()

	ue.Suci, err = util.SupiToSuci(ue.Supi, ue.Plmn)
	if err != nil {
		ue.Log.Errorln("SupiToSuci returned:", err)
		return fmt.Errorf("failed to derive suci")
	}
	mobileId5GS := nasType.MobileIdentity5GS{
		Len:    uint16(len(ue.Suci)), // suci
		Buffer: ue.Suci,
	}

	ue.Log.Traceln("Generating Registration Request Message")
	nasPdu := nasTestpacket.GetRegistrationRequest(nasMessage.RegistrationType5GSInitialRegistration,
		mobileId5GS, nil, ueSecurityCapability, nil, nil, nil)

	m := formUuMessage(common.REG_REQUEST_EVENT, nasPdu)
	SendToSimUe(ue, m)
	ue.Log.Traceln("Sent Registration Request Message to SimUe")
	return nil
}

func HandleAuthResponseEvent(ue *context.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Authentication Response Event")

	msg := intfcMsg.(*common.UeMessage)
	// First process the corresponding Auth Request
	ue.Log.Traceln("Processing corresponding Authentication Request Message")
	authReq := msg.NasMsg.AuthenticationRequest

	ue.NgKsi = nasConvert.SpareHalfOctetAndNgksiToModels(authReq.SpareHalfOctetAndNgksi)

	rand := authReq.GetRANDValue()
	autn := authReq.GetAUTN()
	resStat := ue.DeriveRESstarAndSetKey(autn[:], rand[:], SN_NAME)

	// TODO: Parse Auth Request IEs and update the RealUE Context

	// Now generate NAS Authentication Response
	ue.Log.Traceln("Generating Authentication Reponse Message")
	nasPdu := nasTestpacket.GetAuthenticationResponse(resStat, "")

	m := formUuMessage(common.AUTH_RESPONSE_EVENT, nasPdu)
	SendToSimUe(ue, m)
	ue.Log.Traceln("Sent Authentication Reponse Message to SimUe")
	return nil
}

func HandleSecModCompleteEvent(ue *context.RealUe,
	msg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Security Mode Complete Event")

	//TODO: Process corresponding Security Mode Command first

	mobileId5GS := nasType.MobileIdentity5GS{
		Len:    uint16(len(ue.Suci)), // suci
		Buffer: ue.Suci,
	}
	registrationRequestWith5GMM := nasTestpacket.GetRegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration, mobileId5GS, nil,
		ue.GetUESecurityCapability(), ue.Get5GMMCapability(), nil, nil)

	ue.Log.Traceln("Generating Security Mode Complete Message")
	nasPdu := nasTestpacket.GetSecurityModeComplete(registrationRequestWith5GMM)

	nasPdu, err = test.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext,
		true)
	if err != nil {
		ue.Log.Errorln("EncodeNasPduWithSecurity() returned:", err)
		return fmt.Errorf("failed to encrypt security mode complete message")
	}

	m := formUuMessage(common.SEC_MOD_COMPLETE_EVENT, nasPdu)
	SendToSimUe(ue, m)
	ue.Log.Traceln("Sent Security Mode Complete Message to SimUe")
	return nil
}

func HandleRegCompleteEvent(ue *context.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Registration Complete Event")

	//TODO: Process corresponding Registration Accept first
	msg := intfcMsg.(*common.UeMessage).NasMsg.RegistrationAccept

	var guti []uint8
	if msg.GUTI5G != nil {
		guti = msg.GUTI5G.Octet[:]
	}

	_, ue.Guti = nasConvert.GutiToString(guti)

	ue.Log.Traceln("Generating Registration Complete Message")
	nasPdu := nasTestpacket.GetRegistrationComplete(nil)
	nasPdu, err = test.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		ue.Log.Errorln("EncodeNasPduWithSecurity() returned:", err)
		return fmt.Errorf("failed to encrypt registration complete message")
	}

	m := formUuMessage(common.REG_COMPLETE_EVENT, nasPdu)
	SendToSimUe(ue, m)
	ue.Log.Traceln("Sent Registration Complete Message to SimUe")
	return nil
}

func HandleDeregRequestEvent(ue *context.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling UE Initiated Deregistration Request Event")

	if ue.Guti == "" {
		ue.Log.Errorln("guti not allocated")
		return fmt.Errorf("failed to create deregistration request: guti not unallocated")
	}
	gutiNas := nasConvert.GutiToNas(ue.Guti)
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    11, // 5g-guti
		Buffer: gutiNas.Octet[:],
	}

	nasPdu := nasTestpacket.GetDeregistrationRequest(nasMessage.AccessType3GPP,
		SWITCH_OFF, uint8(ue.NgKsi.Ksi), mobileIdentity5GS)
	nasPdu, err = test.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		ue.Log.Errorln("EncodeNasPduWithSecurity() returned:", err)
		return fmt.Errorf("failed to encrypt deregistration request message")
	}

	m := formUuMessage(common.DEREG_REQUEST_UE_ORIG_EVENT, nasPdu)
	SendToSimUe(ue, m)
	ue.Log.Traceln("Sent UE Initiated Deregistration Request message to SimUe")
	return nil
}

func HandlePduSessEstRequestEvent(ue *context.RealUe,
	msg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling PDU Session Establishment Request Event")

	sNssai := models.Snssai{
		Sst: 1,
		Sd:  "010203",
	}
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10,
		nasMessage.ULNASTransportRequestTypeInitialRequest, "internet", &sNssai)

	nasPdu, err = test.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		fmt.Println("Failed to encrypt PDU Session Establishment Request Message", err)
		return
	}

	m := formUuMessage(common.PDU_SESS_EST_REQUEST_EVENT, nasPdu)
	SendToSimUe(ue, m)
	ue.Log.Traceln("Sent PDU Session Establishment Request Message to SimUe")
	return nil
}

func HandlePduSessEstAcceptEvent(ue *context.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling PDU Session Establishment Accept Event")

	msg := intfcMsg.(*common.UeMessage)
	//TODO: create new pdu session var and parse msg to pdu session var
	nasMsg := msg.NasMsg.PDUSessionEstablishmentAccept
	if nasMsg == nil {
		ue.Log.Errorln("PDUSessionEstablishmentAccept is nil")
		return fmt.Errorf("invalid NAS Message")
	}

	var pduAddr net.IP
	pduSessType := nasConvert.PDUSessionTypeToModels(nasMsg.GetPDUSessionType())
	if pduSessType == models.PduSessionType_IPV4 {
		ip := nasMsg.GetPDUAddressInformation()
		pduAddr = net.IPv4(ip[0], ip[1], ip[2], ip[3])
	}

	pduSess := context.NewPduSession(ue, uint64(nasMsg.PDUSessionID.Octet))
	pduSess.PduSessType = pduSessType
	pduSess.SscMode = nasMsg.GetSSCMode()
	pduSess.PduAddress = pduAddr
	pduSess.WriteUeChan = ue.ReadChan
	ue.AddPduSession(int64(pduSess.PduSessId), pduSess)
	ue.Log.Infoln("PDU Session ID:", pduSess.PduSessId)
	ue.Log.Infoln("PDU Session Type:", pduSess.PduSessType)
	ue.Log.Infoln("SSC Mode:", pduSess.SscMode)
	ue.Log.Infoln("PDU Address:", pduAddr.String())

	return nil
}

func HandleDataBearerSetupRequestEvent(ue *context.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Data Bearer Setup Request Event")

	msg := intfcMsg.(*common.UuMessage)
	for _, item := range msg.DBParams {
		/* Currently gNB also adds failed pdu session ids in the list.
		   pdu sessions are marked failed during decoding. real ue simply
		   returns the same list back by marking any failed pdu sessions on
		   its side. This consolidated list can be used by gnb to form
		   PDUSession Resource Setup Response message
		*/
		if item.PduSess.Success {
			pduSess := ue.GetPduSession(item.PduSess.PduSessId)
			if pduSess == nil {
				item.PduSess.Success = false
				continue
			}

			pduSess.WriteGnbChan = item.CommChan

			/* gNb can use this channel to send DL packets for this PDU session */
			item.CommChan = pduSess.ReadDlChan

			go pdusessworker.Init(pduSess)
		}
	}

	rsp := &common.UuMessage{}
	rsp.Event = common.DATA_BEARER_SETUP_RESPONSE_EVENT
	rsp.DBParams = msg.DBParams
	ue.WriteSimUeChan <- rsp
	ue.Log.Infoln("Sent Data Radio Bearer Setup Response event to SimUe")
	return nil
}

func HandleDataPktGenRequestEvent(ue *context.RealUe,
	msg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Data Packet Generation Request Event")

	for _, v := range ue.PduSessions {
		v.ReadCmdChan <- msg
	}

	return nil
}

func HandleDataPktGenSuccessEvent(ue *context.RealUe,
	msg common.InterfaceMessage) (err error) {
	ue.WriteSimUeChan <- msg
	return nil
}

func HandleDlInfoTransferEvent(ue *context.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Downlink Nas Transport Event")

	msg := intfcMsg.(*common.UuMessage)
	for _, pdu := range msg.NasPdus {
		nasMsg, err := test.NASDecode(ue, nas.GetSecurityHeaderType(pdu), pdu)
		if err != nil {
			ue.Log.Errorln("Failed to decode dowlink NAS Message due to", err)
			return err
		}
		msgType := nasMsg.GmmHeader.GetMessageType()
		ue.Log.Infoln("Received Message Type:", msgType)

		if msgType == nas.MsgTypeDLNASTransport {
			ue.Log.Info("Payload contaner type:",
				nasMsg.GmmMessage.DLNASTransport.SpareHalfOctetAndPayloadContainerType)
			payload := nasMsg.GmmMessage.DLNASTransport.PayloadContainer
			if payload.Len == 0 {
				return fmt.Errorf("payload container length is 0")
			}
			buffer := payload.Buffer[:payload.Len]
			m := nas.NewMessage()
			err := m.PlainNasDecode(&buffer)
			if err != nil {
				ue.Log.Errorln("PlainNasDecode returned:", err)
				return fmt.Errorf("failed to decode payload container")
			}
			nasMsg = m
			msgType = nasMsg.GsmHeader.GetMessageType()

		}

		m := &common.UeMessage{}

		// The MSB out of the 32 bytes represents event type, which in this case
		// is N1_EVENT
		m.Event = common.EventType(msgType) | common.N1_EVENT
		m.NasMsg = nasMsg

		// Simply notify SimUe about the received nas message. Later SimUe will
		// asynchrously send next event to RealUE informing about what to do with
		// the received NAS message
		SendToSimUe(ue, m)
		ue.Log.Infoln("Notified SimUe for message type:", msgType)
	}
	return nil
}

func HandleServiceRequestEvent(ue *context.RealUe,
	msg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Service Request Event")

	nasPdu, err := realue_nas.GetServiceRequest(ue)
	if err != nil {
		return fmt.Errorf("failed to handle service request event:", err)
	}

	// TS 24.501 Section 4.4.6 - Protection of Initial NAS signalling messages
	nasPdu, err = test.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtected, true)
	if err != nil {
		return fmt.Errorf("failed to encode with security:", err)
	}

	m := formUuMessage(common.SERVICE_REQUEST_EVENT, nasPdu)
	SendToSimUe(ue, m)
	ue.Log.Traceln("Sent Service Request Message to SimUe")
	return nil
}
