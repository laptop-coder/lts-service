import { createSignal } from "solid-js";
import { useAuth } from "../lib/auth";
import { useNavigate } from "@solidjs/router";

const Login = () => {
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const auth = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await auth.login(email(), password());
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Не удалось войти в аккаунт",
      );
    } finally {
      setLoading(false);
      navigate("/");
    }
  };

  return (
    <div class="min-h-[80vh] flex items-center justify-center px-4">
      <div class="w-full max-w-md">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-gray-800">Вход в аккаунт</h1>
          <p class="text-gray-500 mt-2">Добро пожаловать!</p>
        </div>

        <form
          onSubmit={handleSubmit}
          class="bg-white rounded-2xl shadow-lg p-6 space-y-5"
        >
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              value={email()}
              onInput={(e) => setEmail(e.currentTarget.value)}
              placeholder="email@example.ru"
              class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
              required
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              Пароль
            </label>
            <input
              type="password"
              value={password()}
              onInput={(e) => setPassword(e.currentTarget.value)}
              placeholder="••••••••"
              class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition"
              required
            />
          </div>

          {error() && (
            <div class="bg-red-50 text-red-600 p-3 rounded-xl text-sm border border-red-200">
              {error()}
            </div>
          )}

          <div class="flex gap-3 pt-2">
            <button
              type="submit"
              disabled={loading()}
              class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
            >
              {loading() ? "Вход..." : "Войти"}
            </button>
          </div>

          <p class="text-center text-sm text-gray-500 mt-4">
            Нет аккаунта?{" "}
            <a
              href="/register"
              class="text-blue-600 hover:text-blue-700 hover:underline"
            >
              Запросить пригласительную ссылку ученика {/*TODO*/}
            </a>
          </p>
        </form>
      </div>
    </div>
  );
};

export default Login;
