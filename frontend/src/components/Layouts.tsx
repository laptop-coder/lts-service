import { type Component, type JSX } from "solid-js";
import { useNavigate } from "@solidjs/router";
import { useAuth } from "../lib/auth";

interface Props {
  children?: JSX.Element;
}

export const PublicRoute: Component<Props> = (props) => {
  const navigate = useNavigate();
  const auth = useAuth();

  const handleLogout = async () => {
    await auth.logout();
    navigate("/login");
  };

  return (
    <div class="min-h-screen bg-gray-100">
      <header class="bg-white shadow">
        <div class="container mx-auto px-4 py-4 flex justify-between items-center">
          <h1 class="text-xl font-bold">Сервис поиска потерянных вещей</h1>
          <div class="flex items-center gap-4">
            <span>
              {auth.user()?.firstName} {auth.user()?.lastName}
            </span>
            <button
              onClick={handleLogout}
              class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600"
            >
              Выйти
            </button>
          </div>
        </div>
      </header>
      <main class="container mx-auto px-4 py-8">{props?.children}</main>
    </div>
  );
};

export const ProtectedRoute: Component<Props> = (props) => {
  const navigate = useNavigate();
  const auth = useAuth();

  return (
    <>
      {auth.isLoading() ? (
        <div class="flex justify-center items-center h-screen">Загрузка...</div>
      ) : auth.isAuthenticated() ? (
        <PublicRoute>{props?.children}</PublicRoute>
      ) : (
        navigate("/login")
      )}
    </>
  );
};
