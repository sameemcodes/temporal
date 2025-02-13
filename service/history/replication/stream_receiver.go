// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package replication

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.temporal.io/server/api/adminservice/v1"
	repicationpb "go.temporal.io/server/api/replication/v1"
	"go.temporal.io/server/common"
	"go.temporal.io/server/common/channel"
	"go.temporal.io/server/common/log/tag"
	"go.temporal.io/server/common/primitives/timestamp"
	ctasks "go.temporal.io/server/common/tasks"
)

type (
	ClusterShardKey struct {
		ClusterName string
		ShardID     int32
	}
	ClusterShardKeyPair struct {
		Source ClusterShardKey
		Target ClusterShardKey
	}

	Stream         BiDirectionStream[*adminservice.StreamWorkflowReplicationMessagesRequest, *adminservice.StreamWorkflowReplicationMessagesResponse]
	StreamReceiver struct {
		ProcessToolBox

		status         int32
		sourceShardKey ClusterShardKey
		targetShardKey ClusterShardKey
		taskTracker    ExecutableTaskTracker
		shutdownChan   channel.ShutdownOnce

		sync.Mutex
		streamCreationTime time.Time
		stream             Stream
	}
)

func NewClusterShardKey(
	ClusterName string,
	ClusterShardID int32,
) ClusterShardKey {
	return ClusterShardKey{
		ClusterName: ClusterName,
		ShardID:     ClusterShardID,
	}
}

func NewStreamReceiver(
	processToolBox ProcessToolBox,
	sourceShardKey ClusterShardKey,
	targetShardKey ClusterShardKey,
) *StreamReceiver {
	taskTracker := NewExecutableTaskTracker(processToolBox.Logger)
	return &StreamReceiver{
		ProcessToolBox: processToolBox,

		status:         common.DaemonStatusInitialized,
		sourceShardKey: sourceShardKey,
		targetShardKey: targetShardKey,
		taskTracker:    taskTracker,
		shutdownChan:   channel.NewShutdownOnce(),

		streamCreationTime: time.Now().UTC(),
		stream: newStream(
			processToolBox,
			sourceShardKey,
			targetShardKey,
		),
	}
}

// Start starts the processor
func (r *StreamReceiver) Start() {
	if !atomic.CompareAndSwapInt32(
		&r.status,
		common.DaemonStatusInitialized,
		common.DaemonStatusStarted,
	) {
		return
	}

	go r.sendEventLoop()
	go r.recvEventLoop()

	r.Logger.Info("StreamReceiver started.")
}

// Stop stops the processor
func (r *StreamReceiver) Stop() {
	if !atomic.CompareAndSwapInt32(
		&r.status,
		common.DaemonStatusStarted,
		common.DaemonStatusStopped,
	) {
		return
	}

	r.shutdownChan.Shutdown()
	r.stream.Close()
	r.taskTracker.Cancel()

	r.Logger.Info("StreamReceiver shutting down.")
}

func (r *StreamReceiver) IsValid() bool {
	return atomic.LoadInt32(&r.status) == common.DaemonStatusStarted
}

func (r *StreamReceiver) sendEventLoop() {
	defer r.Stop()
	timer := time.NewTicker(r.Config.ReplicationStreamSyncStatusDuration())
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			timer.Reset(r.Config.ReplicationStreamSyncStatusDuration())
			r.Lock()
			stream := r.stream
			r.Unlock()
			r.ackMessage(stream)
		case <-r.shutdownChan.Channel():
			return
		}
	}
}

func (r *StreamReceiver) recvEventLoop() {
	defer r.Stop()

	for !r.shutdownChan.IsShutdown() {
		r.Lock()
		streamCreationTime := r.streamCreationTime
		stream := r.stream
		r.Unlock()

		_ = r.processMessages(stream)
		delay := streamCreationTime.Add(r.Config.ReplicationStreamMinReconnectDuration()).Sub(time.Now().UTC())
		if delay > 0 {
			select {
			case <-time.After(delay):
			case <-r.shutdownChan.Channel():
			}
		}

		r.Lock()
		r.streamCreationTime = time.Now().UTC()
		r.stream = newStream(
			r.ProcessToolBox,
			r.sourceShardKey,
			r.targetShardKey,
		)
		r.Unlock()
	}
}

func (r *StreamReceiver) ackMessage(
	stream Stream,
) {
	watermarkInfo := r.taskTracker.LowWatermark()
	if watermarkInfo == nil {
		return
	}
	if err := stream.Send(&adminservice.StreamWorkflowReplicationMessagesRequest{
		Attributes: &adminservice.StreamWorkflowReplicationMessagesRequest_SyncReplicationState{
			SyncReplicationState: &repicationpb.SyncReplicationState{
				LastProcessedMessageId:   watermarkInfo.Watermark,
				LastProcessedMessageTime: timestamp.TimePtr(watermarkInfo.Timestamp),
			},
		},
	}); err != nil {
		r.Logger.Error("StreamReceiver unable to send message, err", tag.Error(err))
	}
}

func (r *StreamReceiver) processMessages(
	stream Stream,
) error {
	streamRespChen, err := stream.Recv()
	if err != nil {
		r.Logger.Error("StreamReceiver unable to recv message, err", tag.Error(err))
		return err
	}
	for streamResp := range streamRespChen {
		if streamResp.Err != nil {
			r.Logger.Error("StreamReceiver recv stream encountered unexpected err", tag.Error(streamResp.Err))
			return streamResp.Err
		}
		tasks := r.ConvertTasks(
			r.targetShardKey.ClusterName, // data come from target
			streamResp.Resp.GetMessages().ReplicationTasks...,
		)
		highWatermark := streamResp.Resp.GetMessages().LastTaskId
		highWatermarkTime := timestamp.TimeValue(streamResp.Resp.GetMessages().LastTaskTime)
		r.taskTracker.TrackTasks(WatermarkInfo{
			Watermark: highWatermark,
			Timestamp: highWatermarkTime,
		}, tasks...)
		for _, task := range tasks {
			r.ProcessToolBox.TaskScheduler.Submit(task)
		}
	}
	r.Logger.Error("StreamReceiver encountered channel close")
	return nil
}

func newStream(
	processToolBox ProcessToolBox,
	sourceShardKey ClusterShardKey,
	targetShardKey ClusterShardKey,
) Stream {
	var clientProvider BiDirectionStreamClientProvider[*adminservice.StreamWorkflowReplicationMessagesRequest, *adminservice.StreamWorkflowReplicationMessagesResponse] = &streamClientProvider{
		processToolBox: processToolBox,
		sourceShardKey: sourceShardKey,
		targetShardKey: targetShardKey,
	}
	return NewBiDirectionStream(
		clientProvider,
		processToolBox.MetricsHandler,
		processToolBox.Logger,
	)
}

type streamClientProvider struct {
	processToolBox ProcessToolBox
	sourceShardKey ClusterShardKey
	targetShardKey ClusterShardKey
}

var _ BiDirectionStreamClientProvider[*adminservice.StreamWorkflowReplicationMessagesRequest, *adminservice.StreamWorkflowReplicationMessagesResponse] = (*streamClientProvider)(nil)

func (p *streamClientProvider) Get(
	ctx context.Context,
) (BiDirectionStreamClient[*adminservice.StreamWorkflowReplicationMessagesRequest, *adminservice.StreamWorkflowReplicationMessagesResponse], error) {
	return NewStreamBiDirectionStreamClientProvider(p.processToolBox.ClientBean).Get(ctx, p.sourceShardKey, p.targetShardKey)
}

type noopSchedulerMonitor struct {
}

func newNoopSchedulerMonitor() *noopSchedulerMonitor {
	return &noopSchedulerMonitor{}
}

func (m *noopSchedulerMonitor) Start()                    {}
func (m *noopSchedulerMonitor) Stop()                     {}
func (m *noopSchedulerMonitor) RecordStart(_ ctasks.Task) {}
