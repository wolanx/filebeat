package loki

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/prometheus/common/model"
	"google.golang.org/grpc"
)

type lokiClient struct {
	host     string
	isHttp   bool
	client   *http.Client
	grpc     logproto.PusherClient
	observer outputs.Observer
	index    string
}

func (c *lokiClient) String() string {
	return "loki"
}

func (c *lokiClient) Connect() error {
	c.client = &http.Client{
		Timeout: 2 * time.Second,
	}

	conn, err := grpc.Dial(c.host, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.grpc = logproto.NewPusherClient(conn)

	return nil
}

func (c *lokiClient) Close() error {
	c.client = nil
	c.grpc = nil
	return nil
}

func (c *lokiClient) Publish(ctx context.Context, batch publisher.Batch) error {
	st := c.observer
	events := batch.Events()
	st.NewBatch(len(events))

	var rest []publisher.Event
	var err error

	if c.isHttp {
		rest, err = c.PublicBatchHttp(ctx, events) // http
	} else {
		rest, err = c.PublicBatchGrpc(ctx, events) // grpc
	}

	dropped := len(rest)
	if dropped == 0 {
		batch.ACK()
	} else {
		batch.RetryEvents(rest)
	}

	st.Dropped(dropped)
	st.Acked(len(events) - dropped)

	return err
}

func fields2labels(fields common.MapStr) model.LabelSet {
	fields.Delete("agent")
	fields.Delete("message")
	//ret := common.MapStr{}
	ret := model.LabelSet{}

	for k, v := range fields.Flatten() {
		label := strings.ReplaceAll(strings.ReplaceAll(k, "-", "_"), ".", "_")
		ret[model.LabelName(label)] = model.LabelValue(fmt.Sprintf("%s", v))
	}
	return ret
}

func (c *lokiClient) PublicBatchGrpc(ctx context.Context, events []publisher.Event) ([]publisher.Event, error) {
	var rest []publisher.Event
	for _, event := range events {
		content := &event.Content
		fields := content.Fields
		msg, _ := fields.GetValue("message")

		logs := []logproto.Entry{
			{Timestamp: content.Timestamp, Line: msg.(string)},
		}
		req := &logproto.PushRequest{Streams: []logproto.Stream{
			{Labels: fields2labels(fields).String(), Entries: logs},
		}}
		_, err := c.grpc.Push(ctx, req)
		if err != nil {
			rest = append(rest, event)
		}
	}
	if len(rest) > 0 {
		return rest, errors.New("loki send failure")
	}
	return nil, nil
}

func (c *lokiClient) PublicBatchHttp(_ context.Context, events []publisher.Event) ([]publisher.Event, error) {
	var rest []publisher.Event
	for _, event := range events {
		content := &event.Content
		ns := content.Timestamp.UnixNano()
		fields := content.Fields

		err := c.publishOneHttp(ns, fields)
		if err != nil {
			rest = append(rest, event)
		}
	}
	if len(rest) > 0 {
		return rest, errors.New("loki send failure")
	}
	return nil, nil
}

func (c *lokiClient) publishOneHttp(ns int64, fields common.MapStr) error {
	msg, _ := fields.GetValue("message")

	var values [][]interface{}
	values = append(values, []interface{}{strconv.FormatInt(ns, 10), msg})

	body := common.MapStr{}
	group := common.MapStr{}
	group["stream"] = fields2labels(fields)
	group["values"] = values
	body["streams"] = []common.MapStr{group}

	req, _ := http.NewRequest("POST", "http://"+c.host+"/loki/api/v1/push", strings.NewReader(body.String()))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}
