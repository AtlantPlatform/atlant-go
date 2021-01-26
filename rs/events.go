// Copyright 2017-21 Digital Asset Exchange Limited. All rights reserved.
// Use of this source code is governed by BSD 3-Clause "New" or "Revised"
// License (BSD 3) that can be found in the LICENSE file.

package rs

import (
	"github.com/AtlantPlatform/atlant-go/proto"
)

// EventType stores the code for event announcement type
type EventType int

const (
	// EventUnknown - code for unknown announcement (0)
	EventUnknown EventType = EventType(proto.ANNOUNCETYPE_UNKNOWN)
	// EventBeatTick - code for beat tick announcement (1)
	EventBeatTick EventType = EventType(proto.ANNOUNCETYPE_BEATTICK)
	// EventBeatInfo - code for beat information announcement (2)
	EventBeatInfo EventType = EventType(proto.ANNOUNCETYPE_BEATINFO)
	// EventRecordUpdate - code for announcement of record update (3)
	EventRecordUpdate EventType = EventType(proto.ANNOUNCETYPE_RECORDUPDATE)
	// EventStopAnnounce - code for stopping announcements
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

// EventFromTopic returns type of event from string
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

// EventAnnounce is a storage for serializable event announcement
type EventAnnounce struct {
	Type     EventType      `json:"type"`
	Announce proto.Announce `json:"announce"`
}
