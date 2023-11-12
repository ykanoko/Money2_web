import {
  Anchor,
  Button,
  Card,
  Group,
  NumberInput,
  SegmentedControl,
  Stack,
  Table,
  Text,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import {
  AddIncomeRecordRequest,
  GetMoneyRecord,
  addIncomeRecord,
  addExpenseRecord,
  addPairExpenseRecord,
  getMoneyRecord,
} from "../api/api";
import { useLocalStorage } from "@mantine/hooks";
import { User } from "../types/user";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { notifications } from "@mantine/notifications";

// DO:AddIncomeRecordRequestだけ用いるのでよいのか
type Form = AddIncomeRecordRequest & {
  payType: "1" | "2" | "3";
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
      payType: "3",
      user_id: 0,
      amount: 0,
    },
  });
  useEffect(() => {
    if (token === "") {
      return;
    }
    getMoneyRecord(token)
      .then((res) => {
        setMoneyRecord(res);
      })
      .catch((err) => {
        console.log(err);
        navigate("/login");
      });
  }, [token, navigate]);

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
      addExpenseRecord(
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
    } else if (values.payType === "3") {
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
        <Anchor href="/login">Login</Anchor>
        <Anchor href="/register">Register</Anchor>
      </Stack>
    );
  }
  return (
    <Stack maw={600} mx="auto" pb={80} px={8}>
      <Title order={2} mt={"md"}>
        Add Money Record
      </Title>
      <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
        <Stack maw={450} mx="auto">
          <SegmentedControl
            color="blue"
            data={[
              { label: "収入", value: "1" },
              { label: "支出", value: "2" },
              { label: "合計支出", value: "3" },
            ]}
            {...form.getInputProps("payType")}
          />
          <SegmentedControl
            color="blue"
            data={[
              { label: user1?.name, value: user1?.id?.toString() },
              { label: user2?.name, value: user2?.id?.toString() },
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
        Balance
      </Title>
      <Group position="center" grow>
        <Card shadow="xs" withBorder>
          <Text size="sm" mt="xs" c="dimmed">
            <Text size="sm" mt="xs" c="dark" span fw={600}>
              {user1?.name}
            </Text>
            の残金
          </Text>
          <Text fw={500} size="lg" mt="md" align="right">
            {moneyRecord?.pair_status.balance_user1} 円
          </Text>
        </Card>
        <Card shadow="xs" withBorder>
          <Text size="sm" mt="xs" c="dimmed">
            <Text size="sm" mt="xs" c="dark" span fw={600}>
              {user2?.name}
            </Text>
            の残金
          </Text>
          <Text fw={500} size="lg" mt="md" align="right">
            {moneyRecord?.pair_status.balance_user2} 円
          </Text>
        </Card>
      </Group>

      <Title order={4} mt={"md"}>
        Pair Status
      </Title>
      <Group position="center">
        <Card shadow="xs" withBorder>
          <Text size="sm" mt="xs" c="dimmed">
            精算方法
          </Text>
          <Text size="sm" mt="xs" c="dimmed">
            <Text fw={500} c="dark" size="lg" mt="md" span>
              {moneyRecord?.pair_status.pay_user}
            </Text>
            が
            <Text fw={500} c="dark" size="lg" mt="md" span>
              {moneyRecord?.pair_status.pay_amount}
            </Text>{" "}
            円払う
          </Text>
        </Card>
      </Group>

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
