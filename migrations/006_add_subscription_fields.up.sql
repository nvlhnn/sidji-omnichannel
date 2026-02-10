-- Add subscription fields to organizations table
ALTER TABLE organizations
ADD COLUMN plan VARCHAR(50) NOT NULL DEFAULT 'starter',
ADD COLUMN subscription_status VARCHAR(20) NOT NULL DEFAULT 'active',
ADD COLUMN ai_credits_limit INTEGER NOT NULL DEFAULT 10,
ADD COLUMN ai_credits_used INTEGER NOT NULL DEFAULT 0,
ADD COLUMN billing_cycle_start TIMESTAMP WITH TIME ZONE DEFAULT NOW();

-- Add check constraints
ALTER TABLE organizations ADD CONSTRAINT check_plan CHECK (plan IN ('starter', 'growth', 'scale', 'enterprise'));
ALTER TABLE organizations ADD CONSTRAINT check_subscription_status CHECK (subscription_status IN ('active', 'past_due', 'canceled', 'trial'));

-- Create index for billing queries
CREATE INDEX idx_organizations_subscription ON organizations(plan, subscription_status);
