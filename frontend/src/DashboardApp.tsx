import { useCallback, useEffect, useMemo, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import type { Key } from "react";
import type { FormInstance, TableColumnsType } from "antd";
import {
  Alert,
  Button,
  Card,
  Col,
  DatePicker,
  Drawer,
  Form,
  Input,
  InputNumber,
  Layout,
  Menu,
  Popconfirm,
  Row,
  Select,
  Space,
  Spin,
  Statistic,
  Switch,
  Table,
  Tag,
  Typography,
  message,
} from "antd";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import {
  AlertOutlined,
  BookOutlined,
  CreditCardOutlined,
  DatabaseOutlined,
  FileTextOutlined,
  LogoutOutlined,
  RobotOutlined,
  ShoppingOutlined,
  UserOutlined,
} from "@ant-design/icons";
import "./App.css";

const { Header, Sider, Content } = Layout;
const { Title, Text } = Typography;

type CurrentUser = {
  id: number;
  name: string;
  email: string;
  role: string;
};

type DashboardOverview = {
  users: number;
  transactions: number;
  products: number;
  decks: number;
  cards: number;
  ai_generation_history: number;
  issues: number;
};

type EntityKey =
  | "users"
  | "transactions"
  | "products"
  | "decks"
  | "cards"
  | "ai_generation_history"
  | "issues";

type RecordData = Record<string, unknown> & { id?: number };

type LookupOption = { value: number; label: string };

type FieldType =
  | "text"
  | "textarea"
  | "number"
  | "boolean"
  | "datetime"
  | "select";

type FieldConfig = {
  key: string;
  label: string;
  type: FieldType;
  required?: boolean;
  options?: Array<{ label: string; value: string | number }>;
  relation?: "users" | "products" | "decks";
};

type EntityConfig = {
  label: string;
  icon: React.ReactNode;
  columns: string[];
  fields: FieldConfig[];
};

const API_URL = import.meta.env.VITE_API_URL ?? "/api";
const TOKEN_KEY = "kilas_admin_token";

const ENTITY_CONFIG: Record<EntityKey, EntityConfig> = {
  users: {
    label: "Users",
    icon: <UserOutlined />,
    columns: [
      "id",
      "email",
      "username",
      "provider",
      "avatar_url",
      "tokens",
      "last_login_date",
      "login_streak",
      "subscription_until",
      "language",
      "role",
      "created_at",
      "updated_at",
    ],
    fields: [
      { key: "email", label: "Email", type: "text", required: true },
      { key: "username", label: "Username", type: "text", required: true },
      { key: "password", label: "Password", type: "text" },
      { key: "provider", label: "Provider", type: "text" },
      { key: "avatar_url", label: "Avatar URL", type: "text" },
      { key: "tokens", label: "Tokens", type: "number" },
      { key: "last_login_date", label: "Last Login Date", type: "datetime" },
      { key: "login_streak", label: "Login Streak", type: "number" },
      {
        key: "subscription_until",
        label: "Subscription Until",
        type: "datetime",
      },
      { key: "language", label: "Language", type: "text" },
      {
        key: "role",
        label: "Role",
        type: "select",
        options: [
          { label: "admin", value: "admin" },
          { label: "user", value: "user" },
        ],
      },
      { key: "created_at", label: "Created At", type: "datetime" },
      { key: "updated_at", label: "Updated At", type: "datetime" },
    ],
  },
  transactions: {
    label: "Transactions",
    icon: <CreditCardOutlined />,
    columns: [
      "id",
      "user_id",
      "product_id",
      "amount",
      "tokens",
      "status",
      "payment_url",
      "created_at",
      "user",
      "product",
    ],
    fields: [
      {
        key: "user_id",
        label: "User",
        type: "select",
        required: true,
        relation: "users",
      },
      {
        key: "product_id",
        label: "Product",
        type: "select",
        required: true,
        relation: "products",
      },
      { key: "amount", label: "Amount", type: "number", required: true },
      { key: "tokens", label: "Tokens", type: "number", required: true },
      { key: "status", label: "Status", type: "text" },
      { key: "payment_url", label: "Payment URL", type: "text" },
      { key: "created_at", label: "Created At", type: "datetime" },
    ],
  },
  products: {
    label: "Products",
    icon: <ShoppingOutlined />,
    columns: [
      "id",
      "name",
      "price",
      "quantity",
      "type",
      "is_listed",
      "description",
      "created_at",
      "updated_at",
    ],
    fields: [
      { key: "name", label: "Name", type: "text", required: true },
      { key: "price", label: "Price", type: "number", required: true },
      { key: "quantity", label: "Quantity", type: "number", required: true },
      { key: "type", label: "Type", type: "text" },
      { key: "is_listed", label: "Is Listed", type: "boolean" },
      { key: "description", label: "Description", type: "textarea" },
      { key: "created_at", label: "Created At", type: "datetime" },
      { key: "updated_at", label: "Updated At", type: "datetime" },
    ],
  },
  decks: {
    label: "Decks",
    icon: <BookOutlined />,
    columns: [
      "id",
      "user_id",
      "title",
      "description",
      "is_public",
      "tags",
      "clone_count",
      "is_ai_generated",
      "created_at",
      "updated_at",
      "cards",
      "card_count",
    ],
    fields: [
      {
        key: "user_id",
        label: "User",
        type: "select",
        required: true,
        relation: "users",
      },
      { key: "title", label: "Title", type: "text", required: true },
      { key: "description", label: "Description", type: "textarea" },
      { key: "is_public", label: "Is Public", type: "boolean" },
      { key: "tags", label: "Tags", type: "text" },
      { key: "clone_count", label: "Clone Count", type: "number" },
      { key: "is_ai_generated", label: "AI Generated", type: "boolean" },
      { key: "created_at", label: "Created At", type: "datetime" },
      { key: "updated_at", label: "Updated At", type: "datetime" },
    ],
  },
  cards: {
    label: "Cards",
    icon: <FileTextOutlined />,
    columns: [
      "id",
      "deck_id",
      "front",
      "back",
      "front_image_url",
      "back_image_url",
      "interval",
      "repetitions",
      "ease_factor",
      "stability",
      "difficulty",
      "due_date",
      "is_ai_created",
      "created_at",
      "updated_at",
    ],
    fields: [
      {
        key: "deck_id",
        label: "Deck",
        type: "select",
        required: true,
        relation: "decks",
      },
      { key: "front", label: "Front", type: "textarea", required: true },
      { key: "back", label: "Back", type: "textarea", required: true },
      { key: "front_image_url", label: "Front Image URL", type: "text" },
      { key: "back_image_url", label: "Back Image URL", type: "text" },
      { key: "interval", label: "Interval", type: "number" },
      { key: "repetitions", label: "Repetitions", type: "number" },
      { key: "ease_factor", label: "Ease Factor", type: "number" },
      { key: "stability", label: "Stability", type: "number" },
      { key: "difficulty", label: "Difficulty", type: "number" },
      { key: "due_date", label: "Due Date", type: "datetime" },
      { key: "is_ai_created", label: "AI Created", type: "boolean" },
      { key: "created_at", label: "Created At", type: "datetime" },
      { key: "updated_at", label: "Updated At", type: "datetime" },
    ],
  },
  ai_generation_history: {
    label: "AI Generation History",
    icon: <RobotOutlined />,
    columns: ["id", "user_id", "text", "card_count", "created_at"],
    fields: [
      {
        key: "user_id",
        label: "User",
        type: "select",
        required: true,
        relation: "users",
      },
      { key: "text", label: "Text", type: "textarea", required: true },
      {
        key: "card_count",
        label: "Card Count",
        type: "number",
        required: true,
      },
      { key: "created_at", label: "Created At", type: "datetime" },
    ],
  },
  issues: {
    label: "Issues",
    icon: <AlertOutlined />,
    columns: [
      "id",
      "reporter_name",
      "reporter_email",
      "transaction_id",
      "category",
      "title",
      "description",
      "status",
      "priority",
      "admin_notes",
      "created_at",
      "updated_at",
    ],
    fields: [
      {
        key: "reporter_name",
        label: "Reporter Name",
        type: "text",
        required: true,
      },
      {
        key: "reporter_email",
        label: "Reporter Email",
        type: "text",
        required: true,
      },
      {
        key: "transaction_id",
        label: "Transaction ID",
        type: "text",
      },
      { key: "category", label: "Category", type: "text", required: true },
      { key: "title", label: "Title", type: "text", required: true },
      {
        key: "description",
        label: "Description",
        type: "textarea",
        required: true,
      },
      {
        key: "status",
        label: "Status",
        type: "select",
        options: [
          { label: "open", value: "open" },
          { label: "in_review", value: "in_review" },
          { label: "resolved", value: "resolved" },
          { label: "rejected", value: "rejected" },
        ],
      },
      {
        key: "priority",
        label: "Priority",
        type: "select",
        options: [
          { label: "low", value: "low" },
          { label: "medium", value: "medium" },
          { label: "high", value: "high" },
        ],
      },
      { key: "admin_notes", label: "Admin Notes", type: "textarea" },
      { key: "created_at", label: "Created At", type: "datetime" },
      { key: "updated_at", label: "Updated At", type: "datetime" },
    ],
  },
};

function App() {
  const [form] = Form.useForm();

  const [token, setToken] = useState<string>(
    () => localStorage.getItem(TOKEN_KEY) ?? "",
  );
  const [authChecking, setAuthChecking] = useState<boolean>(
    Boolean(localStorage.getItem(TOKEN_KEY)),
  );

  const [user, setUser] = useState<CurrentUser | null>(null);
  const [overview, setOverview] = useState<DashboardOverview | null>(null);
  const navigate = useNavigate();
  const location = useLocation();

  const [activeEntity, setActiveEntity] = useState<EntityKey>(() => {
    const fromPath = location.pathname.split("/")[1] as EntityKey | undefined;
    if (fromPath && Object.keys(ENTITY_CONFIG).includes(fromPath)) {
      return fromPath;
    }
    return "users";
  });
  const [rows, setRows] = useState<RecordData[]>([]);

  const [searchText, setSearchText] = useState("");
  const [debouncedSearchText, setDebouncedSearchText] = useState("");
  const [loading, setLoading] = useState(false);
  const [tableLoading, setTableLoading] = useState(false);
  const [error, setError] = useState("");

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [editingRow, setEditingRow] = useState<RecordData | null>(null);

  const [lookupUsers, setLookupUsers] = useState<LookupOption[]>([]);
  const [lookupProducts, setLookupProducts] = useState<LookupOption[]>([]);
  const [lookupDecks, setLookupDecks] = useState<LookupOption[]>([]);

  const [loginEmail, setLoginEmail] = useState("sekilas@kilas.my.id");
  const [loginPassword, setLoginPassword] = useState("admin12345");

  const authHeaders = useMemo(
    () => ({
      Authorization: `Bearer ${token}`,
    }),
    [token],
  );

  const jsonHeaders = useMemo(
    () => ({
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    }),
    [token],
  );

  useEffect(() => {
    const fromPath = location.pathname.split("/")[1] as EntityKey | undefined;
    if (
      fromPath &&
      Object.keys(ENTITY_CONFIG).includes(fromPath) &&
      fromPath !== activeEntity
    ) {
      setActiveEntity(fromPath);
    }
  }, [activeEntity, location.pathname]);

  const currentConfig = ENTITY_CONFIG[activeEntity];

  const fetchMe = useCallback(async () => {
    const res = await fetch(`${API_URL}/auth/me`, { headers: authHeaders });
    if (!res.ok) throw new Error("Session expired. Please login again.");
    const data = (await res.json()) as CurrentUser;
    setUser(data);
  }, [authHeaders]);

  const fetchOverview = useCallback(async () => {
    const res = await fetch(`${API_URL}/admin/summary`, {
      headers: authHeaders,
    });
    if (!res.ok) throw new Error("Failed to fetch dashboard summary");
    const data = (await res.json()) as DashboardOverview;
    setOverview(data);
  }, [authHeaders]);

  const fetchRows = useCallback(
    async (entity: EntityKey, q = "") => {
      setTableLoading(true);
      const qs = new URLSearchParams({ limit: "50" });
      if (q.trim()) qs.set("q", q.trim());
      const res = await fetch(`${API_URL}/admin/${entity}?${qs.toString()}`, {
        headers: authHeaders,
      });
      setTableLoading(false);
      if (!res.ok)
        throw new Error(
          `Failed to fetch ${ENTITY_CONFIG[entity].label.toLowerCase()}`,
        );
      const data = (await res.json()) as RecordData[];
      setRows(data);
    },
    [authHeaders],
  );

  const fetchLookupUsers = useCallback(
    async (q = "") => {
      const qs = new URLSearchParams({ limit: "50" });
      if (q.trim()) qs.set("q", q.trim());
      const res = await fetch(`${API_URL}/admin/users?${qs.toString()}`, {
        headers: authHeaders,
      });
      if (!res.ok) return;
      const users = (await res.json()) as Array<{
        id: number;
        username?: string;
        email?: string;
      }>;
      setLookupUsers(
        users.map((u) => ({
          value: u.id,
          label: `${u.username ?? "user"} (#${u.id}) - ${u.email ?? ""}`,
        })),
      );
    },
    [authHeaders],
  );

  const fetchLookupProducts = useCallback(
    async (q = "") => {
      const qs = new URLSearchParams({ limit: "50" });
      if (q.trim()) qs.set("q", q.trim());
      const res = await fetch(`${API_URL}/admin/products?${qs.toString()}`, {
        headers: authHeaders,
      });
      if (!res.ok) return;
      const products = (await res.json()) as Array<{
        id: number;
        name?: string;
      }>;
      setLookupProducts(
        products.map((p) => ({
          value: p.id,
          label: `${p.name ?? "product"} (#${p.id})`,
        })),
      );
    },
    [authHeaders],
  );

  const fetchLookupDecks = useCallback(
    async (q = "") => {
      const qs = new URLSearchParams({ limit: "50" });
      if (q.trim()) qs.set("q", q.trim());
      const res = await fetch(`${API_URL}/admin/decks?${qs.toString()}`, {
        headers: authHeaders,
      });
      if (!res.ok) return;
      const decks = (await res.json()) as Array<{ id: number; title?: string }>;
      setLookupDecks(
        decks.map((d) => ({
          value: d.id,
          label: `${d.title ?? "deck"} (#${d.id})`,
        })),
      );
    },
    [authHeaders],
  );

  useEffect(() => {
    if (!token) {
      setAuthChecking(false);
      return;
    }

    setAuthChecking(true);
    Promise.all([fetchMe(), fetchOverview()])
      .catch((e: Error) => {
        setError(e.message);
        handleLogout();
      })
      .finally(() => setAuthChecking(false));
  }, [fetchMe, fetchOverview, token]);

  useEffect(() => {
    const t = setTimeout(() => {
      setDebouncedSearchText(searchText);
    }, 250);
    return () => clearTimeout(t);
  }, [searchText]);

  useEffect(() => {
    if (!token || authChecking) return;
    fetchRows(activeEntity, debouncedSearchText).catch((e: Error) =>
      setError(e.message),
    );
  }, [activeEntity, authChecking, debouncedSearchText, fetchRows, token]);

  const filteredRows = useMemo(() => rows, [rows]);

  const tableColumns = useMemo<TableColumnsType<RecordData>>(() => {
    const baseColumns = currentConfig.columns.map((col) => {
      const uniqueValues = Array.from(
        new Set(rows.map((r) => String(r[col] ?? "-"))),
      ).slice(0, 20);

      return {
        title: prettify(col),
        dataIndex: col,
        key: col,
        ellipsis: true,
        sorter: (a: RecordData, b: RecordData) =>
          String(a[col] ?? "").localeCompare(String(b[col] ?? ""), undefined, {
            numeric: true,
          }),
        filters: uniqueValues.map((v) => ({ text: v, value: v })),
        onFilter: (value: Key | boolean, record: RecordData) =>
          String(record[col] ?? "") === String(value),
        render: (value: unknown) => renderCell(col, value),
      };
    });

    return [
      ...baseColumns,
      {
        title: "Actions",
        key: "actions",
        fixed: "right" as const,
        render: (_: unknown, row: RecordData) => (
          <Space>
            <Button size="small" onClick={() => openEditDrawer(row)}>
              Edit
            </Button>
            <Popconfirm
              title="Delete this record?"
              onConfirm={() => handleDelete(row)}
            >
              <Button size="small" danger>
                Delete
              </Button>
            </Popconfirm>
          </Space>
        ),
      },
    ];
  }, [currentConfig.columns, rows]);

  async function handleLogin(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const res = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: loginEmail, password: loginPassword }),
      });

      if (!res.ok) throw new Error("Invalid email or password");
      const data = (await res.json()) as {
        access_token: string;
        user: CurrentUser;
      };

      localStorage.setItem(TOKEN_KEY, data.access_token);
      setToken(data.access_token);
      setUser(data.user);
      message.success("Logged in");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  function handleLogout() {
    localStorage.removeItem(TOKEN_KEY);
    setToken("");
    setUser(null);
    setOverview(null);
    setRows([]);
    setDrawerOpen(false);
    setEditingRow(null);
    setSearchText("");
    setAuthChecking(false);
  }

  function openCreateDrawer() {
    setEditingRow(null);
    form.resetFields();
    const defaults: Record<string, unknown> = {};
    for (const field of currentConfig.fields) {
      if (field.type === "boolean") defaults[field.key] = false;
    }
    form.setFieldsValue(defaults);
    setDrawerOpen(true);
  }

  function openEditDrawer(row: RecordData) {
    setEditingRow(row);
    form.resetFields();
    form.setFieldsValue(normalizeFormValues(currentConfig.fields, row));
    setDrawerOpen(true);
  }

  async function handleSave(formInstance: FormInstance) {
    const values = await formInstance.validateFields();
    const payload = serializeFormPayload(currentConfig.fields, values);

    if (!editingRow?.id && activeEntity === "users" && !payload.password) {
      message.error("Password is required when creating users");
      return;
    }

    const url = editingRow?.id
      ? `${API_URL}/admin/${activeEntity}/${editingRow.id}`
      : `${API_URL}/admin/${activeEntity}`;

    const method = editingRow?.id ? "PUT" : "POST";
    setLoading(true);
    setError("");

    try {
      const res = await fetch(url, {
        method,
        headers: jsonHeaders,
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        const body = (await res.json()) as { error?: string };
        throw new Error(body.error ?? "Failed to save");
      }

      message.success(
        editingRow ? "Updated successfully" : "Created successfully",
      );
      setDrawerOpen(false);
      setEditingRow(null);

      await Promise.all([
        fetchRows(activeEntity, debouncedSearchText),
        fetchOverview(),
      ]);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Save failed");
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(row: RecordData) {
    if (!row.id) return;

    setLoading(true);
    setError("");
    try {
      const res = await fetch(`${API_URL}/admin/${activeEntity}/${row.id}`, {
        method: "DELETE",
        headers: jsonHeaders,
      });
      if (!res.ok) {
        const body = (await res.json()) as { error?: string };
        throw new Error(body.error ?? "Delete failed");
      }

      message.success("Deleted successfully");
      await Promise.all([
        fetchRows(activeEntity, debouncedSearchText),
        fetchOverview(),
      ]);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Delete failed");
    } finally {
      setLoading(false);
    }
  }

  function getRelationOptions(
    relation: FieldConfig["relation"],
  ): LookupOption[] {
    if (relation === "users") return lookupUsers;
    if (relation === "products") return lookupProducts;
    if (relation === "decks") return lookupDecks;
    return [];
  }

  async function handleLookupSearch(
    relation: FieldConfig["relation"],
    q: string,
  ) {
    if (relation === "users") {
      await fetchLookupUsers(q);
      return;
    }
    if (relation === "products") {
      await fetchLookupProducts(q);
      return;
    }
    if (relation === "decks") {
      await fetchLookupDecks(q);
    }
  }

  if (authChecking) {
    return (
      <main className="boot-shell">
        <Spin size="large" />
      </main>
    );
  }

  if (!token || !user) {
    return (
      <main className="login-shell">
        <Card className="login-card" bordered={false}>
          <Title level={3} style={{ marginTop: 0 }}>
            Kilas Admin CMS
          </Title>
          <Text type="secondary">Sign in to manage Kilas resources.</Text>
          <form className="login-form" onSubmit={handleLogin}>
            <Input
              value={loginEmail}
              onChange={(e) => setLoginEmail(e.target.value)}
              placeholder="Email"
              type="email"
              required
            />
            <Input.Password
              value={loginPassword}
              onChange={(e) => setLoginPassword(e.target.value)}
              placeholder="Password"
              required
            />
            <Button htmlType="submit" type="primary" loading={loading} block>
              Login
            </Button>
          </form>
          {error ? (
            <Alert style={{ marginTop: 12 }} type="error" message={error} />
          ) : null}
        </Card>
      </main>
    );
  }

  return (
    <Layout className="dashboard-layout">
      <Sider width={240} className="dashboard-sider">
        <div className="brand">
          <DatabaseOutlined />
          <span>Kilas Admin</span>
        </div>

        <Menu
          mode="inline"
          theme="dark"
          selectedKeys={[activeEntity]}
          items={(Object.keys(ENTITY_CONFIG) as EntityKey[]).map((entity) => ({
            key: entity,
            icon: ENTITY_CONFIG[entity].icon,
            label: ENTITY_CONFIG[entity].label,
            onClick: () => {
              navigate(`/${entity}`);
            },
          }))}
        />

        <Button
          className="logout-button"
          icon={<LogoutOutlined />}
          onClick={handleLogout}
        >
          Logout
        </Button>
      </Sider>

      <Layout>
        <Header className="dashboard-header" style={{ padding: "0 16px" }}>
          <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
            <Text type="secondary">Signed in as</Text>
            <div className="header-user">
              {user.name} ({user.email})
            </div>
          </div>
        </Header>

        <Content className="dashboard-content">
          {error ? (
            <Alert
              type="error"
              message={error}
              showIcon
              style={{ marginBottom: 12 }}
            />
          ) : null}

          <Row gutter={[12, 12]}>
            <Col span={6}>
              <Card>
                <Statistic title="Users" value={overview?.users ?? 0} />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic
                  title="Transactions"
                  value={overview?.transactions ?? 0}
                />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic title="Products" value={overview?.products ?? 0} />
              </Card>
            </Col>
            <Col span={6}>
              <Card>
                <Statistic title="Issues" value={overview?.issues ?? 0} />
              </Card>
            </Col>
          </Row>

          <Card style={{ marginTop: 12 }}>
            <div className="table-toolbar">
              <div>
                <Title level={4} style={{ margin: 0 }}>
                  {currentConfig.label}
                </Title>
                <Text type="secondary">You can search/filter/sort/edit</Text>
              </div>
              <Space>
                <Input.Search
                  allowClear
                  placeholder={`Search ${currentConfig.label.toLowerCase()} table`}
                  value={searchText}
                  onChange={(e) => setSearchText(e.target.value)}
                  style={{ width: 320 }}
                />
                <Button type="primary" onClick={openCreateDrawer}>
                  New
                </Button>
              </Space>
            </div>

            <Spin spinning={tableLoading}>
              <Table<RecordData>
                rowKey={(row) => String(row.id ?? JSON.stringify(row))}
                columns={tableColumns}
                dataSource={filteredRows}
                scroll={{ x: 1800 }}
                pagination={{ pageSize: 10, showSizeChanger: true }}
              />
            </Spin>
          </Card>
        </Content>
      </Layout>

      <Drawer
        open={drawerOpen}
        onClose={() => setDrawerOpen(false)}
        title={
          editingRow
            ? `Edit ${currentConfig.label}`
            : `Create ${currentConfig.label}`
        }
        width={620}
      >
        <Form form={form} layout="vertical" onFinish={() => handleSave(form)}>
          {currentConfig.fields.map((field) => {
            const required = Boolean(field.required);

            if (field.type === "boolean") {
              return (
                <Form.Item
                  key={field.key}
                  label={field.label}
                  name={field.key}
                  valuePropName="checked"
                >
                  <Switch />
                </Form.Item>
              );
            }

            if (field.type === "number") {
              return (
                <Form.Item
                  key={field.key}
                  label={field.label}
                  name={field.key}
                  rules={
                    required
                      ? [
                          {
                            required: true,
                            message: `${field.label} is required`,
                          },
                        ]
                      : undefined
                  }
                >
                  <InputNumber style={{ width: "100%" }} />
                </Form.Item>
              );
            }

            if (field.type === "datetime") {
              return (
                <Form.Item key={field.key} label={field.label} name={field.key}>
                  <DatePicker
                    showTime
                    style={{ width: "100%" }}
                    format="YYYY-MM-DD HH:mm:ss"
                  />
                </Form.Item>
              );
            }

            if (field.type === "select") {
              const options = field.relation
                ? getRelationOptions(field.relation)
                : (field.options ?? []);
              return (
                <Form.Item
                  key={field.key}
                  label={field.label}
                  name={field.key}
                  rules={
                    required
                      ? [
                          {
                            required: true,
                            message: `${field.label} is required`,
                          },
                        ]
                      : undefined
                  }
                >
                  <Select
                    showSearch
                    placeholder={`Select ${field.label}`}
                    options={options}
                    optionFilterProp="label"
                    filterOption={field.relation ? false : true}
                    onSearch={
                      field.relation
                        ? (q) => {
                            handleLookupSearch(field.relation, q).catch(
                              () => undefined,
                            );
                          }
                        : undefined
                    }
                    onDropdownVisibleChange={
                      field.relation
                        ? (open) => {
                            if (open) {
                              handleLookupSearch(field.relation, "").catch(
                                () => undefined,
                              );
                            }
                          }
                        : undefined
                    }
                    allowClear
                  />
                </Form.Item>
              );
            }

            if (field.type === "textarea") {
              return (
                <Form.Item
                  key={field.key}
                  label={field.label}
                  name={field.key}
                  rules={
                    required
                      ? [
                          {
                            required: true,
                            message: `${field.label} is required`,
                          },
                        ]
                      : undefined
                  }
                >
                  <Input.TextArea rows={4} />
                </Form.Item>
              );
            }

            return (
              <Form.Item
                key={field.key}
                label={field.label}
                name={field.key}
                rules={
                  required
                    ? [
                        {
                          required: true,
                          message: `${field.label} is required`,
                        },
                      ]
                    : undefined
                }
              >
                <Input />
              </Form.Item>
            );
          })}

          <Space>
            <Button type="primary" htmlType="submit" loading={loading}>
              {editingRow ? "Save changes" : "Create"}
            </Button>
            <Button onClick={() => setDrawerOpen(false)}>Cancel</Button>
          </Space>
        </Form>
      </Drawer>
    </Layout>
  );
}

function normalizeFormValues(
  fields: FieldConfig[],
  row: RecordData,
): Record<string, unknown> {
  const values: Record<string, unknown> = {};

  for (const field of fields) {
    const raw = row[field.key];
    if (field.type === "datetime") {
      values[field.key] = raw ? dayjs(String(raw)) : undefined;
    } else if (field.type === "boolean") {
      values[field.key] = Boolean(raw);
    } else {
      values[field.key] = raw;
    }
  }

  return values;
}

function serializeFormPayload(
  fields: FieldConfig[],
  values: Record<string, unknown>,
): Record<string, unknown> {
  const payload: Record<string, unknown> = {};

  for (const field of fields) {
    const raw = values[field.key];
    if (raw === undefined || raw === null || raw === "") continue;

    if (field.type === "datetime") {
      payload[field.key] = serializeDate(raw);
      continue;
    }

    payload[field.key] = raw;
  }

  return payload;
}

function serializeDate(value: unknown): string {
  if (dayjs.isDayjs(value)) {
    return (value as Dayjs).format("YYYY-MM-DD HH:mm:ss");
  }

  const parsed = dayjs(String(value));
  return parsed.isValid()
    ? parsed.format("YYYY-MM-DD HH:mm:ss")
    : String(value);
}

function prettify(key: string): string {
  return key.replaceAll("_", " ").replace(/\b\w/g, (ch) => ch.toUpperCase());
}

function renderCell(column: string, value: unknown) {
  if (value === null || value === undefined || value === "")
    return <Text type="secondary">-</Text>;

  if (column === "status") {
    const status = String(value);
    const color =
      status === "resolved"
        ? "green"
        : status === "rejected"
          ? "red"
          : status === "in_review"
            ? "orange"
            : "blue";
    return <Tag color={color}>{status}</Tag>;
  }

  if (typeof value === "boolean") {
    return <Tag color={value ? "green" : "default"}>{String(value)}</Tag>;
  }

  if (typeof value === "object") {
    const stringified = JSON.stringify(value);
    return stringified.length > 100
      ? `${stringified.slice(0, 100)}…`
      : stringified;
  }

  const text = String(value);
  return text.length > 100 ? `${text.slice(0, 100)}…` : text;
}

export default App;
