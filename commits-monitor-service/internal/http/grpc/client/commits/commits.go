package commits

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cmds "commits-monitor-service/internal/http/grpc/protos/commits"
)

type CommitsMetaDataServiceClient struct {
	ServiceUrl string
}

func NewCommitsMetaDataServiceClient(serviceUrl string) *CommitsMetaDataServiceClient {
	return &CommitsMetaDataServiceClient{
		ServiceUrl: serviceUrl,
	}
}

func (rmdsc CommitsMetaDataServiceClient) GetCommitFetchHistory(repoName string) (*cmds.CommitFetchHistoryResponse, error) {
	conn, err := grpc.NewClient(rmdsc.ServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &cmds.CommitFetchHistoryResponse{}, err
	}
	defer conn.Close()

	c := cmds.NewCommitsServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := c.GetCommitFetchHistory(ctx, &cmds.CommitFetchHistoryRequest{
		RepositoryName: repoName,
	})
	if err != nil {
		return &cmds.CommitFetchHistoryResponse{}, err
	}
	return response, nil
}
