// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"sync"

	"github.com/omec-project/gnbsim/common"
	"github.com/omec-project/gnbsim/factory"
	gnbctx "github.com/omec-project/gnbsim/gnodeb/context"
	"github.com/omec-project/gnbsim/logger"
	realuectx "github.com/omec-project/gnbsim/realue/context"

	"github.com/omec-project/nas/security"
	"github.com/sirupsen/logrus"
)

func init() {
	SimUeTable = make(map[string]*SimUe)
}

type LastIndexesProvision struct {
	Ausf int `json:"ausf"`
}
type SequenceNumberProvision struct {
	Sqn         string               `json:"sqn"`
	SqnScheme   string               `json:"sqnScheme"`
	LastIndexes LastIndexesProvision `json:"protectionParameterId"`
}

type SubscriberProvision struct {
	AuthenticationMethod          string                  `json:"authenticationMethod"`
	EncPermanentKey               string                  `json:"encPermanentKey"`
	ProtectionParameterId         string                  `json:"protectionParameterId"`
	SequenceNumber                SequenceNumberProvision `json:"sequenceNumber"`
	AuthenticationManagementField string                  `json:"authenticationManagementField"`
	AlgorithmId                   string                  `json:"algorithmId"`
	EncOpcKey                     string                  `json:"encOpcKey"`
	EncTopcKey                    string                  `json:"encTopcKey"`
	VectorGenerationInHss         bool                    `json:"vectorGenerationInHss"`
	N5gcAuthMethod                string                  `json:"n5gcAuthMethod"`
	RgAuthenticationInd           bool                    `json:"rgAuthenticationInd"`
	Supi                          string                  `json:"supi"`
}

// SimUe controls the flow of messages between RealUe and GnbUe as per the test
// profile. It is the central entry point for all events
type SimUe struct {
	Supi      string
	GnB       *gnbctx.GNodeB
	RealUe    *realuectx.RealUe
	Procedure common.ProcedureType
	WaitGrp   sync.WaitGroup

	// SimUe writes messages to Scenario routine on this channel
	WriteScenarioChan chan *common.InterfaceMessage

	// SimUe writes messages to RealUE on this channel
	WriteRealUeChan chan common.InterfaceMessage

	// SimUe writes messages to GnbUE on this channel
	WriteGnbUeChan chan common.InterfaceMessage

	// SimUe reads messages from other entities on this channel
	// Entities can be RealUe, GnbUe etc.
	ReadChan chan common.InterfaceMessage

	/* logger */
	Log *logrus.Entry
}

var SimUeTable map[string]*SimUe

func NewSimUe(supi string, ueModel string, gnb *gnbctx.GNodeB, result chan *common.InterfaceMessage) *SimUe {
	ueProfile, err := factory.AppConfig.Configuration.GetUeProfile(ueModel)
	if err != nil {
		return nil
	}
	simue := SimUe{}
	simue.GnB = gnb
	simue.Supi = supi
	simue.ReadChan = make(chan common.InterfaceMessage, 5)
	// TODO select prefered security algorithms
	simue.RealUe = realuectx.NewRealUe(supi,
		security.AlgCiphering128NEA0, security.AlgIntegrity128NIA2,
		simue.ReadChan, ueProfile.Plmn, ueProfile.Key, ueProfile.Opc, ueProfile.Nas.SeqNum,
		ueProfile.Nas.Dnn, ueProfile.Nas.SNssai)
	simue.WriteRealUeChan = simue.RealUe.ReadChan
	simue.WriteScenarioChan = result

	simue.Log = logger.SimUeLog.WithField(logger.FieldSupi, supi)

	simue.Log.Traceln("Created new SimUe context")
	SimUeTable[supi] = &simue
	return &simue
}

func GetSimUe(supi string) *SimUe {
	simue, found := SimUeTable[supi]
	if found == false {
		return nil
	}
	return simue
}
