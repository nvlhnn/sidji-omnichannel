package services

import (
	"errors"
)

var (
	ErrSubscriptionLimitReached = errors.New("subscription limit reached for this plan")
)

type SubscriptionLimits struct {
	MaxUsers    int // -1 for unlimited
	MaxChannels int // -1 for unlimited
	MaxAIReply  int // -1 for unlimited
}

func GetSubscriptionLimits(plan string) SubscriptionLimits {
	switch plan {
	case "growth":
		return SubscriptionLimits{
			MaxUsers:    3,
			MaxChannels: -1, // Unlimited
			MaxAIReply:  1000,
		}
	case "scale":
		return SubscriptionLimits{
			MaxUsers:    10,
			MaxChannels: -1, // Unlimited
			MaxAIReply:  -1, // Unlimited
		}
	case "starter":
		fallthrough
	default:
		return SubscriptionLimits{
			MaxUsers:    1,
			MaxChannels: 1,
			MaxAIReply:  10,
		}
	}
}

// Helper to check if a count exceeds the limit
func CheckLimit(currentCount, limit int) error {
	if limit == -1 {
		return nil
	}
	if currentCount >= limit {
		return ErrSubscriptionLimitReached
	}
	return nil
}

// IsCompliance checks if the organization is within its plan limits
func IsCompliance(plan string, userCount, channelCount int) bool {
	limits := GetSubscriptionLimits(plan)
	
	if limits.MaxUsers != -1 && userCount > limits.MaxUsers {
		return false
	}
	
	if limits.MaxChannels != -1 && channelCount > limits.MaxChannels {
		return false
	}
	
	return true
}
