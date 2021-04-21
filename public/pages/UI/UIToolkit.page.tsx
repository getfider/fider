import "./UIToolkit.page.scss"

import React, { useState } from "react"
import { PostStatus, UserStatus } from "@fider/models"
import {
  Heading,
  Button,
  List,
  UserName,
  ListItem,
  Toggle,
  Avatar,
  ShowTag,
  Segment,
  Segments,
  ShowPostStatus,
  Moment,
  Loader,
  Form,
  Input,
  TextArea,
  RadioButton,
  Select,
  Field,
  SelectOption,
  ButtonClickEvent,
  Message,
  Hint,
} from "@fider/components"
import { User, UserRole, Tag } from "@fider/models"
import { notify, Failure } from "@fider/services"
import { DropDown, DropDownItem } from "@fider/components"
import { FaSearch, FaRegLightbulb, FaCogs } from "react-icons/fa"

const jonSnow: User = {
  id: 0,
  name: "Jon Snow",
  role: UserRole.Administrator,
  status: UserStatus.Active,
  avatarURL: "/avatars/letter/0/Jon%20Snow",
}

const aryaStark: User = {
  id: 0,
  name: "Arya Snow",
  role: UserRole.Visitor,
  status: UserStatus.Active,
  avatarURL: "/avatars/letter/0/Arya%20Snow",
}

const easyTag: Tag = { id: 2, slug: "easy", name: "easy", color: "FB3A62", isPublic: true }
const hardTag: Tag = { id: 3, slug: "hard", name: "hard", color: "fbca04", isPublic: false }

const visibilityPublic = { label: "Public", value: "public" }
const visibilityPrivate = { label: "Private", value: "private" }

const UIToolkitPage = () => {
  const [error, setError] = useState<Failure | undefined>(undefined)

  const notifyError = async () => {
    notify.error("Something went wrong...")
  }

  const notifySuccess = async () => {
    notify.success("Congratulations! It worked!")
  }

  const notifyStatusChange = (opt?: SelectOption) => {
    if (opt) {
      notify.success(opt.value)
    }
  }

  const showLoading = async (e: ButtonClickEvent) => {
    return e.preventEnable()
  }

  const forceError = async () => {
    setError({
      errors: [
        { field: "title", message: "Title is mandatory" },
        { field: "description", message: "Error #1" },
        { field: "description", message: "Error #2" },
        { field: "status", message: "Status is mandatory" },
      ],
    })
  }

  const renderText = (item?: DropDownItem) => {
    if (item) {
      return `${item.label} (value: ${item.value})`
    }
    return <span>No country is selected</span>
  }

  const renderControl = (item?: DropDownItem) => {
    if (item) {
      return item.render
    }
    return <span>...</span>
  }

  return (
    <div id="p-ui-toolkit" className="page container">
      <div className="color-scale">
        <div className="color gray-900" />
        <div className="color gray-700" />
        <div className="color gray-500" />
        <div className="color gray-400" />
        <div className="color gray-300" />
        <div className="color gray-200" />
        <div className="color gray-100" />
        <div className="color gray-50" />
      </div>
      <div className="color-scale">
        <div className="color green-900" />
        <div className="color green-700" />
        <div className="color green-500" />
        <div className="color green-400" />
        <div className="color green-300" />
        <div className="color green-200" />
        <div className="color green-100" />
        <div className="color green-50" />
      </div>
      <div className="color-scale">
        <div className="color red-900" />
        <div className="color red-700" />
        <div className="color red-500" />
        <div className="color red-400" />
        <div className="color red-300" />
        <div className="color red-200" />
        <div className="color red-100" />
        <div className="color red-50" />
      </div>
      <div className="color-scale">
        <div className="color blue-900" />
        <div className="color blue-700" />
        <div className="color blue-500" />
        <div className="color blue-400" />
        <div className="color blue-300" />
        <div className="color blue-200" />
        <div className="color blue-100" />
        <div className="color blue-50" />
      </div>
      <div className="color-scale">
        <div className="color yellow-900" />
        <div className="color yellow-700" />
        <div className="color yellow-500" />
        <div className="color yellow-400" />
        <div className="color yellow-300" />
        <div className="color yellow-200" />
        <div className="color yellow-100" />
        <div className="color yellow-50" />
      </div>

      <h1>Heading 1</h1>
      <h2>Heading 2</h2>
      <h3>Heading 3</h3>
      <h4>Heading 4</h4>
      <h5>Heading 5</h5>
      <h6>Heading 6</h6>
      <p>General Text Paragraph</p>
      <p className="info">Info Text</p>

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
          <Avatar user={jonSnow} /> <UserName user={jonSnow} />
        </ListItem>
        <ListItem>
          <Avatar user={aryaStark} /> <UserName user={aryaStark} />
        </ListItem>
      </List>

      <Heading title="Page Heading" icon={FaCogs} subtitle="This is a page heading" />

      <Heading title="Section Heading" icon={FaRegLightbulb} subtitle="This is a page heading" size="small" dividing={true} />

      <h1>Buttons</h1>
      <List>
        <ListItem>
          <Button size="large">
            <FaRegLightbulb /> Large Icon
          </Button>
          <Button size="large">Large Default</Button>
          <Button color="positive" size="large">
            Large Positive
          </Button>
          <Button color="danger" size="large">
            Large Danger
          </Button>
          <Button color="cancel" size="large">
            Large Cancel
          </Button>
        </ListItem>

        <ListItem>
          <Button size="normal">
            <FaRegLightbulb /> Normal Icon
          </Button>
          <Button size="normal">Normal Default</Button>
          <Button color="positive" size="normal">
            Normal Positive
          </Button>
          <Button color="danger" size="normal">
            Normal Danger
          </Button>
          <Button color="cancel" size="normal">
            Normal Cancel
          </Button>
        </ListItem>

        <ListItem>
          <Button size="small">
            <FaRegLightbulb /> Small Icon
          </Button>
          <Button size="small">Small Default</Button>
          <Button color="positive" size="small">
            Small Positive
          </Button>
          <Button color="danger" size="small">
            Small Danger
          </Button>
          <Button color="cancel" size="small">
            Small Cancel
          </Button>
        </ListItem>

        <ListItem>
          <Button size="tiny">
            <FaRegLightbulb /> Tiny Icon
          </Button>
          <Button size="tiny">Tiny Default</Button>
          <Button color="positive" size="tiny">
            Tiny Positive
          </Button>
          <Button color="danger" size="tiny">
            Tiny Danger
          </Button>
          <Button color="cancel" size="tiny">
            Tiny Cancel
          </Button>
        </ListItem>

        <ListItem>
          <Button size="mini">
            <FaRegLightbulb /> Mini Icon
          </Button>
          <Button size="mini">Mini Default</Button>
          <Button color="positive" size="mini">
            Mini Positive
          </Button>
          <Button color="danger" size="mini">
            Mini Danger
          </Button>
          <Button color="cancel" size="mini">
            Mini Cancel
          </Button>
        </ListItem>

        <ListItem>
          <Button href="#">
            <FaRegLightbulb /> Link
          </Button>
          <Button href="#">Link</Button>
          <Button href="#" color="positive">
            Link
          </Button>
          <Button href="#" color="danger">
            Link
          </Button>
        </ListItem>

        <ListItem>
          <Button disabled={true}>
            <FaRegLightbulb /> Default
          </Button>
          <Button disabled={true}>Default</Button>
          <Button disabled={true} color="positive">
            Positive
          </Button>
          <Button disabled={true} color="danger">
            Danger
          </Button>
        </ListItem>

        <ListItem>
          <Button onClick={showLoading}>
            <FaRegLightbulb /> Loading
          </Button>
          <Button onClick={showLoading}>Loading</Button>
          <Button color="positive" onClick={showLoading}>
            Loading
          </Button>
          <Button color="danger" onClick={showLoading}>
            Loading
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
          <ShowPostStatus status={PostStatus.Open} />
        </ListItem>
        <ListItem>
          <ShowPostStatus status={PostStatus.Planned} />
        </ListItem>
        <ListItem>
          <ShowPostStatus status={PostStatus.Started} />
        </ListItem>
        <ListItem>
          <ShowPostStatus status={PostStatus.Duplicate} />
        </ListItem>
        <ListItem>
          <ShowPostStatus status={PostStatus.Completed} />
        </ListItem>
        <ListItem>
          <ShowPostStatus status={PostStatus.Declined} />
        </ListItem>
      </List>

      <h1>Tags</h1>
      <List>
        <ListItem>
          <ShowTag tag={easyTag} />
          <ShowTag tag={hardTag} />
          <ShowTag tag={easyTag} circular={true} />
          <ShowTag tag={hardTag} circular={true} />
        </ListItem>
      </List>

      <h1>Notification</h1>
      <List>
        <ListItem>
          <Button onClick={notifySuccess}>Success</Button>
          <Button onClick={notifyError}>Error</Button>
        </ListItem>
      </List>

      <h1>Moment</h1>
      <List>
        <ListItem>
          Relative: <Moment date="2017-06-03T16:55:06.815042Z" format="relative" />
        </ListItem>
        <ListItem>
          Short: <Moment date="2017-06-03T16:55:06.815042Z" format="short" />
        </ListItem>
        <ListItem>
          Full: <Moment date="2017-06-03T16:55:06.815042Z" format="full" />
        </ListItem>
      </List>

      <h1>Loader</h1>
      <Loader />

      <h1>Message Box</h1>
      <Message showIcon={true} type="error">
        Something went wrong.
      </Message>
      <Message showIcon={true} type="warning">
        Be careful!
      </Message>
      <Message showIcon={true} type="success">
        Your order has been confirmed.
      </Message>

      <h1>Hints</h1>
      <Hint permanentCloseKey="ui-toolkip-example">Did you know that you can close this permanently?</Hint>
      <Hint>You can&apos;t close this one :)</Hint>

      <h1>Form</h1>
      <Form error={error}>
        <Input label="Title" field="title">
          <p className="info">This is the explanation for the field above.</p>
        </Input>
        <Input label="Disabled!" field="unamed" disabled={true} value={"you can't change this!"} />
        <Input label="Name" field="name" placeholder={"Your name goes here..."} />
        <Input label="Subdomain" field="subdomain" suffix="fider.io" />
        <Input label="Email" field="email" suffix={<Button color="positive">Sign in</Button>} />
        <TextArea label="Description" field="description" minRows={5}>
          <p className="info">This textarea resizes as you type.</p>
        </TextArea>
        <Input field="age" placeholder="This field doesn't have a label" />

        <div className="row">
          <div className="col-md-3">
            <Input label="Title1" field="title1" />
          </div>
          <div className="col-md-3">
            <Input label="Title2" field="title2" />
          </div>
          <div className="col-md-3">
            <Input label="Title3" field="title3" />
          </div>
          <div className="col-md-3">
            <RadioButton label="Visibility" field="visibility" defaultOption={visibilityPublic} options={[visibilityPrivate, visibilityPublic]} />
          </div>
        </div>

        <Select
          label="Status"
          field="status"
          options={[
            { value: "open", label: "Open" },
            { value: "started", label: "Started" },
            { value: "planned", label: "Planned" },
          ]}
          onChange={notifyStatusChange}
        />

        <Field label="Number">
          <DropDown
            items={[
              { label: "One", value: "1" },
              { label: "Two", value: "2" },
              { label: "Three", value: "3" },
            ]}
            defaultValue={"1"}
            placeholder="Select a number"
          />
        </Field>

        <Field label="Country (custom render text)">
          <DropDown
            items={[
              { label: "Brazil", value: "br" },
              { label: "United States", value: "us" },
              { label: "Ireland", value: "ie" },
            ]}
            defaultValue={"1"}
            renderText={renderText}
            placeholder="Select a number"
          />
        </Field>

        <Field label="Color (custom render control)">
          <DropDown
            items={[
              { label: "Green", value: "green", render: <span style={{ color: "green" }}>Green</span> },
              { label: "Red", value: "red", render: <span style={{ color: "red" }}>Red</span> },
              { label: "Blue", value: "blue", render: <span style={{ color: "blue" }}>Blue</span> },
            ]}
            placeholder="Select a color"
            inline={true}
            style="simple"
            header="What color do you like the most?"
            renderControl={renderControl}
          />
        </Field>

        <Button onClick={forceError}>Save</Button>
      </Form>

      <Segment>
        <h1>Search</h1>
        <Input field="search" placeholder="Search..." icon={FaSearch} />
      </Segment>
    </div>
  )
}

export default UIToolkitPage
