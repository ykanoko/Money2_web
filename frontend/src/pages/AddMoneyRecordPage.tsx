import {
  Anchor,
  Button,
  Group,
  NumberInput,
  SegmentedControl,
  Stack,
  Table,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import {
  AddIncomeRecordRequest,
  GetMoneyRecord,
  addIncomeRecord,
  addPairExpenseRecord,
  getMoneyRecord,
} from "../api/api";
import { useLocalStorage } from "@mantine/hooks";
import { User } from "../types/user";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { notifications } from "@mantine/notifications";

type Form = AddIncomeRecordRequest & {
  payType: "1" | "2";
};
export default function AddMoneyRecordPage() {
  const [user1] = useLocalStorage<User | null>({
    key: "user1",
    defaultValue: null,
  });
  const [user2] = useLocalStorage<User | null>({
    key: "user2",
    defaultValue: null,
  });
  const [token] = useLocalStorage({ key: "token", defaultValue: "" });
  const [moneyRecord, setMoneyRecord] = useState<GetMoneyRecord>();

  const navigate = useNavigate();
  const form = useForm<Form>({
    initialValues: {
      payType: "2",
      user_id: 0,
      amount: 0,
    },
  });
  useEffect(() => {
    if (token === "") {
      return;
    }
    getMoneyRecord(token).then((res) => {
      setMoneyRecord(res);
    }).catch(err=>{
      console.log(err)
      navigate("/login")
    });
  }, [token,navigate]);

  // DO:loginできてなかったらログインページに移動させる
  // DO:amountが0でerror出す
  function handleSubmit(values: Form) {
    values.user_id = Number(values.user_id);
    if (values.payType === "1") {
      addIncomeRecord({ amount: values.amount, user_id: values.user_id }, token)
        .then(() => {
          navigate(0);
        })
        .catch((err) => {
          console.log(err);
          notifications.show({
            title: "Error!",
            message: "failed to add money record",
            color: "red",
          });
        });
    } else if (values.payType === "2") {
      addPairExpenseRecord(
        { amount: values.amount, user_id: values.user_id },
        token,
      )
        .then(() => {
          navigate(0);
        })
        .catch((err) => {
          console.log(err);
          notifications.show({
            title: "Error!",
            message: "failed to add money record",
            color: "red",
          });
        });
    }
  }

  const pair_status_row = (
    <tr>
      <td>{moneyRecord?.pair_status.balance_user1}</td>
      <td>{moneyRecord?.pair_status.balance_user2}</td>
      <td>{moneyRecord?.pair_status.pay_user}</td>
      <td>{moneyRecord?.pair_status.pay_amount}</td>
    </tr>
  );
  const rows = moneyRecord?.money_records.slice(0, 10).map((money_record) => (
    <tr key={money_record.money2_id}>
      <td>{money_record.money2_id}</td>
      <td>{money_record.date}</td>
      <td>{money_record.type}</td>
      <td>{money_record.user}</td>
      <td>{money_record.amount}</td>
    </tr>
  ));

  if (!user1 || !user2) {
    return (
      <Stack maw={300} mx="auto" mt={"md"}>
        You are not logged in
        <Anchor href="/register">Register</Anchor>
      </Stack>
    );
  }
  return (
    <Stack maw={600} mx="auto">
      <Title order={2} mt={"md"}>
        Add Money Record
      </Title>
      <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
        <Stack maw={300} mx="auto">
          <SegmentedControl
            color="blue"
            data={[
              { label: "収入", value: "1" },
              { label: "合計支出", value: "2" },
              // { label: '個人支出', value: "3" },
            ]}
            {...form.getInputProps("payType")}
          />
          <SegmentedControl
            color="blue"
            data={[
              { label: user1?.name, value: user1?.id.toString() },
              { label: user2?.name, value: user2?.id.toString() },
            ]}
            {...form.getInputProps("user_id")}
          />
          <NumberInput
            withAsterisk
            label="Amount"
            {...form.getInputProps("amount")}
          />
          <Group position="right" mt="md">
            <Button type="submit">Submit</Button>
          </Group>
          <Anchor href="/money_record">Money Record</Anchor>
          <Anchor href="/login">Login</Anchor>
        </Stack>
      </form>
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
    </Stack>
  );
}
