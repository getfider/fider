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

- `p-<page_name>` is the HTML ID of each page component. This is truly unique and should be used to provide page isolated style. E.g.: `p-home`, `p-user-settings`;

- `c-<component_name>` is the main class for each component. Inner components should follow its parent component name as such `c-<component_name>-<inner_component_name>`. E.g: `c-toggle` and `c-toggle-label`;

- `m-<state>` is used to alter the style of a component and it's always related to its component. So `m-disable` of `c-toggle` is different than `m-disable` of `c-button`, even though they have same name.

- `l-<name>` is used to general layout styling. It's does not need to be related to any component or page.

- `is-<state>`, `has-<state>` or simply `<state>` are global style modifiers that have a broader impact.

- `js-<name>` is used to JavaScript hooks, should be avoided as much as possible;

```

```
