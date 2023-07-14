# Money2_web

家計簿 Web アプリ：複数人での金銭の管理を容易に行うことができる。

## Project Structure

- Backend(GO):

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

  <!-- ![hero](images/ER%E5%9B%B3.png) -->

  - vanilla css -> lives in styles/

- Frontend(React):
  - [axios](https://github.com/axios/axios) + [React Query](https://tanstack.com/query/v4/?from=reactQueryV3&original=https://react-query-v3.tanstack.com/)
    -> see lib/

## Getting Started

### Environment Variables

API : [RESAS API](https://opendata.resas-portal.go.jp/)

```
NEXT_PUBLIC_API_KEY={FILL_ME_IN}
```

### Run the development server:

```bash
npm install
npm run dev
```

### Component Structure

```sh
├── components/
│ ├── common/  # common components of the app
│ ├── feature/ # components scoped to a specific feature
│ └── layout/ # layout components, header
└── config/
└── hooks/
└── lib/
└── pages/
└── public/
└── stores/
└── styles/
└── types/
```

## Credits

[RESAS API](https://opendata.resas-portal.go.jp/)
