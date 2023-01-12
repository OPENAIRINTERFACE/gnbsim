// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package realue

import (
	"fmt"
	"net"

	"github.com/omec-project/gnbsim/common"
	realuectx "github.com/omec-project/gnbsim/realue/context"
	"github.com/omec-project/gnbsim/realue/worker/pdusessworker"

	realue_nas "github.com/omec-project/gnbsim/realue/nas"

	"github.com/omec-project/nas"
	"github.com/omec-project/nas/nasConvert"
	"github.com/omec-project/nas/nasMessage"
	"github.com/omec-project/nas/nasTestpacket"
	"github.com/omec-project/nas/nasType"
	"github.com/omec-project/openapi/models"
)

// TODO Remove the hardcoding
const (
	SN_NAME                        string = "5G:mnc093.mcc208.3gppnetwork.org"
	SWITCH_OFF                     uint8  = 0
	REQUEST_TYPE_EXISTING_PDU_SESS uint8  = 0x02
)

func HandleRegistrationAccept(ue *realuectx.RealUe,
	msg *nasMessage.RegistrationAccept) error {

	var guti []uint8
	if msg.GUTI5G != nil {
		guti = msg.GUTI5G.Octet[:]
	}

	_, ue.Guti = nasConvert.GutiToString(guti)
	return nil
}

func HandleDeregRequestEvent(ue *realuectx.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	if ue.Guti == "" {
		ue.Log.Errorln("guti not allocated")
		return fmt.Errorf(
			"failed to create deregistration request: guti not unallocated",
		)
	}
	gutiNas := nasConvert.GutiToNas(ue.Guti)
	mobileIdentity5GS := nasType.MobileIdentity5GS{
		Len:    11, // 5g-guti
		Buffer: gutiNas.Octet[:],
	}

	nasPdu := nasTestpacket.GetDeregistrationRequest(nasMessage.AccessType3GPP,
		SWITCH_OFF, uint8(ue.NgKsi.Ksi), mobileIdentity5GS)
	nasPdu, err = realue_nas.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		ue.Log.Errorln("EncodeNasPduWithSecurity() returned:", err)
		return fmt.Errorf("failed to encrypt deregistration request message")
	}

	m := FormUuMessage(
		common.N1_ENCODED_EVENT+common.NAS_5GMM_DEREGISTRATION_REQUEST_UE_ORIG,
		nasPdu,
	)
	// LG SendToSimUe(ue, m)
	ue.Log.Traceln("TODO LG To avoid comment", m)
	ue.Log.Traceln("Sent UE Initiated Deregistration Request message to SimUe")
	return nil
}

func HandlePduSessEstRequestEvent(ue *realuectx.RealUe,
	msg common.InterfaceMessage) (err error) {

	// sNssai := models.Snssai{
	// 	Sst: 1,
	// 	Sd:  "010203",
	// }
	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionEstablishmentRequest(10,
		nasMessage.ULNASTransportRequestTypeInitialRequest, ue.Dnn, ue.SNssai)

	nasPdu, err = realue_nas.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		fmt.Println(
			"Failed to encrypt PDU Session Establishment Request Message",
			err,
		)
		return
	}

	m := FormUuMessage(
		common.N1_ENCODED_EVENT+common.NAS_5GSM_PDU_SESSION_ESTABLISHMENT_REQUEST,
		nasPdu,
	)
	// LG SendToSimUe(ue, m)
	ue.Log.Traceln("TODO LG To avoid comment", m)
	return nil
}

func HandlePduSessEstAcceptEvent(ue *realuectx.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
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

	pduSess := realuectx.NewPduSession(ue, int64(nasMsg.PDUSessionID.Octet))
	pduSess.PduSessType = pduSessType
	pduSess.SscMode = nasMsg.GetSSCMode()
	pduSess.PduAddress = pduAddr
	pduSess.WriteRealUeChan = ue.ReadChan
	ue.AddPduSession(int64(pduSess.PduSessId), pduSess)
	ue.Log.Infoln("PDU Session ID:", pduSess.PduSessId)
	ue.Log.Infoln("PDU Session Type:", pduSess.PduSessType)
	ue.Log.Infoln("SSC Mode:", pduSess.SscMode)
	ue.Log.Infoln("PDU Address:", pduAddr.String())

	return nil
}

func HandlePduSessReleaseRequestEvent(ue *realuectx.RealUe,
	msg common.InterfaceMessage) (err error) {

	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionReleaseRequest(10)

	nasPdu, err = realue_nas.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		fmt.Println(
			"Failed to encrypt PDU Session Release Request Message",
			err,
		)
		return
	}

	m := FormUuMessage(
		common.N1_ENCODED_EVENT+common.NAS_5GSM_PDU_SESSION_RELEASE_REQUEST,
		nasPdu,
	)
	// LG SendToSimUe(ue, m)
	ue.Log.Traceln("TODO LG To avoid comment", m)
	return nil
}

func HandlePduSessReleaseCompleteEvent(ue *realuectx.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	nasMsg := msg.NasMsg.PDUSessionReleaseCommand
	if nasMsg == nil {
		ue.Log.Errorln("PDUSessionReleaseCommand is nil")
		return fmt.Errorf("invalid NAS Message")
	}

	pduSessId := nasMsg.PDUSessionID.Octet
	ue.Log.Infoln("PDU Session Release Command, PDU Session ID:", pduSessId)

	pduSess, err := ue.GetPduSession(int64(pduSessId))
	if err != nil {
		return fmt.Errorf("failed to fetch PDU session:%v", err)
	}

	quitMsg := &common.UeMessage{}
	quitMsg.Event = common.QUIT_EVENT
	pduSess.ReadCmdChan <- quitMsg

	nasPdu := nasTestpacket.GetUlNasTransport_PduSessionReleaseComplete(
		pduSessId,
		REQUEST_TYPE_EXISTING_PDU_SESS,
		"",
		nil,
	)

	nasPdu, err = realue_nas.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true)
	if err != nil {
		fmt.Println(
			"Failed to encrypt PDU Session Release Request Message",
			err,
		)
		return
	}

	m := FormUuMessage(
		common.N1_ENCODED_EVENT+common.NAS_5GSM_PDU_SESSION_RELEASE_COMPLETE,
		nasPdu,
	)
	// LG SendToSimUe(ue, m)
	ue.Log.Traceln("TODO LG To avoid comment", m)
	return nil
}

func HandleDataBearerSetupRequestEvent(ue *realuectx.RealUe,
	intfcMsg common.InterfaceMessage) (response *common.UuMessage, err error) {

	msg := intfcMsg.(*common.UuMessage)
	for _, item := range msg.DBParams {
		/* Currently gNB also adds failed pdu session ids in the list.
		   pdu sessions are marked failed during decoding. real ue simply
		   returns the same list back by marking any failed pdu sessions on
		   its side. This consolidated list can be used by gnb to form
		   PDUSession Resource Setup/Failed To Setup Response list
		*/
		if item.PduSess.Success {
			pduSess, err := ue.GetPduSession(item.PduSess.PduSessId)
			if err != nil {
				ue.Log.Warnln("Failed to fetch PDU Session:", err)
				item.PduSess.Success = false
				continue
			}

			if !pduSess.Launched {
				pduSess.Launched = true
				ue.WaitGrp.Add(1)
				go pdusessworker.Init(pduSess, &ue.WaitGrp)
			}

			initMsg := &common.UeMessage{}
			initMsg.Event = common.INIT_EVENT
			initMsg.CommChan = item.CommChan
			pduSess.ReadCmdChan <- initMsg

			/* gNb can use this channel to send DL packets for this PDU session */
			item.CommChan = pduSess.ReadDlChan
		}
	}

	response = &common.UuMessage{}
	response.Event = common.DATA_BEARER_SETUP_RESPONSE_EVENT
	response.DBParams = msg.DBParams
	response.TriggeringEvent = msg.TriggeringEvent
	/* LG COMMENT kept it temp for memo
	ue.WriteSimUeChan <- response
	*/
	return
}

func HandleDataPktGenRequestEvent(ue *realuectx.RealUe,
	msg common.InterfaceMessage) (err error) {

	if len(ue.PduSessions) == 0 {
		err = fmt.Errorf("Can't generate traffic, no PDU sessions")
	}
	for _, v := range ue.PduSessions {
		v.ReadCmdChan <- msg
	}

	return
}

func HandleDataPktGenSuccessEvent(ue *realuectx.RealUe,
	msg common.InterfaceMessage) (err error) {
	ue.WriteSimUeChan <- msg
	return nil
}

func HandleConnectionReleaseRequestEvent(ue *realuectx.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {
	msg := intfcMsg.(*common.UuMessage)

	for _, pdusess := range ue.PduSessions {
		pdusess.ReadCmdChan <- msg
	}

	return nil
}

func HandleErrorEvent(ue *realuectx.RealUe,
	intfcMsg common.InterfaceMessage) (err error) {

	// LG COMMENT SendToSimUe(ue, intfcMsg)
	return nil
}

func HandleQuitEvent(
	ue *realuectx.RealUe,
	intfcMsg common.InterfaceMessage,
) (err error) {
	ue.WriteSimUeChan = nil
	for _, pdusess := range ue.PduSessions {
		pdusess.ReadCmdChan <- intfcMsg
	}
	ue.PduSessions = nil
	ue.WaitGrp.Wait()
	ue.Log.Infoln("Real UE terminated")
	return nil
}

func HandleDlInfoTransferEvent(
	ue *realuectx.RealUe,
	nasPdu []byte,
) (*nas.Message, error) {

	nasMsg, err := realue_nas.NASDecode(
		ue,
		nas.GetSecurityHeaderType(nasPdu),
		nasPdu,
	)
	if err != nil {
		ue.Log.Errorln("Failed to decode dowlink NAS Message due to", err)
		return nil, err
	}
	msgType := nasMsg.GmmHeader.GetMessageType()
	ue.Log.Infoln("Received Message Type:", msgType)

	if msgType == nas.MsgTypeDLNASTransport {
		ue.Log.Info(
			"Payload contaner type:",
			nasMsg.GmmMessage.DLNASTransport.SpareHalfOctetAndPayloadContainerType,
		)
		payload := nasMsg.GmmMessage.DLNASTransport.PayloadContainer
		if payload.Len == 0 {
			return nasMsg, fmt.Errorf("payload container length is 0")
		}
		buffer := payload.Buffer[:payload.Len]
		m := nas.NewMessage()
		err := m.PlainNasDecode(&buffer)
		if err != nil {
			ue.Log.Errorln("PlainNasDecode returned:", err)
			return nasMsg, fmt.Errorf("failed to decode payload container")
		}
		nasMsg = m
		msgType = nasMsg.GsmHeader.GetMessageType()

	}
	return nasMsg, err
}

func HandleServiceRequestEvent(ue *realuectx.RealUe,
	msg common.InterfaceMessage) (err error) {

	nasPdu, err := realue_nas.GetServiceRequest(ue)
	if err != nil {
		return fmt.Errorf("failed to handle service request event: %v", err)
	}

	// TS 24.501 Section 4.4.6 - Protection of Initial NAS signalling messages
	nasPdu, err = realue_nas.EncodeNasPduWithSecurity(ue, nasPdu,
		nas.SecurityHeaderTypeIntegrityProtected, true)
	if err != nil {
		return fmt.Errorf("failed to encode with security: %v", err)
	}

	m := FormUuMessage(
		common.N1_ENCODED_EVENT+common.NAS_5GMM_SERVICE_REQUEST,
		nasPdu,
	)
	// LG COMMENT SendToSimUe(ue, m)
	ue.Log.Traceln("TODO LG To avoid comment", m)
	return nil
}

func HandleNwDeregAcceptEvent(
	ue *realuectx.RealUe,
	msg common.InterfaceMessage,
) (err error) {
	ue.Log.Traceln("Generating Dereg Accept Message")
	nasPdu := nasTestpacket.GetDeregistrationAccept()

	nasPdu, err = realue_nas.EncodeNasPduWithSecurity(
		ue,
		nasPdu,
		nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext,
		true,
	)
	if err != nil {
		ue.Log.Errorln("EncodeNasPduWithSecurity() returned:", err)
		return fmt.Errorf("failed to encrypt security mode complete message")
	}

	m := FormUuMessage(
		common.N1_ENCODED_EVENT+common.NAS_5GMM_DEREGISTRATION_ACCEPT_UE_TERM,
		nasPdu,
	)
	// LG COMMENT SendToSimUe(ue, m)
	ue.Log.Traceln("Sent Dereg Accept UE Terminated Message to SimUe", m)
	return nil
}
