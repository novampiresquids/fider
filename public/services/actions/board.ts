import { http, Result, querystring } from "@fider/services"
import { Board } from "@fider/models"

export const getAllBoards = async (): Promise<Result<Board[]>> => {
  return await http.get<Board[]>("/api/v1/boards")
}

export interface SearchBoardsParams {
  query?: string
  view?: string
  limit?: number
  tags?: string[]
}

export const searchBoards = async (params: SearchBoardsParams): Promise<Result<Board[]>> => {
  return await http.get<Board[]>(
    `/api/v1/board${querystring.stringify({
      tags: params.tags,
      query: params.query,
      view: params.view,
      limit: params.limit,
    })}`
  )
}

export const deleteBoard = async (postNumber: number, text: string): Promise<Result> => {
  return http
    .delete(`/api/v1/board/${postNumber}`, {
      text,
    })
    .then(http.event("post", "delete"))
}

interface CreateBoardResponse {
  id: number
  number: number
  title: string
  slug: string
}

export const createBoard = async (title: string, description: string): Promise<Result<CreateBoardResponse>> => {
  return http.post<CreateBoardResponse>(`/api/v1/board`, { tenantName: title, welcome: description }).then(http.event("post", "create"))
}

export const updateBoard = async (postNumber: number, title: string, description: string): Promise<Result> => {
  return http.put(`/api/v1/board/${postNumber}`, { title, description }).then(http.event("post", "update"))
}
