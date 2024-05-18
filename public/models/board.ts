import { User } from "./identity"

export interface Board {
  id: number
  name: string
  status: number
  isPrivate: boolean
  invitation: string
  welcomeMessage: string
}

export interface BoardResponse {
  user: User
  text: string
  respondedAt: Date
  original?: {
    number: number
    title: string
    slug: string
    status: string
  }
}
