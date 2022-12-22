// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package realue

import (
	"github.com/omec-project/gnbsim/common"
	realuectx "github.com/omec-project/gnbsim/realue/context"
	"github.com/omec-project/gnbsim/util/test"
)

func Init(ue *realuectx.RealUe) error {

	ue.AuthenticationSubs = test.GetAuthSubscription(ue.Key, ue.Opc, "", ue.SeqNum)

	// LG HandleEvents(ue)
	return nil
}

/* LG func HandleEvents(ue *realuectx.RealUe) (err error) {

	for msg := range ue.ReadChan {
		event := msg.GetEventType()
		ue.Log.Infoln("Handling:", event)

		switch event {
		case common.N1_ENCODE_EVENT + common.NAS_5GMM_DEREGISTRATION_REQUEST_UE_ORIG:
			err = HandleDeregRequestEvent(ue, msg)
		case common.DL_INFO_TRANSFER_EVENT:
			err = HandleDlInfoTransferEvent(ue, msg)
		case common.N1_ENCODE_EVENT + common.NAS_5GSM_PDU_SESSION_ESTABLISHMENT_REQUEST:
			err = HandlePduSessEstRequestEvent(ue, msg)
		case common.N1_ENCODE_EVENT + common.NAS_5GSM_PDU_SESSION_RELEASE_REQUEST:
			err = HandlePduSessReleaseRequestEvent(ue, msg)
		case common.N1_ENCODE_EVENT + common.NAS_5GSM_PDU_SESSION_RELEASE_COMPLETE:
			err = HandlePduSessReleaseCompleteEvent(ue, msg)
		case common.N1_RECV_SDU_EVENT + common.NAS_5GSM_PDU_SESSION_ESTABLISHMENT_ACCEPT:
			err = HandlePduSessEstAcceptEvent(ue, msg)
		case common.DATA_BEARER_SETUP_REQUEST_EVENT:
			err = HandleDataBearerSetupRequestEvent(ue, msg)
		case common.DATA_PKT_GEN_REQUEST_EVENT:
			err = HandleDataPktGenRequestEvent(ue, msg)
		case common.DATA_PKT_GEN_SUCCESS_EVENT:
			err = HandleDataPktGenSuccessEvent(ue, msg)
		case common.N1_ENCODE_EVENT + common.NAS_5GMM_SERVICE_REQUEST:
			err = HandleServiceRequestEvent(ue, msg)
		case common.CONNECTION_RELEASE_REQUEST_EVENT:
			err = HandleConnectionReleaseRequestEvent(ue, msg)
		case common.N1_ENCODE_EVENT + common.NAS_5GMM_DEREGISTRATION_ACCEPT_UE_TERM:
			err = HandleNwDeregAcceptEvent(ue, msg)
		case common.ERROR_EVENT:
			HandleErrorEvent(ue, msg)
		case common.QUIT_EVENT:
			HandleQuitEvent(ue, msg)
			return nil
		default:
			ue.Log.Warnln("Event", event, "is not supported")
		}

		if err != nil {
			ue.Log.Errorln("real ue failed:", event, ":", err)
			msg := &common.UeMessage{}
			msg.Error = err
			msg.Event = common.ERROR_EVENT
			err = nil
			HandleErrorEvent(ue, msg)
		}
	}

	return nil
} */

func FormUuMessage(event common.EventType, nasPdu []byte) *common.UuMessage {
	msg := &common.UuMessage{}
	msg.Event = event
	msg.NasPdus = append(msg.NasPdus, nasPdu)
	return msg
}

/*  LG func SendToSimUe(ue *realuectx.RealUe,
	msg common.InterfaceMessage) {

	ue.Log.Traceln("Sending", msg.GetEventType(), "to SimUe")
	ue.WriteSimUeChan <- msg
}
*/
