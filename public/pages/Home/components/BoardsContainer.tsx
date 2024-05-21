import "./BoardsContainer.scss"

import React from "react"

import { Board, CurrentUser } from "@fider/models"
import { Loader } from "@fider/components"
import { actions, navigator, querystring } from "@fider/services"
// import IconSearch from "@fider/assets/images/heroicons-search.svg"
// import IconX from "@fider/assets/images/heroicons-x.svg"
import { ListBoards } from "./ListBoards"
import { t, Trans } from "@lingui/macro"

interface PostsContainerProps {
  user?: CurrentUser
  boards: Board[]
}

interface BoardsContainerState {
  loading: boolean
  boards?: Board[]
  view: string
  tags: string[]
  query: string
  limit?: number
}

export class BoardsContainer extends React.Component<PostsContainerProps, BoardsContainerState> {
  constructor(props: PostsContainerProps) {
    super(props)

    this.state = {
      boards: this.props.boards,
      loading: false,
      view: querystring.get("view"),
      query: querystring.get("query"),
      tags: querystring.getArray("tags"),
      limit: querystring.getNumber("limit"),
    }
  }

  private changeFilterCriteria<K extends keyof BoardsContainerState>(obj: Pick<BoardsContainerState, K>, reset: boolean): void {
    this.setState(obj, () => {
      const query = this.state.query.trim().toLowerCase()
      navigator.replaceState(
        querystring.stringify({
          tags: this.state.tags,
          query,
          view: this.state.view,
          limit: this.state.limit,
        })
      )

      this.searchPosts(query, this.state.view, this.state.limit, this.state.tags, reset)
    })
  }

  private timer?: number
  private async searchPosts(query: string, view: string, limit: number | undefined, tags: string[], reset: boolean) {
    window.clearTimeout(this.timer)
    this.setState({ boards: reset ? undefined : this.state.boards, loading: true })
    this.timer = window.setTimeout(() => {
      actions.searchBoards({ query, view, limit, tags }).then((response) => {
        if (response.ok && this.state.loading) {
          this.setState({ loading: false, boards: response.data })
        }
      })
    }, 500)
  }

  // private handleSearchFilterChanged = (query: string) => {
    // this.changeFilterCriteria({ query }, true)
  // }

  // private clearSearch = () => {
    // this.changeFilterCriteria({ query: "" }, true)
  // }

  private showMore = (event: React.MouseEvent<HTMLElement> | React.TouchEvent<HTMLElement>): void => {
    event.preventDefault()
    this.changeFilterCriteria({ limit: (this.state.limit || 30) + 10 }, false)
  }

  private getShowMoreLink = (): string | undefined => {
    if (this.state.boards && this.state.boards.length >= (this.state.limit || 30)) {
      return querystring.set("limit", (this.state.limit || 30) + 10)
    }
  }

  public render() {
    const showMoreLink = this.getShowMoreLink()

    return (
      <div className="c-posts-container">
        <div className="c-posts-container__header mb-4">
          {!this.state.query && (
            <div className="c-posts-container__filter-col">
            </div>
          )}
          <div className="c-posts-container__search-col">
            {/* <Input
              field="query"
              icon={this.state.query ? IconX : IconSearch}
              onIconClick={this.state.query ? this.clearSearch : undefined}
              placeholder={t({ id: "home.postscontainer.query.placeholder", message: "Search" })}
              value={this.state.query}
              onChange={this.handleSearchFilterChanged}
            /> */}
          </div>
        </div>
        <ListBoards
          boards={this.state.boards}
          emptyText={t({ id: "home.boardcontainer.label.noresults", message: "You are not a member of any boards yet. After you create a board or make a contribution to an existing board, it will be available here." })}
        />
        {this.state.loading && <Loader />}
        {showMoreLink && (
          <div className="my-4 ml-4">
            <a href={showMoreLink} className="text-primary-base text-medium hover:underline" onClick={this.showMore}>
              <Trans id="home.postscontainer.label.viewmore">View more posts</Trans>
            </a>
          </div>
        )}
      </div>
    )
  }
}
