/**
 * Admin SMS Orders API endpoints
 */

import { apiClient } from '../client'

export interface SmsOrder {
  order_no: string
  service_type: string
  phone_number: string
  sms_content: string
  status: string
  pending_at?: string
  created_at: string
}

export interface SmsOrderListResponse {
  items: SmsOrder[]
  total: number
}

export interface BatchGenerateResponse {
  items: SmsOrder[]
}

export async function list(
  page: number = 1,
  pageSize: number = 20,
  status?: string,
  options?: { signal?: AbortSignal }
): Promise<SmsOrderListResponse> {
  const { data } = await apiClient.get<SmsOrderListResponse>('/admin/sms-orders', {
    params: {
      page,
      page_size: pageSize,
      ...(status ? { status } : {})
    },
    signal: options?.signal
  })
  return data
}

export async function batchGenerate(count: number, serviceType: string = 'claude'): Promise<BatchGenerateResponse> {
  const { data } = await apiClient.post<BatchGenerateResponse>('/admin/sms-orders/generate', {
    count,
    service_type: serviceType
  })
  return data
}

const smsOrdersAPI = { list, batchGenerate }
export default smsOrdersAPI
