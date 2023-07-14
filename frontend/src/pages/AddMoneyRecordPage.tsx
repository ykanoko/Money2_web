import { Box, Button, Group, NumberInput, PasswordInput, TextInput, Title } from "@mantine/core";
import { useForm } from "@mantine/form";
import { POSTAddIncomeRecordRequest, POSTRegisterRequest, registerUser } from "../api/api";

type Form = POSTAddIncomeRecordRequest & {
  payType:0 | 1
}
export default function AddMoneyRecordPage() {
  const form = useForm<Form>({
    initialValues: {
      payType:0,
      user_id:0,
      amonut:0
    },

  });
  function handleSubmit(values: Form){
// registerUser(values)
  }

  return (
    <Box maw={300} mx="auto">
    <Title order={2}>Add Money Record</Title>

    <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
      <NumberInput
        withAsterisk
        label="Amount"
        {...form.getInputProps('amount')}
      />
      <Group position="right" mt="md">
        <Button type="submit">Submit</Button>
      </Group>
    </form>
  </Box>
  );
}
