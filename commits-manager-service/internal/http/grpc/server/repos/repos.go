package repos

import (
	"commits-manager-service/internal/constants"
	"commits-manager-service/internal/constants/models"
	"commits-manager-service/internal/http/grpc/protos/repos"
	"commits-manager-service/internal/storage/db"
	"context"
)

type ReposMetaDataServer struct {
	repos.UnimplementedRepositoriesServiceServer
	RepositoryPersistence db.GitReposRepository
}

func (rmds *ReposMetaDataServer) GetRepositories(ctx context.Context, req *repos.GetRepositoriesRequest) (*repos.GetRepositoriesResponse, error) {
	repositories, err := rmds.RepositoryPersistence.GetAllRepositories(1000000, 0)
	if err != nil {
		return nil, err
	}
	return &repos.GetRepositoriesResponse{
		Repositories: Convert(repositories),
	}, nil
}

func (rmds *ReposMetaDataServer) GetRepositoryNames(ctx context.Context, req *repos.GetRepositoryNamesRequest) (*repos.GetRepositoryNamesResponse, error) {
	repositories, err := rmds.RepositoryPersistence.GetAllRepositoryNames()
	if err != nil {
		return nil, err
	}
	return &repos.GetRepositoryNamesResponse{
		Repositories: repositories,
	}, nil
}

func (rmds *ReposMetaDataServer) GetReposFetchHistory(ctx context.Context, req *repos.GetReposFetchHistoryRequest) (*repos.GetReposFetchHistoryResponse, error) {
	reposFetchData, err := rmds.RepositoryPersistence.GetLastReposFetchHistory()
	if err != nil {
		return nil, err
	}
	return &repos.GetReposFetchHistoryResponse{
		LastFetchTime: reposFetchData.FetchedAt.UTC().Format(constants.ISO_8601_TIME_LAYOUT),
		LastPage:      int32(reposFetchData.LastPage),
	}, nil
}

func Convert(repositories []*models.Repository) []*repos.Repository {
	convertedRepos := make([]*repos.Repository, 0)
	for i := range repositories {
		convertedRepos = append(convertedRepos, &repos.Repository{
			Name:            repositories[i].Name,
			Description:     repositories[i].Description,
			Url:             repositories[i].URL,
			Language:        repositories[i].Language,
			ForksCount:      int32(repositories[i].ForksCount),
			StarsCount:      int32(repositories[i].StarsCount),
			OpenIssuesCount: int32(repositories[i].OpenIssuesCount),
			WatchersCount:   int32(repositories[i].WatchersCount),
		})
	}
	return convertedRepos
}
