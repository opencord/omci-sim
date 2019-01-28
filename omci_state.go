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
	"errors"
	"fmt"
)

type OnuOmciState struct {
	gemPortId     uint16
	mibUploadCtr  uint16
	uniGInstance  uint8
	tcontInstance uint8
	pptpInstance  uint8
	state         istate
}

type istate int

// TODO - Needs to reflect real ONU/OMCI state
const (
	INCOMPLETE istate = iota
	DONE
)

var OnuOmciStateMap = map[OnuKey]*OnuOmciState{}

func NewOnuOmciState() *OnuOmciState {
	return &OnuOmciState{gemPortId: 0, mibUploadCtr: 0, uniGInstance: 1, tcontInstance: 0, pptpInstance: 1}
}

func GetOnuOmciState(intfId uint32, onuId uint32) istate {
	key := OnuKey{intfId, onuId}
	if onu, ok := OnuOmciStateMap[key]; ok {
		return onu.state
	} else {
		return INCOMPLETE
	}
}

func GetGemPortId(intfId uint32, onuId uint32) (uint16, error) {
	key := OnuKey{intfId, onuId}
	if OnuOmciState, ok := OnuOmciStateMap[key]; ok {
		return OnuOmciState.gemPortId, nil
	}
	errmsg := fmt.Sprintf("Failed to find a key in OnuOmciStateMap key{intfid:%d, onuid:%d}", intfId, onuId)
	return 0, errors.New(errmsg)
}
