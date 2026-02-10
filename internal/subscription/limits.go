package subscription

import "errors"

var (
	ErrSubscriptionLimitReached = errors.New("subscription limit reached for this plan")
)

type SubscriptionLimits struct {
	MaxUsers    int // -1 for unlimited
	MaxChannels int // -1 for unlimited
	MaxAIReply  int // -1 for unlimited
	MaxMessages int // -1 for unlimited
}

func GetSubscriptionLimits(plan string) SubscriptionLimits {
	switch plan {
	case "growth":
		return SubscriptionLimits{
			MaxUsers:    3,
			MaxChannels: -1, // Unlimited
			MaxAIReply:  1000,
			MaxMessages: 5000,
		}
	case "scale":
		return SubscriptionLimits{
			MaxUsers:    10,
			MaxChannels: -1, // Unlimited
			MaxAIReply:  -1, // Unlimited
			MaxMessages: -1, // Unlimited
		}
	case "enterprise":
		return SubscriptionLimits{
			MaxUsers:    -1,
			MaxChannels: -1,
			MaxAIReply:  -1,
			MaxMessages: -1,
		}
	case "starter":
		fallthrough
	default:
		return SubscriptionLimits{
			MaxUsers:    1,
			MaxChannels: 1,
			MaxAIReply:  10,
			MaxMessages: 1000,
		}
	}
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

// CheckLimit checks if a specific limit is reached
func CheckLimit(plan string, currentCount int, limitType string) error {
	limits := GetSubscriptionLimits(plan)
	var max int
	
	switch limitType {
	case "user":
		max = limits.MaxUsers
	case "channel":
		max = limits.MaxChannels
	case "ai_reply":
		max = limits.MaxAIReply
	default:
		return nil
	}
	
	if max != -1 && currentCount >= max {
		return ErrSubscriptionLimitReached
	}
	
	return nil
}
