/**
 * SMS Order Query API
 * Public endpoint - no authentication required
 */

import { apiClient } from './client'

export interface SmsOrderResponse {
  order_no: string
  service_type: string
  phone_number: string
  sms_content: string
  status: string
  pending_at?: string
  created_at: string
  updated_at: string
}

/**
 * Query SMS order by order number (read-only).
 */
export async function queryByOrderNo(orderNo: string): Promise<SmsOrderResponse> {
  const response = await apiClient.get<SmsOrderResponse>(
    `/sms-orders/${encodeURIComponent(orderNo)}`
  )
  return response.data
}

/**
 * Trigger state-machine refresh on the server:
 *  - status=created → assigns a phone number via hero-sms (with retry).
 *  - status=pending → polls hero-sms for the verification code.
 *  - terminal states are returned as-is.
 *
 * Long-running: the server may retry up to 6 times when fetching a phone number,
 * so the request can take up to ~30 seconds in the worst case.
 */
export async function refreshOrder(orderNo: string): Promise<SmsOrderResponse> {
  const response = await apiClient.post<SmsOrderResponse>(
    `/sms-orders/${encodeURIComponent(orderNo)}/refresh`
  )
  return response.data
}
