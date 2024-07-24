package repos

import (
	"commits-manager-service/internal/constants"
	"commits-manager-service/internal/constants/models"
	"commits-manager-service/internal/http/grpc/protos/repos"
	"commits-manager-service/internal/storage/db"
	"context"
)

type ReposMetaDataServer struct {
	repos.UnimplementedRepositoryMetaDataServiceServer
	MetaDataPersistemce   db.MetadataPersistence
	RepositoryPersistence db.RepositoryPersistence
}

func (rmds *ReposMetaDataServer) GetRepositories(ctx context.Context, req *repos.GetRepositoriesRequest) (*repos.GetRepositoriesResponse, error) {
	repositories, err := rmds.RepositoryPersistence.GetAllRepositories()
	if err != nil {
		return nil, err
	}
	return &repos.GetRepositoriesResponse{
		Repositories: Convert(repositories),
	}, nil
}

func (rmds *ReposMetaDataServer) AllRepositoriesMetaData(ctx context.Context, req *repos.AllReposMetaDataRequest) (*repos.AllReposMetaDataResponse, error) {
	lastFetchTime, err := rmds.MetaDataPersistemce.GetLastReposFetchTime()
	if err != nil {
		return nil, err
	}
	return &repos.AllReposMetaDataResponse{
		LastFetchTime: lastFetchTime.UTC().Format(constants.ISO_8601_TIME_LAYOUT),
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
