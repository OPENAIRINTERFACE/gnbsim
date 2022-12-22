// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package gnbcpueworker

import (
	"github.com/omec-project/gnbsim/common"
	gnbctx "github.com/omec-project/gnbsim/gnodeb/context"
	"github.com/omec-project/ngap/ngapType"
)

func Init(gnbue *gnbctx.GnbCpUe) {
	HandleEvents(gnbue)
}

func HandleEvents(gnbue *gnbctx.GnbCpUe) (err error) {

	for msg := range gnbue.ReadChan {
		evt := msg.GetEventType()
		gnbue.Log.Infoln("Handling event:", evt)

		switch msg.GetEventType() {
		case common.CONNECTION_REQUEST_EVENT:
			HandleConnectRequest(gnbue, msg)
		case common.N2_SEND_SDU_EVENT:
			SendToPeer(gnbue, msg)
		case common.N1_SEND_SDU_EVENT | common.NAS_5GMM_REGISTRATION_REQUEST, common.N1_SEND_SDU_EVENT + common.NAS_5GMM_SERVICE_REQUEST:
			HandleInitialUEMessage(gnbue, msg)
		case common.UL_INFO_TRANSFER_EVENT:
			HandleUlInfoTransfer(gnbue, msg)
		case common.DATA_BEARER_SETUP_RESPONSE_EVENT:
			HandleDataBearerSetupResponse(gnbue, msg)
		case common.N2_RECV_SDU_EVENT | common.DOWNLINK_NAS_TRANSPORT_EVENT:
			HandleDownlinkNasTransport(gnbue, msg)
		case common.N2_RECV_SDU_EVENT | common.INITIAL_CTX_SETUP_REQUEST_EVENT:
			HandleInitialContextSetupRequest(gnbue, msg)
		case common.PDU_SESS_RESOURCE_SETUP_REQUEST_EVENT:
			HandlePduSessResourceSetupRequest(gnbue, msg)
		case common.PDU_SESS_RESOURCE_RELEASE_COMMAND_EVENT:
			HandlePduSessResourceReleaseCommand(gnbue, msg)
		case common.UE_CTX_RELEASE_COMMAND_EVENT:
			HandleUeCtxReleaseCommand(gnbue, msg)
		case common.TRIGGER_AN_RELEASE_EVENT:
			HandleRanConnectionRelease(gnbue, msg)
		case common.QUIT_EVENT:
			HandleQuitEvent(gnbue, msg)
			return
		default:
			gnbue.Log.Infoln("Event", evt, "is not supported")
		}

		// TODO: Need to return and handle errors from handlers
	}
	return nil
}

func SendToSimUe(gnbue *gnbctx.GnbCpUe, event common.EventType, ngapPdu *ngapType.NGAPPDU, ngapProcedureCode int64, nasPdu *ngapType.NASPDU) {
	gnbue.Log.Traceln("Sending event", event, "to SimUe")
	uemsg := common.N1N2Message{}
	uemsg.Event = event
	uemsg.NgapPdu = ngapPdu
	uemsg.NgapProcedureCode = ngapProcedureCode
	uemsg.NasPdu = nasPdu
	gnbue.WriteUeChan <- &uemsg
}
