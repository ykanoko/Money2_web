import {
  Anchor,
  Button,
  Group,
  PasswordInput,
  Stack,
  TextInput,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { POSTRegisterRequest, registerUser } from "../api/api";
import { useLocalStorage } from "@mantine/hooks";
import { useNavigate } from "react-router-dom";

export default function RegisterPage() {
  const [, setPairID] = useLocalStorage({ key: "pair_id", defaultValue: "" });
  // TODO:現在は1組のみの登録を想定しているため、storageに入れた
  const navigate = useNavigate();
  const form = useForm<POSTRegisterRequest>({
    initialValues: {
      user1_name: "",
      user2_name: "",
      password: "",
    },
  });
  function handleSubmit(values: POSTRegisterRequest) {
    registerUser(values).then((res) => {
      setPairID(res?.data.pair_id);
      navigate("/login");
    });
  }
  // DO:Submitボタンを押したら、その後一定時間押せないようにする？
  // TODO:メアドに紐づけたい？、googleアカウント等？

  // DO:sign outボタン
  // DO:validation 項目埋まっていない時は、submitできないようにできていない（通知付き）
  return (
    <>
      <Stack maw={600} mx="auto">
        <Title order={2} mt={"md"}>
          Register
        </Title>
        <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
          <Stack maw={300} mx="auto">
            <TextInput
              withAsterisk
              label="User1 name"
              {...form.getInputProps("user1_name")}
            />
            <TextInput
              withAsterisk
              label="User2 name"
              {...form.getInputProps("user2_name")}
            />
            <PasswordInput
              withAsterisk
              label="Password"
              {...form.getInputProps("password")}
            />

            <Group position="right" mt="md">
              <Button type="submit">Submit</Button>
            </Group>
            <Anchor href="/login">Login</Anchor>
          </Stack>
        </form>
      </Stack>
    </>
  );
}
