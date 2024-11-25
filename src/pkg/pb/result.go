package pb

import (
	"context"
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//nolint:stylecheck,gochecknoglobals // Consistent with ProtoBuf.
//goland:noinspection GoSnakeCaseUsage
var (
	ResultStatus_OK = []ResultStatus{
		ResultStatus_RESULT_SUCCESS,
	}
	ResultStatus_ERROR = []ResultStatus{
		ResultStatus_RESULT_FAILURE,
		ResultStatus_RESULT_TIMEOUT,
	}
)

type ResultTime struct {
	Start time.Time
	Took  time.Duration
}

type resultFun func(context.Context) error

//nolint:revive // Context should not be the first argument as this is a decorator.
func Timeit(fun resultFun, ctx context.Context) (ResultTime, error) {
	start := time.Now()
	err := fun(ctx)
	return ResultTime{Start: start, Took: time.Since(start)}, err
}

func EmptyResultTime() ResultTime {
	return ResultTime{Start: time.Now(), Took: time.Duration(0)}
}

func NewResult(action uint32) *Result {
	return &Result{
		Action: action,
		Status: ResultStatus_RESULT_UNKNOWN,
	}
}

func (r *Result) Succeeded() bool {
	return r.GetStatus() == ResultStatus_RESULT_SUCCESS
}

func (r *Result) Failed() bool {
	return !r.Succeeded()
}

func (r *Result) Store(ret ResultTime, err error) {
	r.Time = timestamppb.New(ret.Start.Round(time.Microsecond))
	r.Latency = durationpb.New(ret.Took)
	r.SetStatus(err)
}

//nolint:revive // Context should not be the first argument as this is a decorator.
func (r *Result) Timeit(fun resultFun, ctx context.Context) error {
	ret, err := Timeit(fun, ctx)
	r.Store(ret, err)
	return err
}

func (r *Result) Took() time.Duration {
	return r.GetLatency().AsDuration()
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
func (r *Result) SetStatus(err error) {
	if err == nil {
		r.Status = ResultStatus_RESULT_SUCCESS
		return
	}

	r.Error = err.Error()

	if errors.Is(err, context.DeadlineExceeded) {
		r.Status = ResultStatus_RESULT_TIMEOUT
	} else {
		r.Status = ResultStatus_RESULT_FAILURE
	}
}
