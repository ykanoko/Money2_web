import axios from 'axios';
import { server } from '../common/constants';

export async function getMoneyRecords() {
  try {
    const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2ODkyMzQ4OTl9._YoGR4HiR4-U-Kl63APp25Y9l1Hc6Ej5F_EeZK6EOoU';
    const headers = {
      Authorization: `Bearer ${token}`
    };

    const response = await axios.get(server + '/money_records', { headers });
    console.log(response.data)
    return response.data
  } catch (error) {
    console.error('GETリクエストが失敗しました:', error);
  }
}