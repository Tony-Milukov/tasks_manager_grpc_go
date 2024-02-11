package status

import (
	"sso_3.0/internal/domain/models"
	api "sso_3.0/proto/gen"
)

func GetStatus(status *models.Status) *api.Status {
	if status != nil {
		return &api.Status{
			Title:       status.Title,
			Description: status.Description,
			Id:          int64(status.Id),
		}
	}
	return nil
}
func GetStatuses(statuses []*models.Status) []*api.Status {
	var protoStatuses []*api.Status

	if statuses != nil && len(statuses) > 0 {
		for _, status := range statuses {
			protoStatuses = append(protoStatuses, GetStatus(status))
		}
	}
	return protoStatuses
}
