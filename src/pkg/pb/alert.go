package pb

//nolint:stylecheck,gochecknoglobals // Consistent with ProtoBuf.
//goland:noinspection GoSnakeCaseUsage
var (
	AlertType_PING = []AlertType{
		AlertType_PING_MISSING,
	}
	AlertType_PROBE = []AlertType{
		AlertType_PROBE_SLOW,
		AlertType_PROBE_FAILURE,
		AlertType_PROBE_TIMEOUT,
	}
	AlertStatus_ALL = []AlertStatus{
		AlertStatus_ALERT_OPEN,
		AlertStatus_ALERT_CLOSED_AUTO,
		AlertStatus_ALERT_CLOSED_MANUAL,
	}
	AlertStatus_OPEN = []AlertStatus{
		AlertStatus_ALERT_OPEN,
	}
	AlertStatus_CLOSED = []AlertStatus{
		AlertStatus_ALERT_CLOSED_AUTO,
		AlertStatus_ALERT_CLOSED_MANUAL,
	}
)
