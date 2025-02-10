// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package notifyevent

import (
	"github.com/omec-project/pcf/logger"
)

type NotifyListener struct{}

func (l NotifyListener) Listen(event interface{}) {
	switch event := event.(type) {
	case SendSMpolicyUpdateNotifyEvent:
		event.Handle()
	case SendSMpolicyTerminationNotifyEvent:
		event.Handle()
	default:
		logger.NotifyEventLog.Warnf("registered an invalid user event: %T\n", event)
	}
}
