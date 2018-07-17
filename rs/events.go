// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package rs

import (
	"github.com/AtlantPlatform/atlant-go/proto"
)

type EventType int

const (
	EventUnknown      EventType = EventType(proto.ANNOUNCETYPE_UNKNOWN)
	EventBeatTick     EventType = EventType(proto.ANNOUNCETYPE_BEATTICK)
	EventBeatInfo     EventType = EventType(proto.ANNOUNCETYPE_BEATINFO)
	EventRecordUpdate EventType = EventType(proto.ANNOUNCETYPE_RECORDUPDATE)
	EventStopAnnounce EventType = 999
)

func (e EventType) String() string {
	switch e {
	case EventBeatTick:
		return "beat-tick"
	case EventBeatInfo:
		return "beat-info"
	case EventRecordUpdate:
		return "record-update"
	case EventStopAnnounce:
		return "stop-announce"
	default:
		return "unknown"
	}
}

func EventFromTopic(topic string) EventType {
	switch topic {
	case EventBeatTick.String():
		return EventBeatTick
	case EventBeatInfo.String():
		return EventBeatInfo
	case EventRecordUpdate.String():
		return EventRecordUpdate
	default:
		return EventUnknown
	}
}

type EventAnnounce struct {
	Type     EventType      `json:"type"`
	Announce proto.Announce `json:"announce"`
}
