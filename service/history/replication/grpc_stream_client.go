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

	"google.golang.org/grpc/metadata"

	"go.temporal.io/server/api/adminservice/v1"
	"go.temporal.io/server/client"
	"go.temporal.io/server/client/history"
)

type (
	StreamBiDirectionStreamClientProvider struct {
		clientBean client.Bean
	}
)

func NewStreamBiDirectionStreamClientProvider(
	clientBean client.Bean,
) *StreamBiDirectionStreamClientProvider {
	return &StreamBiDirectionStreamClientProvider{
		clientBean: clientBean,
	}
}

func (p *StreamBiDirectionStreamClientProvider) Get(
	ctx context.Context,
	sourceShardKey ClusterShardKey,
	targetShardKey ClusterShardKey,
) (BiDirectionStreamClient[*adminservice.StreamWorkflowReplicationMessagesRequest, *adminservice.StreamWorkflowReplicationMessagesResponse], error) {
	adminClient, err := p.clientBean.GetRemoteAdminClient(targetShardKey.ClusterName)
	if err != nil {
		return nil, err
	}
	ctx = metadata.NewOutgoingContext(ctx, history.EncodeClusterShardMD(
		history.ClusterShardID{
			ClusterName: sourceShardKey.ClusterName,
			ShardID:     sourceShardKey.ShardID,
		},
		history.ClusterShardID{
			ClusterName: targetShardKey.ClusterName,
			ShardID:     targetShardKey.ShardID,
		},
	))
	return adminClient.StreamWorkflowReplicationMessages(ctx)
}
