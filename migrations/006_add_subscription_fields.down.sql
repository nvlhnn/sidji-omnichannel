ALTER TABLE organizations
DROP CONSTRAINT IF EXISTS check_plan,
DROP CONSTRAINT IF EXISTS check_subscription_status,
DROP COLUMN IF EXISTS plan,
DROP COLUMN IF EXISTS subscription_status,
DROP COLUMN IF EXISTS ai_credits_limit,
DROP COLUMN IF EXISTS ai_credits_used,
DROP COLUMN IF EXISTS billing_cycle_start;

DROP INDEX IF EXISTS idx_organizations_subscription;
