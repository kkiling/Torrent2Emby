package videocontent

import (
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent"
	desc "github.com/kkiling/torrent2emby/pkg/gen/torrent2emby"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapContentIDReq(id *desc.ContentID) videocontent.ContentID {
	if id == nil || (id.MovieId == nil && id.TvShow == nil) {
		return videocontent.ContentID{}
	}

	var result videocontent.ContentID

	if id.MovieId != nil {
		result.MovieID = id.MovieId
	}

	if id.TvShow != nil {
		result.TVShow = &videocontent.TVShowID{
			ID:           id.TvShow.Id,
			SeasonNumber: uint8(id.TvShow.SeasonNumber),
		}
	}

	return result
}
func mapContentID(id *videocontent.ContentID) *desc.ContentID {
	var result desc.ContentID
	if id.MovieID != nil {
		result.MovieId = id.MovieID
	}

	if id.TVShow != nil {
		result.TvShow = &desc.TVShowID{
			Id:           id.TVShow.ID,
			SeasonNumber: uint32(id.TVShow.SeasonNumber),
		}
	}

	return &result
}

func mapVideoContent(result videocontent.VideoContent) *desc.VideoContent {
	return &desc.VideoContent{
		Id:             result.ID.String(),
		CreatedAt:      timestamppb.New(result.CreatedAt),
		ContentId:      mapContentID(&result.ContentID),
		DeliveryStatus: mapDeliveryStatus(result.DeliveryStatus),
	}
}

func mapDeliveryStatus(deliveryStatus videocontent.DeliveryStatus) desc.DeliveryStatus {
	switch deliveryStatus {
	case videocontent.DeliveryStatusFailed:
		return desc.DeliveryStatus_DeliveryStatusFailed
	case videocontent.DeliveryStatusInProgress:
		return desc.DeliveryStatus_DeliveryStatusInProgress
	case videocontent.DeliveryStatusDelivered:
		return desc.DeliveryStatus_DeliveryStatusDelivered
	case videocontent.DeliveryStatusUpdating:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	case videocontent.DeliveryStatusDeleting:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	case videocontent.DeliveryStatusDeleted:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	default:
		return desc.DeliveryStatus_DeliveryStatusUnknown
	}
}

func mapDeliveryStep(step videocontent.StepDelivery) desc.TVShowDeliveryStatus {
	switch step {
	case videocontent.GenerateSearchQuery:
		return desc.TVShowDeliveryStatus_GenerateSearchQuery
	case videocontent.SearchTorrents:
		return desc.TVShowDeliveryStatus_SearchTorrents
	case videocontent.WaitingUserChoseTorrent:
		return desc.TVShowDeliveryStatus_WaitingUserChoseTorrent
	case videocontent.GetMagnetLink:
		return desc.TVShowDeliveryStatus_GetMagnetLink
	case videocontent.AddTorrentToTorrentClient:
		return desc.TVShowDeliveryStatus_AddTorrentToTorrentClient
	case videocontent.PrepareFileMatches:
		return desc.TVShowDeliveryStatus_PrepareFileMatches
	case videocontent.WaitingChoseFileMatches:
		return desc.TVShowDeliveryStatus_WaitingChoseFileMatches
	case videocontent.WaitingTorrentDownloadComplete:
		return desc.TVShowDeliveryStatus_WaitingTorrentDownloadComplete
	case videocontent.CreateVideoContentCatalogs:
		return desc.TVShowDeliveryStatus_CreateVideoContentCatalogs
	case videocontent.DeterminingNeedConvertFiles:
		return desc.TVShowDeliveryStatus_DeterminingNeedConvertFiles
	case videocontent.StartMergeVideoFiles:
		return desc.TVShowDeliveryStatus_StartMergeVideoFiles
	case videocontent.WaitingMergeVideoFiles:
		return desc.TVShowDeliveryStatus_WaitingMergeVideoFiles
	case videocontent.CopyVideoFiles:
		return desc.TVShowDeliveryStatus_CopyVideoFiles
	case videocontent.SetVideoFileGroup:
		return desc.TVShowDeliveryStatus_SetVideoFileGroup
	case videocontent.SetMediaMetaData:
		return desc.TVShowDeliveryStatus_SetMediaMetaData
	case videocontent.SendDeliveryNotification:
		return desc.TVShowDeliveryStatus_SendDeliveryNotification
	default:
		return desc.TVShowDeliveryStatus_TVShowDeliveryStatusUnknown
	}
}

func mapTVShowDeliveryState(data *videocontent.TVShowDeliveryData) *desc.TVShowDeliveryData {
	return &desc.TVShowDeliveryData{
		SearchQuery: data.SearchQuery,
		TorrentSearch: lo.Map(data.TorrentSearch.Result, func(item videocontent.TorrentSearch, _ int) *desc.TorrentSearch {
			return &desc.TorrentSearch{
				Title:     item.Title,
				Href:      item.Href,
				Size:      item.Size,
				Seeds:     item.Seeds,
				Leeches:   item.Leeches,
				Downloads: item.Downloads,
				AddedDate: item.AddedDate,
			}
		}),
	}
}
