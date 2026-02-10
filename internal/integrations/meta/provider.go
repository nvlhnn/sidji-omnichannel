package meta

import (
	"context"
	"fmt"

	"github.com/sidji-omnichannel/internal/models"
)

type MetaProvider struct {
	client *MetaClient
}

func NewMetaProvider(client *MetaClient) *MetaProvider {
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
	case models.ChannelFacebook:
		return p.sendFacebook(channel, contact, message)
	default:
		return "", fmt.Errorf("unsupported channel type for meta provider: %s", channel.Type)
	}
}

func (p *MetaProvider) sendWhatsApp(channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	var resp *WhatsAppSendResponse
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
	var resp *InstagramSendResponse
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

func (p *MetaProvider) sendFacebook(channel *models.Channel, contact *models.Contact, message *models.Message) (string, error) {
	var resp *FacebookSendResponse
	var err error

	// Use Page Access Token if available in config, otherwise use channel.AccessToken
	accessToken := channel.AccessToken
	// TODO: Pull page token from config if it's stored there specifically

	switch message.MessageType {
	case models.MessageTypeText:
		resp, err = p.client.SendFacebookText(
			contact.FacebookID,
			message.Content,
			accessToken,
		)
	case models.MessageTypeImage:
		resp, err = p.client.SendFacebookImage(
			contact.FacebookID,
			message.MediaURL,
			accessToken,
		)
	default:
		resp, err = p.client.SendFacebookText(
			contact.FacebookID,
			message.Content,
			accessToken,
		)
	}

	if err != nil {
		return "", err
	}

	return resp.MessageID, nil
}
