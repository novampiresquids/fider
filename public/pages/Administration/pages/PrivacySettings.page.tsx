import React from "react"
import { Toggle, Form, Field } from "@fider/components"
import { actions, notify, Fider } from "@fider/services"
import { AdminBasePage } from "@fider/pages/Administration/components/AdminBasePage"

interface PrivacySettingsPageState {
  isPrivate: boolean
}

export default class PrivacySettingsPage extends AdminBasePage<any, PrivacySettingsPageState> {
  public id = "p-admin-privacy"
  public name = "privacy"
  public title = "Privacy"
  public subtitle = "Manage your board privacy"

  constructor(props: any) {
    super(props)

    this.state = {
      isPrivate: Fider.session.tenant.isPrivate,
    }
  }

  private toggle = async (active: boolean) => {
    this.setState(
      () => ({
        isPrivate: active,
      }),
      async () => {
        const response = await actions.updateTenantPrivacy(this.state.isPrivate, Fider.session.tenant.id)
        if (response.ok) {
          notify.success("Your privacy settings have been saved.")
        }
      }
    )
  }

  public content() {
    return (
      <Form>
        <Field label="Private Site">
          <Toggle disabled={!Fider.session.user.isAdministrator} active={this.state.isPrivate} onToggle={this.toggle} />
          <p className="text-muted mt-1">
            A private site prevents unauthenticated users from viewing or interacting with its content. <br /> When enabled, only already registered users,
            invited users and users from trusted OAuth providers will have access to this site.
          </p>
        </Field>
      </Form>
    )
  }
}
