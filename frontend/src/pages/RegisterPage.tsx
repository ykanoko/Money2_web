import { Anchor, Box, Button, Group, PasswordInput, TextInput, Title } from "@mantine/core";
import { useForm } from "@mantine/form";
import { POSTRegisterRequest, registerUser } from "../api/api";
import { useLocalStorage } from "@mantine/hooks";
import { useNavigate } from "react-router-dom";

export default function RegisterPage() {
  const [value, setValue] = useLocalStorage({ key: 'pair_id', defaultValue: '' });
  const navigate = useNavigate()
  const form = useForm<POSTRegisterRequest>({
    initialValues:{
      user1_name:'',
      user2_name:'',
      password:''
    }

  });
  function handleSubmit(values: POSTRegisterRequest){
registerUser(values).then(res => {
  setValue(res?.data.pair_id)
  navigate('/login')
})
  }
// DO:Submitボタンを押したら、その後一定時間押せないようにする？
// TODO:メアドに紐づけたい？、googleアカウント等？

// DO:sign outボタン
  return (<>
    <Box maw={300} mx="auto">
    <Title order={2}>Register</Title>
    <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
      <TextInput
        withAsterisk
        label="User1 name"
        {...form.getInputProps('user1_name')}
      />
      <TextInput
        withAsterisk
        label="User2 name"
        {...form.getInputProps('user2_name')}
      />
      <PasswordInput
        withAsterisk
        label="Password"
        {...form.getInputProps('password')}
      />

      <Group position="right" mt="md">
        <Button type="submit">Submit</Button>
      </Group>
      <Anchor href="/login">
      Login
    </Anchor>

    </form>
  </Box>
  </>);
}
