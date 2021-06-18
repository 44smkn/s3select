package query_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	mock_services "github.com/44smkn/s3select/mocks/aws/services"
	"github.com/44smkn/s3select/pkg/aws/services"
	"github.com/44smkn/s3select/pkg/config"
	"github.com/44smkn/s3select/pkg/query"
	s3sdk "github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func Test_defaultObjectSelector_Select(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_services.NewMockS3(ctrl)
	es := s3sdk.NewSelectObjectContentEventStream(func(o *s3sdk.SelectObjectContentEventStream) {
		o.Reader = newFakeEventsStreamReader()
		o.StreamCloser = &fakeStreamCloser{}
	})
	output := &s3sdk.SelectObjectContentOutput{
		EventStream: es,
	}
	m.EXPECT().SelectObjectContentWithContext(gomock.Any(), gomock.Any()).Return(output, nil)

	profile := &config.Profile{}
	cloud := &fakeAwsCloud{
		s3: m,
	}
	selector := query.NewDefaultObjectSelector(profile, cloud, zap.L())
	ctx := context.Background()
	meta := &query.ObjectMetadata{
		BucketName: "fakeBucket",
		ObjectKey:  "fakeKey",
	}
	buf := &bytes.Buffer{}
	go func() {
		time.Sleep(10 * time.Second)
		output.GetEventStream().Close()
	}()
	selector.Select(ctx, meta, "Fake Expression", buf)

	if buf.String() != "summary" {
		t.Errorf("failed to ---")
	}
}

type fakeAwsCloud struct {
	s3 services.S3
}

func (c *fakeAwsCloud) S3() services.S3 {
	return c.s3
}

type fakeEventStreamReader struct {
	stream chan s3sdk.SelectObjectContentEventStreamEvent
}

func newFakeEventsStreamReader() *fakeEventStreamReader {
	ch := make(chan s3sdk.SelectObjectContentEventStreamEvent, 10)
	ch <- &s3sdk.RecordsEvent{Payload: []byte("summary")}
	return &fakeEventStreamReader{
		stream: ch,
	}
}

func (esr *fakeEventStreamReader) Events() <-chan s3sdk.SelectObjectContentEventStreamEvent {
	return esr.stream
}

func (esr *fakeEventStreamReader) Close() error {
	close(esr.stream)
	return nil
}

func (esr *fakeEventStreamReader) Err() error {
	return nil
}

type fakeStreamCloser struct{}

func (sc *fakeStreamCloser) Close() error {
	return nil
}
