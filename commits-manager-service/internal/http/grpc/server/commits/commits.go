package commits

import (
	"commits-manager-service/internal/constants"
	"commits-manager-service/internal/http/grpc/protos/commits"
	"commits-manager-service/internal/storage/db"
	"context"
)

type CommitsMetaDataServer struct {
	commits.UnimplementedCommitsServiceServer
	CommitPersistence db.CommitRepository
}

func (cmds *CommitsMetaDataServer) GetCommitFetchHistory(ctx context.Context, req *commits.CommitFetchHistoryRequest) (*commits.CommitFetchHistoryResponse, error) {
	commitsFetchHistory, err := cmds.CommitPersistence.GetLastCommitFetchTime(req.RepositoryName)
	if err != nil {
		return nil, err
	}
	return &commits.CommitFetchHistoryResponse{
		LastFetchTime: commitsFetchHistory.FetchedAt.UTC().Format(constants.ISO_8601_TIME_LAYOUT),
		LastPage:      int32(commitsFetchHistory.LastPage),
	}, nil
}
