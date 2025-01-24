package external

import "context"

type External struct {
	Notif interface {
		SendNotification(context.Context, NotifRequest) error
	}
	Validation interface{}
}

func NewExternal() External {
	return External{
		Notif: &notif{},
	}
}
