package commits

import (
	"context"
	"commits-manager-service/internal/constants"
	"commits-manager-service/internal/http/grpc/protos/commits"
	"commits-manager-service/internal/storage/db"
)

type CommitsMetaDataServer struct {
	commits.UnimplementedCommitsServiceServer
	CommitPersistence     db.CommitRepository
}

func (cmds *CommitsMetaDataServer) GetRepoCommitFetchData(ctx context.Context, req *commits.RepoCommiFetchDataRequest) (*commits.RepoCommitFetchDataResponse, error) {
	lastFetchTime, err := cmds.CommitPersistence.GetLastCommitFetchTime(req.RepositoryName)
	if err != nil {
		return nil, err
	}
	return &commits.RepoCommitFetchDataResponse{
		LastFetchTime: lastFetchTime.UTC().Format(constants.ISO_8601_TIME_LAYOUT),
	}, nil
}

