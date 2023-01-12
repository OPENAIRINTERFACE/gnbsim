// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package gnodeb

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/omec-project/gnbsim/common"
	"github.com/omec-project/gnbsim/factory"
	gnbctx "github.com/omec-project/gnbsim/gnodeb/context"
	"github.com/omec-project/gnbsim/gnodeb/idrange"
	"github.com/omec-project/gnbsim/gnodeb/transport"
	"github.com/omec-project/gnbsim/gnodeb/worker/gnbcpueworker"
	"github.com/omec-project/gnbsim/logger"
	"github.com/omec-project/gnbsim/util/test"

	"github.com/omec-project/idgenerator"
)

// Init initializes the GNodeB struct var and connects (SCTP only) to the default AMF
func Init(
	gnb *gnbctx.GNodeB,
	simGnbReadChan chan common.InterfaceMessage,
) error {
	gnb.Log = logger.GNodeBLog.WithField(logger.FieldGnb, gnb.GnbName)
	gnb.Log.Traceln("Inititializing GNodeB")
	gnb.Log.Infoln("GNodeB IP:", gnb.GnbN2Ip, "GNodeB Port:", gnb.GnbN2Port)

	gnb.CpTransport = transport.NewGnbCpTransport(gnb)
	gnb.UpTransport = transport.NewGnbUpTransport(gnb)
	err := gnb.UpTransport.Init()
	if err != nil {
		gnb.Log.Errorln("GnbUpTransport.Init returned", err)
		return fmt.Errorf("failed to initialize user plane transport")
	}
	gnb.GnbUes = gnbctx.NewGnbUeDao()
	gnb.GnbPeers = gnbctx.NewGnbPeerDao()
	start, end := idrange.GetIdRange()
	gnb.RanUeNGAPIDGenerator = idgenerator.NewGenerator(
		int64(start),
		int64(end),
	)
	gnb.DlTeidGenerator = idgenerator.NewGenerator(int64(start), int64(end))

	if gnb.Amf == nil {
		// LG TODO hardcoded
		amf, err := factory.AppConfig.Configuration.GetAmf("amf1")
		if err != nil {
			gnb.Log.Errorln("GetAmf returned:", err)
			return err
		}
		if amf.AmfIp == "" {
			// It is important to do this lookup just in time, not at simulation startup
			addrs, err := net.LookupHost(amf.AmfHostName)
			if err != nil {
				return fmt.Errorf(
					"failed to resolve amf host name: %v, err: %s",
					amf.AmfHostName,
					err,
				)
			}
			gnb.Amf = gnbctx.NewGnbAmf(
				addrs[0],
				gnbctx.NGAP_SCTP_PORT,
				simGnbReadChan,
			)
		}
	}

	gnb.Amf.Init()

	err = gnb.CpTransport.ConnectToPeer(gnb.Amf)
	if err != nil {
		gnb.Log.Errorln("ConnectToPeer returned:", err)
		return fmt.Errorf("failed to connect to amf")
	}

	gnb.Log.Tracef("GNodeB Initialized %v ", gnb)
	return nil
}

func QuitGnb(gnb *gnbctx.GNodeB) {
	log.Println("Shutting Down GNodeB:", gnb.GnbName)
	close(gnb.Quit)
}

// SendNGSetup sends the NGSetupRequest to the provided GnbAmf.
func SendNgSdu(gnb *gnbctx.GNodeB, amf *gnbctx.GnbAmf, sdu []byte) error {
	gnb.Log.Traceln("Sending NG SDU")
	err := gnb.CpTransport.SendToPeer(amf, sdu)
	if err != nil {
		gnb.Log.Errorln("SendToPeer returned:", err)
		return fmt.Errorf("failed to send NG SDU")
	}
	return nil
}

// RequestConnection should be called by UE that is willing to connect to this GNodeB
func RequestConnection(
	gnb *gnbctx.GNodeB,
	uemsg *common.UuMessage,
) (chan common.InterfaceMessage, *gnbctx.GnbCpUe, error) {
	ranUeNgapID, err := gnb.AllocateRanUeNgapID()
	if err != nil {
		gnb.Log.Errorln("AllocateRanUeNgapID returned:", err)
		return nil, nil, fmt.Errorf("failed to allocate ran ue ngap id")
	}

	gnbUe := gnbctx.NewGnbCpUe(ranUeNgapID, gnb, gnb.Amf)
	gnb.GnbUes.AddGnbCpUe(ranUeNgapID, gnbUe)

	// TODO: Launching a GO Routine for gNB and handling the waitgroup
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		gnbcpueworker.Init(gnbUe)
	}()
	//Channel on which UE can write message to GnbUe and from which GnbUe will
	//be reading.
	ch := gnbUe.ReadChan
	ch <- uemsg
	return ch, gnbUe, nil
}

func GetInitialUEMessage(
	gnb *gnbctx.GNodeB,
	gnbCpUe *gnbctx.GnbCpUe,
	nasBytes []byte,
) (n2Bytes []byte, err error) {

	msg, err := test.GetInitialUEMessage(gnbCpUe.GnbUeNgapId, nasBytes, "",
		gnb.SupportedTaList[0].Tac, gnb.NrCgiCellList[0])
	if err != nil {
		gnb.Log.Errorln("GetInitialUEMessage failed:", err)
		return nil, err
	}
	return msg, nil
}

func GetUplinkNASTransport(
	gnb *gnbctx.GNodeB,
	gnbCpUe *gnbctx.GnbCpUe,
	nasBytes []byte,
) (n2Bytes []byte, err error) {

	msg, err := test.GetUplinkNASTransport(
		gnbCpUe.AmfUeNgapId,
		gnbCpUe.GnbUeNgapId,
		nasBytes,
	)
	if err != nil {
		gnb.Log.Errorln("GetInitialUEMessage failed:", err)
		return nil, err
	}
	return msg, nil
}

func GetInitialContextSetupResponse(
	gnb *gnbctx.GNodeB,
	gnbCpUe *gnbctx.GnbCpUe,
) (n2Bytes []byte, err error) {

	msg, err := test.GetInitialContextSetupResponse(
		gnbCpUe.AmfUeNgapId,
		gnbCpUe.GnbUeNgapId,
	)
	if err != nil {
		gnb.Log.Errorln("GetInitialContextSetupResponse failed:", err)
		return nil, err
	}
	return msg, nil
}
