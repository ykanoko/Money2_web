import axios from 'axios';
import { server } from '../common/constants';


export interface POSTRegisterRequest {
  user1_name: string;
  user2_name: string;
  password: string;
}

export interface AddIncomeRecordRequest {
  user_id: number;
  amount: number;
}
export interface AddPairExpenseRecordRequest {
  user_id: number;
  amount: number;
}
export interface AddIndivisualExpenseRecordRequest {
  user_id: number;
  amount: number;
}

export interface GetMoneyRecordsListResponse {
  money2_id: number;
	date: string;
	type: string;
	user: string;
	amount: number;
}

export type GetPairStatusReponse = {
  balance_user1: number   
	balance_user2 :number   
	pay_user     : string  
	pay_amount   : number 
}

export type GetMoneyRecord = {
  pair_status:GetPairStatusReponse
  money_records: GetMoneyRecordsListResponse[]
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

export async function addIncomeRecord(data: AddIncomeRecordRequest, token:string) {
  try {
    const dat = JSON.stringify(data);
    console.log("addIncomeRecord",dat)
    const headers = {
      'Content-Type': 'application/json' , 
      Authorization: `Bearer ${token}`
    };

    const response = await axios.post(server + '/record_income',dat ,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}

export async function addPairExpenseRecord(data: AddPairExpenseRecordRequest, token:string) {
  try {
    const dat = JSON.stringify(data);

    const headers = { 'Content-Type': 'application/json' , 
    Authorization: `Bearer ${token}`
    };

    const response = await axios.post(server + '/record_pair_expense', dat,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}

export async function addIndivisualExpenseRecord(data: AddIndivisualExpenseRecordRequest, token: string) {
  try {
    const dat = JSON.stringify(data);

    const headers = {
       'Content-Type': 'application/json', 
       Authorization: `Bearer ${token}`
      };

    const response = await axios.post(server + '/record_indivisual_expense',dat ,{headers} );
    console.log(response);
  return response
  } catch (error) {
    console.error('POSTリクエストが失敗しました:', error);
  }
}

export async function getMoneyRecord(token: string) {
    try {
      const headers = {
        Authorization: `Bearer ${token}`
      };
  
      const [pair_status, money_records] = await Promise.all([axios.get(server + '/pair_status', { headers }), axios.get(server + '/money_records', { headers })]) ;
      // DO:これ以降のconsole.log消す
      console.log(money_records.data, pair_status.data)
      const data:GetMoneyRecord = {
        pair_status:pair_status.data,
        money_records:money_records.data   
      }
      return data
    } catch (error) {
      console.error('GETリクエストが失敗しました:', error);
    }
  }