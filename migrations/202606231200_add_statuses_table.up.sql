-- Tenant-configurable custom post statuses. Implements feedback.fider.io/posts/111.
--
-- Each tenant owns a list of statuses they can use to label posts. The 6 existing
-- built-in statuses (Open, Started, Completed, Declined, Planned, Duplicate) are
-- seeded as system rows for every existing tenant via the data migration below,
-- so behavior is unchanged for sites that don't touch the new admin UI.
--
-- "kind" maps the custom name to a fixed semantic so the rest of Fider still
-- knows whether a status counts as open / active / closed-completed / closed-declined
-- / duplicate, regardless of what the tenant labels it.

CREATE TABLE statuses (
    id           SERIAL PRIMARY KEY,
    tenant_id    INTEGER     NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    slug         VARCHAR(50) NOT NULL,
    label        VARCHAR(50) NOT NULL,
    kind         VARCHAR(20) NOT NULL,
    color        VARCHAR(20) NOT NULL DEFAULT 'blue',
    icon         VARCHAR(60) NOT NULL DEFAULT 'lightbulb',
    show_on_home BOOLEAN     NOT NULL DEFAULT TRUE,
    filterable   BOOLEAN     NOT NULL DEFAULT TRUE,
    sort_order   INTEGER     NOT NULL DEFAULT 0,
    is_system    BOOLEAN     NOT NULL DEFAULT FALSE,
    is_active    BOOLEAN     NOT NULL DEFAULT TRUE,
    legacy_enum  INTEGER     NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, slug),
    CHECK (kind IN ('open', 'active', 'closed-completed', 'closed-declined', 'duplicate'))
);

CREATE INDEX statuses_tenant_active_idx ON statuses (tenant_id, is_active);
CREATE INDEX statuses_tenant_legacy_idx ON statuses (tenant_id, legacy_enum);

-- Seed the 6 built-in statuses for every existing tenant. PostDeleted (6) is
-- intentionally not seeded — it's a hidden tombstone state, not user-settable.
INSERT INTO statuses (tenant_id, slug, label, kind, color, icon, show_on_home, filterable, sort_order, is_system, legacy_enum)
SELECT id, 'open',      'Open',      'open',             'blue',   'lightbulb',        FALSE, TRUE, 10, TRUE, 0 FROM tenants
UNION ALL
SELECT id, 'planned',   'Planned',   'active',           'blue',   'thumbsup',         TRUE,  TRUE, 20, TRUE, 4 FROM tenants
UNION ALL
SELECT id, 'started',   'Started',   'active',           'blue',   'sparkles-outline', TRUE,  TRUE, 30, TRUE, 1 FROM tenants
UNION ALL
SELECT id, 'completed', 'Completed', 'closed-completed', 'green',  'check-circle',     TRUE,  TRUE, 40, TRUE, 2 FROM tenants
UNION ALL
SELECT id, 'declined',  'Declined',  'closed-declined',  'red',    'thumbsdown',       TRUE,  TRUE, 50, TRUE, 3 FROM tenants
UNION ALL
SELECT id, 'duplicate', 'Duplicate', 'duplicate',        'yellow', 'duplicate',        TRUE,  TRUE, 60, TRUE, 5 FROM tenants;
