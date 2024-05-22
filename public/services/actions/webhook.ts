import { http, Result, StringObject } from "@fider/services"
import { WebhookData, WebhookPreviewResult, WebhookTriggerResult, WebhookType } from "@fider/models"

export const createWebhook = async (data: WebhookData, boardId: number): Promise<Result<{ id: number }>> => {
  return await http.post(`/fider/_api/board/${boardId}/admin/webhook`, data)
}

export const updateWebhook = async (id: number, data: WebhookData, boardId: number): Promise<Result> => {
  return await http.put(`/fider/_api/board/${boardId}/admin/webhook/${id}`, data)
}

export const deleteWebhook = async (id: number, boardId: number): Promise<Result> => {
  return await http.delete(`/fider/fider/_api/board/${boardId}/admin/webhook/${id}`)
}

export const testWebhook = async (id: number, boardId: number): Promise<Result<WebhookTriggerResult>> => {
  return await http.get(`/_api/board/${boardId}/admin/webhook/test/${id}`)
}

export const previewWebhook = async (type: WebhookType, url: string, content: string, boardId: number): Promise<Result<WebhookPreviewResult>> => {
  return await http.post(`/fider/_api/board/${boardId}/admin/webhook/preview`, { type, url, content })
}

export const getWebhookHelp = async (type: WebhookType, boardId: number): Promise<Result<StringObject>> => {
  return await http.get(`/fider/_api/board/${boardId}/admin/webhook/props/${type}`)
}
