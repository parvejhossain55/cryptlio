# Cryplio Database Schema — Entity Mapping

This document maps SRS requirements to database tables.

## Entity-to-Table Mapping (SRS §7.1)

| SRS Entity | Table | Notes |
|------------|-------|-------|
| User | `users` | Core identity with KYC level, status |
| User Profile | `users` (bio, avatar) + `public_user_profiles` view | Public-facing aggregated data |
| User Statistics | `user_stats` | Denormalized aggregates for performance |
| User Sessions | `user_sessions` | Device tracking, revocation (SRS FR-115) |
| Trade Ad | `trade_ads` | Maker advertisements with escrow lock |
| Trade | `trades` | Execution records with status lifecycle |
| Trade Messages | `trade_messages` | Real-time chat (SRS FR-223) |
| Trade Attachments | `trade_attachments` | File upload support (SRS FR-224) |
| Wallet | `wallets` | Per-crypto wallet with locked balance |
| Transaction | `wallet_transactions` | Full transaction ledger |
| Dispute | `disputes` | Dispute resolution workflow (SRS 3.4) |
| Dispute Messages | `dispute_messages` | 3-way chat (buyer/seller/admin) |
| KYC Record | `kyc_records` | Document storage, verification status |
| Feedback | `trade_feedback` | Post-trade ratings (SRS FR-244) |
| Referral | `referrals` | Track referrals + commission (SRS 3.7) |

## Additional Supporting Tables

| Table | Purpose | SRS Requirement |
|-------|---------|----------------|
| `payment_methods` | Lookup: bkash, nagad, bank, wise, etc. | SRS 5.3 |
| `crypto_assets` | Lookup: USDT, ETH, BTC with network info | SRS 5.2 |
| `fiat_currencies` | Lookup: USD, BDT, EUR, etc. | SRS 5.3 |
| `dispute_reasons` | Standardized dispute categories | SRS FR-401 |
| `fee_tiers` | Volume-based platform fee tiers | SRS 3.7/3.8 |
| `admin_actions` | Admin audit trail | SRS FR-809 |
| `audit_logs` | Generic audit for all tables | SRS 4.2 Security |
| `platform_config` | Key-value platform configuration | SRS FR-804 |
| `announcements` | System-wide announcements | SRS FR-806 |
| `user_blocks` | User blocking feature | SRS FR-136 |
| `rate_limit_counts` | Rate limiting enforcement | SRS 4.2 Security |
| `login_attempts` | Brute-force tracking | SRS FR-116 |
| `api_request_logs` | API analytics | Monitoring requirement |
| `notifications` | In-app notification queue | SRS FR-501/502 |
| `notification_preferences` | User notification settings | SRS FR-505 |
| `merchant_applications` | Merchant verification workflow | SRS FR-601-603 |
| `merchant_analytics` | Daily merchant rollup metrics | SRS FR-604 |
| `password_reset_tokens` | Secure password reset | SRS FR-117 |
| `email_verification_tokens` | Email verification | SRS FR-102 |

## Key Relationships

### User → Trade Ads
```
users (1) ──→ (N) trade_ads
```
- One user can create many ads
- `trade_ads.user_id` references `users.user_id`
- ON DELETE CASCADE: user deleted → ads removed

### Trade Ad → Trades
```
trade_ads (1) ──→ (N) trades
```
- One ad can spawn multiple trade attempts
- `trades.ad_id` references `trade_ads.ad_id`
- ON DELETE RESTRICT: prevent ad deletion if active trades exist

### Trade → Messages
```
trades (1) ──→ (N) trade_messages
```
- Trade chat messages
- `trade_messages.trade_id` ON DELETE CASCADE
- Deleting trade removes all chat history (per data retention policy 2 years)

### Trade → Feedback
```
trades (1) ──→ (1) trade_feedback
```
- One feedback from each party (two rows)
- `trade_feedback.trade_id` UNIQUE constraint ensures one rating per user pair

### Trade → Dispute
```
trades (1) ──→ (0..1) disputes
```
- Optional one dispute per trade
- `disputes.trade_id` UNIQUE constraint

### Dispute → Messages
```
disputes (1) ──→ (N) dispute_messages
```
- 3-way chat (buyer, seller, admin)
- ON DELETE CASCADE

### User → Wallets
```
users (1) ──→ (N) wallets
```
- One wallet per crypto per user
- Compound unique: `(user_id, crypto_id)`

### Wallet → Transactions
```
wallets (1) ──→ (N) wallet_transactions
```
- Complete ledger of all wallet activity

### Trade → Escrow Implications

Trades don't directly reference wallets. Instead, wallet transactions reference trades via `reference_id`.

Escrow locking creates `wallet_transactions` with:
- `type = 'escrow_lock'` — when trade starts
- `type = 'escrow_release'` — when seller releases

## ENUM Types & Status Values

### user_status
```
'pending'   → Email unverified, cannot trade (FR-121)
'active'    → Full trading privileges
'suspended' → Temporarily blocked by admin
'banned'    → Permanently banned
'deleted'   → Soft-deleted for GDPR (purged after 30 days)
```

### kyc_level (SRS FR-121–FR-124)
```
'level_0' → No verification (browse only)
'level_1' → Email verified ($500/day limit)
'level_2' → ID + selfie ($10,000/day limit)
'level_3' → Address proof + enhanced due diligence (unlimited)
```

### trade_status (SRS 3.2.3)
```
'pending'   → Trade created, awaiting escrow lock
'active'    → Escrow locked, buyer must pay
'paid'      → Buyer marked as paid (timer started)
'released'  → Seller released escrow, trade complete
'cancelled' → Trade cancelled before completion
'disputed'  → Dispute raised
'completed' → Final state (alias of 'released')
'expired'   → Payment window expired, auto-cancelled
```

### ad_type
```
'buy'  → Buyer advertisement (user wants to buy crypto)
'sell' → Seller advertisement (user wants to sell crypto)
```

### price_type (SRS FR-202)
```
'fixed'   → Fixed fiat price per crypto unit
'floating' → Percentage above/below market rate
```

### dispute_status (SRS 3.4)
```
'pending'     → Dispute opened, awaiting admin assignment
'assigned'    → Admin assigned, review in progress
'under_review'→ Admin actively reviewing evidence
'resolved'    → Ruling issued, funds distributed
'appealed'    → Party appealed decision
'closed'      → Appeal process complete
```

### dispute_resolution (SRS FR-407)
```
'release_to_buyer' → Escrow released to buyer (payment confirmed)
'return_to_seller' → Escrow returned to seller (no payment received)
'partial_split'    → Partial release/return (compromise)
'cancel'          → Trade cancelled, funds returned to buyer
```

### admin_action_type (SRS FR-801–FR-811)
```
'user_suspend'        → Temporarily suspend user
'user_ban'            → Permanently ban user
'user_unban'          → Lift ban
'dispute_resolve'     → Resolve dispute with ruling
'withdrawal_approve'  → Approve large withdrawal
'withdrawal_reject'   → Reject withdrawal request
'kyc_review'          → Review KYC submission
'announcement_post'   → Post platform announcement
'fee_update'          → Update platform fee config
'merchant_approve'    → Approve merchant application
'merchant_reject'     → Reject merchant application
'config_change'       → Platform configuration change
'bulk_message'        → Send bulk notification
'report_generate'     → Generate admin report
```

## Data Retention & GDPR Compliance

### Soft Deletes
- `users.deleted_at` — user account deletion (GDPR right to erasure)
- `trade_ads.deleted_at` — ad deletion (soft)
- `trades.deleted_at` — trade record deletion (rare, for compliance)
- `trade_messages.deleted_at` — message deletion
- Actual data purged after retention period via cron job

### Automated Purge Queries

```sql
-- Purge soft-deleted users after 30-day grace period
DELETE FROM users WHERE deleted_at < NOW() - INTERVAL '30 days';

-- Purge chat messages older than 2 years (SRS 7.2)
DELETE FROM trade_messages WHERE created_at < NOW() - INTERVAL '2 years';
DELETE FROM dispute_messages WHERE created_at < NOW() - INTERVAL '2 years';

-- Archive old trades to cold storage (PostgreSQL table partitioning recommended)
CREATE TABLE trades_archive PARTITION OF trades FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');
```

### Nothing Ever Truly Deletes

Financial records (`trades`, `wallet_transactions`, `disputes`, `audit_logs`) are **never** physically deleted. They are marked deleted but retained per regulatory requirements (SRS 7.2):
- Trade records: 7 years
- KYC documents: 5 years
- Audit logs: Permanent

## Indexing Strategy

### High-Selectivity Indexes
- `users.email` — login lookups
- `users.username` — public profile lookups
- `wallet_transactions.tx_hash` — blockchain lookup by hash

### Composite Query Indexes
- `active_trade_ads` view uses: `(type, status, is_public, crypto_id, fiat_id)` for listing ads
- `trades` uses: `(status, created_at)` for finding active trades
- `notifications` uses: `(user_id, is_read, created_at DESC)` for inbox queries

### GIN Indexes (JSONB/Arrays)
- `trade_ads.payment_methods` — array containment queries (`@>`)
- `disputes.evidence_links` — JSONB search
- `audit_logs` JSONB fields (for change tracking)

### Partial Indexes
- `users.deleted_at` WHERE `deleted_at IS NOT NULL` — only soft-deleted rows
- `trades` WHERE `status IN ('pending', 'active', 'paid')` — active trades only

## Views

### `active_trade_ads`
Joins ads with user reputation stats. Used by public ad listing API.

### `completed_trades`
Joins trades with both buyer and seller feedback. Used in user history and analytics.

### `public_user_profiles`
Aggregated user stats excluding sensitive data. Used in public profiles.

### `open_disputes`
Filtered list of disputes needing admin attention. Used in admin dashboard.

## Migration Order

Migrations MUST be applied in numerical order due to foreign key dependencies:

000 (types) → 001 (users) → 002 (lookups) → 003 (kyc) → 004 (ads) → 005 (trades) → 006 (feedback) → 007 (disputes) → 008 (wallets) → 009 (notifications) → 010 (referrals) → 011 (merchant) → 012 (audit/admin) → 013 (monitoring) → 014 (views) → 015 (triggers) → 016 (seed)

**Rollback note:** Dropping tables in reverse order.

## Future Extensions

Tables/features to add in later phases:
- `blockchain_links` — escrow contract tracking (not needed if off-chain escrow)
- `two_factor_backup_codes` — 2FA backup code storage
- `api_keys` — Merchant API key management
- `subscription_payments` — Recurring merchant billing
- `chat_conversations` — General messaging system (beyond trade-specific)
- `support_tickets` — User support system
- `kyc_verification_log` — External KYC provider webhook logs
- `withdrawal_whitelist` — User-approved withdrawal addresses (SRS FR-309)
- `sms_logs` — SMS delivery tracking
- `email_logs` — Email delivery tracking
- `system_metrics` — Time-series metrics for Grafana

---

Questions? Refer to `pkg/database/migrations/README.md` for migration best practices.
