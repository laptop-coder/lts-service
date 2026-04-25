import { createSignal, createEffect } from "solid-js";
import { api } from "../lib/api";

const RequestStudentInvite = () => {
  const [email, setEmail] = createSignal("");
  const [error, setError] = createSignal("");
  const [loading, setLoading] = createSignal(false);
  const [done, setDone] = createSignal(false); // email was sent


  let emailInputRef: HTMLInputElement | undefined;
  const focusEmailInput = () => {
    if (emailInputRef) {
      emailInputRef.focus();
    }
  };

  createEffect(() => {
    focusEmailInput();
  });

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    if (!email().trim()) return;
    setError("");
    setLoading(true);

    const formData = new URLSearchParams();
    formData.append("email", email());

    try {
      await api.post<{}>("/invite/request/student", formData);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Ошибка запроса пригласительной ссылки",
      );
    } finally {
      setLoading(false);
      setDone(true);
    }
  };

  return (
    <div class="min-h-screen py-8 px-4">
      <div class="max-w-2xl mx-auto">
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-gray-800">
            Создание аккаунта ученика
          </h1>
          <p class="text-gray-500 mt-2">Запрос пригласительной ссылки</p>
        </div>

        {loading() ? (
          <div class="text-center py-12 text-gray-500">Загрузка...</div>
        ) : done() ? (
          <div class="bg-green-100 border border-green-400 text-green-800 px-4 py-3 rounded-xl flex items-center gap-3">
            <svg
              class="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
              ></path>
            </svg>
            <span>
              Письмо отправлено на{" "}
              <span class="font-semibold underline">{email()}</span>! Перейдите
              по ссылке в письме для завершения регистрации.
            </span>
          </div>
        ) : (
          <form
            onSubmit={handleSubmit}
            class="bg-white rounded-2xl shadow-lg p-6 space-y-5"
          >
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Email *
              </label>
              <input
              ref={emailInputRef}
                type="email"
                value={email()}
                placeholder="email@example.ru"
                onInput={(e) => setEmail(e.currentTarget.value)}
                class="w-full px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition disabled:cursor-not-allowed"
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
                {loading() ? "Отправка..." : "Отправить"}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default RequestStudentInvite;
