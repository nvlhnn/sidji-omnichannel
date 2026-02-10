package integrations

import (
	"context"
	"fmt"

	"github.com/sidji-omnichannel/internal/integrations/meta"
	"github.com/sidji-omnichannel/internal/models"
)

type MetaProvider struct {
	client *meta.MetaClient
}

func NewMetaProvider(client *meta.MetaClient) *MetaProvider {
	return &MetaProvider{client: client}
}

func (p *MetaProvider) GetProviderName() string {
	return "meta"
}

func (p *MetaProvider) SendMessage(ctx context.Context, channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	switch channel.Type {
	case models.ChannelWhatsApp:
		return p.sendWhatsApp(channel, contact, message)
	case models.ChannelInstagram:
		return p.sendInstagram(channel, contact, message)
	default:
		return "", fmt.Errorf("unsupported channel type for meta provider: %s", channel.Type)
	}
}

func (p *MetaProvider) sendWhatsApp(channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	var resp *meta.WhatsAppSendResponse
	var err error

	switch message.MessageType {
	case models.MessageTypeText:
		resp, err = p.client.SendWhatsAppText(
			channel.PhoneNumberID,
			contact.WhatsAppID,
			message.Content,
			channel.AccessToken,
		)
	case models.MessageTypeImage:
		resp, err = p.client.SendWhatsAppImage(
			channel.PhoneNumberID,
			contact.WhatsAppID,
			message.MediaURL,
			message.Content, // caption
			channel.AccessToken,
		)
	case models.MessageTypeDocument:
		resp, err = p.client.SendWhatsAppDocument(
			channel.PhoneNumberID,
			contact.WhatsAppID,
			message.MediaURL,
			message.MediaFileName,
			message.Content, // caption
			channel.AccessToken,
		)
	default:
		resp, err = p.client.SendWhatsAppText(
			channel.PhoneNumberID,
			contact.WhatsAppID,
			message.Content,
			channel.AccessToken,
		)
	}

	if err != nil {
		return "", err
	}

	if len(resp.Messages) > 0 {
		return resp.Messages[0].ID, nil
	}

	return "", fmt.Errorf("no message ID returned from Meta")
}

func (p *MetaProvider) sendInstagram(channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	var resp *meta.InstagramSendResponse
	var err error

	switch message.MessageType {
	case models.MessageTypeText:
		resp, err = p.client.SendInstagramText(
			channel.IGUserID,
			contact.InstagramID,
			message.Content,
			channel.AccessToken,
		)
	case models.MessageTypeImage:
		resp, err = p.client.SendInstagramImage(
			channel.IGUserID,
			contact.InstagramID,
			message.MediaURL,
			channel.AccessToken,
		)
	default:
		resp, err = p.client.SendInstagramText(
			channel.IGUserID,
			contact.InstagramID,
			message.Content,
			channel.AccessToken,
		)
	}

	if err != nil {
		return "", err
	}

	return resp.MessageID, nil
}
