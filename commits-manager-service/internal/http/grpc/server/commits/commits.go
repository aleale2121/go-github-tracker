package commits

import (
	"context"
	"commits-manager-service/internal/constants"
	"commits-manager-service/internal/http/grpc/protos/commits"
	"commits-manager-service/internal/storage/db"
)

type CommitsMetaDataServer struct {
	commits.UnimplementedCommitsMetaDataServiceServer
	MetaDataPersistemce db.MetadataPersistence
}

func (cmds *CommitsMetaDataServer) GetRepoCommitMetaData(ctx context.Context, req *commits.RepoCommitMetaDataRequest) (*commits.RepoCommitMetaDataResponse, error) {
	lastFetchTime, err := cmds.MetaDataPersistemce.GetLastCommitFetchTime(req.RepositoryName)
	if err != nil {
		return nil, err
	}
	return &commits.RepoCommitMetaDataResponse{
		LastFetchTime: lastFetchTime.UTC().Format(constants.ISO_8601_TIME_LAYOUT),
	}, nil
}

