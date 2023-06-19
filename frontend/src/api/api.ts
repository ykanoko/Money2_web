import axios from 'axios';

export async function getMoneyRecords() {
  try {
    const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2ODc0MDExOTd9.nWg_CMWBoPcT8NP7C5VAFPAU8drtW11yPURGfv1Dzi0';
    const headers = {
      Authorization: `Bearer ${token}`
    };

    const response = await axios.get('http://127.0.0.1:9000/money_records', { headers });
    return response.data
  } catch (error) {
    console.error('GETリクエストが失敗しました:', error);
  }
}