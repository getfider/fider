import "./DesignSystem.page.scss"

import React, { useState } from "react"
import { PostStatus, UserStatus } from "@fider/models"
import {
  PageTitle,
  Button,
  UserName,
  Toggle,
  Avatar,
  ShowTag,
  ShowPostStatus,
  Moment,
  Loader,
  Form,
  Input,
  TextArea,
  RadioButton,
  Select,
  SelectOption,
  ButtonClickEvent,
  Message,
  Hint,
  AvatarStack,
  SocialSignInButton,
  Dropdown,
  Icon,
  Checkbox,
  ImageUploader,
} from "@fider/components"
import { User, UserRole, Tag } from "@fider/models"
import { notify, Failure } from "@fider/services"
import { HStack, VStack } from "@fider/components/layout"
import IconLightBulb from "@fider/assets/images/heroicons-light-bulb.svg"
import IconSearch from "@fider/assets/images/heroicons-search.svg"

const jonSnow: User = {
  id: 0,
  name: "Jon Snow",
  role: UserRole.Administrator,
  status: UserStatus.Active,
  avatarURL:
    "https://images.unsplash.com/photo-1491528323818-fdd1faba62cc?ixlib=rb-1.2.1&ixqx=1JzWlMeJDF&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80",
}

const aryaStark: User = {
  id: 0,
  name: "Arya Snow",
  role: UserRole.Visitor,
  status: UserStatus.Active,
  avatarURL: "https://images.unsplash.com/photo-1550525811-e5869dd03032?ixlib=rb-1.2.1&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80",
}

const robStark: User = {
  id: 0,
  name: "Robert Stark",
  role: UserRole.Visitor,
  status: UserStatus.Active,
  avatarURL:
    "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?ixlib=rb-1.2.1&ixqx=1JzWlMeJDF&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2.25&w=256&h=256&q=80",
}

const easyTag: Tag = { id: 1, slug: "easy", name: "Easy", color: "82c460", isPublic: true }
const normalTag: Tag = { id: 2, slug: "normal", name: "Normal", color: "ebb134", isPublic: false }
const hardTag: Tag = { id: 3, slug: "hard", name: "Hard", color: "9c3630", isPublic: false }

const visibilityPublic = { label: "Public", value: "public" }
const visibilityPrivate = { label: "Private", value: "private" }

const DesignSystemPage = () => {
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

  const showLoading = (e: ButtonClickEvent): void => {
    return e.preventEnable()
  }

  const forceError = async () => {
    setError({
      errors: [
        { field: "title", message: "Title is mandatory" },
        { field: "description", message: "Error #1" },
        { field: "description", message: "Error #2" },
        { field: "status", message: "Status is mandatory" },
        { field: "logo", message: "Logo is mandatory" },
        { field: "gdpr", message: "You have to agree!" },
      ],
    })
  }

  return (
    <div id="p-ui-toolkit" className="page container">
      <h2 className="text-display2 mb-2">1. Colors</h2>

      <div className="color-scale">
        <div className="color gray-900" />
        <div className="color gray-800" />
        <div className="color gray-700" />
        <div className="color gray-600" />
        <div className="color gray-500" />
        <div className="color gray-400" />
        <div className="color gray-300" />
        <div className="color gray-200" />
        <div className="color gray-100" />
        <div className="color gray-50" />
      </div>
      <div className="color-scale">
        <div className="color green-900" />
        <div className="color green-800" />
        <div className="color green-700" />
        <div className="color green-600" />
        <div className="color green-500" />
        <div className="color green-400" />
        <div className="color green-300" />
        <div className="color green-200" />
        <div className="color green-100" />
        <div className="color green-50" />
      </div>
      <div className="color-scale">
        <div className="color red-900" />
        <div className="color red-800" />
        <div className="color red-700" />
        <div className="color red-600" />
        <div className="color red-500" />
        <div className="color red-400" />
        <div className="color red-300" />
        <div className="color red-200" />
        <div className="color red-100" />
        <div className="color red-50" />
      </div>
      <div className="color-scale">
        <div className="color blue-900" />
        <div className="color blue-800" />
        <div className="color blue-700" />
        <div className="color blue-600" />
        <div className="color blue-500" />
        <div className="color blue-400" />
        <div className="color blue-300" />
        <div className="color blue-200" />
        <div className="color blue-100" />
        <div className="color blue-50" />
      </div>
      <div className="color-scale">
        <div className="color yellow-900" />
        <div className="color yellow-800" />
        <div className="color yellow-700" />
        <div className="color yellow-600" />
        <div className="color yellow-500" />
        <div className="color yellow-400" />
        <div className="color yellow-300" />
        <div className="color yellow-200" />
        <div className="color yellow-full" />
        <div className="color yellow-50" />
      </div>
      <div className="color-scale">
        <div className="color bg-primary-dark" />
        <div className="color bg-primary-base" />
        <div className="color bg-primary-light" />
      </div>

      <h2 className="text-display2 mb-3 mt-6">2. Text</h2>

      <VStack>
        <span className="text-display2">text-display2</span>
        <span className="text-display">text-display</span>
        <span className="text-title">text-title</span>
        <span className="text-body">text-body</span>
        <span className="text-sm">text-sm</span>
        <span className="text-xs">text-xs</span>
        <span className="text-2xs">text-2xs</span>
        <span className="text-muted">text-muted</span>
        <span className="text-category">text-category</span>
      </VStack>

      <h2 className="text-display2 mb-3 mt-6">3. Avatars</h2>

      <VStack>
        <HStack>
          <Avatar user={jonSnow} /> <UserName user={jonSnow} />
        </HStack>
        <HStack>
          <Avatar user={aryaStark} /> <UserName user={aryaStark} />
        </HStack>
        <HStack>
          <AvatarStack users={[jonSnow, aryaStark, robStark]} />
        </HStack>
      </VStack>

      <h2 className="text-display2 mb-3 mt-6">4. Page Title</h2>

      <PageTitle title="Page Title" subtitle="This is a page subtitle" />

      <h2 className="text-display2 mb-3 mt-6">5. Buttons</h2>

      <VStack>
        <HStack>
          <Button variant="primary" size="large">
            Primary
          </Button>
          <Button size="large">
            <Icon sprite={IconLightBulb} /> <span>Secondary</span>
          </Button>
          <Button size="large">
            <Icon sprite={IconLightBulb} />
          </Button>
          <Button size="large">Secondary</Button>
          <Button variant="tertiary" size="large">
            Tertiary
          </Button>
          <Button variant="danger" size="large">
            Danger
          </Button>
        </HStack>

        <HStack>
          <Button variant="primary">Primary</Button>
          <Button>
            <Icon sprite={IconLightBulb} /> <span>Secondary</span>
          </Button>
          <Button>
            <Icon sprite={IconLightBulb} />
          </Button>
          <Button>Secondary</Button>
          <Button variant="tertiary">Tertiary</Button>
          <Button variant="danger">Danger</Button>
        </HStack>

        <HStack>
          <Button variant="primary" size="small">
            Primary
          </Button>
          <Button size="small">
            <Icon sprite={IconLightBulb} /> <span>Secondary</span>
          </Button>
          <Button size="small">
            <Icon sprite={IconLightBulb} />
          </Button>
          <Button size="small">Secondary</Button>
          <Button variant="tertiary" size="small">
            Tertiary
          </Button>
          <Button variant="danger" size="small">
            Danger
          </Button>
        </HStack>

        <HStack>
          <Button href="#" variant="primary">
            Link
          </Button>
          <Button href="#">
            <Icon sprite={IconLightBulb} /> <span>Link</span>
          </Button>
          <Button href="#">
            <Icon sprite={IconLightBulb} />
          </Button>
          <Button>Link</Button>
          <Button variant="tertiary" href="#">
            Link
          </Button>
          <Button href="#" variant="danger">
            Link
          </Button>
        </HStack>

        <HStack>
          <Button disabled={true} variant="primary">
            Primary
          </Button>
          <Button disabled={true}>
            <Icon sprite={IconLightBulb} /> <span>Secondary</span>
          </Button>
          <Button disabled={true}>
            <Icon sprite={IconLightBulb} />
          </Button>
          <Button disabled={true}>Secondary</Button>
          <Button disabled={true} variant="tertiary">
            Tertiary
          </Button>
          <Button disabled={true} variant="danger">
            Danger
          </Button>
        </HStack>

        <HStack>
          <Button variant="primary" onClick={showLoading}>
            Loading
          </Button>
          <Button onClick={showLoading}>
            <Icon sprite={IconLightBulb} /> <span>Loading</span>
          </Button>
          <Button onClick={showLoading}>
            <Icon sprite={IconLightBulb} />
          </Button>
          <Button onClick={showLoading}>Loading</Button>
          <Button variant="tertiary" onClick={showLoading}>
            Loading
          </Button>
          <Button variant="danger" onClick={showLoading}>
            Loading
          </Button>
        </HStack>

        <HStack>
          <SocialSignInButton option={{ displayName: "GitHub", provider: "github" }} />
          <SocialSignInButton option={{ displayName: "Facebook", provider: "facebook" }} />
          <SocialSignInButton option={{ displayName: "Google", provider: "google" }} />
        </HStack>
      </VStack>

      <h2 className="text-display2 mb-3 mt-6">6. Toggle</h2>

      <VStack>
        <Toggle active={true} label="Active" />
        <Toggle active={false} label="Inactive" />
        <Toggle active={true} disabled={true} label="Disabled" />
      </VStack>

      <h2 className="text-display2 mb-3 mt-6">7. Statuses</h2>

      <VStack>
        <ShowPostStatus status={PostStatus.Open} />
        <ShowPostStatus status={PostStatus.Planned} />
        <ShowPostStatus status={PostStatus.Started} />
        <ShowPostStatus status={PostStatus.Duplicate} />
        <ShowPostStatus status={PostStatus.Completed} />
        <ShowPostStatus status={PostStatus.Declined} />
      </VStack>

      <h2 className="text-display2 mb-3 mt-6">8. Tags</h2>

      <VStack>
        <HStack>
          <ShowTag tag={easyTag} />
          <ShowTag tag={normalTag} />
          <ShowTag tag={hardTag} />
        </HStack>
        <HStack>
          <ShowTag tag={easyTag} circular={true} />
          <ShowTag tag={normalTag} circular={true} />
          <ShowTag tag={hardTag} circular={true} />
        </HStack>
      </VStack>

      <h2 className="text-display2 mb-3 mt-6">9. Notification</h2>

      <HStack>
        <Button onClick={notifySuccess}>Success</Button>
        <Button onClick={notifyError}>Error</Button>
      </HStack>

      <h2 className="text-display2 mb-3 mt-6">10. Moment</h2>

      <VStack>
        <div>
          Relative: <Moment date="2017-06-03T16:55:06.815042Z" format="relative" />
        </div>
        <div>
          Short: <Moment date="2017-06-03T16:55:06.815042Z" format="short" />
        </div>
        <div>
          Full: <Moment date="2017-06-03T16:55:06.815042Z" format="full" />
        </div>
      </VStack>

      <h2 className="text-display2 mb-3 mt-6">11. Loader</h2>

      <Loader />

      <h2 className="text-display2 mb-3 mt-6">12. Message</h2>

      <Message showIcon={true} type="error">
        Something went wrong.
      </Message>
      <Message showIcon={true} type="warning">
        Be careful!
      </Message>
      <Message showIcon={true} type="success">
        Your order has been confirmed.
      </Message>

      <h2 className="text-display2 mb-3 mt-6">13. Hints</h2>

      <Hint permanentCloseKey="ui-toolkip-example">Did you know that you can close this permanently?</Hint>
      <Hint>You can&apos;t close this one :)</Hint>

      <h2 className="text-display2 mb-3 mt-6">14. Form</h2>

      <Form error={error}>
        <Input label="Title" field="title">
          <p className="text-muted">This is the explanation for the field above.</p>
        </Input>
        <Input label="Disabled!" field="unamed" disabled={true} value={"you can't change this!"} />
        <Input label="Name" field="name" placeholder={"Your name goes here..."} />
        <Input label="Subdomain" field="subdomain" suffix=".fider.io" />
        <Input label="Email" field="email" suffix={<Button variant="primary">Sign in</Button>} />
        <TextArea label="Description" field="description" minRows={5}>
          <p className="text-muted">This textarea resizes as you type.</p>
        </TextArea>
        <Input field="age" placeholder="This field doesn't have a label" />

        <div className="grid gap-4 grid-cols-4">
          <Input label="Title1" field="title1" />
          <Input label="Title2" field="title2" />
          <Input label="Title3" field="title3" />
          <RadioButton label="Visibility" field="visibility" defaultOption={visibilityPublic} options={[visibilityPrivate, visibilityPublic]} />
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

        <ImageUploader label="Logo" field="logo" onChange={console.log} />

        <Checkbox field="gdpr">I agree with this policy.</Checkbox>

        <Button onClick={forceError}>Save</Button>
      </Form>

      <h2 className="text-display2 mb-3 mt-6">15. Dropdown</h2>

      <Dropdown renderHandle={<Avatar user={aryaStark} />}>
        <div className="p-2 text-bold uppercase">Arya Stark</div>
        <Dropdown.ListItem>Account Settings</Dropdown.ListItem>
        <Dropdown.ListItem>My Profile</Dropdown.ListItem>
        <Dropdown.Divider />
        <Dropdown.ListItem>Sign Out</Dropdown.ListItem>
      </Dropdown>

      <h2 className="text-display2 mb-3 mt-6">16. Search</h2>

      <Input field="search" placeholder="Search..." icon={IconSearch} />
    </div>
  )
}

export default DesignSystemPage
