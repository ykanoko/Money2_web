import { Table } from "@mantine/core"
import { getMoneyRecords } from "../api/api"
import {useEffect, useState} from "react"

type MoneyRecord = {
    id: number
    date: string
    type: string
	user         :string 
	amount       :number   
	balance_user1: number   
	balance_user2 :number   
	pay_user     : string  
	pay_amount   : number 

}
const MoneyRecordsPage = () => {
    // usestaと打つ
    const [moneyRecords, setMoneyRecords] = useState<MoneyRecord[]>([])
useEffect( () => {
    getMoneyRecords().then(res =>{
        setMoneyRecords(res)
        console.log(res)
    })
}, [])

const rows = moneyRecords.map((moneyRecord) => (
    <tr key={moneyRecord.id}>
      <td>{moneyRecord.id}</td>
      <td>{moneyRecord.date}</td>
      <td>{moneyRecord.type}</td>
      <td>{moneyRecord.user}</td>
      <td>{moneyRecord.amount}</td>
      <td>{moneyRecord.balance_user1}</td>
      <td>{moneyRecord.balance_user2}</td>
      <td>{moneyRecord.pay_user}</td>
      <td>{moneyRecord.pay_amount}</td>
    </tr>
  ));

return (
    <Table>
      <thead>
        <tr>
          <th>ID</th>
          <th>日付</th>
          <th>種類</th>
          <th>名前</th>
          <th>金額</th>
          <th>残金1</th>
          <th>残金2</th>
          <th>払うべき人</th>
          <th>払うべき金額</th>
          {/* <th>Date</th>
          <th>Type</th>
          <th>User</th>
          <th>Amount</th>
          <th>Balance User1</th>
          <th>Balance User2</th>
          <th>Pay User</th>
          <th>Pay Amount</th> */}

        </tr>
      </thead>
      <tbody>{rows}</tbody>
    </Table>
  )
}

export default MoneyRecordsPage