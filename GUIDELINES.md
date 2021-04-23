# Server Development

**Models Naming**

- `app/models/action.<something>`: actions are based on user interaction for POST/PUT/PATCH requests. Actions map 1-to-1 with an Command. E.g: `action.CreateNewUser`;
- `app/models/dto.<something>`: A simple object used for data transfer between various packages/services. E.g: `dto.NewUserInfo`;
- `app/models/entity.<something>`: An object that is mapped to a database table. E.g: `entity.User`;
- `app/models/cmd.<something>` something that must be done and potentially return some value. E.g.: `cmd.HttpRequest`, `cmd.LogDebug`, `cmd.SendMail`, `cmd.CreateNewUser`;
- `app/models/query.<something>` get some information from somewhere. E.g.: `query.GetUserById`, `query.GetAllPosts`;

# UI Development

## React / CSS / HTML Convention

**Folder Structure**

```javascript
public/components // Shared/Basic Components

public/pages // Pages

public/pages/Home // Home page component folder
	-> index.ts // exporter
	-> Home.page.scss // Page specific styles
	-> Home.page.spec.tsx // Page Component unit tests
	-> Home.page.tsx // Page Component
	-> ./components // Inner components of home page
```

**CSS Naming**

Fider uses a combination of BEM and Utility Classes.

- `p-<page_name>` is the HTML ID of each page component. This is truly unique and should be used to provide page isolated style. E.g.: `p-home`, `p-user-settings`;

- `c-<component_name>` is the "Block" class for each component. Elements should follow its parent component name as such `c-<component_name>__<element>`. E.g: `c-toggle` and `c-toggle__label`;

- `{block}--<state>` is used to alter the style of its. . E.g: `c-toggle` and `c-toggle--checked`;

- `is-<state>`, `has-<state>` are global style modifiers that have a broader impact.

- Utility classes do not have a preffix.
