import axios from 'axios';
import { server } from '../common/constants';

export async function getMoneyRecords() {
  try {
    const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2ODg4MTA5ODN9.72VBvKD3uVu3RNpGfAHr94pW-T2EHCToWWbYM2a5Uq4';
    const headers = {
      Authorization: `Bearer ${token}`
    };

    const response = await axios.get(server + '/money_records', { headers });
    return response.data
  } catch (error) {
    console.error('GETリクエストが失敗しました:', error);
  }
}