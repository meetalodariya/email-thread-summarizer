import {
  BrowserRouter,
  Navigate,
  Outlet,
  Route,
  Routes,
  useLocation,
} from "react-router";
import { Toaster } from "@/components/ui/toaster";
import { Inbox } from "@/pages/inbox";
import { Register } from "@/pages/Register";
import { useAuth } from "@/providers/auth";
import { AuthProvider } from "@/providers/auth/provider";

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/auth" element={<Register />} />
          <Route
            path="/auth/google/callback"
            element={<Register isCallback={true} />}
          />

          <Route element={<RequireAuth />}>
            <Route path="/" element={<Navigate to="/inbox" replace />} />
            <Route index path="/inbox" element={<Inbox />} />
          </Route>
        </Routes>
        <Toaster />
      </AuthProvider>
    </BrowserRouter>
  );
}

function RequireAuth() {
  const auth = useAuth();
  const location = useLocation();

  if (!auth?.user) {
    return <Navigate to="/auth" state={{ from: location }} replace />;
  }

  return <Outlet />;
}

export default App;
