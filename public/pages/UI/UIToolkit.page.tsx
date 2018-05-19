import * as React from "react";
import { CurrentUser, IdeaStatus } from "@fider/models";
import { Heading, Button, List, UserName, ListItem, Toggle, Gravatar, ShowTag, Segment, Segments, ShowIdeaStatus } from "@fider/components";
import { User, UserRole, Tag } from "@fider/models";
import { notify } from "@fider/services";

const jonSnow: User = {
  id: 0,
  name: "Jon Snow",
  role: UserRole.Administrator
};

const aryaStark: User = {
  id: 0,
  name: "Arya Snow",
  role: UserRole.Visitor
};

const easyTag: Tag = { id: 2, slug: "easy", name: "Easy", color: "00a3a5", isPublic: true };
const hardTag: Tag = { id: 3, slug: "hard", name: "Hard", color: "ad43ec", isPublic: false };

export class UIToolkitPage extends React.Component<{}, {}> {
  public render() {
    return (
      <div id="p-ui-toolkit" className="page container">
        <h1>Heading 1</h1>
        <h2>Heading 2</h2>
        <h3>Heading 3</h3>
        <h4>Heading 4</h4>
        <h5>Heading 5</h5>
        <p className="primary">Primary Text</p>
        <p className="text">General Text</p>
        <p className="info">Info Text</p>
        <p className="primary-hover clickable">Primary Text (on hover) and clickable</p>

        <Segment>
          <h2>The title</h2>
          <p>The content goes here</p>
        </Segment>

        <Segments>
          <Segment>
            <p>First Segment</p>
          </Segment>
          <Segment>
            <p>Second Segment</p>
          </Segment>
          <Segment>
            <p>Third Segment</p>
          </Segment>
        </Segments>

        <List>
          <ListItem>
            <Gravatar user={jonSnow} /> <UserName user={jonSnow} />
          </ListItem>
          <ListItem>
            <Gravatar user={aryaStark} /> <UserName user={aryaStark} />
          </ListItem>
        </List>

        <Heading title="Page Heading" icon="settings" subtitle="This is a page heading" />

        <Heading
          title="Section Heading"
          icon="lightbulb"
          subtitle="This is a page heading"
          level={3}
          dividing={true}
          iconClassName="primary"
        />

        <h1>Buttons</h1>
        <List>
          <ListItem>
            <Button size="large">Large Default</Button>
            <Button color="positive" size="large">
              Large Positive
            </Button>
            <Button color="danger" size="large">
              Large Danger
            </Button>
          </ListItem>

          <ListItem>
            <Button size="normal">Normal Default</Button>
            <Button color="positive" size="normal">
              Normal Positive
            </Button>
            <Button color="danger" size="normal">
              Normal Danger
            </Button>
          </ListItem>

          <ListItem>
            <Button size="small">Small Default</Button>
            <Button color="positive" size="small">
              Small Positive
            </Button>
            <Button color="danger" size="small">
              Small Danger
            </Button>
          </ListItem>

          <ListItem>
            <Button size="tiny">Tiny Default</Button>
            <Button color="positive" size="tiny">
              Tiny Positive
            </Button>
            <Button color="danger" size="tiny">
              Tiny Danger
            </Button>
          </ListItem>

          <ListItem>
            <Button size="mini">Mini Default</Button>
            <Button color="positive" size="mini">
              Mini Positive
            </Button>
            <Button color="danger" size="mini">
              Mini Danger
            </Button>
          </ListItem>

          <ListItem>
            <Button disabled={true}>Default</Button>
            <Button disabled={true} color="positive">
              Positive
            </Button>
            <Button disabled={true} color="danger">
              Danger
            </Button>
          </ListItem>
        </List>

        <h1>Toggle</h1>
        <List>
          <ListItem>
            <Toggle active={true} label="Active" />
          </ListItem>
          <ListItem>
            <Toggle active={false} label="Inactive" />
          </ListItem>
          <ListItem>
            <Toggle active={true} disabled={true} label="Disabled" />
          </ListItem>
        </List>

        <h1>Statuses</h1>
        <List>
          <ListItem>
            <ShowIdeaStatus status={IdeaStatus.Open} />
          </ListItem>
          <ListItem>
            <ShowIdeaStatus status={IdeaStatus.Planned} />
          </ListItem>
          <ListItem>
            <ShowIdeaStatus status={IdeaStatus.Started} />
          </ListItem>
          <ListItem>
            <ShowIdeaStatus status={IdeaStatus.Duplicate} />
          </ListItem>
          <ListItem>
            <ShowIdeaStatus status={IdeaStatus.Completed} />
          </ListItem>
          <ListItem>
            <ShowIdeaStatus status={IdeaStatus.Declined} />
          </ListItem>
        </List>

        <h1>Tags</h1>
        <List>
          <ListItem>
            <ShowTag tag={easyTag} size="large" />
            <ShowTag tag={hardTag} size="large" />
            <ShowTag tag={easyTag} circular={true} size="large" />
            <ShowTag tag={hardTag} circular={true} size="large" />
          </ListItem>
          <ListItem>
            <ShowTag tag={easyTag} size="normal" />
            <ShowTag tag={hardTag} size="normal" />
            <ShowTag tag={easyTag} circular={true} size="normal" />
            <ShowTag tag={hardTag} circular={true} size="normal" />
          </ListItem>
          <ListItem>
            <ShowTag tag={easyTag} size="small" />
            <ShowTag tag={hardTag} size="small" />
            <ShowTag tag={easyTag} circular={true} size="small" />
            <ShowTag tag={hardTag} circular={true} size="small" />
          </ListItem>
          <ListItem>
            <ShowTag tag={easyTag} size="tiny" />
            <ShowTag tag={hardTag} size="tiny" />
            <ShowTag tag={easyTag} circular={true} size="tiny" />
            <ShowTag tag={hardTag} circular={true} size="tiny" />
          </ListItem>
          <ListItem>
            <ShowTag tag={easyTag} size="mini" />
            <ShowTag tag={hardTag} size="mini" />
            <ShowTag tag={easyTag} circular={true} size="mini" />
            <ShowTag tag={hardTag} circular={true} size="mini" />
          </ListItem>
        </List>

        <h1>Notification</h1>
        <List>
          <ListItem>
            <Button onClick={async () => notify.success("Congratulations! It worked!")}>Success</Button>
            <Button onClick={async () => notify.error("Something went wrong...")}>Error</Button>
          </ListItem>
        </List>
      </div>
    );
  }
}
