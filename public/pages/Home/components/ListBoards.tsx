import React from "react"
import { Board, CurrentUser } from "@fider/models"
import { Markdown } from "@fider/components"
import { HStack, VStack } from "@fider/components/layout"

interface ListBoardsProps {
  boards?: Board[]
  emptyText: string
}

const ListBoardItem = (props: { board: Board; user?: CurrentUser; }) => {
  return (
    <HStack center={true}>
      <VStack className="w-full" spacing={2}>
        <HStack justify="between">
          <a className="text-title hover:text-primary-base" href={`/board/${props.board.id}`}>
            {props.board.name}
          </a>
        </HStack>
        <Markdown className="text-gray-600" maxLength={300} text={props.board.welcomeMessage} style="plainText" />
      </VStack>
    </HStack>
  )
}

export const ListBoards = (props: ListBoardsProps) => {
  if (!props.boards) {
    return null
  }

  if (props.boards.length === 0) {
    return <p className="text-center">{props.emptyText}</p>
  }

  return (
    <VStack spacing={4} divide={true} center={true}>
      {props.boards.map((board) => (
        <ListBoardItem key={board.id} board={board} />
      ))}
    </VStack>
  )
}
