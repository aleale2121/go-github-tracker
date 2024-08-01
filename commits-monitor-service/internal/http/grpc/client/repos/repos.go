package repos

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	rmds "commits-monitor-service/internal/http/grpc/protos/repos"
)

type ReposMetaDataServiceClient struct {
	ServiceUrl string
}

func NewReposMetaDataServiceClient(serviceUrl string) *ReposMetaDataServiceClient {
	return &ReposMetaDataServiceClient{
		ServiceUrl: serviceUrl,
	}
}

func (rmdsc ReposMetaDataServiceClient) GetRepositories() ([]*rmds.Repository, error) {
	conn, err := grpc.NewClient(rmdsc.ServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := rmds.NewRepositoriesServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := c.GetRepositories(ctx, &rmds.GetRepositoriesRequest{})
	if err != nil {
		return nil, err
	}
	return response.Repositories, nil
}

func (rmdsc ReposMetaDataServiceClient) GetRepositoryNames() ([]string, error) {
	conn, err := grpc.NewClient(rmdsc.ServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := rmds.NewRepositoriesServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := c.GetRepositoryNames(ctx, &rmds.GetRepositoryNamesRequest{})
	if err != nil {
		return nil, err
	}
	return response.Repositories, nil
}
