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

import "log"


func OmciSim(intfId uint32, onuId uint32, request []byte) ([]byte, error) {
	var resp []byte

	transactionId, deviceId, msgType, class, instance, content, err := ParsePkt(request)
	if err != nil {
		log.Printf("ONU {intfid:%d, onuid:%d} - Cannot parse omci msg", intfId, onuId)
		return resp, &OmciError{"Cannot parse omci msg"}
	}

	log.Printf("ONU {intfid:%d, onuid:%d} - OmciRun - transactionId: %d msgType: %d, ME Class: %d, ME Instance: %d",
		intfId, onuId, transactionId, msgType, class, instance)

	key := OnuKey{intfId, onuId}
	if _, ok := OnuOmciStateMap[key]; !ok {
		OnuOmciStateMap[key] = NewOnuOmciState()
	}

	if _, ok := Handlers[msgType]; !ok {
		log.Printf("ONU {intfid:%d, onuid:%d} - Ignore omci msg (msgType %d not handled)", intfId, onuId, msgType)
		return resp, &OmciError{"Unimplemented omci msg"}
	}

	resp, err = Handlers[msgType](class, content, key)
	if err != nil {
		log.Printf("ONU {intfid:%d, onuid:%d} - Unable to send a successful response, error:%s", intfId, onuId, err)
		return resp, nil
	}

	// In the OMCI message, first 2-bytes is the Transaction Correlation ID
	resp[0] = byte(transactionId >> 8)
	resp[1] = byte(transactionId & 0xFF)
	resp[2] = 0x2<<4 | byte(msgType) // Upper nibble 0x2 is fixed (0010), Lower nibbles defines the msg type (i.e., mib-upload, mib-upload-next, etc)
	resp[3] = deviceId

	// for create, get and set
	if ((msgType & 0xFF) != MibUploadNext) && ((msgType & 0xFF) != MibReset) && ((msgType & 0xFF) != MibUpload) {
		// Common fields for create, get, and set
		resp[4] = byte(class >> 8)
		resp[5] = byte(class & 0xFF)
		resp[6] = byte(instance >> 8)
		resp[7] = byte(instance & 0xFF)
		resp[8] = 0 // Result: Command Processed Successfully

		// Hardcoding class specific values for Get
		if (class == 0x82) && ((msgType & 0x0F) == Get) {
			resp[9] = 0
			resp[10] = 0x78

		} else if (class == 0x2F) && ((msgType & 0x0F) == Get) {
			resp[9] = 0x0F
			resp[10] = 0xB8
		} else if (class == 0x138) && ((msgType & 0x0F) == Get) {
			resp[9] = content[0] // 0xBE
			resp[10] = 0x00
		}
	}

	log.Printf("OMCI-SIM Response %+x\n", resp)

	return resp, nil
}
