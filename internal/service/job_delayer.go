package service

type DeleteBannerByFtTgMsg struct {
	featureID, tagID *int64
}

func NewDeleteBannerByFtTgMsg(featureID, tagID *int64) DeleteBannerByFtTgMsg {
	return DeleteBannerByFtTgMsg{
		featureID: featureID,
		tagID:     tagID,
	}
}

type JobDelayer interface {
	DeleteBannerFtTg(message DeleteBannerByFtTgMsg) error
}

type RabbitMQJobDelayer struct {
}

func NewRabbitMQJobDelayer() *RabbitMQJobDelayer {
	return &RabbitMQJobDelayer{}
}

func (r *RabbitMQJobDelayer) DeleteBannerFtTg(message DeleteBannerByFtTgMsg) error {
	panic(message)
}
