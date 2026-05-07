/**
 * SMS Order Query API
 * Public endpoint - no authentication required
 */

import { apiClient } from './client'

export interface SmsOrderResponse {
  order_no: string
  phone_number: string
  sms_content: string
  status: string
  created_at: string
  updated_at: string
}

/**
 * Query SMS order by order number
 */
export async function queryByOrderNo(orderNo: string): Promise<SmsOrderResponse> {
  const response = await apiClient.get<SmsOrderResponse>(`/sms-orders/${encodeURIComponent(orderNo)}`)
  return response.data
}
