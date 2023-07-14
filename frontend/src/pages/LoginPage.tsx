import { Anchor, Box, Button, Group, NumberInput, PasswordInput, Title } from "@mantine/core";
import { POSTLoginRequest, loginUser } from "../api/api";
import { useForm } from "@mantine/form";
import { useLocalStorage } from "@mantine/hooks";
import { useNavigate } from "react-router-dom";

// rafceを打つ
export default function LoginPage() {
  const [value, setValue] = useLocalStorage({ key: 'pair_id', defaultValue: '' });
  const [token, setToken] = useLocalStorage({ key: 'token', defaultValue: '' });
  const navigate = useNavigate()
  const form = useForm<POSTLoginRequest>({
    initialValues:{
      pair_id: 0,
      password:''
    }
  });

  function handleSubmit(values: POSTLoginRequest){
loginUser(values).then(res => {
  setToken(res?.data.token)
  navigate('/')
})
  }

  return (<>
    <Box maw={300} mx="auto">
    <Title order={2}>Login</Title>
    <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
      {value !== ''? <p>Your pair ID is {value}</p>:null}
      <NumberInput
        withAsterisk
        label="Pair ID"
        {...form.getInputProps('pair_id')}
      />
      <PasswordInput
        withAsterisk
        label="Password"
        {...form.getInputProps('password')}
      />

      <Group position="right" mt="md">
        <Button type="submit">Submit</Button>
      </Group>
      <Anchor href="/register">
      Register
    </Anchor>
    </form>
  </Box>
  </>);
}
