import {
  Anchor,
  Box,
  Button,
  Group,
  NumberInput,
  PasswordInput,
  Stack,
  Title,
} from "@mantine/core";
import { POSTLoginRequest, loginUser } from "../api/api";
import { useForm } from "@mantine/form";
import { useLocalStorage } from "@mantine/hooks";
import { useNavigate } from "react-router-dom";
import { User } from "../types/user";

// rafceを打つ
export default function LoginPage() {
  const [pairID] = useLocalStorage({ key: "pair_id", defaultValue: "" });
  const [, setToken] = useLocalStorage({ key: "token", defaultValue: "" });
  const [, setUser1] = useLocalStorage<User | null>({
    key: "user1",
    defaultValue: null,
  });
  const [, setUser2] = useLocalStorage<User | null>({
    key: "user2",
    defaultValue: null,
  });
  // TODO:現在は1組のみの登録を想定しているため、storageに入れた
  const navigate = useNavigate();
  const form = useForm<POSTLoginRequest>({
    initialValues: {
      pair_id: 0,
      password: "",
    },
  });

  function handleSubmit(values: POSTLoginRequest) {
    loginUser(values).then((res) => {
      setToken(res?.data.token);
      setUser1(res?.data.user1);
      setUser2(res?.data.user2);
      navigate("/");
    });
  }
  // DO:validation 項目埋まっていない時は、submitできないようにできていない（通知付き）

  return (
    <>
      <Box maw={600} mx="auto">
        <Title order={2} mt={"md"}>
          Login
        </Title>
        <form onSubmit={form.onSubmit((values) => handleSubmit(values))}>
          <Stack maw={300} mx="auto">
            {pairID ? <p>Your pair ID is {pairID}</p> : null}
            <NumberInput
              withAsterisk
              label="Pair ID"
              {...form.getInputProps("pair_id")}
            />
            <PasswordInput
              withAsterisk
              label="Password"
              {...form.getInputProps("password")}
            />

            <Group position="right" mt="md">
              <Button type="submit">Submit</Button>
            </Group>
            <Anchor href="/register">Register</Anchor>
          </Stack>
        </form>
      </Box>
    </>
  );
}
