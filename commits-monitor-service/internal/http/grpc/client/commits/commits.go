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

func (rmdsc CommitsMetaDataServiceClient) GetRepoLastFetchTime(repoName string) (string, error) {
	conn, err := grpc.NewClient(rmdsc.ServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	c := cmds.NewCommitsMetaDataServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := c.GetRepoCommitMetaData(ctx, &cmds.RepoCommitMetaDataRequest{
		RepositoryName: repoName,
	})
	if err != nil {
		return "", err
	}
	return response.LastFetchTime, nil
}
