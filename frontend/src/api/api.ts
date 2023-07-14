import axios from 'axios';
import { server } from '../common/constants';

export async function getMoneyRecords() {
  try {
    const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2ODkyMzQ4OTl9._YoGR4HiR4-U-Kl63APp25Y9l1Hc6Ej5F_EeZK6EOoU';
    const headers = {
      Authorization: `Bearer ${token}`
    };

    const response = await axios.get(server + '/money_records', { headers });
    // DO:これ以降のconsole.log消す
    console.log(response.data)
    return response.data
  } catch (error) {
    console.error('GETリクエストが失敗しました:', error);
  }
}

export interface POSTRegisterRequest {
  user1_name: string;
  user2_name: string;
  password: string;
}

export async function registerUser(data: POSTRegisterRequest) {
  try {
    const dat = JSON.stringify(data);
    console.log(dat)
    const headers = { 'Content-Type': 'application/json' 
    };

    const response = await axios.post(server + '/register',dat ,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}


export interface POSTLoginRequest {
  pair_id: number;
  password: string;
}
export async function loginUser(data: POSTLoginRequest) {
  try {
    const dat = JSON.stringify(data);
    console.log(dat)
    const headers = { 'Content-Type': 'application/json' 
    };

    const response = await axios.post(server + '/login',dat ,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}

export interface POSTAddIncomeRecordRequest {
  user_id: number;
  amonut: number;
}
export interface POSTAddPairExpenseRecordRequest {
  user_id: number;
  amonut: number;
}
export interface POSTAddIndivisualExpenseRecordRequest {
  user_id: number;
  amonut: number;
}

export async function addIncomeRecord(data: POSTAddIncomeRecordRequest) {
  try {
    const dat = JSON.stringify(data);
    console.log(dat)
    const headers = { 'Content-Type': 'application/json' 
    };

    const response = await axios.post(server + '/record_income',dat ,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}

export async function addPairExpenseRecord(data: POSTAddPairExpenseRecordRequest) {
  try {
    const dat = JSON.stringify(data);
    console.log(dat)
    const headers = { 'Content-Type': 'application/json' 
    };

    const response = await axios.post(server + '/record_pair_expense',dat ,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}

export async function addIndivisualExpenseRecord(data: POSTAddIndivisualExpenseRecordRequest) {
  try {
    const dat = JSON.stringify(data);
    console.log(dat)
    const headers = { 'Content-Type': 'application/json' 
    };

    const response = await axios.post(server + '/record_indivisual_expense',dat ,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}