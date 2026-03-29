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
      navigate("/");
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <form
      onSubmit={handleSubmit}
      class="max-w-md mx-auto bg-white rounded-lg shadow-md p-6 space-y-4"
    >
      <input
        type="email"
        value={email()}
        onInput={(e) => setEmail(e.currentTarget.value)}
        placeholder="Email"
        required
      />
      <input
        type="password"
        value={password()}
        onInput={(e) => setPassword(e.currentTarget.value)}
        placeholder="Пароль"
        required
      />
      {error() && <div class="error">{error()}</div>}
      <button type="submit" disabled={loading()}>
        {loading() ? "Вход..." : "Войти"}
      </button>
    </form>
  );
};

export default Login;
