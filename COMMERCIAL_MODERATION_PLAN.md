# Commercial Content Moderation Restructure Plan

## Overview

Moving Fider's content moderation feature from open source (AGPL) to commercial licensing using an "open core" model. The commercial code will reside in a `commercial/` folder with a restrictive license, while the open source core provides infrastructure and gracefully degrades when commercial features aren't licensed.

## Key Architectural Decisions

### 1. Webpack Integration (Simple Solution)
- Update `webpack.config.js` to scan both folders:
  ```js
  paths: [
    ...glob.sync(`./public/**/*.{html,tsx}`, { nodir: true }),
    ...glob.sync(`./commercial/**/*.{html,tsx}`, { nodir: true })
  ]
  ```
- Commercial code gets bundled but legal license prevents usage
- No complex conditional compilation needed

### 2. Backend Package Separation Strategy

**Handlers - Dynamic Route Registration**
- Open source: Register stub routes returning "upgrade" messages
- Commercial: Override routes with real handlers via dynamic registration
- License validation controls active routes

**Services - Stub + Override Pattern**  
- Open source: Register stub bus handlers returning "not licensed" errors
- Commercial: Override with real implementations via bus registration
- Same command/query types, different implementations

## What Moves to Commercial (Strong Legal Protection)

### Commercial Folder Structure
```
commercial/
├── LICENSE (restrictive commercial license)
├── pages/
│   └── Administration/
│       └── ContentModeration.page.tsx (~300 lines)
├── components/
│   └── ModerationIndicator.tsx
├── handlers/
│   ├── moderation.go (ModerationPage, GetModerationItems, GetModerationCount)
│   └── apiv1/
│       ├── moderation.go (all API endpoints)
│       └── moderation_test.go
├── services/
│   └── moderation.go (approve/decline business logic)
└── init.go (service/route registration)
```

### Frontend (Moves to Commercial)
- Complete moderation admin UI (ContentModeration.page.tsx + styles)
- Moderation indicator component showing pending counts
- All moderation-specific UI components

### Backend (Moves to Commercial)
- All HTTP handlers (`ModerationPage`, `GetModerationItems`, `GetModerationCount`)
- All API endpoints (`/api/v1/admin/moderation/*` routes)
- Business logic implementations (approve/decline/verify/block functions)
- Route registrations for moderation endpoints
- Tests for commercial functionality

## What Stays Open Source (Minimal Infrastructure)

### Database/Models (Cannot Move - Bus System Dependency)
- `app/models/cmd/moderation.go` - Command type definitions
- `app/models/query/moderation.go` - Query type definitions
- Database schema, migrations (`is_moderation_enabled`, `is_approved` columns)
- Entity definitions (`ModerationItem` structs)

### Settings Infrastructure
- Privacy settings toggle (shows "upgrade" message if unlicensed)
- Basic tenant property (`isModerationEnabled`)
- License validation service

### Content Flow Logic
- Logic that marks content as needing approval
- Basic "content requires moderation" checks throughout codebase
- Content queuing when moderation is enabled

### Locale Files
- All `moderation.*` translation strings (used by upgrade messages too)

## Technical Implementation Details

### 1. Handler Registration Pattern
```go
// Open source routes.go - stub routes
ui.Get("/admin/moderation", upgradeHandler("content-moderation"))
ui.Get("/_api/admin/moderation/items", upgradeHandler("content-moderation"))

// Commercial init.go - overrides routes
func init() {
    if license.IsCommercialFeatureEnabled("content-moderation") {
        web.RegisterRoute("GET", "/admin/moderation", handlers.ModerationPage())
        web.RegisterRoute("GET", "/_api/admin/moderation/items", handlers.GetModerationItems())
        // ... all other moderation routes
    }
}
```

### 2. Service Registration Pattern
```go
// Open source postgres.go - register stubs
func (s *Service) Init() {
    bus.AddHandler(approvePostStub)   // Returns "feature not licensed"
    bus.AddHandler(declinePostStub)
    bus.AddHandler(getModerationItemsStub)
    // ... other stubs
}

// Commercial service init - override with real handlers
func (cs *CommercialService) Init() {
    if license.IsCommercialFeatureEnabled("content-moderation") {
        bus.AddHandler(approvePost)       // Real implementation
        bus.AddHandler(declinePost)
        bus.AddHandler(getModerationItems)
        // ... other real handlers
    }
}
```

### 3. License Validation Service
```go
// app/services/license.go
type LicenseService interface {
    IsCommercialFeatureEnabled(feature string) bool
}

// Implementation checks for valid commercial license
// Controls route registration and service overrides
```

## Implementation Phases

### Phase 1: Setup Commercial Infrastructure
1. Create `commercial/` folder structure
2. Add restrictive LICENSE file to commercial folder
3. Update webpack.config.js to scan commercial folder
4. Create license validation service interface

### Phase 2: Move Frontend Components
1. Move `ContentModeration.page.tsx` to `commercial/pages/Administration/`
2. Move `ModerationIndicator.tsx` to `commercial/components/`
3. Test that webpack builds both locations correctly
4. Add license checks in components to show upgrade messages

### Phase 3: Implement Backend Service Separation  
1. Create stub implementations for all moderation commands/queries
2. Register stubs in open source postgres service
3. Move real implementations to `commercial/services/moderation.go`
4. Implement commercial service registration with license checks

### Phase 4: Implement Route Separation
1. Replace direct handler calls with upgrade handlers in routes.go
2. Move real handlers to `commercial/handlers/`
3. Implement dynamic route registration in commercial init
4. Test route overriding works correctly

### Phase 5: Testing & Validation
1. Move tests to commercial folder
2. Test open source build without commercial features
3. Test commercial build with license validation
4. Verify graceful degradation and upgrade messaging
5. Test that forking open source works without commercial parts

## Key Benefits

### Strong Commercial Protection
- Recreating moderation requires rebuilding entire UI + API + business logic
- Substantial engineering effort (days/weeks) to replicate functionality
- Clear legal separation under different licenses

### Clean Architecture  
- Open source provides database infrastructure and settings
- Commercial enhances with actual moderation functionality
- No broken states - content flows normally when moderation disabled
- True open core model: commercial builds on open source foundation

### Simple Build Process
- Single repository with clear folder separation
- Standard webpack build process
- License provides protection, not technical hiding
- Easy development workflow

## Migration Considerations

### Existing Moderation Data
- All existing moderation settings and data remain compatible
- Database schema stays in open source (infrastructure)
- Only the management interface becomes commercial

### User Experience
- **With License**: Full moderation functionality as before
- **Without License**: Moderation simply disabled, standard user flow
- **Upgrade Path**: Clear messaging and sales funnel for commercial features

### Development Workflow  
- Developers can see all code (legal license controls usage)
- Standard build process works for both open source and commercial
- Clean separation makes it easy to add more commercial features

## Success Metrics

1. **Legal Protection**: Someone forking open source cannot easily recreate moderation
2. **Functional Separation**: Open source works perfectly without moderation
3. **Build Compatibility**: Both open source and commercial versions build successfully  
4. **Clean Boundaries**: Clear understanding of what's core vs commercial
5. **Scalability**: Pattern works for future commercial features

This plan provides genuine open core protection while maintaining clean architecture and manageable implementation complexity.