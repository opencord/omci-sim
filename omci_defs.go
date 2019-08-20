/*
 * Copyright 2018-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

//
// OMCI definitions
//

// OmciMsgType represents a OMCI message-type
type OmciMsgType byte

func (t OmciMsgType) PrettyPrint() string {
	switch t {
	case Create:
		return "Create"
	case Delete:
		return "Delete"
	case Set:
		return "Set"
	case Get:
		return "Get"
	case GetAllAlarms:
		return "GetAllAlarms"
	case GetAllAlarmsNext:
		return "GetAllAlarmsNext"
	case MibUpload:
		return "MibUpload"
	case MibUploadNext:
		return "MibUploadNext"
	case MibReset:
		return "MibReset"
	case AlarmNotification:
		return "AlarmNotification"
	case AttributeValueChange:
		return "AttributeValueChange"
	case Test:
		return "Test"
	case StartSoftwareDownload:
		return "StartSoftwareDownload"
	case DownloadSection:
		return "DownloadSection"
	case EndSoftwareDownload:
		return "EndSoftwareDownload"
	case ActivateSoftware:
		return "ActivateSoftware"
	case CommitSoftware:
		return "CommitSoftware"
	case SynchronizeTime:
		return "SynchronizeTime"
	case Reboot:
		return "Reboot"
	case GetNext:
		return "GetNext"
	case TestResult:
		return "TestResult"
	case GetCurrentData:
		return "GetCurrentData"
	case SetTable:
		return "SetTable"
	default:
		log.Warnf("Cant't convert state %v to string", t)
		return string(t)
	}
}

const (
	// Message Types
	_                                 = iota
	Create                OmciMsgType = 4
	Delete                OmciMsgType = 6
	Set                   OmciMsgType = 8
	Get                   OmciMsgType = 9
	GetAllAlarms          OmciMsgType = 11
	GetAllAlarmsNext      OmciMsgType = 12
	MibUpload             OmciMsgType = 13
	MibUploadNext         OmciMsgType = 14
	MibReset              OmciMsgType = 15
	AlarmNotification     OmciMsgType = 16
	AttributeValueChange  OmciMsgType = 17
	Test                  OmciMsgType = 18
	StartSoftwareDownload OmciMsgType = 19
	DownloadSection       OmciMsgType = 20
	EndSoftwareDownload   OmciMsgType = 21
	ActivateSoftware      OmciMsgType = 22
	CommitSoftware        OmciMsgType = 23
	SynchronizeTime       OmciMsgType = 24
	Reboot                OmciMsgType = 25
	GetNext               OmciMsgType = 26
	TestResult            OmciMsgType = 27
	GetCurrentData        OmciMsgType = 28
	SetTable              OmciMsgType = 29 // Defined in Extended Message Set Only
)

const (
	// Managed Entity Class values
	EthernetPMHistoryData OmciClass = 24
	ONUG                  OmciClass = 256
	ANIG                  OmciClass = 263
	GEMPortNetworkCTP     OmciClass = 268
)

// OMCI Managed Entity Class
type OmciClass uint16

// OMCI Message Identifier
type OmciMessageIdentifier struct {
	Class    OmciClass
	Instance uint16
}

type OmciContent [32]byte

type OmciMessage struct {
	TransactionId uint16
	MessageType   OmciMsgType
	DeviceId      uint8
	MessageId     OmciMessageIdentifier
	Content       OmciContent
}

func ParsePkt(pkt []byte) (uint16, uint8, OmciMsgType, OmciClass, uint16, OmciContent, error) {
	var m OmciMessage

	r := bytes.NewReader(pkt)

	if err := binary.Read(r, binary.BigEndian, &m); err != nil {
		log.WithFields(log.Fields{
			"Packet": pkt,
			"omciMsg": fmt.Sprintf("%x", pkt),
		}).Errorf("Failed to read packet: %s", err)
		return 0, 0, 0, 0, 0, OmciContent{}, errors.New("Failed to read packet")
	}
	/*    Message Type = Set
	      0... .... = Destination Bit: 0x0
	      .1.. .... = Acknowledge Request: 0x1
	      ..0. .... = Acknowledgement: 0x0
	      ...0 1000 = Message Type: Set (8)
	*/

	log.WithFields(log.Fields{
		"TransactionId": m.TransactionId,
		"MessageType": m.MessageType.PrettyPrint(),
		"MeClass": m.MessageId.Class,
		"MeInstance": m.MessageId.Instance,
		"Conent": m.Content,
		"Packet": pkt,
	}).Tracef("Parsing OMCI Packet")

	return m.TransactionId, m.DeviceId, m.MessageType & 0x1F, m.MessageId.Class, m.MessageId.Instance, m.Content, nil
}
