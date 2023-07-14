# Money2_web

家計簿 Web アプリ：複数人での金銭の管理を容易に行うことができる。

## Project Structure

- Backend(GO)

  ```mermaid
  erDiagram

  types ||--o{ money2 : ""
  users ||--o{ money2 : ""
  pairs ||--o{ money2 : ""
  users ||--|| pairs: ""

  users {
  	int id
    varchar name
    decimal balance
  }

  pairs {
  	int id
  	binary password
  	int user1_id
  	int user2_id
  	decimal calculation_user1
  	datetime created_at
  }

  money2 {
    int id
  	int pair_id
    int type_id
  	int user_id
  	int amount
  	datetime created_at
  }

  types {
    int id
  	varchar name
  }
  ```

- Frontend(React)
