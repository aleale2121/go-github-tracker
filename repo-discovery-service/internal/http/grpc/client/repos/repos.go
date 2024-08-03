package repos

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	rs "repos-discovery-service/internal/http/grpc/protos/repos"
)

type RepositoriesServiceClient struct {
	ServiceUrl string
}

func NewRepositoriesServiceClient(serviceUrl string) *RepositoriesServiceClient {
	return &RepositoriesServiceClient{
		ServiceUrl: serviceUrl,
	}
}

func (rmdsc RepositoriesServiceClient) GetRepositoryNames() ([]string, error) {
	conn, err := grpc.NewClient(rmdsc.ServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	c := rs.NewRepositoriesServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := c.GetRepositoryNames(ctx, &rs.GetRepositoryNamesRequest{})
	if err != nil {
		return []string{}, err
	}
	return response.Repositories, nil
}

func (rmdsc RepositoriesServiceClient) GetReposFetchHistory() (*rs.GetReposFetchHistoryResponse, error) {
	conn, err := grpc.NewClient(rmdsc.ServiceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &rs.GetReposFetchHistoryResponse{}, err
	}
	defer conn.Close()

	c := rs.NewRepositoriesServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response, err := c.GetReposFetchHistory(ctx, &rs.GetReposFetchHistoryRequest{})
	if err != nil {
		return &rs.GetReposFetchHistoryResponse{}, err
	}
	return &rs.GetReposFetchHistoryResponse{
		LastFetchTime: response.LastFetchTime,
		LastPage:      response.LastPage,
	}, nil
}
