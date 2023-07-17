import { Anchor, Stack, Table, Title } from "@mantine/core";
import { GetMoneyRecord, getMoneyRecord } from "../api/api";
import { useEffect, useState } from "react";
import { useLocalStorage } from "@mantine/hooks";
import { User } from "../types/user";

export default function MoneyRecordPage() {
  const [moneyRecord, setMoneyRecord] = useState<GetMoneyRecord>();
  const [token] = useLocalStorage({ key: "token", defaultValue: "" });
  const [user1] = useLocalStorage<User | null>({
    key: "user1",
    defaultValue: null,
  });
  const [user2] = useLocalStorage<User | null>({
    key: "user2",
    defaultValue: null,
  });

  useEffect(() => {
    if (token === "") {
      return;
    }
    getMoneyRecord(token).then((res) => {
      setMoneyRecord(res);
    });
  }, [token]);

  const pair_status_row = (
    <tr>
      <td>{moneyRecord?.pair_status.balance_user1}</td>
      <td>{moneyRecord?.pair_status.balance_user2}</td>
      <td>{moneyRecord?.pair_status.pay_user}</td>
      <td>{moneyRecord?.pair_status.pay_amount}</td>
    </tr>
  );

  const rows = moneyRecord?.money_records.map((money_record) => (
    <tr key={money_record.money2_id}>
      <td>{money_record.money2_id}</td>
      <td>{money_record.date}</td>
      <td>{money_record.type}</td>
      <td>{money_record.user}</td>
      <td>{money_record.amount}</td>
    </tr>
  ));

  return (
    <Stack maw={600} mx="auto">
      <Title order={2} mt={"md"}>
        MoneyRecord
      </Title>
      <Title order={4} mt={"md"}>
        Pair Status
      </Title>
      <Table striped withBorder>
        <thead>
          <tr>
            <th>{user1?.name}の残金</th>
            <th>{user2?.name}の残金</th>
            <th>払うべき人</th>
            <th>払うべき金額</th>
          </tr>
        </thead>
        <tbody>{pair_status_row}</tbody>
      </Table>
      <Title order={4} mt={"md"}>
        Money Records
      </Title>
      <Table striped withBorder>
        <thead>
          <tr>
            <th>ID</th>
            <th>日付</th>
            <th>種類</th>
            <th>名前</th>
            <th>金額</th>
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
      <Anchor href="/">Add Money Record</Anchor>
    </Stack>
  );
}
