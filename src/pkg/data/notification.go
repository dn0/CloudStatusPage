package data

type Alerts []*Alert

type Incidents []*Incident

type Notification struct {
	Alerts    Alerts
	Incidents Incidents
}

func (alerts Alerts) ToNotification() *Notification {
	return &Notification{Alerts: alerts}
}

func (incs Incidents) ToNotification() *Notification {
	return &Notification{Incidents: incs}
}
