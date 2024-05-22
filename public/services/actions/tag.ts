import { http, Result } from "@fider/services/http"
import { Tag } from "@fider/models"

export const createTag = async (boardId: number, name: string, color: string, isPublic: boolean): Promise<Result<Tag>> => {
  return http.post<Tag>(`/fider/api/v1/board/${boardId}/tags`, { name, color, isPublic }).then(http.event("tag", "create"))
}

export const updateTag = async (boardId: number, slug: string, name: string, color: string, isPublic: boolean): Promise<Result<Tag>> => {
  return http.put<Tag>(`/fider/api/v1/board/${boardId}/tags/${slug}`, { name, color, isPublic }).then(http.event("tag", "update"))
}

export const deleteTag = async (boardId: number, slug: string): Promise<Result> => {
  return http.delete(`/fider/api/v1/board/${boardId}/tags/${slug}`).then(http.event("tag", "delete"))
}

export const assignTag = async (boardId: number, slug: string, postNumber: number): Promise<Result> => {
  return http.post(`/fider/api/v1/board/${boardId}/posts/${postNumber}/tags/${slug}`).then(http.event("tag", "assign"))
}

export const unassignTag = async (boardId: number, slug: string, postNumber: number): Promise<Result> => {
  return http.delete(`/fider/api/v1/board/${boardId}/posts/${postNumber}/tags/${slug}`).then(http.event("tag", "unassign"))
}
