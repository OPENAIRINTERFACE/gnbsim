// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package simue

import (
	"fmt"

	"github.com/omec-project/gnbsim/common"
	simuectx "github.com/omec-project/gnbsim/simue/context"
)

func HandleRegRequestEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Infoln("Initiating Registration Procedure")
	msg := &common.UeMessage{}
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_REGISTRATION_REQUEST
	SendToRealUe(ue, msg)
	return nil
}

func HandleRegRequestEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.N1_SEND_SDU_EVENT + common.NAS_5GMM_REGISTRATION_REQUEST
	SendToGnbUe(ue, msg)
	return nil
}

func HandleRegCompleteEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := &common.UeMessage{}
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_REGISTRATION_COMPLETE
	SendToRealUe(ue, msg)
	return nil
}

func HandleRegCompleteEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Registration Complete to the network")
	return nil
}

func HandleAuthRequestEvent(ue *simuectx.SimUe,
	intfMsg common.InterfaceMessage) (err error) {

	msg := intfMsg.(*common.UeMessage)
	SendToScenario(ue, msg)
	return nil
}

func HandleAuthResponseEvent(ue *simuectx.SimUe,
	intfMsg common.InterfaceMessage) (err error) {

	msg := &common.UeMessage{}
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_AUTHENTICATION_RESPONSE
	SendToRealUe(ue, msg)
	return nil
}

func HandleAuthResponseEncodedEvent(ue *simuectx.SimUe,
	intfMsg common.InterfaceMessage) (err error) {

	msg := intfMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Registration Complete to the network")
	return nil
}

func HandleSecModCommandEvent(ue *simuectx.SimUe,
	intfMsg common.InterfaceMessage) (err error) {

	msg := intfMsg.(*common.UeMessage)
	SendToScenario(ue, msg)
	return nil
}

func HandleSecModCompleteEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Security Mode Complete Event")
	msg := &common.UeMessage{}
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_SECURITY_MODE_COMPLETE
	SendToRealUe(ue, msg)
	return nil
}

func HandleSecModCompleteEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Security Mode Complete Event")
	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Security Mode Complete to the network")
	return nil
}

func HandleSecModRejectEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Security Mode Reject Event")
	msg := &common.UeMessage{}
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_SECURITY_MODE_REJECT
	SendToRealUe(ue, msg)
	return nil
}

func HandleSecModRejectEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Security Mode Reject Event")

	msg := intfcMsg.(*common.UuMessage)

	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Security Mode Reject to the network")
	return nil
}

func HandleDeregRequestEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_DEREGISTRATION_REQUEST_UE_ORIG
	SendToRealUe(ue, msg)
	return nil
}

func HandleDeregRequestEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Deregistration Request to the network")
	return nil
}

func HandleDeregAcceptEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	return nil
}

func HandlePduSessEstRequestEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GSM_PDU_SESSION_ESTABLISHMENT_REQUEST
	SendToRealUe(ue, msg)
	return nil
}

func HandlePduSessEstRequestEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	return nil
}

func HandlePduSessEstAcceptEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	SendToScenario(ue, msg)
	return nil
}

func HandlePduSessEstRejectEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	SendToScenario(ue, msg)
	return nil
}

func HandlePduSessReleaseRequestEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GSM_PDU_SESSION_RELEASE_REQUEST
	SendToRealUe(ue, msg)
	return nil
}

func HandlePduSessReleaseRequestEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	return nil
}

func HandlePduSessReleaseCommandEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	SendToScenario(ue, msg)
	return nil
}

func HandlePduSessReleaseCompleteEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UeMessage)
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GSM_PDU_SESSION_RELEASE_COMPLETE
	SendToRealUe(ue, msg)
	return nil
}

func HandlePduSessReleaseCompleteEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := intfcMsg.(*common.UuMessage)
	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	return nil
}

func HandleDlInfoTransferEvent(ue *simuectx.SimUe,
	msg common.InterfaceMessage) (err error) {

	SendToRealUe(ue, msg)
	return nil
}

func HandleDataBearerSetupRequestEvent(ue *simuectx.SimUe,
	msg common.InterfaceMessage) (err error) {

	SendToRealUe(ue, msg)
	return nil
}

func HandleDataBearerSetupResponseEvent(ue *simuectx.SimUe,
	msg common.InterfaceMessage) (err error) {

	SendToGnbUe(ue, msg)
	return nil
}

func HandleDataBearerReleaseRequestEvent(ue *simuectx.SimUe,
	msg common.InterfaceMessage) (err error) {
	// This event is sent by gNB component after it has sent
	// PDU Session Resource Release Complete over N2, However the PDU Sesson
	// routines in the RealUE will be terminated while processing PDU Session
	// Release Complete which will also release the communication links
	// (go channels) with the gNB
	//Current Procedure is complete. Move to next one
	return nil
}

func HandleDataPktGenSuccessEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	//Current Procedure is complete. Move to next one
	return nil
}

func HandleDataPktGenFailureEvent(ue *simuectx.SimUe,
	msg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("HandleDataPktGenFailureEvent")
	SendToScenario(ue, msg)
	return nil
}

func HandleServiceRequestEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	msg := &common.UeMessage{}
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_SERVICE_REQUEST
	SendToRealUe(ue, msg)
	return nil
}

func HandleServiceRequestEncodedEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	err = ConnectToGnb(ue)
	if err != nil {
		return fmt.Errorf("failed to connect gnb %v:", err)
	}

	SendToGnbUe(ue, intfcMsg)

	ue.Log.Traceln("Sent Service Request Event to the network")
	return nil
}

func HandleConnectionReleaseRequestEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Errorln("TODO HandleConnectionReleaseRequestEvent")

	msg := intfcMsg.(*common.UuMessage)
	ue.WriteGnbUeChan = nil
	SendToRealUe(ue, msg)
	return nil
}

func HandleNwDeregRequestEvent(ue *simuectx.SimUe, intfcMsg common.InterfaceMessage) (err error) {

	//msg := intfcMsg.(*common.UeMessage)
	ue.Log.Errorln("TODO HandleNwDeregRequestEvent")
	return nil
}

func HandleNwDeregAcceptEvent(ue *simuectx.SimUe, intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Dereg Accept Event")
	msg := &common.UeMessage{}
	msg.Event = common.N1_ENCODE_EVENT + common.NAS_5GMM_DEREGISTRATION_ACCEPT_UE_TERM
	SendToRealUe(ue, msg)
	return nil
}

func HandleNwDeregAcceptDecodedEvent(ue *simuectx.SimUe, intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Traceln("Handling Dereg Accept Event")

	msg := intfcMsg.(*common.UuMessage)

	msg.Event = common.UL_INFO_TRANSFER_EVENT
	SendToGnbUe(ue, msg)
	ue.Log.Traceln("Sent Dereg Accept to the network")
	return nil
}

func HandleErrorEvent(ue *simuectx.SimUe,
	intfcMsg common.InterfaceMessage) (err error) {

	ue.Log.Errorln("TODO HandleErrorEvent")
	//ue.Log.Traceln("debug3")
	//SendToScenario(ue, common.PROC_FAIL_EVENT, intfcMsg.GetErrorMsg())

	msg := &common.UuMessage{}
	msg.Event = common.QUIT_EVENT
	HandleQuitEvent(ue, msg)
	return nil
}

func HandleQuitEvent(ue *simuectx.SimUe,
	msg common.InterfaceMessage) (err error) {
	if ue.WriteGnbUeChan != nil {
		SendToGnbUe(ue, msg)
	}
	SendToRealUe(ue, msg)
	ue.WriteRealUeChan = nil
	ue.WaitGrp.Wait()
	ue.Log.Infoln("Sim UE terminated")
	return nil
}
