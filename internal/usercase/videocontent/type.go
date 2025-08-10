package videocontent

import (
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/common"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/content"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/delivery"
	"github.com/kkiling/torrent2emby/internal/usercase/videocontent/runners/tvshowdeliverystate"
)

type TVShowID = common.TVShowID
type ContentID = common.ContentID
type VideoContent = content.VideoContent

type TorrentState = delivery.TorrentState

const (
	TorrentStateError       = delivery.TorrentStateError
	TorrentStateUploading   = delivery.TorrentStateUploading
	TorrentStateDownloading = delivery.TorrentStateDownloading
	TorrentStateStopped     = delivery.TorrentStateStopped
	TorrentStateQueued      = delivery.TorrentStateQueued
	TorrentStateUnknown     = delivery.TorrentStateUnknown
)

type DeliveryStatus = content.DeliveryStatus

const (
	DeliveryStatusFailed     = content.DeliveryStatusFailed
	DeliveryStatusInProgress = content.DeliveryStatusInProgress
	DeliveryStatusDelivered  = content.DeliveryStatusDelivered
	DeliveryStatusUpdating   = content.DeliveryStatusUpdating
	DeliveryStatusDeleting   = content.DeliveryStatusDeleting
	DeliveryStatusDeleted    = content.DeliveryStatusDeleted
)

type CreateVideoContentParams = content.CreateVideoContentParams
type TVShowDeliveryState = tvshowdeliverystate.State
type TVShowDeliveryData = tvshowdeliverystate.TVShowDeliveryData
type TorrentSearch = delivery.TorrentSearch
type ContentMatches = delivery.ContentMatches
type Track = delivery.Track
type FileInfo = delivery.FileInfo
type ChoseTorrentOptions = tvshowdeliverystate.ChoseTorrentOptions
type ChoseFileMatchesOptions = tvshowdeliverystate.ChoseFileMatchesOptions

type StepDelivery = tvshowdeliverystate.StepDelivery

const (
	GenerateSearchQuery            = tvshowdeliverystate.GenerateSearchQuery
	SearchTorrents                 = tvshowdeliverystate.SearchTorrents
	WaitingUserChoseTorrent        = tvshowdeliverystate.WaitingUserChoseTorrent
	GetMagnetLink                  = tvshowdeliverystate.GetMagnetLink
	AddTorrentToTorrentClient      = tvshowdeliverystate.AddTorrentToTorrentClient
	PrepareFileMatches             = tvshowdeliverystate.PrepareFileMatches
	WaitingChoseFileMatches        = tvshowdeliverystate.WaitingChoseFileMatches
	WaitingTorrentDownloadComplete = tvshowdeliverystate.WaitingTorrentDownloadComplete
	CreateVideoContentCatalogs     = tvshowdeliverystate.CreateVideoContentCatalogs
	DeterminingNeedConvertFiles    = tvshowdeliverystate.DeterminingNeedConvertFiles
	StartMergeVideoFiles           = tvshowdeliverystate.StartMergeVideoFiles
	WaitingMergeVideoFiles         = tvshowdeliverystate.WaitingMergeVideoFiles
	CopyVideoFiles                 = tvshowdeliverystate.CopyVideoFiles
	SetVideoFileGroup              = tvshowdeliverystate.SetVideoFileGroup
	GetCatalogsSize                = tvshowdeliverystate.GetCatalogsSize
	SetMediaMetaData               = tvshowdeliverystate.SetMediaMetaData
	SendDeliveryNotification       = tvshowdeliverystate.SendDeliveryNotification
)
