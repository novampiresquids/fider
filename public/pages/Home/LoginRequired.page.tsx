import "./Home.page.scss"

import React from "react"
import { PoweredByFider, Header, TenantLogo } from "@fider/components"


export interface HomePageState {
  title: string
}

const LoginRequired = () => {

  // const defaultInvitation = t({
    // id: "home.form.defaultinvitation",
    // message: "Enter your suggestion here...",
  // })

  return (
    <>
      <Header />
      <div id="" className="page container">
        <div className="w-max-7xl mx-auto">
          <div className="h-20 text-center mb-4 mt-8"> <TenantLogo size={100} useFiderIfEmpty={true} /> </div>
          <div className="text-center"> Sign in to create your own feedback boards </div>
          <PoweredByFider slot="home-footer" className="lg:hidden xl:hidden mt-8" />
        </div>
      </div>
    </>
  )
}

export default LoginRequired
