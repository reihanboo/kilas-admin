import { Suspense, lazy } from "react";
import { Spin } from "antd";

const DashboardApp = lazy(() => import("./DashboardApp"));

function App() {
  return (
    <Suspense
      fallback={
        <main className="boot-shell">
          <Spin size="large" />
        </main>
      }
    >
      <DashboardApp />
    </Suspense>
  );
}

export default App;
