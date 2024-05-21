import "./Home.page.scss"

import React, { useState } from "react"
import { Board } from "@fider/models"
import { Markdown, PoweredByFider, Header } from "@fider/components"
import { BoardsContainer } from "./components/BoardsContainer"
import { useFider } from "@fider/hooks"
import { VStack } from "@fider/components/layout"

import { t } from "@lingui/macro"
import { BoardInput } from "./components/BoardInput"

export interface HomePageProps {
  boards: Board[]
}

export interface HomePageState {
  title: string
}

const HomePage = (props: HomePageProps) => {
  const fider = useFider()
  const [, setTitle] = useState("")

  const defaultWelcomeMessage = t({
    id: "home.board.defaultwelcomemessage",
    message: `Create a new feedback board`,
  })

  // const defaultInvitation = t({
    // id: "home.form.defaultinvitation",
    // message: "Enter your suggestion here...",
  // })

  return (
    <>
      <Header />
      <div id="p-home" className="page container">
        <div className="p-home__welcome-col">
          <VStack spacing={2}>
            <Markdown text={fider.session.tenant?.welcomeMessage || defaultWelcomeMessage} style="full" />
            <BoardInput placeholder="Board title" onTitleChanged={setTitle} />
            <PoweredByFider slot="home-input" className="sm:hidden md:hidden lg:block" />
          </VStack>
        </div>
        <div className="p-home__posts-col">
            {/* <SimilarPosts title={title} tags={props.tags} /> */}
            <BoardsContainer boards={props.boards} />
          <PoweredByFider slot="home-footer" className="lg:hidden xl:hidden mt-8" />
        </div>
      </div>
    </>
  )
}

export default HomePage
